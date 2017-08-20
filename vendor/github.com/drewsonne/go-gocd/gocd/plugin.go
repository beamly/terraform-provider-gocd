package gocd

import (
	"context"
	"fmt"
	"net/url"
)

// PluginsService exposes calls for interacting with Plugin objects in the GoCD API.
type PluginsService service

// PluginsResponseLinks describes the HAL _link resource for the api response object for a collection of agent objects.
//go:generate gocd-response-links-generator -type=PluginsResponseLinks,PluginLinks
type PluginsResponseLinks struct {
	Self *url.URL `json:"self"`
	Doc  *url.URL `json:"doc"`
}

// PluginLinks describes the HAL _link resource for the api response object for a collection of agent objects.
type PluginLinks struct {
	Self *url.URL `json:"self"`
	Doc  *url.URL `json:"doc"`
	Find *url.URL `json:"find"`
}

// PluginsResponse describes the response obejct for a plugin API call.
type PluginsResponse struct {
	Links    PluginsResponseLinks `json:"_links"`
	Embedded struct {
		PluginInfo []*Plugin `json:"plugin_info"`
	} `json:"_embedded"`
}

// Plugin describes a single plugin resource.
type Plugin struct {
	Links                     PluginLinks               `json:"_links"`
	ID                        string                    `json:"id"`
	Name                      string                    `json:"name"`
	DisplayName               string                    `json:"display_name"`
	Version                   string                    `json:"version"`
	Type                      string                    `json:"type"`
	PluggableInstanceSettings PluggableInstanceSettings `json:"pluggable_instance_settings"`
}

// PluggableInstanceSettings describes plugin configuration
type PluggableInstanceSettings struct {
	Configurations []PluginConfiguration `json:"configurations"`
	View           PluginView            `json:"view"`
}

// PluginView describes any view attached to a plugin.
type PluginView struct {
	Template string `json:"template"`
}

// List retrieves all plugins
func (ps *PluginsService) List(ctx context.Context) (*PluginsResponse, *APIResponse, error) {
	pr := PluginsResponse{}
	_, resp, err := ps.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/plugin_info",
		ResponseBody: &pr,
		APIVersion:   apiV2,
	})

	return &pr, resp, err
}

// Get retrieves information about a specific plugin.
func (ps *PluginsService) Get(ctx context.Context, name string) (*Plugin, *APIResponse, error) {
	p := &Plugin{}
	_, resp, err := ps.client.getAction(ctx, &APIClientRequest{
		Path:         fmt.Sprintf("admin/plugin_info/%s", name),
		ResponseBody: &p,
		APIVersion:   apiV2,
	})

	return p, resp, err
}
