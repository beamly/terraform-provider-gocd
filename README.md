# terraform-provider-gocd 0.1.22

[![GoDoc](https://godoc.org/github.com/beamly/terraform-provider-gocd/gocd?status.svg)](https://godoc.org/github.com/beamly/terraform-provider-gocd/gocd)
[![Build Status](https://travis-ci.org/beamly/terraform-provider-gocd.svg?branch=master)](https://travis-ci.org/beamly/terraform-provider-gocd)
[![codecov](https://codecov.io/gh/beamly/terraform-provider-gocd/branch/master/graph/badge.svg)](https://codecov.io/gh/beamly/terraform-provider-gocd)
[![codebeat badge](https://codebeat.co/badges/50ee50a6-46a9-497d-a85b-d929e6f00f65)](https://codebeat.co/projects/github-com-beamly-terraform-provider-gocd-master-0ff3887a-10ea-4399-acfe-b990029412b6)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fbeamly%2Fterraform-provider-gocd.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fbeamly%2Fterraform-provider-gocd?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/beamly/terraform-provider-gocd)](https://goreportcard.com/report/github.com/beamly/terraform-provider-gocd)

Terraform provider for GoCD Server

## Installation

    $ brew tap beamly/tap
    $ brew install terraform-provider-gocd
    $ tf-install-provider gocd

__NOTE__: `terraform` does not currently provide a way to easily install 3rd party providers. Until this is implemented,
the `tf-install-provider` utility can be used to copy the provider binary to the correct location.

__NOTE__: Given the requirements by hashicorp for SLA's to be set for importing providers into the `hashicorp-provider` github organisation, and as I am but one person on this, there's little chance this will become an officially sanctioned terraform provider. See [Terraform Provider Development Program](https://www.terraform.io/guides/terraform-provider-development-program.html):

> The expectation is to resolve all critical issues within 48 hours and all other issues within 5 business days. 

## Components

 - Data
     - [`gocd_task_definition`](#gocd_task_definition)
     - [`gocd_job_definition`](#gocd_job_definition)
 - Resources
     - [`gocd_pipeline`](#gocd_pipeline)
     - [`gocd_pipeline_template`](#gocd_pipeline_template)
     - [`gocd_pipeline_stage`](#gocd_pipeline_stage)
     - [`gocd_environment`](#gocd_environment)
     - [`gocd_environment_association`](#gocd_environment_association)

## Data


 - [`gocd_task_definition`](#gocd_task_definition)
 - [`gocd_job_definition`](#gocd_job_definition)

### gocd\_task\_definition

Generates json strings for GoCD task definitions

#### Example Usage

```hcl
data "gocd_task_definition" "my-task" {
  type = "exec"
  command = "terraform"
  arguments = ["init"]
}
```

#### Argument Reference

Task definition as defined in the [GoCD API](https://api.gocd.org/current/#the-task-object).

 - `type` - (Required) The type of a task. Can be one of exec, ant, nant, rake, fetch, pluggable_task.
 - `run_if` - (Optional) The run_if condition specifies when a task should be allowed to run. Can be one of passed, failed, any.

##### Task `type = "exec"`
 - `command` - (Required) The name of the executable. If the executable is not on PATH, you may also specify the full path.
 - `arguments`- (Optional) The list of arguments to be passed to the executable.
 - `working_directory` - (Optional) The directory in which the executable is to be executed.

##### Task `type = "ant"`
 - `build_file` - (Required) The path to Ant build file.
 - `target` - (Required) The Ant target(s) to run.
 - `working_directory` - (Optional) The directory in which the executable is to be executed.

##### Task `type = "fetch"`
 - `pipeline` - (Optional) The name of direct upstream pipeline or ancestor pipeline of one of the upstream pipelines on which the pipeline of the job depends on. .
 - `stage` - (Required) The name of the stage to fetch artifacts from.
 - `job`  - (Required) The name of the job to fetch artifacts from.
 - `source` - (Required) The path of the artifact directory or file of a specific job, relative to the sandbox directory. If the directory or file does not exist, the job is failed.
 - `is_source_a_file` - (Optional) Whether source is a file or directory.
 - `destination` - (Optional) The path of the directory where the artifact is fetched to.

### gocd\_job\_definition

Generates json strings for GoCD job definitions

#### Example Usage

```hcl
data "gocd_job_definition" "my-job" {
  name = "my-job"
  tasks = []
  environment_variables = [{
    name = "HOME"
    value = "/home/go"
  }]
}

output "my-job" {
  value = "${data.gocd_job_definition.my-job.json}"
}
```

#### Argument Reference

 - `name` - (Required) The name of the job.
 - `tasks` - (Required) A list of json strings defining a task definition for this job
 - `run_instance_count` - (Optional) The number of jobs to run. If set to null (default), one job will be created. If set to the literal string all, the job will be run on all agents. If set to a positive integer, the specified number of jobs will be created. Can be one of null, Integer, all.
 - `timeout` - (Optional) The time period(in minute) after which the job will be terminated by go if it has not generated any output.
 - `environment_variables` - (Optional) The list of environment variables defined here are set on the agents and can be used within your tasks. Each `environment_variables` block supports fields documented below.
 - `resources` - (Optional) The list of (String) resources that specifies the resource which the job requires to build. MUST NOT be specified along with elastic_profile_id.
 - `tabs` - (Optional) The list of tabs which let you add custom tabs within the job details page. Each `tabs` block supports fields documented below.
 - `artifacts` - (Optional) The list of artifacts specifies what files the agent will publish to the server. Each `artifacts` block supports fields documented below.
 - `properties` - (Optional) The list of properties of the build from XML files or artifacts created during your build. Each `properties` block supports fields documented below.
 - `elastic_profile_id` - (Optional) The id of the elastic profile, specifying this attribute would run the job on an elastic agent asociated with this profile. MUST NOT be specified along with resources. Since v16.10.0.

The `environment_variables` block supports:

 - `name` - (Required) The name of the environment variable.
 - `value` - (Optional) The value of the environment variable. One of `value` or `encrypted_value` must be set.
 - `encrypted_value` - (Optional) The encrypted value of the environment variable. One of `value` or `encrypted_value` must be set.
 - `secure` - Whether environment variable is secure or not. When set to `true`, encrypts the value if one is specified. The default value is `false`.

The `tabs` block supports:

 - `name` - (Required) The name of the tab which will appear in the Job detail page.
 - `path` - (Required) The relative path of a file in the server artifact destination directory of the job that will be render in the tab.

The `artifacts` block supports:

 - `type` - (Required) The type of the artifact. Can be one of test, build.
 - `source` - (Required) The file or folder to publish to the server.
 - `destination` - (Optional) The destination is relative to the artifacts folder of the current job instance on the server side. If it is not specified, the artifact will be stored in the root of the artifacts directory.

The `properties` block supports:

 - `name` - (Required) The name of the property.
 - `source` - (Optional) The relative path to the XML file containing the data that you want to use to create the property.
 - `xpath` - (Optional) The xpath that will be used to create property.

#### Attributes Reference

 - `json` - JSON encoded string of the job definition


## Resources

 - [`gocd_pipeline`](#gocd_pipeline)
 - [`gocd_pipeline_template`](#gocd_pipeline_template)
 - [`gocd_pipeline_stage`](#gocd_pipeline_stage)
 - [`gocd_environment`](#gocd_environment)
 - [`gocd_environment_association`](#gocd_environment_association)

### gocd\_pipeline

Provides support for creating pipelines in GoCD.

#### Example Usage

```hcl
resource "gocd_pipeline" "build" {
  name = "build"
  group = "terraform-provider-gocd"
  label_template = "0.0.$${COUNT}"
  materials = [
    {
      type = "git"
      attributes {
        name = "terraform-provider-gocd"
        url = "https://github.com/beamly/terraform-provider-gocd.git"
      }
    }
  ]
}
```

#### Argument Reference

 - `name` - (Required) The name of the pipeline.
 - `group` - (Required) The name of the pipeline group to deploy into.
 - `materials` - (Required) The list of materials to be used by pipeline. At least one material must be specified. Each `materials` block supports fields documented below.
 - `label_template` - (Optional)  The label template to customise the pipeline instance label.
 - `enable_pipeline_locking` - (Optional)  Whether pipeline is locked to run single instance or not.
 - `template` - (Optional)  The name of the template used by pipeline. A `gocd_pipeline_stage` can not be assigned to a `gocd_pipeline` it `template` is set.
 - `parameters` - (Optional) A [map](https://www.terraform.io/docs/configuration/variables.html#maps) of parameters to be used for substitution in a pipeline or pipeline template.
 - `environment_variables` - (Optional) The list of environment variables that will be passed to all tasks (commands) that are part of this pipeline. Each `environment_variables` block supports fields documented below.

The `environment_variables` block supports:

 - `name` - (Required) The name of the environment variable.
 - `value` - (Optional) The value of the environment variable. One of `value` or `encrypted_value` must be set.
 - `encrypted_value` - (Optional) The encrypted value of the environment variable. One of `value` or `encrypted_value` must be set.
 - `secure` - Whether environment variable is secure or not. When set to `true`, encrypts the value if one is specified. The default value is `false`.

Type `materials` block supports:

 - `type` (Required) The type of a material. Can be one of git, dependency.
 - `attributes` (Required) A [map](https://www.terraform.io/docs/configuration/variables.html#maps) of attributes for each material type. See the [GoCD API Documentation](https://api.gocd.org/current/#the-pipeline-material-object) for each material type attributes.


#### Attributes Reference

 - `version` - The current version of the resource configuration in GoCD.

#### Import

Pipelines can be imported using the pipeline name, e.g.

    $ terraform import gocd_pipeline.pipeline pipeline-name

### gocd\_pipeline\_template

Provides support for creating pipeline templates in GoCD.

#### Example Usage

```hcl
resource "gocd_pipeline_template" "terraform-builder" {
  name = "terraform-build-template"
}
```

#### Argument Reference

 - `name` - (Required) The name of the pipeline template.

#### Attributes Reference

 - `version` - The current version of the resource configuration in GoCD.

#### Import

Pipeline Templates can be imported using the pipeline template name, e.g.

    $ terraform import gocd_pipeline_template.pipeline-template pipeline-template-name

### gocd\_pipeline\_stage

Provides support for creating stages for pipelines or pipeline templates in GoCD.

#### Example Usage

```hcl
resource "gocd_pipeline_stage" "build" {
  name = "plan"
  pipeline = "plan"
  jobs = [
  <<JOB
 {
  "name": "plan",
  "tasks": [{
    "type": "exec",
    "attributes": {
      "run_if": ["passed"],
      "command": "terraform",
      "arguments": ["plan"]
    }
  }]
 }
  JOB
  ]
}
```

#### Import

Pipeline Stages can be imported using a type identifier (`template` or `pipeline`), the pipeline or template name, and the stage name, e.g.

    $ terraform import gocd_pipeline_stage.pipeline-stage pipeline/pipeline-name/pipeline-stage
    $ ...
    $ terraform import gocd_pipeline_stage.template-stage template/template-name/template-stage

### gocd\_environment

Provides support for creating environments in GoCD.

#### Example Usage

```hcl
resource "gocd_environment" "testing" {
  name = "testing"
}
```

#### Argument Rference

 - `name` - (Required) Name of the environment to create.

#### Attributes Reference

 - `version` - The current version of the resource configuration in GoCD.

#### Import

Environments can be imported using the environment name, e.g.

    $ terraform import gocd_environment.environment-name environment-name

### gocd\_environment\_association

Provides support for associating pipelines and environments in GoCD.

__NOTE:__ There is an intention to support agents and environment variables in the future.

#### Example Usage

```hcl
resource "gocd_environment_association" "build-in-testing" {
  environment = "${gocd_environment.testing.name}"
  pipeline = "${gocd_pipeline.build.name}"
}

resource "gocd_environment" "testing" {
  name = "testing"
}

resource "gocd_pipeline" "build" {
  name = "build"
  # ...
}
```

#### Argument Reference

 - `environment` - (Required) The name of the environment which the resource is being associated to.
 - `pipeline` - (Required) The name of the pipeline to associate to the environment


#### Attributes Reference

 - `version` - The current version of the resource configuration in GoCD.

#### Import

Environment Associations can be imported using the environment name, a type identifier, and type value, e.g.

    $ terraform import gocd_environment_association.association-name environment-name/p/pipeline-name

__NOTE__: Currently only pipeline (`p`) type identifiers are supported.

## Development Setup

Get project and install dependencies using glide:

    $ go get github.com/beamly/terraform-provider-gocd
    $ cd $GOPATH/src/github.com/beamly/terraform-provider-gocd/
    $ make before_install

Then you can run tests as follows:

    $ make testacc

## Demo

You will need docker and terraform >= 0.10.0 installed for this demo to work.

Either build the provider with `go build` or download it from the gihub repository. If you download it, make sure the binary is in the current folder.

	$ go build

Spin up the test gocd server, with endpoint at http://127.0.0.1:8153/go/

    $ make provision-test-gocd && sh ./scripts/wait-for-test-server.sh

Then initialise and apply the configuration.

    $ terraform init
    $ terraform apply

When you're finished, run:

    $ make teardown-test-gocd

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fbeamly%2Fterraform-provider-gocd.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fbeamly%2Fterraform-provider-gocd?ref=badge_large)
