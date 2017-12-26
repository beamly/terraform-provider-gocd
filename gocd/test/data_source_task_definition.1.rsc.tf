data "gocd_task_definition" "test" {
  type = "ant"

  run_if = [
    "failed",
  ]

  working_directory = "/tmp/ant/"
  target            = "test-ant"
  build_file        = "./ant.xml"
}
