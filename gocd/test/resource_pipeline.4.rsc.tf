resource "gocd_pipeline_template" "test-template4" {
  name = "test-template4"
}

resource "gocd_pipeline" "test-pipeline" {
  name     = "test-pipeline"
  group    = "ecsagent"
  template = "${gocd_pipeline_template.test-template4.name}"

  parameters {
    Image = "base"
  }

  materials = [
    {
      type = "git"

      attributes {
        url    = "git@github.com:org/gocd-ecsagents.git"
        branch = "master"

        //        auto_update = true
        filter {
          ignore = [
            "gocd-agents/Dockerfile.base",
            "Makefile",
            "gocd-agents/files/base/",
          ]
        }
      }
    },
  ]
}
