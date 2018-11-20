package gocd

import (
	"fmt"
	"github.com/beamly/go-gocd/gocd"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testDataSourceTaskDefinition(t *testing.T) {
	for i := 0; i <= 7; i++ {
		t.Run(
			fmt.Sprintf("gocd_task_definition.%d", i),
			DataSourceTaskDefinition(t, i,
				fmt.Sprintf("data_source_task_definition.%d.rsc.tf", i),
				fmt.Sprintf("data_source_task_definition.%d.rsp.json", i),
			),
		)
	}

	t.Run("UnexpectedTaskType", dataSourceTaskTypeFail)

}

func dataSourceTaskTypeFail(t *testing.T) {
	task := gocd.Task{
		Attributes: gocd.TaskAttributes{},
	}

	rd := (&schema.Resource{Schema: map[string]*schema.Schema{
		"build_file": {Type: schema.TypeString, Required: true},
		"target":     {Type: schema.TypeString, Required: true},
	}}).Data(&terraform.InstanceState{
		Attributes: map[string]string{
			"build_file": "mock-build-file",
			"target":     "mock-target",
		},
	})

	dataSourceGocdTaskBuildExec(&task, rd)

	assert.Equal(t, "mock-build-file", task.Attributes.BuildFile)
	assert.Equal(t, "mock-target", task.Attributes.Target)

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
