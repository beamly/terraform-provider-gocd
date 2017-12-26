resource "gocd_pipeline" "pipeline1" {
  group = "test-auto-update"
  name  = "pipeline1"

  materials = [{
    type = "git"

    attributes {
      url         = "https://github.com/gocd/gocd"
      auto_update = false
    }
  }]
}

resource "gocd_pipeline" "pipeline2" {
  group = "test-auto-update"
  name  = "pipeline2"

  materials = [{
    type = "git"

    attributes {
      url         = "https://github.com/gocd/gocd"
      auto_update = true
    }
  }]
}
