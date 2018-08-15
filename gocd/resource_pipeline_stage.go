package gocd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beamly/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
)

const STAGE_TYPE_PIPELINE = "pipeline"
const STAGE_TYPE_PIPELINE_TEMPLATE = "template"

func resourcePipelineStage() *schema.Resource {
	stringArg := &schema.Schema{Type: schema.TypeString}
	optionalBoolArg := &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	return &schema.Resource{
		Create: resourcePipelineStageCreate,
		Read:   resourcePipelineStageRead,
		Update: resourcePipelineStageUpdate,
		Delete: resourcePipelineStageDelete,
		Exists: resourcePipelineStageExists,
		Importer: &schema.ResourceImporter{
			State: resourcePipelineStageImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"clean_working_directory": optionalBoolArg,
			"never_cleanup_artifacts": optionalBoolArg,
			"fetch_materials": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"jobs": {
				Type:             schema.TypeList,
				Required:         true,
				Elem:             stringArg,
				DiffSuppressFunc: supressJSONDiffs,
				//ValidateFunc:     validation.ValidateJsonString,
			},
			"manual_approval": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"success_approval"},
			},
			"success_approval": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"manual_approval"},
			},
			"authorization_users": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"success_approval", "authorization_roles"},
				Elem:          stringArg,
			},
			"authorization_roles": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"success_approval", "authorization_users"},
				Elem:          stringArg,
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
			"pipeline": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"pipeline_template"},
				Optional:      true,
				ForceNew:      true,
			},
			"pipeline_template": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"pipeline"},
				Optional:      true,
				ForceNew:      true,
			},
		},
	}
}

func resourcePipelineStageImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var pType, pipeline, name string
	var err error

	if pType, pipeline, name, err = parseGoCDPipelineStageId(d.Id()); err != nil {
		return nil, err
	}

	d.Set("name", name)

	if err := resourcePipelineStageSetPTypeName(d, pType, pipeline); err != nil {
		return nil, err
	}

	if err := resourcePipelineStageRead(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourcePipelineStageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	var pType, pipeline, name string
	var err error

	if pType, pipeline, name, err = parseGoCDPipelineStageId(d.Id()); err != nil {
		return false, err
	}

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if stage, err := retrieveStage(pType, name, pipeline, client); (err == nil) && (stage != nil) {
		d.SetId(fmt.Sprintf("%s/%s/%s", pType, pipeline, stage.Name))
		return true, nil
	}

	return false, err
}

func resourcePipelineStageCreate(d *schema.ResourceData, meta interface{}) (err error) {
	var existing *gocd.StageContainer
	var updated *gocd.StageContainer

	stage := gocd.Stage{
		Approval: &gocd.Approval{},
	}

	ingestStageConfig(d, &stage)

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	pipelineName, pType := pipelineNameType(d)

	if existing, err = getStageContainer(pType, pipelineName, client); err != nil {
		return err
	}

	(*existing).SetStages(cleanPlaceHolderStage((*existing).GetStages()))
	(*existing).AddStage(&stage)

	if updated, err = updateStageContainer(pType, existing, client); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", pType, (*updated).GetName(), stage.Name))

	return err
}

func resourcePipelineStageRead(d *schema.ResourceData, meta interface{}) error {
	var stage *gocd.Stage

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	pType, pipeline, name, err := parseGoCDPipelineStageId(d.Id())
	if stage, err = retrieveStage(pType, name, pipeline, client); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", pType, pipeline, stage.Name))

	d.Set("name", stage.Name)
	d.Set("fetch_materials", stage.FetchMaterials)
	d.Set("clean_working_directory", stage.CleanWorkingDirectory)
	d.Set("never_cleanup_artifacts", stage.NeverCleanupArtifacts)
	d.Set("resources", stage.Resources)

	d.Set(
		"environment_variables",
		ingestEnvironmentVariables(stage.EnvironmentVariables),
	)

	if appr := stage.Approval; appr != nil {
		if appr.Type == "manual" {
			d.Set("manual_approval", true)
		} else if appr.Type == "success_approval" {
			d.Set("success_approval", true)
		}
		if auth := stage.Approval.Authorization; auth != nil {
			if users := auth.Users; users != nil {
				d.Set("authorization_users", users)
			}
			if roles := auth.Roles; roles != nil {
				d.Set("authorization_roles", roles)
			}
		}
	}

	var s string
	if jobs := stage.Jobs; len(jobs) > 0 {
		stringJobs := []string{}
		for _, job := range jobs {
			if s, err = job.JSONString(); err != nil {
				return err
			}
			stringJobs = append(stringJobs, s)
		}

		d.Set("jobs", stringJobs)
	}

	if err := resourcePipelineStageSetPTypeName(d, pType, pipeline); err != nil {
		return err
	}

	return nil
}

