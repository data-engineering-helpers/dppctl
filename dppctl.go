//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/dppctl.go
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
    "gopkg.in/yaml.v3"
	
	"github.com/data-engineering-helpers/dppctl/utilities"
	"github.com/data-engineering-helpers/dppctl/service"
)

const AppVersion = "0.0.1-alpha.1"

var (
	bucketName string
	versionFlag bool
	specFilepath string
)

func init() {
	flag.StringVar(&bucketName, "bucket", "baldwins",
		"The `name` of the S3 bucket to list item from.")

	flag.StringVar(&specFilepath, "f",  "depl/aws-dev-sample.yaml",
		"The `name` of the deployment YAML specification file.")

	flag.BoolVar(&versionFlag, "v", false, "Shows the current version")
}

func main() {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("[test-dppctl] ")
	log.SetFlags(0)

	//
	flag.Parse()
	if versionFlag {
      log.Println(AppVersion)
      os.Exit(0)
    }

	// STS - Caller identity (IAM)
	stsStruct, err := service.AWSGetCallerIdentity()
	if err != nil {
		log.Print(err)
	}
	log.Println("AWS IAM/caller identity:")
	log.Println(stsStruct)


	// Specification of the deployment
	depl_spec, err := utilities.ReadSpecFile(specFilepath)
	if err != nil {
		log.Print(err)
	}
	log.Println("Parsed spec file: ", depl_spec)

	depl_spec_struct, err := yaml.Marshal(&depl_spec)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println("Spec file dump:\n", string(depl_spec_struct))

	//
	bucketName := depl_spec.StorageContainer.Name
	bucketPrefix := depl_spec.StorageContainer.Prefix
	
	// AWS S3
	object_list, err := service.AWSS3List(bucketName, bucketPrefix)
	if err != nil {
		log.Print(err)
	}

	log.Println("List of objects within the following bucket:", bucketName)
	for _, object_metadata := range object_list {
		log.Println(object_metadata)
	}

	// CodeArtifact
	domain_list, err := service.AWSCodeArtifact()
	if err != nil {
		log.Print(err)
	}
	log.Println("List of domains within the CodeArtifact service:")
	for _, domain := range domain_list {
		log.Println(domain)
	}

	// MWAA/Airflow
	mwaaEnv := depl_spec.Airflow.Domain
	webServerHostname, cliToken, _,
		err := service.AWSAirflowCreateLoginToken(mwaaEnv)
	if err != nil {
		log.Print(err)
	}
	mwaaStr := fmt.Sprintf("hostname=%s token=%s", webServerHostname, cliToken)
	log.Println("MWAA/Airflow CLI token created: ", mwaaStr)

	//
	//command := "version"
	command := "dags list -o json"
	stdoutStr, err := service.AWSAirflowCLI(webServerHostname, cliToken,
		command)
	if err != nil {
		log.Print(err)
	}
	//log.Println("MWAA/Airflow CLI response: ", stdoutStr)

	// Parse the output when the command is "dags list -o json"
	mwaaDagMetadata, err := utilities.ParseAWSMWAADagListOutput(stdoutStr)
	if err != nil {
		log.Print(err)
	}
	log.Println("MWAA/Airflow DAGs: ", mwaaDagMetadata)
}


