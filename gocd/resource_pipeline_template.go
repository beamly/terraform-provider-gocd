package gocd

import (
	"context"
	"errors"
	"github.com/beamly/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

const PLACEHOLDER_NAME = "TERRAFORM_PLACEHOLDER"

func resourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineTemplateCreate,
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
		},
	}
}

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

	//stages := extractStages(d)
	// As a pipeline must be created with a stage, when we first create the pipeline, add a dummy placeholder stage.
	// This will be cleaned up by any stage creation actions.

	placeholderStages := []*gocd.Stage{
		stagePlaceHolder(),
	}
	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	pt, _, err := client.PipelineTemplates.Create(context.Background(), name, placeholderStages)
	return readPipelineTemplate(d, pt, err)
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

	return nil
}
