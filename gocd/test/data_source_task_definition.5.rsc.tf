data "gocd_task_definition" "test" {
  type = "pluggable"

  run_if = [
    "passed",
  ]

  plugin_id      = "plugin.id"
  plugin_version = "plugin.version"

  configuration {
    "key"   = "key1"
    "value" = "value1"
  }

  configuration {
    "key"   = "key2"
    "value" = "value2"
  }
}
