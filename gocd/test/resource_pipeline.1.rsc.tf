resource "gocd_pipeline_template" "test-pipeline" {
  name = "template1-terraform"
}

resource "gocd_pipeline" "test-pipeline" {
  name                    = "pipeline1-terraform"
  group                   = "testing"
  template                = "${gocd_pipeline_template.test-pipeline.name}"
  lock_behavior           = "lockOnFailure"
  label_template          = "build-$${COUNT}"

  materials = [
    {
      type = "git"

      attributes {
        name        = "gocd-github"
        url         = "git@github.com:gocd/gocd"
        branch      = "feature/my-addition"
        destination = "gocd-dir"
      }
    },
  ]
}
