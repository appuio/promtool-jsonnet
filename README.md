# promtool-jsonnet

[![Build](https://img.shields.io/github/workflow/status/appuio/promtool-jsonnet/Test)][build]
![Go version](https://img.shields.io/github/go-mod/go-version/appuio/promtool-jsonnet)
[![Version](https://img.shields.io/github/v/release/appuio/promtool-jsonnet)][releases]
[![GitHub downloads](https://img.shields.io/github/downloads/appuio/promtool-jsonnet/total)][releases]

[build]: https://github.com/appuio/promtool-jsonnet/actions?query=workflow%3ATest
[releases]: https://github.com/appuio/promtool-jsonnet/releases

## Usage

### Run tests

```sh
export PJ_JSONNET_PATH=~"`pwd`/jsonnet/"
make ensure-prometheus
go run . --test-file ~/path/to/your/promtool/tests.jsonnet  --add-yaml-file ~/path/to/supplemental/yaml/file.yml
```
