//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/service/aws.go
//
package service

import (
	"context"
	"fmt"
	"time"
	"log"
	"flag"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	awsConfig aws.Config
	region    string
)

func init() {
	flag.StringVar(
		&region, "region", "eu-west-1",	"The `region` of the AWS project.")

    // Using the SDK's default configuration, loading additional config
    // and credentials values from the environment variables, shared
    // credentials, and shared configuration files
    cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }
	awsConfig = cfg
}

func AWSGetCallerIdentity() (string, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	client := sts.NewFromConfig(awsConfig)
	params := &sts.GetCallerIdentityInput{}
	resp, err := client.GetCallerIdentity(ctx, params)
	if err != nil {
		log.Fatalf("expect no error, got %v", err)
	}

	sts_identity := fmt.Sprintf("UserId=%s Account=%s Arn=%s",
		aws.ToString(resp.UserId),
		aws.ToString(resp.Account),
		aws.ToString(resp.Arn))

    //
    return sts_identity, nil
}

func AWSDynamodDB() ([]string, error) {
    messages := []string {}

    // Using the Config value, create the DynamoDB client
    svc := dynamodb.NewFromConfig(awsConfig)

    // Build the request with its input parameters
	params := &dynamodb.ListTablesInput{
        Limit: aws.Int32(5),
    }
    resp, err := svc.ListTables(context.TODO(), params)
    if err != nil {
        log.Fatalf("failed to list tables, %v", err)
    }

    for _, tableName := range resp.TableNames {
		messages = append(messages, tableName)
    }

    //
    return messages, nil
}

func AWSCodeArtifact() ([]string, error) {
	// References:
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/codeartifact/types/types.go
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/codeartifact/api_op_ListDomains.go
	
    messages := []string {}

    // Using the Config value, create the CodeArtifact client
    svc := codeartifact.NewFromConfig(awsConfig)

    // Build the request with its input parameters
	params := &codeartifact.ListDomainsInput{}
    resp, err := svc.ListDomains(context.TODO(), params)
    if err != nil {
        log.Fatalf("failed to list domains, %v", err)
    }

	// https://docs.aws.amazon.com/codeartifact/latest/APIReference/API_DomainSummary.html
    for _, domain := range resp.Domains {
		message := fmt.Sprintf("Name=%s Status=%s",
			aws.ToString(domain.Name), domain.Status)
		messages = append(messages, message)
    }

    //
    return messages, nil
}

func AWSS3List(bucketName string) ([]string, error) {
    messages := []string {}

    //
    if bucketName == "" {
        return messages, errors.New("empty bucket name")
    }

    // Create an Amazon S3 service client
    svc := s3.NewFromConfig(awsConfig)

    // Get the first page of results for ListObjectsV2 for a bucket
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}
    output, err := svc.ListObjectsV2(context.TODO(), params)
    if err != nil {
		log.Fatal(err)
    }

    for _, object := range output.Contents {
		message := fmt.Sprintf("Key=%s size=%d",
			aws.ToString(object.Key), object.Size)
		messages = append(messages, message)
    }

    //
    return messages, nil
}