func resourcePipelineStageUpdate(d *schema.ResourceData, meta interface{}) error {
	var existing *gocd.StageContainer
	var pType, pipeline, name string
	var stage *gocd.Stage
	var err error

	if pType, pipeline, name, err = parseGoCDPipelineStageId(d.Id()); err != nil {
		return err
	}
	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if stage, err = retrieveStage(pType, name, pipeline, client); stage == nil {
		if err != nil {
			return err
		}
		return fmt.Errorf("Could not find stage `%s` in pipeline/template `%s`", name, pipeline)
	}

	// If we are updating, make sure we are only adding jobs which we are responsible for. This avoids conflicts
	// and encourages state to be managed by tf.
	stage.Jobs = []*gocd.Job{}
	ingestStageConfig(d, stage)

	// Retrieve Pipeline object so we have the latest version
	if existing, err = getStageContainer(pType, pipeline, client); err != nil {
		return err
	}

	(*existing).SetStage(stage)
	if _, err = updateStageContainer(pType, existing, client); err != nil {
		return err
	}

	return nil
}

func resourcePipelineStageDelete(d *schema.ResourceData, meta interface{}) (err error) {
	var updated *gocd.StageContainer
	var existing *gocd.StageContainer

	pipeline, pType := pipelineNameType(d)

	client := meta.(*gocd.Client)

	// Make delete operation atomic
	client.Lock()
	defer client.Unlock()

	// Retrieve Pipeline object so we have the latest version
	if existing, err = getStageContainer(pType, pipeline, client); err != nil {
		return err
	}

	stages := cleanPlaceHolderStage((*existing).GetStages())

	cleanedStages := []*gocd.Stage{}
	stageName := d.Get("name").(string)
	for _, stage := range stages {
		if stage.Name != stageName {
			cleanedStages = append(cleanedStages, stage)
		}
	}

	if len(cleanedStages) == 0 {
		cleanedStages = append(cleanedStages, stagePlaceHolder())
	}

	(*existing).SetStages(cleanedStages)

	// Perform stage update
	if updated, err = updateStageContainer(pType, existing, client); err != nil {
		return err
	}

	if stage := (*updated).GetStage(stageName); stage != nil {
		return fmt.Errorf("could not delete stage `%s` as it does not exist", stageName)
	}

	return nil
}

func retrieveStage(pType string, stageName string, pipeline string, client *gocd.Client) (stage *gocd.Stage, err error) {
	var existing *gocd.StageContainer

	if existing, err = getStageContainer(pType, pipeline, client); err != nil {
		return nil, err
	}

	if stage = (*existing).GetStage(stageName); stage != nil {
		return stage, nil
	}

	return nil, nil
}

func pipelineNameType(d *schema.ResourceData) (pipelineName string, pType string) {
	if pipelineI, hasPipeline := d.GetOk("pipeline"); hasPipeline {
		return pipelineI.(string), STAGE_TYPE_PIPELINE
	}
	return d.Get("pipeline_template").(string), STAGE_TYPE_PIPELINE_TEMPLATE
}

func updateStageContainer(pType string, existing *gocd.StageContainer, client *gocd.Client) (updated *gocd.StageContainer, err error) {
	var sc gocd.StageContainer
	ctx := context.Background()
	if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		sc, _, err = client.PipelineTemplates.Update(ctx, (*existing).GetName(), (*existing).(*gocd.PipelineTemplate))
	} else if pType == STAGE_TYPE_PIPELINE {
		sc, _, err = client.PipelineConfigs.Update(ctx, (*existing).GetName(), (*existing).(*gocd.Pipeline))
	}
	return &sc, err
}

