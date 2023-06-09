//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/workflow/workflow.go
//
package workflow

import (
	"fmt"
	"log"
	
	"github.com/data-engineering-helpers/dppctl/utilities"
	"github.com/data-engineering-helpers/dppctl/service"
)

func Check(deplSpec utilities.SpecFile) {
	// /////////////////////////////////
	// STS - Caller identity (IAM)
	// /////////////////////////////////
	stsStruct, err := service.AWSGetCallerIdentity()
	if err != nil {
		log.Print(err)
	}
	log.Println("AWS IAM/caller identity:")
	log.Println(stsStruct)

	//
	bucketName := deplSpec.StorageContainer.Name
	bucketPrefix := deplSpec.StorageContainer.Prefix
	
	// /////////////////////////////////
	// AWS S3
	// /////////////////////////////////
	object_list, err := service.AWSS3List(bucketName, bucketPrefix)
	if err != nil {
		log.Print(err)
	}

	log.Println("List of objects within the following bucket:", bucketName)
	for _, object_metadata := range object_list {
		log.Println(object_metadata)
	}

	// /////////////////////////////////
	// CodeArtifact
	// /////////////////////////////////
	// List of the domains
	//
	caDomainName := deplSpec.ArtifactRepo.Domain
	caDomainOwner := deplSpec.ArtifactRepo.AccountId
	caRepoName := deplSpec.ArtifactRepo.Name
	caFormatStr := deplSpec.ArtifactRepo.Format
	packageName := deplSpec.Container.Module.Name
	packageVersion := deplSpec.Container.Module.Version

	caFormat, err := service.AWSCodeArtifactFormatFromString(caFormatStr)
    if err != nil {
		errMsg := fmt.Sprintf("The %s CodeArtifact format is not known",
		caFormatStr)
        log.Fatalf(errMsg, err)
    }

	pkgVersions, err := service.AWSCodeArtifactListPackageVersions(caDomainName,
		caDomainOwner, caRepoName, caFormat, packageName)
    if err != nil {
		errMsg := fmt.Sprintf("No versioned package can be retrieved from CodeArtifact repository for Domain-name=%s Domain-owner=%s Repo-name=%s Format=%s Pkg-name=%s",
		caDomainName, caDomainOwner, caRepoName, caFormat, packageName)
        log.Fatalf(errMsg, err)
    }
	
	log.Println("List of versioned packages within the CodeArtifact repository:")
	for _, pkgVersion := range pkgVersions {
		log.Println(pkgVersion)
	}

	//
	pkgDetails, err := service.AWSCodeArtifactDescribePackageVersion(caDomainName,
		caDomainOwner, caRepoName, caFormat, packageName, packageVersion)
    if err != nil {
		errMsg := fmt.Sprintf("No versioned package can be retrieved from CodeArtifact repository for Domain-name=%s Domain-owner=%s Repo-name=%s Format=%s Pkg-name=%s Pkg-version=%s",
			caDomainName, caDomainOwner, caRepoName, caFormat,
			packageName, packageVersion)
        log.Fatalf(errMsg + ": %v", err)
    }
	
	log.Println("Details for the versioned package within the CodeArtifact repository:", pkgDetails)

	// /////////////////////////////////
	// Elastic Container Registry (ECR)
	// /////////////////////////////////

	/*
	// List of the repositories
	ecrRepoList, err := service.AWSECRListRepositories()
    if err != nil {
		errMsg := fmt.Sprintf("No repository can be retrieved from ECR service")
        log.Fatalf(errMsg, err)
    }
	
	log.Println("List of repositories within the ECR service:")
	for _, ecrRepo := range ecrRepoList {
		log.Println(ecrRepo)
	}
	*/

	// List of the image tags
	//
	ecrRepoName := deplSpec.ContainerRepo.Name
	ecrImgList, err := service.AWSECRListImages(ecrRepoName)
    if err != nil {
		errMsg := fmt.Sprintf("No image can be retrieved from ECR service for the %s repository",
			ecrRepoName)
        log.Fatalf(errMsg, err)
    }
	
	log.Println("List of repositories within the ECR service for", ecrRepoName)
	for _, ecrImg := range ecrImgList {
		log.Println(ecrImg)
	}
	
	// Description of the images
	//
	ecrImgDetailList, err := service.AWSECRDescribeImages(ecrRepoName)
    if err != nil {
		errMsg := fmt.Sprintf("No image can be retrieved from ECR service for the %s repository",
			ecrRepoName)
        log.Fatalf(errMsg, err)
    }
	
	log.Println("List of image details within the ECR service for", ecrRepoName)
	for _, ecrImg := range ecrImgDetailList {
		log.Println(ecrImg)
	}
	
	// /////////////////////////////////
	// MWAA/Airflow
	// /////////////////////////////////
	// Create a one-time MWAA CLI token
	mwaaEnv := deplSpec.Airflow.Domain
	webServerHostname, cliToken, _,
		err := service.AWSAirflowCreateLoginToken(mwaaEnv)
	if err != nil {
		log.Print(err)
	}
	mwaaStr := fmt.Sprintf("hostname=%s token=%s", webServerHostname, cliToken)
	log.Println("MWAA/Airflow CLI token created: ", mwaaStr)

	// Invoke the MWAA CLI API for the specific command (here, the list of DAGs)
	//command := "version"
	command := "dags list -o json"
	stdoutStr, err := service.AWSAirflowCLI(webServerHostname, cliToken,
		command)
	if err != nil {
		log.Print(err)
	}
	//log.Println("MWAA/Airflow CLI response: ", stdoutStr)

	// Parse the output when the command is "dags list -o json"
	mwaaDagMetadataList, err := utilities.ParseAWSMWAADagListOutput(stdoutStr)
	if err != nil {
		log.Print(err)
	}
	//log.Println("MWAA/Airflow DAGs: ", mwaaDagMetadataList)

	// Retrieve the DAGs, for which the name is matching the given pattern
	namePattern := deplSpec.Airflow.Dag.NamePattern
	log.Println("Retrieving all the Airflow DAG matching the following name pattern:",
		namePattern)

	dagList, err := utilities.ExtractMatchingAWSMWAADagList(mwaaDagMetadataList,
		namePattern)
	
	log.Println("List of MWAA/Airflow DAGs mathcing the name pattern:", dagList)
}

