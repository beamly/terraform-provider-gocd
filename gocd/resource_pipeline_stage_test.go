package gocd

import (
	"context"
	"fmt"
	"github.com/beamly/go-gocd/gocd"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testResourceStage(t *testing.T) {
	t.Run("Basic", testResourceStageBasic)
	t.Run("FetchMaterials", testResourceStageFetchMaterials)
	t.Run("FetchMaterialsTrueToFalse", testResourceStageFetchMaterialsTrueToFalse)
	t.Run("CleanWorkingDirectory", testResourceStageCleanWorkingDirectory)
	t.Run("NeverCleanupArtefacts", testResourceStageNeverCleanupArtefacts)
	t.Run("Import", testResourcePipelineStageImportBasic)
	t.Run("PTypeName", testResourcePipelineStagePtypeName)
	t.Run("Helpers", testResourcePipelineStageHelpers)
	t.Run("ExistsFails", testResourcePipelineStageExistsFails)
	//t.Run("Update", testResourcePipelineStageUpdate)
}

func testResourcePipelineStageExistsFails(t *testing.T) {
	ds := resourcePipelineStage().Data(&terraform.InstanceState{})
	exists, err := resourcePipelineStageExists(ds, nil)
	assert.False(t, exists)
	assert.EqualError(t, err, "could not parse the provided id ``")
}

//func testResourcePipelineStageUpdate(t *testing.T) {
//	ds := resourcePipelineStage().Data(&terraform.InstanceState{
//		ID: "template/test-pipeline/test-stage",
//	})
//	err := resourcePipelineStageUpdate(ds, )
//	assert.EqualError(t, err, "not implemented")
//}

func testResourcePipelineStagePtypeName(t *testing.T) {
	t.Run("Template", testResourcePipelineStagePtypeNameTemplate)
	t.Run("Pipeline", testResourcePipelineStagePtypeNamePipeline)
	t.Run("Fail", testResourcePipelineStagePtypeNameFail)
}

func testResourcePipelineStagePtypeNameFail(t *testing.T) {
	ds := (&schema.Resource{Schema: map[string]*schema.Schema{
		"pipeline":          {Type: schema.TypeString, Optional: true},
		"pipeline_template": {Type: schema.TypeString, Optional: true},
	}}).Data(&terraform.InstanceState{})
	err := resourcePipelineStageSetPTypeName(ds, "unknown-type", "test-pipeline")
	assert.EqualError(t, err, "Unexpected pipeline type `unknown-type`")

	p, pOk := ds.GetOk("pipeline_template")
	assert.False(t, pOk)
	assert.Empty(t, p)

	pt, ptOk := ds.GetOk("pipeline")
	assert.False(t, ptOk)
	assert.Empty(t, pt)
}

func testResourcePipelineStagePtypeNamePipeline(t *testing.T) {
	ds := (&schema.Resource{Schema: map[string]*schema.Schema{
		"pipeline":          {Type: schema.TypeString, Optional: true},
		"pipeline_template": {Type: schema.TypeString, Optional: true},
	}}).Data(&terraform.InstanceState{})
	err := resourcePipelineStageSetPTypeName(ds, STAGE_TYPE_PIPELINE, "test-pipeline")
	assert.Nil(t, err)

	p, pOk := ds.GetOk("pipeline_template")
	assert.False(t, pOk)
	assert.Empty(t, p)

	pt, ptOk := ds.GetOk("pipeline")
	assert.True(t, ptOk)
	assert.Equal(t, pt, "test-pipeline")
}

func testResourcePipelineStagePtypeNameTemplate(t *testing.T) {
	ds := (&schema.Resource{Schema: map[string]*schema.Schema{
		"pipeline":          {Type: schema.TypeString, Optional: true},
		"pipeline_template": {Type: schema.TypeString, Optional: true},
	}}).Data(&terraform.InstanceState{})
	err := resourcePipelineStageSetPTypeName(ds, STAGE_TYPE_PIPELINE_TEMPLATE, "test-pipeline-template")
	assert.Nil(t, err)

	p, pOk := ds.GetOk("pipeline")
	assert.False(t, pOk)
	assert.Empty(t, p)

	pt, ptOk := ds.GetOk("pipeline_template")
	assert.True(t, ptOk)
	assert.Equal(t, pt, "test-pipeline-template")
}

func testResourceStageBasic(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdStageDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_stage.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineStageExists("gocd_pipeline_stage.test-stage", true, false, false),
				),
			},
		},
	})
}

