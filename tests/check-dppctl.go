// File: https://github.com/data-engineering-helpers/dppctl/blob/main/tests/check-dppctl.go

package main

import (
  "context"
  "fmt"
  "log"

  "github.com/data-engineering-helpers/dppctl/service"

  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
  // Set properties of the predefined Logger, including
  // the log entry prefix and a flag to disable printing
  // the time, source file, and line number.
  log.SetPrefix("[test-dppctl] ")
  log.SetFlags(0)

  // AWS S3
  file_list, err := service.S3List("baldwins")
  if err != nil {
    log.Print(err)
  }

  for _, file_metadata := range file_list {
    log.Println(file_metadata)
  }

    // Using the SDK's default configuration, loading additional config
    // and credentials values from the environment variables, shared
    // credentials, and shared configuration files
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-1"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    // Using the Config value, create the DynamoDB client
    svc := dynamodb.NewFromConfig(cfg)

    // Build the request with its input parameters
    resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
        Limit: aws.Int32(5),
    })
    if err != nil {
        log.Fatalf("failed to list tables, %v", err)
    }

    fmt.Println("Tables:")
    for _, tableName := range resp.TableNames {
        fmt.Println(tableName)
    }
}


