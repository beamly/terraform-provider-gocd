package gocd

import (
	"context"
	"fmt"
	"strings"
)

// PipelineConfigsService describes the HAL _link resource for the api response object for a pipelineconfig
type PipelineConfigsService service

// PipelineConfig describes the configuration for a pipeline
type PipelineConfig struct{}

// PipelineConfigRequest describes a request object for creating or updating pipelines
type PipelineConfigRequest struct {
	Group    string    `json:"group"`
	Pipeline *Pipeline `json:"pipeline"`
}

// Update a pipeline configuration
func (pcs *PipelineConfigsService) Update(ctx context.Context, group string, name string, p *Pipeline) (*Pipeline, *APIResponse, error) {

	pt := &PipelineConfigRequest{
		Group:    group,
		Pipeline: p,
	}

	req, err := pcs.client.NewRequest("PUT", "admin/pipelines/"+name, pt, apiV4)
	if err != nil {
		return nil, nil, err
	}

	req.HTTP.Header.Set("If-Match", fmt.Sprintf("\"%s\"", p.Version))

	pc := Pipeline{}
	resp, err := pcs.client.Do(ctx, req, &pc, responseTypeJSON)
	return &pc, resp, err

}

// Create a pipeline configuration
func (pcs *PipelineConfigsService) Create(ctx context.Context, group string, p *Pipeline) (*Pipeline, *APIResponse, error) {
	pt := &PipelineConfigRequest{
		Group:    group,
		Pipeline: p,
	}
	pc := Pipeline{}

	_, resp, err := pcs.client.postAction(ctx, &APIClientRequest{
		Path:         "admin/pipelines",
		APIVersion:   apiV4,
		RequestBody:  pt,
		ResponseBody: &pc,
	})
	if err == nil {
		pc.Version = strings.Replace(resp.HTTP.Header.Get("Etag"), "\"", "", -1)
	}

	return &pc, resp, err
}

// Delete a pipeline configuration
func (pcs *PipelineConfigsService) Delete(ctx context.Context, name string) (string, *APIResponse, error) {
	return pcs.client.deleteAction(ctx, "admin/pipelines/"+name, apiV4)
}
