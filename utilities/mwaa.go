//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/utilities/mwaa.go
//
package utilities

import (
	"errors"
    //"log"
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

