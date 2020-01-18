package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aaronsky/codeowners-coverage/internal"
	"github.com/urfave/cli/v2"
)

var app = cli.App{
	Name:  "codeowners-coverage",
	Usage: "Return codeowners coverage report for a repository",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "format",
			Usage: "Report format",
			Value: "json",
		},
	},
	ArgsUsage: "[path to repository]",
	Action:    executeCommand,
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func executeCommand(c *cli.Context) error {
	path := c.Args().First()

	format := c.String("format")
	// Currently only the 'json' format is supported
	if format != "json" {
		return fmt.Errorf("Only JSON format is supported for now")
	}

	report, err := internal.NewCoverageReport(path)
	if err != nil {
		return err
	}
	json, err := report.ToFormat(internal.ReportFormatJSON)
	if err != nil {
		return err
	}

	fmt.Println(json)

	return nil
}
