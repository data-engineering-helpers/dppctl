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
	awscatypes "github.com/aws/aws-sdk-go-v2/service/codeartifact/types"	
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

func AWSCodeArtifactListDomains() ([]string, error) {
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

func AWSCodeArtifactFormatFromString(format string) (awscatypes.PackageFormat,
	error) {
	switch format {
	case "pypi":
		return awscatypes.PackageFormatPypi, nil
	case "maven":
		return awscatypes.PackageFormatMaven, nil
	}

	errMsg := fmt.Sprintf("The %s CodeArtifact repository format is not known")
	return awscatypes.PackageFormatGeneric, errors.New(errMsg)
}

func AWSCodeArtifactListPackageVersions(domainName string,
	domainOwner string, repoName string, repoFormat awscatypes.PackageFormat,
	packageName string) ([]string, error) {
	// References:
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/codeartifact/types/types.go
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/codeartifact/api_op_ListPackageVersions.go
	
    pkgVersions := []string {}

    // Using the Config value, create the CodeArtifact client
    svc := codeartifact.NewFromConfig(awsConfig)

    // Build the request with its input parameters
	params := &codeartifact.ListPackageVersionsInput{
		Domain: aws.String(domainName),
		DomainOwner: aws.String(domainOwner),
		Format: repoFormat,
		Repository: aws.String(repoName),
		Package: aws.String(packageName),
	}
    resp, err := svc.ListPackageVersions(context.TODO(), params)
    if err != nil {
        log.Fatalf("failed to list the versions of the package, %v", err)
    }

	// https://docs.aws.amazon.com/codeartifact/latest/APIReference/API_PackageVersionSummary.html
    for _, versionStruct := range resp.Versions {
		origin := versionStruct.Origin
		domainEntryPoint := origin.DomainEntryPoint
		externalConnectionName := aws.ToString(domainEntryPoint.ExternalConnectionName)
		repositoryName := aws.ToString(domainEntryPoint.RepositoryName)
		originType := origin.OriginType
		version := aws.ToString(versionStruct.Version)
		revision := aws.ToString(versionStruct.Revision)
		status := versionStruct.Status

		message := fmt.Sprintf("Pkg-name=%s Version=%s Status=%s Revision=%s Origin=(domain-entry-point=%s, repository-name=%s, origin-type=%s)",
			packageName, version, status, revision,
			externalConnectionName, repositoryName, originType)
		//log.Println("Package details:", message)
		pkgVersions = append(pkgVersions, message)
    }

    //
    return pkgVersions, nil
}

func AWSCodeArtifactDescribePackageVersion(domainName string,
	domainOwner string, repoName string, repoFormat awscatypes.PackageFormat,
	packageName string, packageVersion string) (string, error) {
	// References:
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/codeartifact/types/types.go
	// * https://github.com/aws/aws-sdk-go-v2/blob/main/service/codeartifact/api_op_DescribePackageVersion.go
	
    var pkgDetailsStr string

    // Using the Config value, create the CodeArtifact client
    svc := codeartifact.NewFromConfig(awsConfig)

    // Build the request with its input parameters
	params := &codeartifact.DescribePackageVersionInput{
		Domain: aws.String(domainName),
		DomainOwner: aws.String(domainOwner),
		Format: repoFormat,
		Repository: aws.String(repoName),
		Package: aws.String(packageName),
		PackageVersion: aws.String(packageVersion),
	}
    resp, err := svc.DescribePackageVersion(context.TODO(), params)
    if err != nil {
		errMsg := fmt.Sprintf("Failed to describe package for the specific versioned package, %v",
			err)
		log.Println(errMsg)
        return pkgDetailsStr, err
    }

	// https://docs.aws.amazon.com/codeartifact/latest/APIReference/API_PackageVersionDescription.html
	packageVersionDesc := resp.PackageVersion
	displayName := aws.ToString(packageVersionDesc.DisplayName)
	homePage := aws.ToString(packageVersionDesc.HomePage)
	licenses := packageVersionDesc.Licenses
	namespace := aws.ToString(packageVersionDesc.Namespace)
	publishedTime := packageVersionDesc.PublishedTime
	revision := aws.ToString(packageVersionDesc.Revision)
	status := packageVersionDesc.Status
	sourceCodeRepository := aws.ToString(packageVersionDesc.SourceCodeRepository)
	origin := packageVersionDesc.Origin
	domainEntryPoint := origin.DomainEntryPoint
	externalConnectionName := aws.ToString(domainEntryPoint.ExternalConnectionName)
	repositoryName := aws.ToString(domainEntryPoint.RepositoryName)
	originType := origin.OriginType

	pkgDetailsStr = fmt.Sprintf("Pkg-name=%s Display-name=%s Version=%s Status=%s Revision=%s Homepage=%s Namespace=%s Source-code-repo=%s Published-time=%s Licenses=%s Origin=(domain-entry-point=%s, repository-name=%s, origin-type=%s)",
		packageName, displayName, packageVersion, status, revision, homePage,
		namespace, sourceCodeRepository, publishedTime, licenses,
		externalConnectionName, repositoryName, originType)
	
    //
    return pkgDetailsStr, nil
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
	/**
	 * As of 2023, it does not seem possible to target/use the Airflow API
	 * directly on the AWS managed service (MWAA). One has to use
	 * the API backend of the MWAA CLI. That is why the code for
	 * that Go function is not straightforward.
	 * Note that the use of the MWAA CLI API (through `curl`) is itself
	 * convoluted. See also
	 * https://github.com/data-engineering-helpers/dppctl/blob/main/README.md
	 *
	 * References:
	 * Stack Overflow - Is it possible to access the Airflow API in AWS MWAA?
	 *  https://stackoverflow.com/questions/67884770/is-it-possible-to-access-the-airflow-api-in-aws-mwaa
	 * Apache Airflow - Airflow API reference guide: https://airflow.apache.org/docs/apache-airflow/stable/stable-rest-api-ref.html
	 * AWS - Amazon Managed Workflows for Apache Airflow (MWAA) User Guide:
	 *  https://docs.aws.amazon.com/mwaa/index.html
	 * AWS - Accessing the Apache Airflow UI:
	 *  https://docs.aws.amazon.com/mwaa/latest/userguide/access-airflow-ui.html
	 * AWS - Apache Airflow CLI command reference:
	 *  https://docs.aws.amazon.com/mwaa/latest/userguide/airflow-cli-command-reference.html)
	 * GitHub - AWS - Sample code for MWAA:
	 *  https://github.com/aws-samples/amazon-mwaa-examples
	 * GitHub - AWS - Sample code for MWAA - Bash operator script:
	 *  https://github.com/aws-samples/amazon-mwaa-examples/tree/main/dags/bash_operator_script
	 */
	
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
	//request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Type", "application/json")

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
	//log.Println("MWAA response data: ", string(responseData))

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