func testResourceStageCleanWorkingDirectory(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdStageDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_stage.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineStageExists("gocd_pipeline_stage.test-stage", true, false, false),
					r.TestCheckResourceAttr("gocd_pipeline_stage.test-stage", "clean_working_directory", "true"),
					r.TestCheckResourceAttr("gocd_pipeline_stage.test-stage", "fetch_materials", "false"),
					r.TestCheckResourceAttr("gocd_pipeline_stage.test-stage", "never_cleanup_artifacts", "false"),
				),
			},
		},
	})
}

func testResourceStageFetchMaterials(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdStageDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_stage.1.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineStageExists("gocd_pipeline_stage.test-stage", false, true, false),
					r.TestCheckResourceAttr("gocd_pipeline_stage.test-stage", "clean_working_directory", "false"),
					r.TestCheckResourceAttr("gocd_pipeline_stage.test-stage", "fetch_materials", "true"),
					r.TestCheckResourceAttr("gocd_pipeline_stage.test-stage", "never_cleanup_artifacts", "false"),
				),
			},
		},
	})
}

func testResourceStageFetchMaterialsTrueToFalse(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdStageDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_stage.1.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineStageExists("gocd_pipeline_stage.test-stage", false, true, false),
				),
			},
			{
				Config: testFile("resource_pipeline_stage.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineStageExists("gocd_pipeline_stage.test-stage", true, false, false),
				),
			},
		},
	})
}

func testResourceStageNeverCleanupArtefacts(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdStageDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_stage.2.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineStageExists("gocd_pipeline_stage.test-stage", false, false, true),
					r.TestCheckResourceAttr("gocd_pipeline_stage.test-stage", "clean_working_directory", "false"),
					r.TestCheckResourceAttr("gocd_pipeline_stage.test-stage", "fetch_materials", "false"),
					r.TestCheckResourceAttr("gocd_pipeline_stage.test-stage", "never_cleanup_artifacts", "true"),
				),
			},
		},
	})
}

func testCheckPipelineStageExists(resource string, cleanWorkingDir bool, fetchMaterials bool, neverCleanupArtifacts bool) r.TestCheckFunc {
	return func(s *terraform.State) error {
		var pipeline *gocd.PipelineTemplate
		var err error

		rcs := s.RootModule().Resources
		rs, ok := rcs[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pipeline stage name is set")
		}

		if pipeline, _, err = testGocdClient.PipelineTemplates.Get(context.Background(), "test-pipeline-template"); err != nil {
			return err
		}

		if pipeline.Stages[0].CleanWorkingDirectory != cleanWorkingDir {
			return fmt.Errorf("clean_working_directory property not set to %t", cleanWorkingDir)
		}

		if pipeline.Stages[0].FetchMaterials != fetchMaterials {
			return fmt.Errorf("fetch_materials property not set to %t", fetchMaterials)
		}

		if pipeline.Stages[0].NeverCleanupArtifacts != neverCleanupArtifacts {
			return fmt.Errorf("never_cleanup_artifacts property not set to %t", neverCleanupArtifacts)
		}

		return nil
	}
}

func testGocdStageDestroy(s *terraform.State) error {
	root := s.RootModule()
	for _, rs := range root.Resources {
		if rs.Type != "gocd_pipeline_stage" {
			continue
		}

		_, _, err := testGocdClient.PipelineTemplates.Get(context.Background(), rs.Primary.ID)
		//stage := pt.GetStage()
		if err == nil {
			return fmt.Errorf("still exists")
		}
	}

	return nil
}
