package gocd

import (
	r "github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func testEnvironmentAssociation(t *testing.T) {
	t.Run("Import", testResourceEnvironmentAssociationImportBasic)
	t.Run("Basic", testResourceEnvironmentAssociationBasic)
}

func testResourceEnvironmentAssociationBasic(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdEnvironmentAssociationDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_environment_association.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr(
						"gocd_environment.test-environment",
						"name",
						"test-environment",
					),
					r.TestCheckResourceAttr(
						"gocd_pipeline.test-pipeline",
						"name",
						"test-pipeline",
					),
					r.TestCheckResourceAttr(
						"gocd_pipeline.test-pipeline",
						"id",
						"test-pipeline",
					),
					r.TestCheckResourceAttr(
						"gocd_environment_association.test-environment-association",
						"id",
						"test-environment/p/test-pipeline",
					),
				),
			},
		},
	})
}