func getStageContainer(pType string, pipelineName string, client *gocd.Client) (existing *gocd.StageContainer, err error) {
	var sc gocd.StageContainer
	ctx := context.Background()
	if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		sc, _, err = client.PipelineTemplates.Get(ctx, pipelineName)
	} else if pType == STAGE_TYPE_PIPELINE {
		sc, _, err = client.PipelineConfigs.Get(ctx, pipelineName)
	}

	return &sc, err

}

func cleanPlaceHolderStage(stages []*gocd.Stage) []*gocd.Stage {
	cleanStages := []*gocd.Stage{}
	for _, stage := range stages {
		if stage.Name != PLACEHOLDER_NAME {
			cleanStages = append(cleanStages, stage)
		}
	}
	return cleanStages
}

func stagePlaceHolder() *gocd.Stage {
	return &gocd.Stage{
		Name: PLACEHOLDER_NAME,
		Jobs: []*gocd.Job{
			{Name: PLACEHOLDER_NAME},
		},
	}
}

func dataSourceStageParseManuallApproval(data *schema.ResourceData, doc *gocd.Stage) error {
	if doc.Approval == nil {
		doc.Approval = &gocd.Approval{}
	}
	doc.Approval.Type = "manual"
	doc.Approval.Authorization = &gocd.Authorization{}
	if users := data.Get("authorization_users").(*schema.Set).List(); len(users) > 0 {
		doc.Approval.Authorization.Users = decodeConfigStringList(users)
	} else if roles := data.Get("authorization_roles").(*schema.Set).List(); len(roles) > 0 {
		doc.Approval.Authorization.Roles = decodeConfigStringList(roles)
	}
	return nil
}

func dataSourceStageParseJobs(jobs []string, doc *gocd.Stage) error {
	for _, rawjob := range jobs {
		job := &gocd.Job{}
		if err := json.Unmarshal([]byte(rawjob), job); err != nil {
			return err
		}
		doc.Jobs = append(doc.Jobs, job)
	}
	return nil
}

func parseGoCDPipelineStageId(id string) (pType string, pipeline string, stage string, err error) {
	var r *regexp.Regexp
	r, _ = regexp.Compile(`^(template|pipeline)/([^/]+)/([^/]+)$`)

	if matches := r.FindAllStringSubmatch(id, -1); len(matches) == 1 {
		return matches[0][1], matches[0][2], matches[0][3], nil
	}

	return "", "", "", fmt.Errorf("could not parse the provided id `%s`", id)
}

func ingestStageConfig(d *schema.ResourceData, stage *gocd.Stage) {
	stage.Name = d.Get("name").(string)
	if manualApproval, ok := d.GetOk("manual_approval"); ok && manualApproval.(bool) {
		dataSourceStageParseManuallApproval(d, stage)
	} else if d.Get("success_approval").(bool) {
		stage.Approval.Type = "success"
		stage.Approval.Authorization = nil
	} else {
		stage.Approval = nil
	}

	stage.FetchMaterials = d.Get("fetch_materials").(bool)
	stage.CleanWorkingDirectory = d.Get("clean_working_directory").(bool)
	stage.NeverCleanupArtifacts = d.Get("never_cleanup_artifacts").(bool)

	if rJobs, hasJobs := d.GetOk("jobs"); hasJobs {
		if jobs := decodeConfigStringList(rJobs.([]interface{})); len(jobs) > 0 {
			dataSourceStageParseJobs(jobs, stage)
		}
	}

	if rawEnvVars, hasEnvVars := d.GetOk("environment_variables"); hasEnvVars {
		if envVars := rawEnvVars.([]interface{}); len(envVars) > 0 {
			stage.EnvironmentVariables = dataSourceGocdJobEnvVarsRead(envVars)
		}
	}
}

func ingestEnvironmentVariables(environmentVariables []*gocd.EnvironmentVariable) []map[string]interface{} {
	envVarMaps := []map[string]interface{}{}
	for _, rawEnvVar := range environmentVariables {
		envVarMap := map[string]interface{}{
			"name":   rawEnvVar.Name,
			"secure": rawEnvVar.Secure,
		}
		if rawEnvVar.Value != "" {
			envVarMap["value"] = rawEnvVar.Value
		}
		if rawEnvVar.EncryptedValue != "" {
			envVarMap["encrypted_value"] = rawEnvVar.EncryptedValue
		}
		envVarMaps = append(envVarMaps, envVarMap)
	}
	return envVarMaps
}
