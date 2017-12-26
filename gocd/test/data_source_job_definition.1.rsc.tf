data "gocd_task_definition" "test" {
  type = "exec"

  run_if = [
    "passed",
  ]

  command = "/usr/local/bin/terraform"

  arguments = [
    "-debug",
    "version",
  ]

  working_directory = "/tmp/"
}

data "gocd_job_definition" "test" {
  name               = "job-name"
  run_instance_count = 3
  timeout            = 6

  environment_variables {
    name  = "USERNAME"
    value = "myusername"
  }

  environment_variables = [
    {
      name            = "PASSWORD"
      encrypted_value = "$R*YN:LDFIOH"
    },
    {
      name   = "HIDDEN"
      value  = "shown"
      secure = true
    },
  ]

  resources = [
    "alpha",
    "beta",
  ]

  tasks = [
    "${data.gocd_task_definition.test.json}",
  ]

  tabs = [
    {
      name = "Report"
      path = "report1.html"
    },
  ]

  artifacts = [
    {
      type   = "build"
      source = "web.war"
    },
  ]

  properties {
    name   = "coverage.class"
    source = "target/emma/coverage.xml"
    xpath  = "substring-before(//report/data/all/coverage[starts-with(@type,'class')]/@value, '%')"
  }
}
