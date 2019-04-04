package gocd

import (
	"fmt"
	r "github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func testDataSourceStageDefinition(t *testing.T) {
	for i := 0; i <= 2; i++ {
		t.Run(
			fmt.Sprintf("gocd_stage_definition.%d", i),
			dataSourceStageDefinition(t, i,
				fmt.Sprintf("data_source_stage_definition.%d.rsc.tf", i),
				fmt.Sprintf("data_source_stage_definition.%d.rsp.json", i),
			),
		)
	}

}

func dataSourceStageDefinition(t *testing.T, index int, configPath string, expectedPath string) func(t *testing.T) {
	return func(t *testing.T) {
		config := testFile(configPath)
		expected := testFile(expectedPath)
		r.UnitTest(t, r.TestCase{
			Providers: testGocdProviders,
			Steps: testStepComparisonCheck(&TestStepJSONComparison{
				Index:        index,
				ID:           "data.gocd_stage_definition.test",
				Config:       config,
				ExpectedJSON: expected,
			}),
		})
	}
}
