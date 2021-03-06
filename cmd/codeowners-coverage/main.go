// Package main is the command line interface to package coverage
package main

import (
	"log"
	"os"
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
