locals {
  "group" = "test-pipelines"
}

resource "gocd_pipeline" "pipe-A" {
  name      = "pipe-A"
  group     = "${local.group}"
  materials = [{
    type = "git"
    attributes {
      url = "github.com/gocd/gocd"
    }
  }]

  stages = ["${data.gocd_stage_definition.stage-A.json}"]
}

data "gocd_stage_definition" "stage-A" {
  name     = "stage-A"
  jobs     = ["${data.gocd_job_definition.list.json}"]
}

data "gocd_job_definition" "list" {
  name  = "list"
  tasks = ["${data.gocd_task_definition.list.json}"]
}

data "gocd_task_definition" "list" {
  type    = "exec"
  command = "ls"
}

resource "gocd_pipeline" "pipe-B" {
  name      = "pipe-B"
  group     = "${local.group}"
  materials = [{
    type = "dependency"
    attributes {
      pipeline = "${gocd_pipeline.pipe-A.name}"
      stage    = "${data.gocd_stage_definition.stage-A.name}"
    }
  }]

  stages = ["${data.gocd_stage_definition.stage-B.json}"]
}

data "gocd_stage_definition" "stage-B" {
  name     = "stage-B"
  jobs     = ["${data.gocd_job_definition.list.json}"]
}

