package gocd

import (
	"context"
	"github.com/beamly/go-gocd/gocd"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
				ValidateFunc: RegexRuleset(map[string]string{
					`^[a-zA-Z0-9_\-]{1}`:                  "first character of %q (%q) must be alphanumeric, underscore, or dot",
					`^[a-zA-Z0-9_\-]{1}[a-zA-Z0-9_\-.]*$`: "only alphanumeric, underscores, hyphens, or dots allowed in %q (%q)",
				}),
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
				Computed: true,
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
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"materials": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"git",
								"svn",
								"hg",
								"p4",
								"tfs",
								"dependency",
								"package",
								"plugin",
							}, true),
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
									"submodule_folder": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"shallow_clone": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"destination": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"pipeline": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"stage": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"auto_update": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
										Removed:  "The `auto_update` attribute has been disabled until a way to manage updates atomically has been devised.",
									},
									"invert_filter": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"filter": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     materialFilterResource(),
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

func materialFilterResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ignore": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcePipelineCreate(d *schema.ResourceData, meta interface{}) (err error) {
	var p *gocd.Pipeline

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if p, err = extractPipeline(d); err != nil {
		return
	}
	// GoCD does not allow pipelines to have _no_ stages. As this provider has split stages and pipelines into separate
	// resources, we need to add a place holder stage if there are no stages yet configured and we're not using a
	// template.
	if (p.Stages == nil || len(p.Stages) == 0) && p.Template == "" {
		p.Stages = []*gocd.Stage{
			stagePlaceHolder(),
		}
	}

	group := d.Get("group").(string)
	pc, _, err := client.PipelineConfigs.Create(context.Background(), group, p)
	return readPipeline(d, pc, err)
}

func resourcePipelineRead(d *schema.ResourceData, meta interface{}) error {

	d.Set("name", d.Get("name").(string))

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	ctx := context.Background()
	pc, _, err := client.PipelineConfigs.Get(ctx, d.Id())
	if err := readPipeline(d, pc, err); err != nil {
		return err
	}

	pgs, _, err := client.PipelineGroups.List(ctx, "")
	d.Set("group", pgs.GetGroupByPipelineName(d.Id()).Name)
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
	name := d.Id()

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()
	if p, _, err := client.PipelineConfigs.Get(context.Background(), name); err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return false, nil
		} else {
			return false, err
		}
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
		m := gocd.Material{}

		mat := rawMaterial.(map[string]interface{})
		m.Ingest(mat)
		if mattr1, ok1 := mat["attributes"].([]interface{}); ok1 {
			if mattr2, ok2 := mattr1[0].(map[string]interface{}); ok2 {
				if filterI, ok3 := mattr2["filter"]; ok3 {
					if filterSet, ok4 := filterI.(*schema.Set); ok4 {
						filters := filterSet.List()
						if len(filters) > 0 {
							if filter, ok5 := filters[0].(map[string]interface{}); ok5 {
								if ignore, ok6 := filter["ignore"]; ok6 {
									mattr2["filter"] = map[string]interface{}{
										"ignore": decodeConfigStringList(ignore.([]interface{})),
									}
								}
							}
						}
					}
				}
				m.IngestAttributes(mattr2)
			}
		}

		ms = append(ms, m)

	}
	return ms, nil
}

func readPipelineMaterials(d *schema.ResourceData, materials []gocd.Material) error {
	materialImports := make([]interface{}, len(materials))
	for i, m := range materials {
		attrs := m.Attributes.GenerateGeneric()
		filterSet := schema.NewSet(
			schema.HashResource(materialFilterResource()),
			[]interface{}{},
		)
		if filters, ok1 := attrs["filter"]; ok1 {
			if filterI, ok2 := filters.(map[string]interface{}); ok2 {
				if len(filterI) > 0 {
					filterSet.Add(filterI)
				}
				attrs["filter"] = filterSet
			}
		}
		materialImports[i] = map[string]interface{}{
			"type":       m.Type,
			"attributes": []interface{}{attrs},
		}
	}
	if err := d.Set("materials", materialImports); err != nil {
		return err
	}
	return nil
}

//func extractPipelineMaterialFilter(attr interface{}) *gocd.MaterialFilter {
//	filterI := attr.([]interface{})
//	var mf *gocd.MaterialFilter
//	if len(filterI) > 0 {
//		filtersI := filterI[0].(map[string]interface{})
//		filters := filtersI["ignore"].([]interface{})
//		mf = &gocd.MaterialFilter{
//			Ignore: decodeConfigStringList(filters),
//		}
//	}
//	return mf
//}

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
