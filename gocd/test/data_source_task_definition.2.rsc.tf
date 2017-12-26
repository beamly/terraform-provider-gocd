data "gocd_task_definition" "test" {
  type = "nant"

  run_if = [
    "any",
  ]

  working_directory = "/tmp/nant/"
  target            = "test-nant"
  build_file        = "./nant.xml"
  nant_path         = "c:/windows/nant.exe"
}
