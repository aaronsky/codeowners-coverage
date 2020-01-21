package main

import (
	"fmt"

	coverage "github.com/aaronsky/codeowners-coverage"
	"github.com/urfave/cli/v2"
)

// app is the static configuration of the primary CLI
var app = cli.App{
	Name:      "codeowners-coverage",
	Usage:     "Return codeowners coverage report for a repository",
	ArgsUsage: "[path to repository]",
	Action:    executeCommand,
}

// arguments is a type that describes the simple arguments for this CLI
type arguments struct {
	Path string
}

// newArguments constructs an Arguments object from a cli.Context Args object
func newArguments(args cli.Args) (*arguments, error) {
	path := args.First()
	if path == "" {
		return nil, fmt.Errorf("no path was supplied")
	}
	return &arguments{
		Path: path,
	}, nil
}

// executeCommand is the action handler for `app` and is executed by the CLI
func executeCommand(c *cli.Context) error {
	args, err := newArguments(c.Args())
	if err != nil {
		return err
	}

	report, err := coverage.NewCoverageReport(args.Path)
	if err != nil {
		return err
	}

	json, err := report.ToFormat(coverage.ReportFormatJSON)
	if err != nil {
		return err
	}

	fmt.Println(json)

	return nil
}
