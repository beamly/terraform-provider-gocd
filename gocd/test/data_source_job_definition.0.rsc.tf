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
