package gocd

import (
	"encoding/json"
	"github.com/beamly/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
)

// codebeat:disable[LOC]
func dataSourceGocdJobTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGocdJobTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"run_instance_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"environment_variables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type: schema.TypeString,
							// ConflictsWith can only be applied to top level configs.
							// A custom validation will need to be used.
							//ConflictsWith: []string{"encrypted_value"},
							Optional: true,
						},
						"encrypted_value": {
							Type: schema.TypeString,
							// ConflictsWith can only be applied to top level configs.
							// A custom validation will need to be used.
							//ConflictsWith: []string{"value"},
							Optional: true,
						},
						"secure": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
					},
				},
			},
			"resources": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"elastic_profile_id"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"elastic_profile_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"resources"},
			},
			"tabs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"artifacts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"properties": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"xpath": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
// codebeat:enable[LOC]

func dataSourceGocdJobTemplateRead(d *schema.ResourceData, meta interface{}) error {

	tasks := []*gocd.Task{}
	for _, rawTask := range d.Get("tasks").([]interface{}) {
		task := gocd.Task{}
		err := json.Unmarshal([]byte(rawTask.(string)), &task)
		if err != nil {
			return err
		}
		tasks = append(tasks, &task)
	}

	j := gocd.Job{
		Name:  d.Get("name").(string),
		Tasks: tasks,
	}

	if ric, ok := d.GetOk("run_instance_count"); ok {
		j.RunInstanceCount = ric.(int)
	}

	if to, ok := d.GetOk("timeout"); ok {
		j.Timeout = gocd.TimeoutField(to.(int))
	}

	if elasticProfile, ok := d.GetOk("elastic_profile_id"); ok {
		j.ElasticProfileID = elasticProfile.(string)
	}

	if elasticProfile, ok := d.GetOk("elastic_profile_id"); ok {
		j.ElasticProfileID = elasticProfile.(string)
	}

	if envVars, ok := d.Get("environment_variables").([]interface{}); ok && len(envVars) > 0 {
		j.EnvironmentVariables = dataSourceGocdJobEnvVarsRead(envVars)
	}

	if props, ok := d.Get("properties").([]interface{}); ok && len(props) > 0 {
		j.Properties = dataSourceGocdJobPropertiesRead(props)
	}

	if resources := d.Get("resources").(*schema.Set).List(); len(resources) > 0 {
		if rscs := decodeConfigStringList(resources); len(rscs) > 0 {
			j.Resources = rscs
		}
	}

	if resources := d.Get("tabs").([]interface{}); len(resources) > 0 {
		j.Tabs = []*gocd.Tab{}
		for _, rawTab := range resources {
			tabMap := rawTab.(map[string]interface{})
			j.Tabs = append(j.Tabs, &gocd.Tab{
				Name: tabMap["name"].(string),
				Path: tabMap["path"].(string),
			})
		}
	}

	if resources := d.Get("artifacts").([]interface{}); len(resources) > 0 {
		j.Artifacts = []*gocd.Artifact{}
		for _, rawArtifact := range resources {
			artifactMap := rawArtifact.(map[string]interface{})
			j.Artifacts = append(j.Artifacts, &gocd.Artifact{
				Type:        artifactMap["type"].(string),
				Source:      artifactMap["source"].(string),
				Destination: artifactMap["destination"].(string),
			})
		}
	}

	return definitionDocFinish(d, j)
}

func dataSourceGocdJobPropertiesRead(rawProps []interface{}) []*gocd.JobProperty {
	props := []*gocd.JobProperty{}
	for _, propRaw := range rawProps {
		propStruct := &gocd.JobProperty{}
		prop := propRaw.(map[string]interface{})

		if name, ok := prop["name"]; ok {
			propStruct.Name = name.(string)
		}

		if name, ok := prop["source"]; ok {
			propStruct.Source = name.(string)
		}

		if name, ok := prop["xpath"]; ok {
			propStruct.XPath = name.(string)
		}
		props = append(props, propStruct)
	}
	return props
}

func dataSourceGocdJobEnvVarsRead(rawEnvVars []interface{}) []*gocd.EnvironmentVariable {
	envVars := []*gocd.EnvironmentVariable{}
	for _, envVarRaw := range rawEnvVars {
		envVarStruct := &gocd.EnvironmentVariable{}
		envVar := envVarRaw.(map[string]interface{})

		if name, ok := envVar["name"]; ok {
			envVarStruct.Name = name.(string)
		}

		if val, ok := envVar["value"]; ok {
			envVarStruct.Value = val.(string)
		}

		if encrypted, ok := envVar["encrypted_value"]; ok {
			envVarStruct.EncryptedValue = encrypted.(string)
		}

		if secure, ok := envVar["secure"]; ok {
			envVarStruct.Secure = secure.(bool)
		}

		envVars = append(envVars, envVarStruct)
	}

	return envVars

}
