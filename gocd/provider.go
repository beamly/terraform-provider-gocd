package gocd

import (
	"crypto/tls"
	"github.com/beamly/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"net/http"
	"os"
	"strings"
)

func Provider() terraform.ResourceProvider {
	return SchemaProvider()
}

// SchemaProvider describing the required configs to interact with GoCD server. Environment variables can also be set:
//   baseurl        - GOCD_URL
//   username       - GOCD_USERNAME
//   password       - GOCD_PASSWORD
//   skip_ssl_check - GOCD_SKIP_SSL_CHECK
func SchemaProvider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			//"gocd_stage_definition": dataSourceGocdStageTemplate(),
			"gocd_job_definition":  dataSourceGocdJobTemplate(),
			"gocd_task_definition": dataSourceGocdTaskDefinition(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gocd_environment":             resourceEnvironment(),
			"gocd_environment_association": resourceEnvironmentAssociation(),
			"gocd_pipeline_template":       resourcePipelineTemplate(),
			"gocd_pipeline":                resourcePipeline(),
			"gocd_pipeline_stage":          resourcePipelineStage(),
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
	var rUrl, rU, rP, rB interface{}
	var ok, nossl, b bool
	var cfg *gocd.Configuration

	if rUrl, ok = d.GetOk("baseurl"); ok {
		if url, ok = rUrl.(string); !ok || url == "" {
			url = os.Getenv("GOCD_URL")
		}
	}
	log.Printf("[DEBUG] Using GoCD config 'baseurl': %s", url)

	if rU, ok = d.GetOk("username"); ok {
		if u, ok = rU.(string); !ok || u == "" {
			u = os.Getenv("GOCD_USERNAME")
		}
	}
	log.Printf("[DEBUG] Using GoCD config 'username': %s", u)

	if rP, ok = d.GetOk("password"); ok {
		if p, ok = rP.(string); !ok || p == "" {
			p = os.Getenv("GOCD_PASSWORD")
		}
	}
	log.Printf("[DEBUG] Using GoCD config 'password': %s", rP)

	if rB, ok = d.GetOk("skip_ssl_check"); ok {
		if b, ok = rB.(bool); !ok {
			nossl = false
		} else {
			nossl = b
		}
	}
	log.Printf("[DEBUG] Using GoCD config 'skip_ssl_check': %t", nossl)

	cfg = &gocd.Configuration{
		Server:       url,
		Username:     u,
		Password:     p,
		SkipSslCheck: nossl,
	}

	hClient := &http.Client{}

	if strings.HasPrefix(cfg.Server, "https") {
		log.Printf("[DEBUG] GoCD is using https.")
		hClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.SkipSslCheck},
		}
	} else {
		hClient.Transport = http.DefaultTransport
	}

	// Add API logging
	hClient.Transport = logging.NewTransport("GoCD", hClient.Transport)
	gc := gocd.NewClient(cfg, hClient)

	// No-longer supported by go-gocd
	// versionString := terraform.VersionString()
	// gc.params.UserAgent = fmt.Sprintf("(%s %s) Terraform/%s", runtime.GOOS, runtime.GOARCH, versionString)

	return gc, nil

}

func envDefault(e string) schema.SchemaDefaultFunc {
	return schema.MultiEnvDefaultFunc([]string{
		e,
	}, nil)
}
