package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/appuio/promtool-jsonnet/pkg/promjsonnet"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var (
	// these variables are populated by Goreleaser when releasing
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")

	appName     = "promtool-jsonnet"
	appLongName = "Run Jsonnet files against promtool test"

	// envPrefix is the global prefix to use for the keys in environment variables
	envPrefix = "PJ"
)

type commandArgs struct {
	testFilePath        string
	promtoolPath        string
	additionalYamlFiles cli.StringSlice
	jsonnetPath         cli.StringSlice
}

func main() {
	ctx, stop, app := newApp()
	defer stop()
	err := app.RunContext(ctx, os.Args)
	// If required flags aren't set, it will exit before we could set up logging
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func newApp() (context.Context, context.CancelFunc, *cli.App) {
	command := &commandArgs{}
	app := &cli.App{
		Name:     appName,
		Usage:    appLongName,
		Version:  fmt.Sprintf("%s, revision=%s, date=%s", version, commit, date),
		Compiled: compilationDate(),

		EnableBashCompletion: true,

		Action: command.execute,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "test-file", Usage: "Path of the jsonnet test file from which to test queries",
				EnvVars: envVars("TEST_FILE"), Destination: &command.testFilePath, Value: "./test.jsonnet"},
			&cli.StringFlag{Name: "promtool-path", Usage: "Path of the promtool binary",
				EnvVars: envVars("PROMTOOL_PATH"), Destination: &command.promtoolPath, Value: "./.cache/prometheus/promtool"},
			&cli.StringSliceFlag{Name: "add-yaml-file", Usage: "Additional yaml files to include into the test (available to jsonnet as extVar)",
				EnvVars: envVars("ADD_YAML_FILE"), Destination: &command.additionalYamlFiles},
			&cli.StringSliceFlag{Name: "jsonnet-path", Usage: "Additional library paths for jsonnet",
				EnvVars: envVars("JSONNET_PATH"), Destination: &command.jsonnetPath},
		},
		ExitErrHandler: func(context *cli.Context, err error) {
			if err == nil {
				return
			}
			// Don't show stack trace if the error is expected (someone called cli.Exit())
			var exitErr cli.ExitCoder
			if errors.As(err, &exitErr) {
				cli.HandleExitCoder(err)
				return
			}
			cli.OsExiter(1)
		},
	}
	// There is logr.NewContext(...) which returns a context that carries the logger instance.
	// However, since we are configuring and replacing this logger after starting up and parsing the flags,
	// we'll store a thread-safe atomic reference.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	return ctx, stop, app
}

// env combines envPrefix with given suffix delimited by underscore.
func env(suffix string) string {
	if envPrefix == "" {
		return suffix
	}
	return envPrefix + "_" + suffix
}

// envVars combines envPrefix with each given suffix delimited by underscore.
func envVars(suffixes ...string) []string {
	arr := make([]string, len(suffixes))
	for i := range suffixes {
		arr[i] = env(suffixes[i])
	}
	return arr
}

func compilationDate() time.Time {
	compiled, err := time.Parse(time.RFC3339, date)
	if err != nil {
		// an empty Time{} causes cli.App to guess it from binary's file timestamp.
		return time.Time{}
	}
	return compiled
}

func (cmd *commandArgs) execute(cliCtx *cli.Context) error {
	log.Print("Running test for file ", cmd.testFilePath)
	extVars, err := cmd.buildExtVars()
	if err != nil {
		log.Print(err, "Query test setup failed")
	}
	err = promjsonnet.RunTestQueries(cmd.testFilePath, cmd.promtoolPath, extVars, cmd.jsonnetPath.Value())

	if err != nil {
		log.Print(err, "Query test failed")
	}
	log.Print("Done")
	return err
}

func (cmd *commandArgs) buildExtVars() (*map[string]string, error) {
	extVars := make(map[string]string)
	for file := range cmd.additionalYamlFiles.Value() {
		path := cmd.additionalYamlFiles.Value()[file]
		name := filepath.Base(path)
		fileHandle, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer fileHandle.Close()
		fileContent, err := io.ReadAll(fileHandle)
		if err != nil {
			return nil, err
		}

		var parsed map[string]interface{}
		err = yaml.Unmarshal(fileContent, &parsed)
		if err != nil {
			return nil, err
		}
		jsonStr, err := json.Marshal(parsed)
		if err != nil {
			return nil, err
		}
		extVars[name] = string(jsonStr)
	}
	return &extVars, nil
}
