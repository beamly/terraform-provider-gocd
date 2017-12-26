data "gocd_task_definition" "test" {
  type = "exec"

  run_if = [
    "passed",
  ]

  working_directory = "/tmp/"
  command           = "/usr/local/bin/terraform"

  arguments = [
    "-debug",
    "version",
  ]
}
