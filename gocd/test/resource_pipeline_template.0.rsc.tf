resource "gocd_pipeline_template" "test-pipeline" {
  name = "template0-terraform"
  stages = [
    <<STAGE
{
  "name": "test-stage",
  "fetch_materials": false,
  "clean_working_directory": false,
  "never_cleanup_artifacts": false,
  "approval": {
    "type": "success"
  },
  "jobs": [
    {
      "name": "job1",
      "tasks": [
        {
          "type": "exec",
          "attributes": {
            "run_if": [
              "passed"
            ],
            "command": "terraform"
          }
        }
      ]
    }
  ]
}
STAGE
  ]
}
