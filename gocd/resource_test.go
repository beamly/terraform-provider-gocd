package gocd

import "testing"

func TestResource(t *testing.T) {
	t.Run("PipelineTemplate", testResourcePipelineTemplate)
	t.Run("Pipeline", testResourcePipeline)
	t.Run("Environment", testEnvironment)
	t.Run("EnvironmentAssociation", testEnvironmentAssociation)
}
