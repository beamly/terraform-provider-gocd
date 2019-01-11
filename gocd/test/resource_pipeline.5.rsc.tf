## START pipeline_template.gocd-image-build-deploy
# CMD terraform import gocd_pipeline_template.gocd-image-build-deploy "gocd-image-build-deploy"
resource "gocd_pipeline_template" "gocd-image-build-deploy" {
  name = "gocd-image-build-deploy"
}

# CMD terraform import gocd_pipeline_stage.build "build"
resource "gocd_pipeline_stage" "build" {
  name              = "build"
  pipeline_template = "${gocd_pipeline_template.gocd-image-build-deploy.name}"
  fetch_materials   = true

  jobs = [
    "${data.gocd_job_definition.build.json}",
  ]

  environment_variables = [
    {
      name  = "IMAGE"
      value = "#{Image}"
    },
  ]
}

data "gocd_job_definition" "build" {
  name = "build"

  tasks = [
    "${data.gocd_task_definition.gocd-image-build-deploy_build_build_0.json}",
  ]

  resources = [
    "v5.80",
  ]

  environment_variables = [
    {
      name  = "IMAGE"
      value = "#{Image}"
    },
  ]
}

data "gocd_task_definition" "gocd-image-build-deploy_build_build_0" {
  type = "exec"

  run_if = [
    "passed",
  ]

  command = "make"

  arguments = [
    "build",
    "#{Image}\"complex\"$${Env}",
  ]
}

# CMD terraform import gocd_pipeline_stage.clean "clean"
resource "gocd_pipeline_stage" "clean" {
  name              = "clean"
  pipeline_template = "${gocd_pipeline_template.gocd-image-build-deploy.name}"
  fetch_materials   = true

  jobs = [
    "${data.gocd_job_definition.clean.json}",
  ]
}

data "gocd_job_definition" "clean" {
  name = "clean"

  tasks = [
    "${data.gocd_task_definition.gocd-image-build-deploy_clean_clean_0.json}",
  ]

  resources = [
    "v5.80",
  ]
}

data "gocd_task_definition" "gocd-image-build-deploy_clean_clean_0" {
  type = "exec"

  run_if = [
    "passed",
  ]

  command = "docker"

  arguments = [
    "system",
    "prune",
    "-f",
  ]
}

# CMD terraform import gocd_pipeline_stage.deploy "deploy"
resource "gocd_pipeline_stage" "deploy" {
  name              = "deploy"
  pipeline_template = "${gocd_pipeline_template.gocd-image-build-deploy.name}"
  fetch_materials   = true

  jobs = [
    "${data.gocd_job_definition.deploy.json}",
  ]
}

data "gocd_job_definition" "deploy" {
  name = "deploy"

  tasks = [
    "${data.gocd_task_definition.gocd-image-build-deploy_deploy_deploy_0.json}",
  ]

  resources = [
    "v5.80",
  ]
}

data "gocd_task_definition" "gocd-image-build-deploy_deploy_deploy_0" {
  type = "exec"

  command = "make"

  arguments = [
    "deploy",
  ]
}

## END
resource "gocd_pipeline" "terraform-image" {
  name     = "terraform-image"
  group    = "ecsagent"
  template = "${gocd_pipeline_template.gocd-image-build-deploy.name}"

  parameters {
    Image = "terraform"
  }

  environment_variables = [
    {
      name  = "PACKER_ANSIBLE_VERSION"
      value = "2.0.2.0"
    },
    {
      name  = "INSTALL_ANSIBLE_DEPENDENCIES"
      value = "true"
    },
  ]

  materials = [
    {
      type = "git"

      attributes {
        url    = "git@github.com:company/gocd-ecsagents.git"
        branch = "master"

        //        auto_update = true
      }
    },
    {
      type = "dependency"

      attributes {
        pipeline = "${gocd_pipeline.test-pipeline.name}"
        stage    = "${gocd_pipeline_stage.clean.name}"

        //        auto_update = true
      }
    },
  ]
}

resource "gocd_pipeline" "test-pipeline" {
  name     = "base-image"
  group    = "ecsagent"
  template = "${gocd_pipeline_template.gocd-image-build-deploy.name}"

  parameters {
    Image = "base"
  }

  environment_variables = [
    {
      name  = "PACKER_ANSIBLE_VERSION"
      value = "2.0.2.0"
    },
    {
      name  = "INSTALL_ANSIBLE_DEPENDENCIES"
      value = "true"
    },
  ]

  materials = [
    {
      type = "git"

      attributes {
        url = "git@github.com:company/gocd-ecsagents.git"

        filter = [
          "company-gocd-agents/Dockerfile.base",
          "Makefile",
          "company-gocd-agents/files/base/",
        ]

        invert_filter = true
        branch        = "master"

        //        auto_update = true
      }
    },
  ]
}
