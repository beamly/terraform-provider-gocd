resource "gocd_pipeline_template" "test-pipeline" {
  name = "template0-terraform"

  stages = ["${data.gocd_stage_definition.test-stage.json}"]
}

resource "gocd_pipeline" "test-pipeline" {
  name                    = "pipeline0-terraform"
  group                   = "testing"
  template                = "${gocd_pipeline_template.test-pipeline.name}"
  enable_pipeline_locking = true

  materials = [
    {
      type = "git"

      attributes {
        name        = "gocd-src"
        url         = "git@github.com:gocd/gocd"
        branch      = "feature/my-addition"
        destination = "gocd-dir"

        filter = [
          "one",
          "two",
        ]
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

