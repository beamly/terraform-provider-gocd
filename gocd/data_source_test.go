package gocd

import "testing"

func TestDataSource(t *testing.T) {
	t.Run("JobDefinition", testDataSourceJobDefinition)
	t.Run("StageDefinition", testDataSourceStageDefinition)
	t.Run("TaskDefinition", testDataSourceTaskDefinition)
}
