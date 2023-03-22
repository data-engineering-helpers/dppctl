//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/utilities/mwaa.go
//
package utilities

import (
	"errors"
    "fmt"
	"regexp"
	"encoding/json"
)

type MwaaDagMetadata struct {
	DagId string `json:"dag_id"`
	Filepath string `json:"filepath"`
	Owner string `json:"owner"`
	Paused string `json:"paused"`
}

func ParseAWSMWAADagListOutput(rawOutput string) ([]MwaaDagMetadata, error) {
	// Retrieve, from a JSON-formatted string, the list of DAGs as a structure
	var mwaaDagMetadataList []MwaaDagMetadata

	// Extract the DAG list JSON part
	re := regexp.MustCompile(`(\[\{"dag_id":.*\])`)
	match := re.FindStringSubmatch(rawOutput)
	if (len(match) == 0) {
		return mwaaDagMetadataList, errors.New("No match for JSON DAG list")
	}
	dagListStr := match[1]
	
	// Map the string (containing the JSON-formatted DAG list)
	// onto a MWAAResponse structure
	json.Unmarshal([]byte(dagListStr), &mwaaDagMetadataList)

	//
	return mwaaDagMetadataList, nil
}

func ExtractMatchingAWSMWAADagList(mwaaDagMetadataList []MwaaDagMetadata,
	// Retrieve the DAGs, for which the name is matching the given pattern
	namePattern string) ([]MwaaDagMetadata, error) {
	dagList := []MwaaDagMetadata{}

	// Build a RegExp from the given name pattern
	nameRegex := fmt.Sprintf(".*%s.*", namePattern)
	re := regexp.MustCompile(nameRegex)
	
	for _, dag := range mwaaDagMetadataList {
		dagId := dag.DagId

		match := re.FindStringSubmatch(dagId)
		if (len(match) == 0) {
			continue
		}

		dagList = append(dagList, dag)
	}

	//
	return dagList, nil
}

