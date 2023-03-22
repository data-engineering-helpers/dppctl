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
	domain_list, err := service.AWSCodeArtifactListDomains()
	if err != nil {
		log.Print(err)
	}
	log.Println("List of domains within the CodeArtifact service:")
	for _, domain := range domain_list {
		log.Println(domain)
	}

	//
	caDomainName := deplSpec.ArtifactRepo.Domain
	caDomainOwner := deplSpec.ArtifactRepo.AccountId
	caRepoName := deplSpec.ArtifactRepo.Name
	caRepoFormatStr := deplSpec.ArtifactRepo.Format
	packageName := deplSpec.Container.Module.Name
	//packageVersion := deplSpec.Container.Module.Version

	caRepoFormat, err := service.AWSCodeArtifactFormatFromString(caRepoFormatStr)
    if err != nil {
		errMsg := fmt.Sprintf("The %s CodeArtifact repository format is not known")
        log.Fatalf(errMsg, err)
    }
	service.AWSCodeArtifactListPackageVersions(caDomainName, caDomainOwner,
		caRepoName, caRepoFormat, packageName)
	
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

