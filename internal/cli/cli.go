// Package cli contains the brains of the primary cmd of this package
package cli

import (
	"fmt"
	"log"
	"os"

	coverage "github.com/aaronsky/codeowners-coverage"
	"github.com/urfave/cli/v2"
)

var app = cli.App{
	Name:      "codeowners-coverage",
	Usage:     "Return codeowners coverage report for a repository",
	ArgsUsage: "[path to repository]",
	Action:    executeCommand,
}

// CLI runs the CLI app
func CLI() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Arguments is a type that describes the simple arguments for this CLI
type Arguments struct {
	Path string
}

func getArguments(c *cli.Context) (*Arguments, error) {
	path := c.Args().First()
	if path == "" {
		return nil, fmt.Errorf("no path was supplied")
	}
	return &Arguments{
		Path: path,
	}, nil
}

func executeCommand(c *cli.Context) error {
	args, err := getArguments(c)
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
