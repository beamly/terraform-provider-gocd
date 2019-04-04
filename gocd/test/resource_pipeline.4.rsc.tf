resource "gocd_pipeline_template" "test-template4" {
  name = "test-template4"
  stages = ["${data.gocd_stage_definition.test-stage.json}"]
}

resource "gocd_pipeline" "test-pipeline" {
  name     = "test-pipeline"
  group    = "ecsagent"
  template = "${gocd_pipeline_template.test-template4.name}"

  parameters {
    Image = "base"
  }

  materials = [
    {
      type = "git"

      attributes {
        url    = "git@github.com:org/gocd-ecsagents.git"
        branch = "master"

        //        auto_update = true
        filter = [
          "gocd-agents/Dockerfile.base",
          "Makefile",
          "gocd-agents/files/base/",
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

