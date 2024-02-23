package promjsonnet

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
)

func RunTestQueries(filepath string, promtoolPath string, extcodes *map[string]string, jpaths []string) error {
	tmp, err := renderJsonnet(filepath, extcodes, &jpaths)
	if err != nil {
		return err
	}
	return runPromtool(tmp, promtoolPath)
}

func runPromtool(tmp string, promtoolPath string) error {
	cmd := exec.Command(promtoolPath, "test", "rules", tmp)
	var stderr, stdout strings.Builder
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	// Not using t.Log to keep formatting sane
	fmt.Println("STDOUT")
	fmt.Println(stdout.String())
	fmt.Println("STDERR")
	fmt.Println(stderr.String())
	return err
}

func renderJsonnet(tFile string, extcodes *map[string]string, jpaths *[]string) (string, error) {
	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: *jpaths,
	})

	for key := range *extcodes {
		vm.ExtCode(key, (*extcodes)[key])
	}

	ev, err := vm.EvaluateFile(tFile)
	if err != nil {
		return "", err
	}

	filename := filepath.Base(tFile)

	tmp := path.Join("/tmp", fmt.Sprintf("%s.json", filename))
	err = os.WriteFile(tmp, []byte(ev), 0644)
	return tmp, err
}
