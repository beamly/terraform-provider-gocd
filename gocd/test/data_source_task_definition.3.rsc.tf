data "gocd_task_definition" "test" {
  type = "rake"

  run_if = [
    "failed",
  ]

  working_directory = "tmp/rake/"
  target            = "test-rake"
  build_file        = "./rake.rb"
}
