package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aaronsky/codeowners-coverage/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "codeowners-coverage",
		Usage: "Return codeowners coverage report for a repository",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "format",
				Usage: "Report format",
				Value: "json",
			},
			&cli.StringFlag{
				Name:     "token",
				Usage:    "GitHub token",
				Required: true,
				EnvVars:  []string{"GITHUB_TOKEN"},
			},
		},
		Action: func(c *cli.Context) error {
			remote := c.Args().First()
			token := c.String("token")
			if token == "" {
				return fmt.Errorf("A valid Github token is required")
			}

			format := c.String("format")
			// Currently only the 'json' format is supported
			if format != "json" {
				return fmt.Errorf("Only JSON format is supported for now")
			}

			report, err := internal.NewCoverageReport(remote, token)
			if err != nil {
				return err
			}
			json, err := report.ToFormat(format)
			if err != nil {
				return err
			}

			fmt.Println(json)

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
