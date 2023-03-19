//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/tests/check-dppctl.go
//
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
	
	"github.com/data-engineering-helpers/dppctl/service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
)

var (
	awsConfig aws.Config
	bucketName string
	region    string
)

func init() {
	flag.StringVar(&bucketName, "bucket", "baldwins", "The `name` of the S3 bucket to list item from.")
	flag.StringVar(&region, "region", "eu-west-1", "The `region` of your AWS project.")

    // Using the SDK's default configuration, loading additional config
    // and credentials values from the environment variables, shared
    // credentials, and shared configuration files
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }
	awsConfig = cfg
}

func AWSGetCallerIdentity() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	client := sts.NewFromConfig(awsConfig)
	params := &sts.GetCallerIdentityInput{}
	resp, err := client.GetCallerIdentity(ctx, params)
	if err != nil {
		log.Fatalf("expect no error, got %v", err)
	}

    log.Println("AWS IAM role/caller identity:")
	sts_identity := fmt.Sprintf("UserId=%s Account=%s Arn=%s", aws.ToString(resp.UserId), aws.ToString(resp.Account), aws.ToString(resp.Arn))
	log.Println(sts_identity)
}

func AWSDynamodDB() {
    // Using the Config value, create the DynamoDB client
    svc := dynamodb.NewFromConfig(awsConfig)

    // Build the request with its input parameters
    resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
        Limit: aws.Int32(5),
    })
    if err != nil {
        log.Fatalf("failed to list tables, %v", err)
    }

    log.Println("Tables:")
    for _, tableName := range resp.TableNames {
        log.Println(tableName)
    }
}

func AWSCodeArtifact() {
	// References:
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/codeartifact/types/types.go
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/codeartifact/api_op_ListDomains.go
	
    // Using the Config value, create the CodeArtifact client
    svc := codeartifact.NewFromConfig(awsConfig)

    // Build the request with its input parameters
    resp, err := svc.ListDomains(context.TODO(), &codeartifact.ListDomainsInput{
    })
    if err != nil {
        log.Fatalf("failed to list domains, %v", err)
    }

	//
    log.Println("Domains:")
    for _, domain := range resp.Domains {
		// https://docs.aws.amazon.com/codeartifact/latest/APIReference/API_DomainSummary.html
		domainStr := fmt.Sprintf("Name=%s Status=%s", aws.ToString(domain.Name), domain.Status)
        log.Println(domainStr)
    }
}

func main() {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("[test-dppctl] ")
	log.SetFlags(0)

	//
	flag.Parse()
	
	// STS - Caller identity (IAM)
	AWSGetCallerIdentity()

	// AWS S3
	file_list, err := service.S3List(bucketName)
	if err != nil {
		log.Print(err)
	}

	for _, file_metadata := range file_list {
		log.Println(file_metadata)
	}

	// DynamoDB
	AWSDynamodDB()

	// CodeArtifact
	AWSCodeArtifact()
}


