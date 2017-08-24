package gocd

import (
	"context"
	"fmt"
	"net/url"
)

// AgentsService describes the HAL _link resource for the api response object for an agent objects.
type AgentsService service

// AgentsLinks describes the HAL _link resource for the api response object for a collection of agent objects.
//go:generate gocd-response-links-generator -type=AgentsLinks,AgentLinks
type AgentsLinks struct {
	Self *url.URL `json:"self"`
	Doc  *url.URL `json:"doc"`
}

// AgentLinks describes the HAL _link resource for the api response object for a single agent object.
type AgentLinks struct {
	Self *url.URL `json:"self"`
	Doc  *url.URL `json:"doc"`
	Find *url.URL `json:"find"`
}

// AgentsResponse describes the structure of the API response when listing collections of agent object.
type AgentsResponse struct {
	Links    *AgentsLinks `json:"_links,omitempty"`
	Embedded *struct {
		Agents []*Agent `json:"agents"`
	} `json:"_embedded,omitempty"`
}

// Agent describes a single agent object.
type Agent struct {
	UUID             string        `json:"uuid"`
	Hostname         string        `json:"hostname"`
	ElasticAgentID   string        `json:"elastic_agent_id"`
	ElasticPluginID  string        `json:"elastic_plugin_id"`
	IPAddress        string        `json:"ip_address"`
	Sandbox          string        `json:"sandbox"`
	OperatingSystem  string        `json:"operating_system"`
	FreeSpace        int           `json:"free_space"`
	AgentConfigState string        `json:"agent_config_state"`
	AgentState       string        `json:"agent_state"`
	Resources        []string      `json:"resources"`
	Environments     []string      `json:"environments"`
	BuildState       string        `json:"build_state"`
	BuildDetails     *BuildDetails `json:"build_details"`
	Links            *AgentLinks   `json:"_links,omitempty"`
	client           *Client
}

// AgentUpdate describes the structure for the PUT payload when updating an agent
type AgentUpdate struct {
	Hostname         string   `json:"hostname,omitempty"`
	Resources        []string `json:"resources,omitempty"`
	Environments     []string `json:"environments,omitempty"`
	AgentConfigState string   `json:"agent_config_state,omitempty"`
}

// AgentBulkUpdate describes the structure for the PUT payload when updating multiple agents
type AgentBulkUpdate struct {
	Uuids            []string                   `json:"uuids"`
	Operations       *AgentBulkOperationsUpdate `json:"operations,omitempty"`
	AgentConfigState string                     `json:"agent_config_state,omitempty"`
}

// AgentBulkOperationsUpdate describes the structure for a single Operation in AgentBulkUpdate the PUT payload when
// updating multiple agents
type AgentBulkOperationsUpdate struct {
	Environments *AgentBulkOperationUpdate `json:"environments,omitempty"`
	Resources    *AgentBulkOperationUpdate `json:"resources,omitempty"`
}

// AgentBulkOperationUpdate describes an action to be performed on an Environment or Resource during an agent update.
type AgentBulkOperationUpdate struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

// BuildDetails describes the builds being performed on this agent.
type BuildDetails struct {
	Links    *BuildDetailsLinks `json:"_links"`
	Pipeline string             `json:"pipeline"`
	Stage    string             `json:"stage"`
	Job      string             `json:"job"`
}

// BuildDetailsLinks describes the HAL structure for _link objects for the build details.
//go:generate gocd-response-links-generator -type=BuildDetailsLinks
type BuildDetailsLinks struct {
	Job      *url.URL `json:"job"`
	Stage    *url.URL `json:"stage"`
	Pipeline *url.URL `json:"pipeline"`
}

// RemoveLinks sets the `Link` attribute as `nil`. Used when rendering an `Agent` struct to JSON.
func (a *Agent) RemoveLinks() {
	a.Links = nil
}

// List will retrieve all agents, their status, and metadata from the GoCD Server.
func (s *AgentsService) List(ctx context.Context) ([]*Agent, *APIResponse, error) {
	r := AgentsResponse{}
	_, resp, err := s.client.getAction(ctx, &APIClientRequest{
		Path:         "agents",
		ResponseBody: &r,
		APIVersion:   apiV4,
	})

	for _, agent := range r.Embedded.Agents {
		agent.client = s.client
	}

	return r.Embedded.Agents, resp, err
}

// Get will retrieve a single agent based on the provided UUID.
func (s *AgentsService) Get(ctx context.Context, uuid string) (*Agent, *APIResponse, error) {
	return s.handleAgentRequest(ctx, "GET", uuid, nil)
}

// Update will modify the configuration for an existing agents.
func (s *AgentsService) Update(ctx context.Context, uuid string, agent AgentUpdate) (*Agent, *APIResponse, error) {
	return s.handleAgentRequest(ctx, "PATCH", uuid, &agent)
}

// Delete will remove an existing agent. Note: The agent must be disabled, and not currently building to be deleted.
func (s *AgentsService) Delete(ctx context.Context, uuid string) (string, *APIResponse, error) {
	return s.client.deleteAction(ctx, "agents/"+uuid, apiV4)
}

// BulkUpdate will change the configuration for multiple agents in a single request.
func (s *AgentsService) BulkUpdate(ctx context.Context, agents AgentBulkUpdate) (string, *APIResponse, error) {
	a := StringResponse{}
	_, resp, err := s.client.patchAction(ctx, &APIClientRequest{
		Path:         "agents",
		APIVersion:   apiV4,
		ResponseBody: &a,
		RequestBody:  agents,
	})
	return a.Message, resp, err
}

// JobRunHistory will return a list of Jobs run on this agent.
func (s *AgentsService) JobRunHistory(ctx context.Context, uuid string) ([]*Job, *APIResponse, error) {
	a := JobRunHistoryResponse{}
	_, resp, err := s.client.getAction(ctx, &APIClientRequest{
		Path:         fmt.Sprintf("agents/%s/job_run_history", uuid),
		APIVersion:   apiV4,
		ResponseBody: &a,
	})
	return a.Jobs, resp, err
}

// handleAgentRequest handles the flow to perform an HTTP action on an agent resource.
func (s *AgentsService) handleAgentRequest(ctx context.Context, action string, uuid string, body *AgentUpdate) (*Agent, *APIResponse, error) {
	a := Agent{client: s.client}
	_, resp, err := s.client.httpAction(ctx, &APIClientRequest{
		Method:       action,
		Path:         "agents/" + uuid,
		APIVersion:   apiV4,
		RequestBody:  body,
		ResponseBody: &a,
	})

	return &a, resp, err
}
