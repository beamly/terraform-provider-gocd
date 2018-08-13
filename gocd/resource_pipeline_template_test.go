package gocd

import (
	"context"
	"fmt"
	"github.com/beamly/go-gocd/gocd"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func testResourcePipelineTemplate(t *testing.T) {
	t.Run("Basic", testResourcePipelineTemplateBasic)
	t.Run("ImportBasic", testResourcePipelineTemplateImportBasic)
	t.Run("Exists", testResourcePipelineTemplateExists)
	t.Run("PipelineReadHelper", testResourcePipelineTemplateReadHelper)
	t.Run("Missing", testResourcePipelineTemplateMissing)
}

func testResourcePipelineTemplateReadHelper(t *testing.T) {

	t.Run("MissingName", testResourcePipelineTemplateReadHelperMissingName)
}

func testResourcePipelineTemplateReadHelperMissingName(t *testing.T) {

	rd := (&schema.Resource{Schema: map[string]*schema.Schema{}}).Data(&terraform.InstanceState{})
	e := errors.New("mock-error")
	err := readPipelineTemplate(rd, nil, e)

	assert.EqualError(t, err, "mock-error")
}

func testResourcePipelineTemplateExists(t *testing.T) {
	rd := (&schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true},
	}}).Data(&terraform.InstanceState{})

	exists, err := resourcePipelineTemplateExists(rd, nil)
	assert.False(t, exists)
	assert.EqualError(t, err, "`name` can not be empty")
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
				),
			},
		},
	})

}

func testResourcePipelineTemplateMissing(t *testing.T) {

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
				),
			},
			{
				Config: testFile("resource_pipeline_template.0.rsc.tf"),
				SkipFunc: func() (bool, error) {
					cfg := gocd.Configuration{
						Server:   os.Getenv("GOCD_URL"),
						Username: os.Getenv("GOCD_USERNAME"),
						Password: os.Getenv("GOCD_PASSWORD"),
					}
					c := cfg.Client()

					c.PipelineTemplates.Delete(context.Background(), "template0-terraform")

					return false, nil
				},
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineTemplateExists("gocd_pipeline_template.test-pipeline"),
					testCheckPipelineTemplateName(
						"gocd_pipeline_template.test-pipeline", "template0-terraform"),
				),
			},
		},
	})

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
