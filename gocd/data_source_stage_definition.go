package gocd

import (
	"encoding/json"
	"github.com/beamly/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"strconv"
)

func dataSourceGocdStageDefinition() *schema.Resource {
	stringArg := &schema.Schema{Type: schema.TypeString}
	optionalBoolArg := &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	return &schema.Resource{
		Read: dataSourceGocdStageDefinitionRead,
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
			"approval": {
				Optional: true,
				Type:     schema.TypeSet,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"manual",
								"success",
							}, true),
						},
						"authorization": {
							Optional: true,
							Type:     schema.TypeSet,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"users": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     stringArg,
									},
									"roles": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     stringArg,
									},
								},
							},
						},
					},
				},
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
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGocdStageDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	stage := gocd.Stage{
		Name: d.Get("name").(string),
		Approval: &gocd.Approval{
			Type: "success",
			Authorization: &gocd.Authorization{
				Users: []string{},
				Roles: []string{},
			},
		},
	}

	dataSourceStageParseApproval(d, &stage)

	stage.FetchMaterials = d.Get("fetch_materials").(bool)
	stage.CleanWorkingDirectory = d.Get("clean_working_directory").(bool)
	stage.NeverCleanupArtifacts = d.Get("never_cleanup_artifacts").(bool)

	if rJobs, hasJobs := d.GetOk("jobs"); hasJobs {
		if jobs := decodeConfigStringList(rJobs.([]interface{})); len(jobs) > 0 {
			dataSourceStageParseJobs2(jobs, &stage)
		}
	}

	if rawEnvVars, hasEnvVars := d.GetOk("environment_variables"); hasEnvVars {
		if envVars := rawEnvVars.([]interface{}); len(envVars) > 0 {
			stage.EnvironmentVariables = dataSourceGocdJobEnvVarsRead(envVars)
		}
	}

	jsonDoc, err := json.MarshalIndent(stage, "", "  ")
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	jsonString := string(jsonDoc)
	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(jsonString)))

	return nil
}

func dataSourceStageParseApproval(data *schema.ResourceData, doc *gocd.Stage) error {

	if resources := data.Get("approval").(*schema.Set).List(); len(resources) > 0 {
		log.Printf("[DEBUG] XXX resources=%+v\n", resources)

		var users, roles []interface{}

		for _, rawApproval := range resources {
			log.Printf("[DEBUG] XXX rawApproval=%+v\n", rawApproval)

			approvalMap := rawApproval.(map[string]interface{})
			log.Printf("[DEBUG] XXX approvalMap=%+v\n", approvalMap)

			if sAuth := approvalMap["authorization"].(*schema.Set).List(); len(sAuth) > 0 {
				log.Printf("[DEBUG] XXX sAuth=%+v\n", sAuth)

				for _, rawAuth := range sAuth {
					log.Printf("[DEBUG] XXX rawAuth=%+v\n", rawAuth)

					auth := rawAuth.(map[string]interface{})
					log.Printf("[DEBUG] XXX auth=%+v\n", auth)

					users = auth["users"].(*schema.Set).List()
					roles = auth["roles"].(*schema.Set).List()
				}
			}

			doc.Approval = &gocd.Approval{
				Type: approvalMap["type"].(string),
				Authorization: &gocd.Authorization{
					Users: decodeConfigStringList(users),
					Roles: decodeConfigStringList(roles),
				},
			}

			log.Printf("[DEBUG] XXX %+v\n", doc.Approval)

		}
	}

	return nil
}

func dataSourceStageParseJobs2(jobs []string, doc *gocd.Stage) error {
	for _, rawjob := range jobs {
		job := &gocd.Job{}
		if err := json.Unmarshal([]byte(rawjob), job); err != nil {
			return err
		}
		doc.Jobs = append(doc.Jobs, job)
	}
	return nil
}
