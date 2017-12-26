package gocd

import (
	"context"
	"errors"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create:   resourcePipelineCreate,
		Read:     resourcePipelineRead,
		Update:   resourcePipelineUpdate,
		Delete:   resourcePipelineDelete,
		Exists:   resourcePipelineExists,
		Importer: resourcePipelineStateImport(),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label_template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_pipeline_locking": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
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
			"materials": &schema.Schema{
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"attributes": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"branch": {
										Type:             schema.TypeString,
										Optional:         true,
										Computed:         true,
										DiffSuppressFunc: supressMaterialBranchDiff,
									},
									"shallow_clone": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"destination": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"pipeline": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"stage": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"auto_update": {
										Type:     schema.TypeBool,
										Optional: true,
										Removed:  "The `auto_update` attribute has been disabled until a way to manage updates atomically has been devised.",
									},
									"invert_filter": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"filter": {
										Type:     schema.TypeList,
										Optional: true,
										//Elem: &schema.Schema{
										//	Type: schema.TypeString,
										//},
										MaxItems: 1,
										MinItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ignore": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourcePipelineCreate(d *schema.ResourceData, meta interface{}) (err error) {
	var group string
	var p *gocd.Pipeline
	if ptgroup, hasGroup := d.GetOk("group"); hasGroup {
		group = ptgroup.(string)
	}
	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if p, err = extractPipeline(d); err != nil {
		return err
	}
	if (p.Stages == nil || len(p.Stages) == 0) && p.Template == "" {
		p.Stages = []*gocd.Stage{
			stagePlaceHolder(),
		}
	}
	pc, _, err := client.PipelineConfigs.Create(context.Background(), group, p)
	return readPipeline(d, pc, err)
}

func resourcePipelineRead(d *schema.ResourceData, meta interface{}) error {
	var name string
	if pname, hasName := d.GetOk("name"); hasName {
		name = pname.(string)
		d.SetId(name)
		d.Set("name", name)
	}
	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	ctx := context.Background()
	pc, _, err := client.PipelineConfigs.Get(ctx, name)
	if err := readPipeline(d, pc, err); err != nil {
		return err
	}

	pgs, _, err := client.PipelineGroups.List(ctx, "")
	d.Set("group", pgs.GetGroupByPipelineName(name).Name)
	return nil
}

func resourcePipelineUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	var name string
	var p *gocd.Pipeline
	if pname, hasName := d.GetOk("name"); hasName {
		name = pname.(string)
		d.SetId(name)
		d.Set("name", name)
	}

	templateToPipeline, templateChange := isSwitchToTemplate(d)

	if p, err = extractPipeline(d); err != nil {
		return err
	}

	client := meta.(*gocd.Client)
	ctx := context.Background()
	client.Lock()
	defer client.Unlock()

	existing, _, err := client.PipelineConfigs.Get(ctx, name)

	if templateChange && !templateToPipeline {
		p.Stages = nil
	} else if templateToPipeline {
		p.Stages = []*gocd.Stage{stagePlaceHolder()}
	} else {
		p.Stages = existing.Stages
	}

	p.Version = existing.Version
	pc, _, err := client.PipelineConfigs.Update(ctx, name, p)
	return readPipeline(d, pc, err)
}

func resourcePipelineDelete(d *schema.ResourceData, meta interface{}) error {
	var name string
	if pname, hasName := d.GetOk("name"); hasName {
		name = pname.(string)
	}
	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	_, _, err := client.PipelineConfigs.Delete(context.Background(), name)
	return err
}

func resourcePipelineExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	} else {
		return false, errors.New("`name` can not be empty")
	}

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()
	if p, _, err := client.PipelineConfigs.Get(context.Background(), name); err != nil {
		return false, err
	} else {
		return (p.Name == name), nil
	}
}

func resourcePipelineStateImport() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			d.Set("name", d.Id())
			return []*schema.ResourceData{d}, nil
		},
	}
}

func extractPipeline(d *schema.ResourceData) (p *gocd.Pipeline, err error) {
	p = &gocd.Pipeline{}

	if template, hasTemplate := d.GetOk("template"); hasTemplate {
		p.Template = template.(string)
	}

	if pipelineLocking, hasPipelineLocking := d.GetOk("enable_pipeline_locking"); hasPipelineLocking {
		p.EnablePipelineLocking = pipelineLocking.(bool)
	}

	p.Name = d.Get("name").(string)

	rawMaterials := d.Get("materials")
	if materials := rawMaterials.([]interface{}); len(materials) > 0 {
		if p.Materials, err = extractPipelineMaterials(materials); err != nil {
			return nil, err
		}
	}

	rawParameters := d.Get("parameters")
	if parameters := rawParameters.(map[string]interface{}); len(parameters) > 0 {
		p.Parameters = extractPipelineParameters(parameters)
	}

	if envVars, ok := d.Get("environment_variables").([]interface{}); ok && len(envVars) > 0 {
		p.EnvironmentVariables = dataSourceGocdJobEnvVarsRead(envVars)
	}

	return p, nil
}

func extractPipelineParameters(rawProperties map[string]interface{}) []*gocd.Parameter {
	ps := []*gocd.Parameter{}
	for key, value := range rawProperties {
		ps = append(ps, &gocd.Parameter{
			Name:  key,
			Value: value.(string),
		})
	}
	return ps
}

