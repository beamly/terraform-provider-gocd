package gocd

import (
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
)

// SchemaProvider describing the required configs to interact with GoCD server. Environment variables can also be set:
//   baseurl        - GOCD_URL
//   username       - GOCD_USERNAME
//   password       - GOCD_PASSWORD
//   skip_ssl_check - GOCD_SKIP_SSL_CHECK
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"gocd_stage_definition": dataSourceGocdStageTemplate(),
			"gocd_job_definition":   dataSourceGocdJobTemplate(),
			"gocd_task_definition":  dataSourceGocdTaskDefinition(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gocd_pipeline_template": resourcePipelineTemplate(),
		},
		Schema: map[string]*schema.Schema{
			"baseurl": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["gocd_baseurl"],
				DefaultFunc: envDefault("GOCD_URL"),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["username"],
				DefaultFunc: envDefault("GOCD_USERNAME"),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["password"],
				DefaultFunc: envDefault("GOCD_PASSWORD"),
			},
			"skip_ssl_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["skip_ssl_check"],
				DefaultFunc: envDefault("GOCD_SKIP_SSL_CHECK"),
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"baseurl":  "URL for the GoCD Server",
		"username": "User to interact with the GoCD API with.",
		"password": "Password for User for GoCD API interaction.",
	}

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	var url, u, p string
	var cfg *gocd.Configuration

	if url = d.Get("baseurl").(string); url == "" {
		url = os.Getenv("GOCD_URL")
	}

	if u = d.Get("username").(string); u == "" {
		u = os.Getenv("GOCD_USERNAME")
	}
	if p = d.Get("password").(string); p == "" {
		p = os.Getenv("GOCD_PASSWORD")
	}
	nossl := d.Get("skip_ssl_check").(bool)

	cfg = &gocd.Configuration{
		Server:   url,
		Username: u,
		Password: p,
		SslCheck: !nossl,
	}

	return cfg.Client(), nil

}

func envDefault(e string) schema.SchemaDefaultFunc {
	return schema.MultiEnvDefaultFunc([]string{
		e,
	}, nil)
}
