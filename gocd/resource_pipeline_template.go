package gocd

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/beamly/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

const PLACEHOLDER_NAME = "TERRAFORM_PLACEHOLDER"

// codebeat:disable[LOC]
func resourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineTemplateCreate,
		Update: resourcePipelineTemplateUpdate,
		Read:   resourcePipelineTemplateRead,
		Delete: resourcePipelineTemplateDelete,
		Exists: resourcePipelineTemplateExists,
		Importer: &schema.ResourceImporter{
			State: resourcePipelineTemplateImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stages": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: supressJSONDiffs,
			},
		},
	}
}

// codebeat:enable[LOC]

func resourcePipelineTemplateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("name", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourcePipelineTemplateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	} else {
		return false, errors.New("`name` can not be empty")
	}

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if p, _, err := client.PipelineTemplates.Get(context.Background(), name); err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return (p.Name == name), nil
	}
}

func resourcePipelineTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	}

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	pt := gocd.PipelineTemplate{}
	resourcePipelineTemplateParseStages(d, &pt)

	pt2, _, err := client.PipelineTemplates.Create(context.Background(), name, pt.Stages)
	return readPipelineTemplate(d, pt2, err)
}

func resourcePipelineTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	}

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	pt := gocd.PipelineTemplate{
		Name:    name,
		Version: d.Get("version").(string),
	}

	resourcePipelineTemplateParseStages(d, &pt)

	pt2, _, err := client.PipelineTemplates.Update(context.Background(), name, &pt)
	return readPipelineTemplate(d, pt2, err)
}

func resourcePipelineTemplateRead(d *schema.ResourceData, meta interface{}) error {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	}

	var pt *gocd.PipelineTemplate
	var resp *gocd.APIResponse
	var err error
	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if pt, resp, err = client.PipelineTemplates.Get(context.Background(), name); err != nil {
		if resp.HTTP.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	return readPipelineTemplate(d, pt, nil)

}

func resourcePipelineTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	if ptname, hasName := d.GetOk("name"); hasName {
		client := meta.(*gocd.Client)
		client.Lock()
		defer client.Unlock()

		if _, _, err := client.PipelineTemplates.Delete(context.Background(), ptname.(string)); err != nil {
			return err
		}
	}

	return nil
}

func readPipelineTemplate(d *schema.ResourceData, p *gocd.PipelineTemplate, err error) error {

	if err != nil {
		return err
	}

	d.SetId(p.Name)
	d.Set("version", p.Version)

	var s string
	if stages := p.Stages; len(stages) > 0 {
		stringStages := []string{}
		for _, stage := range stages {
			if s, err = stage.JSONString(); err != nil {
				return err
			}
			stringStages = append(stringStages, s)
		}

		d.Set("stages", stringStages)
	}

	return nil
}

func resourcePipelineTemplateParseStages(d *schema.ResourceData, pt *gocd.PipelineTemplate) error {

	if rStages, hasStages := d.GetOk("stages"); hasStages {
		if stages := decodeConfigStringList(rStages.([]interface{})); len(stages) > 0 {
			for _, rawstage := range stages {
				stage := &gocd.Stage{}
				if err := json.Unmarshal([]byte(rawstage), stage); err != nil {
					return err
				}
				pt.Stages = append(pt.Stages, stage)
			}
		}
	}

	return nil
}
