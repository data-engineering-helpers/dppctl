//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/utilities/depl.go
//
package utilities

import (
    "log"
	"io/ioutil"
    "gopkg.in/yaml.v3"
)

type SpecFile struct {
	// Some meta-data for the project
	Metadata struct {
		Env string `yaml:"env"`
		Project string `yaml:"project"`
		GitUrl string `yaml:"git_url"`
	} `yaml:"metadata"`
	
	// Payload/workload: what has to be deployed
    Container struct {
		Module struct {
			Stack string `yaml:"stack"`
			Name string `yaml:"name"`
			Version string `yaml:"version"`
		} `yaml:"module"`

		//
		Dependencies struct {
			Spark struct {
				Version string `yaml:"version"`
			} `yaml:"spark"`
			DeltaSpark struct {
				Version string `yaml:"version"`
			} `yaml:"delta_spark"`
		} `yaml:"dependencies"`
		
	} `yaml:"container"`
	
	// Details of the environment to be deployed

	// Storage container (e.g., AWS S3 bucket, Azure Data Storage, GCS)
	StorageContainer struct {
		Provider string `yaml:"provider"`
		Region string `yaml:"region"`
		AccountId string `yaml:"acct_id"`
		Name string `yaml:"name"`
		Prefix string `yaml:"prefix"`
	} `yaml:"storage_container"`

	// Repository for the software artifacts
	ArtifactRepo struct {
		Provider string `yaml:"provider"`
		Region string `yaml:"region"`
		AccountId string `yaml:"acct_id"`
		Domain string `yaml:"domain"`
		Name string `yaml:"name"`
	} `yaml:"artifact_repo"`

	// Repository for the OCI (e.g., Docker) container images
	ContainerRepo struct {
		Provider string `yaml:"provider"`
		Region string `yaml:"region"`
		AccountId string `yaml:"acct_id"`
		Domain string `yaml:"domain"`
		Name string `yaml:"name"`
	} `yaml:"container_repo"`

	// Airflow service (e.g., AWS MWAA)
	Airflow struct {
		Provider string `yaml:"provider"`
		Region string `yaml:"region"`
		AccountId string `yaml:"acct_id"`
		Domain string `yaml:"domain"`

		//
		Dag struct {
			NamePattern string `yaml:"name_pattern"`
			Tag string `yaml:"tag"`
		} `yaml:"dag"`

		StorageContainer struct {
			Name string `yaml:"name"`
			Prefix string `yaml:"prefix"`
		} `yaml:"storage_container"`
	} `yaml:"airflow"`

	// Compute engine (e.g., Spark on DataBricks, Spark on AWS EMR)
	ComputeEngine struct {
		Provider string `yaml:"provider"`
		Region string `yaml:"region"`
		AccountId string `yaml:"acct_id"`
		Domain string `yaml:"domain"`

		//
		Cluster struct {
			Name string `yaml:"name"`
			Version string `yaml:"version"`
		} `yaml:"cluster"`
	} `yaml:"compute_engine"`


	// Kubernetes service (e.g., AWS EKS)
	Kubernetes struct {
		Provider string `yaml:"provider"`
		Region string `yaml:"region"`
		AccountId string `yaml:"acct_id"`
		Domain string `yaml:"domain"`
		Namespace string `yaml:"Namespace"`
	} `yaml:"kubernetes"`
}

func ReadSpecFile(specFilepath string) (SpecFile, error) {
	t := SpecFile{}

    yamlFile, err := ioutil.ReadFile(specFilepath)
    if err != nil {
        log.Fatalf("Error while reading the specification file: %v", err)
    }
	
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
    
    return t, nil
}

