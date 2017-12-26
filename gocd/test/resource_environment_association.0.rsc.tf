resource "gocd_environment" "test-environment" {
  name = "test-environment"
}

resource "gocd_pipeline" "test-pipeline" {
  name  = "test-pipeline"
  group = "test-group"

  materials = [
    {
      type = "git"

      attributes {
        name   = "gocd-src"
        url    = "git@github.com:gocd/gocd"
        branch = "master"

        //        auto_update = "true"
      }
    },
  ]
}

resource "gocd_environment_association" "test-environment-association" {
  environment = "${gocd_environment.test-environment.name}"
  pipeline    = "${gocd_pipeline.test-pipeline.name}"
}
