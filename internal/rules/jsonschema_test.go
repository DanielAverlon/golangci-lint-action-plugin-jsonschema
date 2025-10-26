package rules_test

import (
	"testing"

	"github.com/golangci/golangci-lint-action-plugin-example/internal/rules"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestJsonSchema(t *testing.T) {
	testDataDir := analysistest.TestData()

	analysistest.Run(t, testDataDir, rules.NoDescriptionCommas, "./src/jsonschema/")
}

func TestJsonSchemaAutoFix(t *testing.T) {
	testDataDir := analysistest.TestData()

	results := analysistest.RunWithSuggestedFixes(t, testDataDir, rules.NoDescriptionCommas, "./src/jsonschema/")
	suggestedFixProvided := false
	for _, result := range results {
		for _, diagnostic := range result.Diagnostics {
			for _, suggestedFix := range diagnostic.SuggestedFixes {
				if len(suggestedFix.TextEdits) != 0 {
					suggestedFixProvided = true
				}
			}
		}
	}

	if !suggestedFixProvided {
		t.Errorf("expected a suggested fix to be provided, but didn't have any in %+v", results)
	}
}
