package gocd

import (
	"fmt"
	r "github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestDataSourceTaskDefinition(t *testing.T) {
	for i := 0; i <= 5; i++ {
		t.Run(
			fmt.Sprintf("gocd_task_definition.%d", i),
			DataSourceTaskDefinition(t, i,
				fmt.Sprintf("data_source_task_definition.%d.rsc.tf", i),
				fmt.Sprintf("data_source_task_definition.%d.rsp.json", i),
			),
		)
	}
}

func DataSourceTaskDefinition(t *testing.T, index int, configPath string, expectedPath string) func(t *testing.T) {
	return func(t *testing.T) {
		config := testFile(configPath)
		expected := testFile(expectedPath)
		r.UnitTest(t, r.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testGocdProviders,
			Steps: testStepComparisonCheck(&TestStepJSONComparison{
				Index:        index,
				ID:           "data.gocd_task_definition.test",
				Config:       config,
				ExpectedJSON: expected,
			}),
		})
	}
}