func extractPipelineMaterials(rawMaterials []interface{}) ([]gocd.Material, error) {
	ms := []gocd.Material{}
	for _, rawMaterial := range rawMaterials {
		mat := rawMaterial.(map[string]interface{})
		m := gocd.Material{}

		if mType, ok := mat["type"]; ok {
			m.Type = mType.(string)
		}

		if mAttributes, ok := mat["attributes"]; ok {

			attr := gocd.MaterialAttributes{}

			rawAttr := mAttributes.([]interface{})[0].(map[string]interface{})
			for attrKey, attrValue := range rawAttr {
				if s, ok := attrValue.(string); ok && s == "" {
					continue
				}
				switch attrKey {
				case "url":
					attr.URL = attrValue.(string)
				case "destination":
					attr.Destination = attrValue.(string)
				case "filter":
					attr.Filter = extractPipelineMaterialFilter(attrValue)
				case "invert_filter":
					attr.InvertFilter = attrValue.(bool)
				case "name":
					attr.Name = attrValue.(string)
				case "auto_update":
					attr.AutoUpdate = attrValue.(bool)
				case "branch":
					attr.Branch = attrValue.(string)
				case "submodule_folder":
					attr.SubmoduleFolder = attrValue.(string)
				case "shallow_clone":
					attr.ShallowClone = attrValue.(bool)
				case "pipeline":
					attr.Pipeline = attrValue.(string)
				case "stage":
					attr.Stage = attrValue.(string)
				default:
					return nil, fmt.Errorf("Unexpected material attribute: `%s:%s`", attrKey, attrValue)
				}
			}

			switch m.Type {
			case "dependency":
				if attr.Name == attr.Pipeline {
					attr.Name = ""
				}
				//case "git":
				//	if attr.Branch == "" {
				//		attr.Branch = "master"
				//	}
			}

			m.Attributes = attr
		}
		ms = append(ms, m)

	}
	return ms, nil
}

func readPipelineMaterials(d *schema.ResourceData, materials []gocd.Material) error {
	materialImports := make([]interface{}, len(materials))
	for i, m := range materials {
		materialMap := make(map[string]interface{})
		materialMap["type"] = m.Type

		attrs := map[string]interface{}{
			"url": m.Attributes.URL,
			//"auto_update":   m.Attributes.AutoUpdate,
			"branch":        m.Attributes.Branch,
			"destination":   m.Attributes.Destination,
			"invert_filter": m.Attributes.InvertFilter,
			"stage":         m.Attributes.Stage,
			"pipeline":      m.Attributes.Pipeline,
			"shallow_clone": m.Attributes.ShallowClone,
		}

		if m.Type == "dependency" {
			if m.Attributes.Name != m.Attributes.Pipeline {
				attrs["name"] = m.Attributes.Name
			}
		} else {
			attrs["name"] = m.Attributes.Name
		}

		filter := make([]map[string]interface{}, 1)
		if m.Attributes.Filter != nil {
			filter[0] = map[string]interface{}{
				"ignore": m.Attributes.Filter.Ignore,
			}
			attrs["filter"] = filter
		}

		materialMap["attributes"] = []map[string]interface{}{attrs}
		materialImports[i] = materialMap
	}
	if err := d.Set("materials", materialImports); err != nil {
		return err
	}
	return nil
}

func extractPipelineMaterialFilter(attr interface{}) *gocd.MaterialFilter {
	filterI := attr.([]interface{})
	var mf *gocd.MaterialFilter
	if len(filterI) > 0 {
		filtersI := filterI[0].(map[string]interface{})
		filters := filtersI["ignore"].([]interface{})
		mf = &gocd.MaterialFilter{
			Ignore: decodeConfigStringList(filters),
		}
	}
	return mf
}

func readPipeline(d *schema.ResourceData, p *gocd.Pipeline, err error) error {
	if err != nil {
		return err
	}

	d.SetId(p.Name)
	if p.Template != "" {
		d.Set("template", p.Template)
	}

	if p.LabelTemplate != "" && p.LabelTemplate != "${COUNT}" {
		d.Set("label_template", p.LabelTemplate)
	}

	d.Set("enable_pipeline_locking", p.EnablePipelineLocking)
	d.Set(
		"environment_variables",
		ingestEnvironmentVariables(p.EnvironmentVariables),
	)

	err = readPipelineMaterials(d, p.Materials)

	if len(p.Parameters) > 0 {
		rawParams := make(map[string]string, len(p.Parameters))
		for _, param := range p.Parameters {
			rawParams[param.Name] = param.Value
		}
		d.Set("parameters", rawParams)
	}

	return err
}

func isSwitchToTemplate(d *schema.ResourceData) (templateToPipeline bool, change bool) {
	change = d.HasChange("template")
	if !change {
		return false, false
	}
	if template, hasTemplate := d.GetOk("template"); hasTemplate {
		return template == "", change
	}
	return templateToPipeline, change
}

func supressMaterialBranchDiff(k, old, new string, d *schema.ResourceData) bool {
	if old == "" && new == "master" || old == "master" && new == "" {
		return true
	}
	return false
}
