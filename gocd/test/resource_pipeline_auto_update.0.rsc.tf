resource "gocd_pipeline" "pipeline1" {
  group = "test-auto-update"
  name  = "pipeline1"

  materials = [{
    type = "git"

    attributes {
      name        = "gocd-src"
      url         = "https://github.com/gocd/gocd"
      auto_update = false
    }
  }]

  stages = ["${data.gocd_stage_definition.test-stage.json}"]
}

resource "gocd_pipeline" "pipeline2" {
  group = "test-auto-update"
  name  = "pipeline2"

  materials = [{
    type = "git"

    attributes {
      name        = "gocd-src"
      url         = "https://github.com/gocd/gocd"
      auto_update = false
    }
  }]

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

