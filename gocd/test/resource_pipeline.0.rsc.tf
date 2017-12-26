resource "gocd_pipeline_template" "test-pipeline" {
  name = "template0-terraform"
}

//resource "gocd_pipeline_stage" "test-stage" {
//  name = "test-stage"
//  fetch_materials = false
//  clean_working_directory = false
//  never_cleanup_artifacts = false
//  success_approval = true
//  pipeline_template = "${gocd_pipeline_template.test-pipeline.name}"
//  jobs = [
//    <<JOB
//      {
//        "name": "job1",
//        "tasks": [
//          {
//            "type": "exec",
//            "attributes": {
//              "run_if": [
//                "passed"
//              ],
//            "command": "terraform"
//            }
//          }
//        ]
//      }
//JOB
//  ]
//}

resource "gocd_pipeline" "test-pipeline" {
  name     = "pipeline0-terraform"
  group    = "testing"
  template = "${gocd_pipeline_template.test-pipeline.id}"

  materials = [
    {
      type = "git"

      attributes {
        name        = "gocd-src"
        url         = "git@github.com:gocd/gocd"
        branch      = "feature/my-addition"
        destination = "gocd-dir"

        //        auto_update = true
        filter {
          ignore = [
            "one",
            "two",
          ]
        }
      }
    },
  ]
}
