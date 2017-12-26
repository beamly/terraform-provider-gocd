data "gocd_task_definition" "test" {
  type = "fetch"

  run_if = [
    "failed",
  ]

  pipeline         = "pipeline1"
  stage            = "stage2"
  job              = "job3"
  source           = "source_artifact/"
  is_source_a_file = false
  destination      = "dest_artifact/"
}
