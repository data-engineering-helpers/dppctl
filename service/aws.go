//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/service/aws.go
//
package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"
	"log"
	"flag"
	"errors"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"bytes"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/codeartifact"
	"github.com/aws/aws-sdk-go-v2/service/mwaa"
	"github.com/aws/smithy-go/middleware"
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

// A Response struct to map the MWAA CLI API response
type MWAAResponse struct {
    StdErr string `json:"stderr"`
    StdOut string `json:"stdout"`
}

func AWSGetCallerIdentity() (string, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

    // Create an Amazon STS service client
	svc := sts.NewFromConfig(awsConfig)

	// Get the IAM role
	params := &sts.GetCallerIdentityInput{}
	output, err := svc.GetCallerIdentity(ctx, params)
	if err != nil {
		log.Fatalf("expect no error, got %v", err)
	}

	sts_identity := fmt.Sprintf("UserId=%s Account=%s Arn=%s",
		aws.ToString(output.UserId),
		aws.ToString(output.Account),
		aws.ToString(output.Arn))

    //
    return sts_identity, nil
}

func AWSS3List(bucketName string, prefix string) ([]string, error) {
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
		Prefix: aws.String(prefix),
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

func AWSAirflowCreateLoginToken(environment string) (string, string,
	middleware.Metadata, error) {
	// References:
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/mwaa/api_op_CreateCliToken.go
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/mwaa/api_op_CreateWebLoginToken.go
	// * https://github.com/aws/smithy-go/blob/main/middleware/metadata.go
	cliToken := ""
	webServerHostname := ""
	var resultMetadata middleware.Metadata

    //
    if environment == "" {
        return webServerHostname, cliToken, resultMetadata,
			errors.New("empty Airflow/MWAA environment")
    }

    // Create an Amazon MWAA (managed Airflow service) client
    svc := mwaa.NewFromConfig(awsConfig)

    // Create a token for login through MWAA CLI API
	params := &mwaa.CreateCliTokenInput{
		Name: aws.String(environment),
	}
    output, err := svc.CreateCliToken(context.TODO(), params)
    if err != nil {
		log.Fatal(err)
    }

	webServerHostname = aws.ToString(output.WebServerHostname)
	cliToken = aws.ToString(output.CliToken)
	resultMetadata = output.ResultMetadata

    //
    return webServerHostname, cliToken, resultMetadata, nil
}

func AWSAirflowCLI(webServerHostname string, cliToken string,
	command string) (string, error) {
	stdoutStr := ""
	
    //
    if command == "" {
        return stdoutStr, errors.New("empty MWAA CLI command")
    }

	api_url := fmt.Sprintf("https://%s/aws_mwaa/cli", webServerHostname)
	body := []byte(fmt.Sprintf(command))
    request, err := http.NewRequest("POST", api_url, bytes.NewBuffer(body))
    if err != nil {
		log.Fatal(err)
    }

	// Add the headers
	request.Header.Add("Content-Type", "text/plain")

	bearerToken := fmt.Sprintf("Bearer %s", cliToken)
	request.Header.Add("Authorization", bearerToken)

	// Call the API
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
	log.Println("MWAA response data: ", string(responseData))

	// Map the HTTP reponse onto a MWAAResponse structure
	var mwaaResponseObject MWAAResponse
	json.Unmarshal(responseData, &mwaaResponseObject)
	stdoutB64Str := mwaaResponseObject.StdOut

	// Base64 decode the `stdout` string
	stdoutData, err := base64.StdEncoding.DecodeString(stdoutB64Str)
	if err != nil {
		log.Fatal("error:", err)
	}

	stdoutStr = string(stdoutData)

	//
	return stdoutStr, nil
}

