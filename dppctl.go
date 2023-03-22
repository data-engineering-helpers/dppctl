//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/dppctl.go
//
package main

import (
	"flag"
	"log"
	"os"
    "gopkg.in/yaml.v3"
	
	"github.com/data-engineering-helpers/dppctl/utilities"
	"github.com/data-engineering-helpers/dppctl/workflow"
)

const AppVersion = "0.0.1-alpha.1"

var (
	versionFlag bool
	specFilepath string
	command string
)

func init() {
	flag.BoolVar(&versionFlag, "v", false, "Shows the current version")

	flag.StringVar(&specFilepath, "f",  "depl/aws-dev-sample.yaml",
		"The `name` of the deployment YAML specification file.")

	flag.StringVar(&command, "c",  "check",
		"The command to perform.")
}

func main() {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("[dppctl] ")
	log.SetFlags(0)

	//
	flag.Parse()
	if versionFlag {
      log.Println(AppVersion)
      os.Exit(0)
    }

	// Specification of the deployment
	deplSpec, err := utilities.ReadSpecFile(specFilepath)
	if err != nil {
		log.Print(err)
	}
	log.Println("Parsed spec file: ", deplSpec)

	deplSpecStruct, err := yaml.Marshal(&deplSpec)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println("Spec file dump:\n", string(deplSpecStruct))

	//
	switch command {
	case "check":
		workflow.Check(deplSpec)
	}
}


