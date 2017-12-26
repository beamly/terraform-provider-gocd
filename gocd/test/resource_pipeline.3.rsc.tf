resource "gocd_pipeline" "test-pipeline3-upstream" {
  name           = "test-pipeline3-upstream"
  group          = "testing"
  label_template = "$${COUNT}"

  materials = [{
    type = "git"

    attributes {
      url    = "https://github.com/drewsonne/terraform-provider-gocd.git"
      branch = "master"

      //      auto_update = true
    }
  }]
}

resource "gocd_pipeline" "test-pipeline3" {
  name           = "test-pipeline3"
  group          = "testing"
  label_template = "$${COUNT}"

  materials = [
    {
      type = "git"

      attributes {
        url    = "https://github.com/drewsonne/terraform-provider-gocd.git"
        branch = "master"

        //        auto_update = true
      }
    },
  ]
}

# CMD terraform import gocd_pipeline_stage.test "test"
resource "gocd_pipeline_stage" "test" {
  name            = "test"
  pipeline        = "${gocd_pipeline.test-pipeline3.name}"
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
