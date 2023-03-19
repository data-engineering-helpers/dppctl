// File: https://github.com/data-engineering-helpers/dppctl/blob/main/service/aws.go

package service

import (
	"context"
	"fmt"
	"log"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3List(bucket_name string) ([]string, error) {
    messages := []string {}

    //
    if bucket_name == "" {
        return messages, errors.New("empty bucket name")
    }

    // Load the Shared AWS Configuration (~/.aws/config)
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
	log.Fatal(err)
    }

    // Create an Amazon S3 service client
    client := s3.NewFromConfig(cfg)

    // Get the first page of results for ListObjectsV2 for a bucket
    output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
	Bucket: aws.String(bucket_name),
    })
    if err != nil {
	log.Fatal(err)
    }

    for _, object := range output.Contents {
	message := fmt.Sprintf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	messages = append(messages, message)
    }

    //
    return messages, nil
}

