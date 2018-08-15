data "gocd_task_definition" "test" {
  type = "exec"

  run_if = [
    "passed",
  ]

  command = "/usr/local/bin/terraform"

  arguments = [
    "-debug",
    "version",
  ]

  working_directory = "tmp/"
}

data "gocd_job_definition" "test" {
  name = "job-name"

  tasks = [
    "${data.gocd_task_definition.test.json}",
  ]
}

resource "gocd_pipeline_stage" "test-stage" {
  name = "test-stage"

  jobs = [
    "${data.gocd_job_definition.test.json}",
  ]

  manual_approval = true

  authorization_roles = [
    "one",
    "two",
  ]

  environment_variables = [
    {
      name  = "IMAGE"
      value = "#{Image}"
    },
  ]

  pipeline_template = "${gocd_pipeline_template.test-pipeline.id}"

  clean_working_directory = false
  fetch_materials = false
  never_cleanup_artifacts = true
}

resource "gocd_pipeline_template" "test-pipeline" {
  name = "test-pipeline-template"
}
