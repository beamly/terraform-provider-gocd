package gocd

import (
	"context"
	"encoding/json"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineTemplateCreate,
		Read:   resourcePipelineTemplateRead,
		Update: resourcePipelineTemplateUpdate,
		Delete: resourcePipelineTemplateDelete,
		Exists: resourcePipelineTemplateExists,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"stages": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: supressJSONDiffs,
				},
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePipelineTemplateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	}
	pt, _, err := meta.(*gocd.Client).PipelineTemplates.Get(context.Background(), name)
	exists := (pt.Name == name) && (err == nil)
	return exists, err
}

func resourcePipelineTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	}

	stages := extractStages(d)
	pt, _, err := meta.(*gocd.Client).PipelineTemplates.Create(context.Background(), name, stages)
	if err != nil {
		return err
	}

	return readPipelineTemplate(d, pt)
}

func resourcePipelineTemplateRead(d *schema.ResourceData, meta interface{}) error {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	}

	pt, resp, err := meta.(*gocd.Client).PipelineTemplates.Get(context.Background(), name)
	if err != nil {
		if resp.HTTP.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	return readPipelineTemplate(d, pt)

}

func resourcePipelineTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	}

	version := d.Get("version")
	stages := extractStages(d)
	pt, _, err := meta.(*gocd.Client).PipelineTemplates.Update(context.Background(), name, version.(string), stages)
	if err != nil {
		return err
	}

	return readPipelineTemplate(d, pt)

}

func resourcePipelineTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	if ptname, hasName := d.GetOk("name"); hasName {
		_, _, err := meta.(*gocd.Client).PipelineTemplates.Delete(context.Background(), ptname.(string))
		if err != nil {
			return err
		}
	}

	return nil
}

func extractStages(d *schema.ResourceData) []*gocd.Stage {
	stages := []*gocd.Stage{}
	for _, rawstage := range d.Get("stages").([]interface{}) {
		stage := gocd.Stage{}
		json.Unmarshal([]byte(rawstage.(string)), &stage)
		stages = append(stages, &stage)
	}
	return stages
}

func readPipelineTemplate(d *schema.ResourceData, p *gocd.PipelineTemplate) error {
	d.SetId(p.Name)

	stages := []string{}

	for _, stage := range p.Stages {
		bdy, err := stage.JSONString()
		if err != nil {
			return err
		}
		stages = append(stages, bdy)
	}

	if err := d.Set("stages", stages); err != nil {
		return err
	}

	d.Set("version", p.Version)

	return nil
}
