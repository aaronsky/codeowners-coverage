# codeowners-coverage

[![Actions Status](https://github.com/aaronsky/codeowners-coverage/workflows/Run%20Tests/badge.svg?branch=master)](https://github.com/wayfair/aaronsky/codeowners-coverage/actions) [![codecov](https://codecov.io/gh/aaronsky/codeowners-coverage/branch/master/graph/badge.svg)](https://codecov.io/gh/aaronsky/codeowners-coverage) [![Go Report Card](https://goreportcard.com/badge/github.com/aaronsky/codeowners-coverage)](https://goreportcard.com/report/github.com/aaronsky/codeowners-coverage) [![GoDoc](https://godoc.org/github.com/aaronsky/codeowners-coverage?status.svg)](https://godoc.org/github.com/aaronsky/codeowners-coverage)

## Installation

```
go get github.com/aaronsky/codeowners-coverage/cmd/codeowners-coverage
```

## Usage

### As a Package 

```go
import (
    "fmt"
    "github.com/aaronsky/codeowners-coverage"
)

func getMyReport(path string) error {
	report, err := coverage.NewCoverageReport(path)
	if err != nil {
		return err
    }

	jsonString, err := report.ToFormat(coverage.ReportFormatJSON)
	if err != nil {
		return err
    }

    fmt.Println(jsonString)
}
```

### CLI

`codeowners-coverage` also has a CLI. It works by loading a local Git repository, parsing its CODEOWNERS file, and crawling the disk for matches. To run, simply provide a path to a Git repository.

```
codeowners-coverage ~/go/src/github.com/docker/compose
```

In the event of a successful navigation, this will print JSON to stdout describing the coverage attributes of the repository. 

## License

This package is licensed under the [MIT License](./LICENSE).
