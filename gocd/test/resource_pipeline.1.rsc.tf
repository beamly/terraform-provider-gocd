resource "gocd_pipeline_template" "test-pipeline" {
  name = "template1-terraform"
  stages = ["${data.gocd_stage_definition.test-stage.json}"]
}

resource "gocd_pipeline" "test-pipeline" {
  name                    = "pipeline1-terraform"
  group                   = "testing"
  template                = "${gocd_pipeline_template.test-pipeline.name}"
  lock_behavior           = "lockOnFailure"
  label_template          = "build-$${COUNT}"

  materials = [
    {
      type = "git"

      attributes {
        name        = "gocd-github"
        url         = "git@github.com:gocd/gocd"
        branch      = "feature/my-addition"
        destination = "gocd-dir"
      }
    },
  ]
}

data "gocd_stage_definition" "test-stage" {
  name = "test"
  jobs = [
    "${data.gocd_job_definition.test-job.json}"
  ]
}

data "gocd_job_definition" "test-job" {
  name = "test"
  tasks = [
    "${data.gocd_task_definition.test.json}"
  ]
}
data "gocd_task_definition" "test" {
  type = "exec"
  command = "echo"
  arguments = [
    "hello",
    "world",
  ]
}

