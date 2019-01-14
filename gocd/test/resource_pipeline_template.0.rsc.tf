resource "gocd_pipeline_template" "test-pipeline" {
  name = "template0-terraform"
  stages = ["${data.gocd_stage_definition.test-stage.json}"]
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

