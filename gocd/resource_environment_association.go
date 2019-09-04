package gocd

import (
	"context"
	"fmt"
	"github.com/beamly/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourceEnvironmentAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentAssociationCreate,
		Read:   resourceEnvironmentAssociationRead,
		Delete: resourceEnvironmentAssociationDelete,
		Exists: resourceEnvironmentAssociationExists,
		Importer: &schema.ResourceImporter{
			State: resourceEnvironmentAssociationImport,
		},
		Schema: map[string]*schema.Schema{
			"environment": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pipeline": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				//	Optional:      true,
				//	ConflictsWith: []string{"agent", "environment_variable"},
			},
			// TODO implement this
			//"agent": {
			//	Type:          schema.TypeString,
			//	Optional:      true,
			//	ConflictsWith: []string{"pipeline", "environment_variable"},
			//},
			//"environment_variable": {
			//	Type:          schema.TypeList,
			//	Optional:      true,
			//	ConflictsWith: []string{"agent", "pipeline"},
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"name": {
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//			"value": {
			//				Type: schema.TypeString,
			//				// ConflictsWith can only be applied to top level configs.
			//				// A custom validation will need to be used.
			//				//ConflictsWith: []string{"encrypted_value"},
			//				Optional: true,
			//			},
			//			"encrypted_value": {
			//				Type: schema.TypeString,
			//				// ConflictsWith can only be applied to top level configs.
			//				// A custom validation will need to be used.
			//				//ConflictsWith: []string{"value"},
			//				Optional: true,
			//			},
			//			"secure": {
			//				Type:     schema.TypeBool,
			//				Default:  false,
			//				Optional: true,
			//			},
			//		},
			//	},
			//},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEnvironmentAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	environment := d.Get("environment").(string)
	pipeline := d.Get("pipeline").(string)

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()
	env, _, err := client.Environments.Patch(context.Background(), environment, &gocd.EnvironmentPatchRequest{
		Pipelines: &gocd.PatchStringAction{
			Add:    []string{pipeline},
			Remove: []string{},
		},
	})
	if err != nil {
		return err
	}
	d.SetId(environmentAssociationId(environment, pipeline, "", ""))
	d.Set("version", env.Version)
	return nil
}

func resourceEnvironmentAssociationRead(d *schema.ResourceData, meta interface{}) error {
	id := strings.Split(d.Id(), "/")
	environment := id[0]
	//associationType := id[1]
	value := id[2]
	client := meta.(*gocd.Client)
	env, _, err := client.Environments.Get(context.Background(), environment)
	if err != nil {
		return err
	}

	for _, p := range env.Pipelines {
		if p.Name == value {
			d.Set("pipeline", p.Name)
		}
	}
	d.Set("version", env.Version)

	return nil
}

func resourceEnvironmentAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	id := strings.Split(d.Id(), "/")
	environment := id[0]
	//associationType := id[1]
	value := id[2]

	client := meta.(*gocd.Client)
	_, _, err := client.Environments.Patch(context.Background(), environment, &gocd.EnvironmentPatchRequest{
		Pipelines: &gocd.PatchStringAction{
			Add:    []string{},
			Remove: []string{value},
		},
	})
	return err
}

func resourceEnvironmentAssociationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id := strings.Split(d.Id(), "/")
	environment := id[0]
	//associationType := id[1]
	value := id[2]

	client := meta.(*gocd.Client)
	env, _, err := client.Environments.Get(context.Background(), environment)
	if err != nil {
		return false, err
	}
	for _, p := range env.Pipelines {
		if p.Name == value {
			return true, nil
		}
	}
	return false, err
}

func resourceEnvironmentAssociationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := strings.Split(d.Id(), "/")
	environment := id[0]
	//associationType := id[1]
	value := id[2]
	d.Set("environment", environment)
	d.Set("pipeline", value)
	return []*schema.ResourceData{d}, nil
}

func environmentAssociationId(env string, pipeline string, agent string, envvar string) string {
	var envAssociationType string
	var value string
	if pipeline != "" {
		envAssociationType = "p"
		value = pipeline
	} else if agent != "" {
		envAssociationType = "a"
		value = pipeline
	} else if envvar != "" {
		envAssociationType = "e"
		value = pipeline
	}

	return fmt.Sprintf("%s/%s/%s", env, envAssociationType, value)
}
