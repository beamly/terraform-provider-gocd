package gocd

import (
	"fmt"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataSourceGocdStageTemplateRead(t *testing.T) {
	t.Run("success", testDataSourceGocdStageTemplateReadSuccess)
	t.Run("fail", testDataSourceGocdStageTemplateReadFail)
}
func testDataSourceGocdStageTemplateReadFail(t *testing.T) {
	s := dataSourceGocdStageTemplate().Schema
	d := map[string]interface{}{
		"success_approval": true,
		"jobs":             []string{"one", "two"},
	}
	rd := schema.TestResourceDataRaw(t, s, d)
	err := dataSourceGocdStageTemplateRead(rd, nil)
	assert.NotNil(t, err)
}
func testDataSourceGocdStageTemplateReadSuccess(t *testing.T) {
	s := dataSourceGocdStageTemplate().Schema
	d := map[string]interface{}{
		"name":             "one",
		"success_approval": true,
		"jobs":             []string{"one", "two"},
	}
	rd := schema.TestResourceDataRaw(t, s, d)
	err := dataSourceGocdStageTemplateRead(rd, nil)
	assert.Nil(t, err)
	assert.Equal(t, `{
  "name": "one",
  "fetch_materials": false,
  "clean_working_directory": false,
  "never_cleanup_artifacts": false,
  "approval": {
    "type": "success"
  }
}`, rd.Get("json"))

}

func TestDataSourceStageTemplate(t *testing.T) {
	for i := 0; i <= 0; i++ {
		t.Run(
			fmt.Sprintf("gocd_stage_definition.%d", i),
			DataSourceStageDefinition(t, i,
				fmt.Sprintf("data_source_stage_definition.%d.rsc.tf", i),
				fmt.Sprintf("data_source_stage_definition.%d.rsp.json", i),
			),
		)
	}
}

func DataSourceStageDefinition(t *testing.T, index int, configPath string, expectedPath string) func(t *testing.T) {
	return func(t *testing.T) {
		config := testFile(configPath)
		expected := testFile(expectedPath)
		r.UnitTest(t, r.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
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
