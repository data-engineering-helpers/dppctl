//
// File: https://github.com/data-engineering-helpers/dppctl/blob/main/dppctl_test.go
//
package main

import (
  "testing"
	
	"github.com/data-engineering-helpers/dppctl/utilities"
)

/**
 * Check that the default/sample deployment specification file
 * is parsed correctly
 */
func TestReadSpecFile(t *testing.T) {
	// Specification of the deployment
    specFilepath := "depl/aws-dev-sample.yaml"
	deplSpec, err := utilities.ReadSpecFile(specFilepath)
	if err != nil {
        t.Fatalf(`utilities.ReadSpecFile() = %q, %v, parsed %#q spec file`,
			deplSpec, err, specFilepath)
	}
	
}

