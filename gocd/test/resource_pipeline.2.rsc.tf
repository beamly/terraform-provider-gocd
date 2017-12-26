resource "gocd_pipeline_template" "test-pipeline" {
  name = "template0-terraform"
}

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
