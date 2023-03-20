//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/tests/check-dppctl.go
//
package main

import (
	"flag"
	"log"
	"os"
	
	"github.com/data-engineering-helpers/dppctl/service"
)

const AppVersion = "0.0.1-alpha.1"

var (
	bucketName string
	versionFlag bool
)

func init() {
	flag.StringVar(&bucketName, "bucket", "baldwins",
		"The `name` of the S3 bucket to list item from.")

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

	// AWS S3
	object_list, err := service.AWSS3List(bucketName)
	if err != nil {
		log.Print(err)
	}

	log.Println("List of objects within the following bucket:", bucketName)
	for _, object_metadata := range object_list {
		log.Println(object_metadata)
	}

	// DynamoDB
	table_list, err := service.AWSDynamodDB()
	log.Println("List of tables in the DynamoDB service:")
	for _, table := range table_list {
		log.Println(table)
	}

	// CodeArtifact
	domain_list, err := service.AWSCodeArtifact()
	log.Println("List of domains within the CodeArtifact service:")
	for _, domain := range domain_list {
		log.Println(domain)
	}

}


