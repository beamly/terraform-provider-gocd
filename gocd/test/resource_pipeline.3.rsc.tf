resource "gocd_pipeline" "test-pipeline3-upstream" {
  name           = "test-pipeline3-upstream"
  group          = "testing"
  label_template = "$${COUNT}"

  materials = [{
    type = "git"

    attributes {
      url    = "https://github.com/beamly/terraform-provider-gocd.git"
      branch = "master"

      //      auto_update = true
    }
  }]

  stages = ["${data.gocd_stage_definition.test.json}"]
}

resource "gocd_pipeline" "test-pipeline3" {
  name           = "test-pipeline3"
  group          = "testing"
  label_template = "$${COUNT}"

  materials = [
    {
      type = "git"

      attributes {
        url    = "https://github.com/beamly/terraform-provider-gocd.git"
        branch = "master"

        //        auto_update = true
      }
    },
  ]

  stages = ["${data.gocd_stage_definition.test.json}"]
}

# CMD terraform import gocd_pipeline_stage.test "test"
data "gocd_stage_definition" "test" {
  name            = "test"
  fetch_materials = true

  jobs = [
    "${data.gocd_job_definition.test.json}",
  ]
}

data "gocd_job_definition" "test" {
  name = "test"

  tasks = [
    "${data.gocd_task_definition.test-pipeline3_test_test_1.json}",
  ]
}

data "gocd_task_definition" "test-pipeline3_test_test_1" {
  type    = "exec"
  run_if  = ["passed"]
  command = "echo"

  arguments = [
    "test",
  ]
}
