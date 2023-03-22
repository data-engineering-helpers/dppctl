//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/workflow/workflow.go
//
package workflow

import (
	"fmt"
	"log"
	"regexp"
	
	"github.com/data-engineering-helpers/dppctl/utilities"
	"github.com/data-engineering-helpers/dppctl/service"
)

func Check(deplSpec utilities.SpecFile) {
	// STS - Caller identity (IAM)
	stsStruct, err := service.AWSGetCallerIdentity()
	if err != nil {
		log.Print(err)
	}
	log.Println("AWS IAM/caller identity:")
	log.Println(stsStruct)

	//
	bucketName := deplSpec.StorageContainer.Name
	bucketPrefix := deplSpec.StorageContainer.Prefix
	
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
	mwaaEnv := deplSpec.Airflow.Domain
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
	//log.Println("MWAA/Airflow DAGs: ", mwaaDagMetadata)

	namePattern := deplSpec.Airflow.Dag.NamePattern
	log.Println("Retrieving all the Airflow DAG following the name pattern: ",
		namePattern)

	// Extract the DAG list JSON part
	nameRegex := fmt.Sprintf(".*%s.*", namePattern)
	re := regexp.MustCompile(nameRegex)
	
	for _, dag := range mwaaDagMetadata {
		dagId := dag.DagId

		match := re.FindStringSubmatch(dagId)
		if (len(match) == 0) {
			continue
		}

		log.Println("MWAA/Airflow specific DAG: ", dag)
	}
	
}

