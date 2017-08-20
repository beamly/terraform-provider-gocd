package gocd

import (
	"fmt"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestResourcePipelineTemplate(t *testing.T) {
	t.Run("Basic", testResourcePipelineTemplateBasic)
	t.Run("ImportBasic", testResourcePipelineTemplateImportBasic)
}

func testResourcePipelineTemplateBasic(t *testing.T) {

	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineTemplateDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_template.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineTemplateExists("gocd_pipeline_template.test-pipeline"),
					testCheckPipelineTemplateName(
						"gocd_pipeline_template.test-pipeline", "template0-terraform"),
					testCheckPipelineTemplate1StageCount("gocd_pipeline_template.test-pipeline"),
				),
			},
			{
				Config: testFile("resource_pipeline_template.1.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineTemplateExists("gocd_pipeline_template.test-pipeline"),
					testCheckPipelineTemplateName(
						"gocd_pipeline_template.test-pipeline", "template1-terraform"),
					testCheckPipelineTemplate2StageCount("gocd_pipeline_template.test-pipeline"),
				),
			},
		},
	})

}

func testCheckPipelineTemplate1StageCount(resource string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		if rs := s.RootModule().Resources[resource].Primary; rs.Attributes["stages.#"] != "1" {
			return fmt.Errorf("Expected 1 stage. Found '%s'", rs.Attributes["stages.#"])
		}
		return nil
	}
}

func testCheckPipelineTemplate2StageCount(resource string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		if rs := s.RootModule().Resources[resource].Primary; rs.Attributes["stages.#"] != "2" {
			return fmt.Errorf("Expected 2 stages. Found '%s'", rs.Attributes["stages.#"])
		}
		return nil
	}
}

func testCheckPipelineTemplateName(resource string, id string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		if rs := s.RootModule().Resources[resource]; rs.Primary.ID != id {
			return fmt.Errorf("Expected id 'template1-terraform', got '%s", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckPipelineTemplateExists(resource string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pipeline template name is set")
		}

		return nil
	}
}
