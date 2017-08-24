package gocd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func dataSourceGocdTaskDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGocdTaskDefinitionRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"run_if": {
				Type:     schema.TypeList,
				MaxItems: 3,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"command": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"arguments": {
				Type:     schema.TypeList,
				MaxItems: 3,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"build_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"working_directory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nant_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_source_a_file": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"job": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stage": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGocdTaskDefinitionRead(d *schema.ResourceData, meta interface{}) error {

	task := gocd.Task{
		Type:       d.Get("type").(string),
		Attributes: gocd.TaskAttributes{},
	}

	if runIf := decodeConfigStringList(d.Get("run_if").([]interface{})); len(runIf) > 0 {
		task.Attributes.RunIf = runIf
	}

	switch task.Type {
	case "exec":
		dataSourceGocdTaskBuildExec(&task, d)
	case "ant":
		dataSourceGocdTaskBuildAnt(&task, d)
	case "nant":
		dataSourceGocdTaskBuildNant(&task, d)
	case "rake":
		dataSourceGocdRakeTemplate(&task, d)
	case "fetch":
		dataSourceGocdFetchTemplate(&task, d)
	case "pluggable":
		dataSourceGocdPluggabeTemplate(&task, d)
	default:
		return fmt.Errorf("unexpected `gocd.Task.Type`: '%s'", task.Type)
	}

	jsonDoc, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	jsonString := string(jsonDoc)
	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(jsonString)))

	return nil
}

// Extract attributes for Pluggable Task
func dataSourceGocdPluggabeTemplate(t *gocd.Task, d *schema.ResourceData) error {
	t.Attributes.PluginConfiguration = &gocd.TaskPluginConfiguration{}
	if pid, ok := d.GetOk("plugin_id"); ok {
		t.Attributes.PluginConfiguration.ID = pid.(string)
	} else {
		return errors.New("Missing pluging id")
	}

	if pv, ok := d.GetOk("plugin_version"); ok {
		t.Attributes.PluginConfiguration.Version = pv.(string)
	} else {
		return errors.New("Missing pluging version")
	}

	configs := []gocd.PluginConfigurationKVPair{}
	if cfg, ok := d.GetOk("configuration"); ok {
		for _, kv := range cfg.([]interface{}) {
			kvm := kv.(map[string]interface{})
			configs = append(configs, gocd.PluginConfigurationKVPair{
				Key:   kvm["key"].(string),
				Value: kvm["value"].(string),
			})
		}
	}
	t.Attributes.Configuration = configs

	return nil
}

// Extract attributes for Fetch Task
func dataSourceGocdFetchTemplate(t *gocd.Task, d *schema.ResourceData) {
	if pipe, ok := d.GetOk("pipeline"); ok {
		t.Attributes.Pipeline = pipe.(string)
	}

	if s, ok := d.GetOk("stage"); ok {
		t.Attributes.Stage = s.(string)
	}

	if j, ok := d.GetOk("job"); ok {
		t.Attributes.Job = j.(string)
	}

	if isaf, ok := d.GetOk("is_source_a_file"); ok && isaf.(bool) {
		t.Attributes.IsSourceAFile = "true"
	} else {
		t.Attributes.IsSourceAFile = "false"
	}

	if d, ok := d.GetOk("destination"); ok {
		t.Attributes.Destination = d.(string)
	}

	if s, ok := d.GetOk("source"); ok {
		t.Attributes.Source = s.(string)
	}

}

// Extract attributes for rake tas
func dataSourceGocdRakeTemplate(t *gocd.Task, d *schema.ResourceData) {
	if bf, ok := d.GetOk("build_file"); ok {
		t.Attributes.BuildFile = bf.(string)
	}

	if template, ok := d.GetOk("target"); ok {
		t.Attributes.Target = template.(string)
	}

	if wd, ok := d.GetOk("working_directory"); ok {
		t.Attributes.WorkingDirectory = wd.(string)
	}

}

// Extract attributes for nant task
func dataSourceGocdTaskBuildNant(t *gocd.Task, d *schema.ResourceData) {
	if bf, ok := d.GetOk("build_file"); ok {
		t.Attributes.BuildFile = bf.(string)
	}

	if template, ok := d.GetOk("target"); ok {
		t.Attributes.Target = template.(string)
	}

	if nantPath, ok := d.GetOk("nant_path"); ok {
		t.Attributes.NantPath = nantPath.(string)
	}

	if wd, ok := d.GetOk("working_directory"); ok {
		t.Attributes.WorkingDirectory = wd.(string)
	}

}

// Extract attributes for ant task
func dataSourceGocdTaskBuildAnt(t *gocd.Task, d *schema.ResourceData) {
	if wd, ok := d.GetOk("working_directory"); ok {
		t.Attributes.WorkingDirectory = wd.(string)
	}

	if bf, ok := d.GetOk("build_file"); ok {
		t.Attributes.BuildFile = bf.(string)
	}

	if template, ok := d.GetOk("target"); ok {
		t.Attributes.Target = template.(string)
	}
}

// Extract attributes for exec task
func dataSourceGocdTaskBuildExec(t *gocd.Task, d *schema.ResourceData) {
	if bf, ok := d.GetOk("build_file"); ok {
		t.Attributes.BuildFile = bf.(string)
	}

	if cmd, ok := d.GetOk("command"); ok {
		t.Attributes.Command = cmd.(string)
	}

	if args := decodeConfigStringList(d.Get("arguments").([]interface{})); len(args) > 0 {
		t.Attributes.Arguments = args
	}

	if wd, ok := d.GetOk("working_directory"); ok {
		t.Attributes.WorkingDirectory = wd.(string)
	}

	if template, ok := d.GetOk("target"); ok {
		t.Attributes.Target = template.(string)
	}

}
