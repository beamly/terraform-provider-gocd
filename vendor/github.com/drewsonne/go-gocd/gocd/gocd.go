package gocd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	// Version of the gocd library in the event that we change it for the user agent.
	libraryVersion = "1"
	// UserAgent to be used when calling the GoCD agent.
	userAgent = "go-gocd/" + libraryVersion
	// Version 1 of the GoCD API.
	apiV1 = "application/vnd.go.cd.v1+json"
	// Version 2 of the GoCD API.
	apiV2 = "application/vnd.go.cd.v2+json"
	// Version 3 of the GoCD API.
	apiV3 = "application/vnd.go.cd.v3+json"
	// Version 4 of the GoCD API.
	apiV4 = "application/vnd.go.cd.v4+json"
)

//Body Response Types
const (
	responseTypeXML  = "xml"
	responseTypeJSON = "json"
)

// StringResponse handles the unmarshaling of the single string response from DELETE requests.
type StringResponse struct {
	Message string `json:"message"`
}

// APIResponse encapsulates the net/http.Response object, a string representing the Body, and a gocd.Request object
// encapsulating the response from the API.
type APIResponse struct {
	HTTP    *http.Response
	Body    string
	Request *APIRequest
}

// APIRequest encapsulates the net/http.Request object, and a string representing the Body.
type APIRequest struct {
	HTTP *http.Request
	Body string
}

// Client struct which acts as an interface to the GoCD Server. Exposes resource service handlers.
type Client struct {
	client *http.Client

	BaseURL  *url.URL
	Username string
	Password string

	UserAgent string

	Agents            *AgentsService
	PipelineGroups    *PipelineGroupsService
	Stages            *StagesService
	Jobs              *JobsService
	PipelineTemplates *PipelineTemplatesService
	Pipelines         *PipelinesService
	PipelineConfigs   *PipelineConfigsService
	Configuration     *ConfigurationService
	Encryption        *EncryptionService
	Plugins           *PluginsService

	common service
	cookie string
}

// PaginationResponse is a struct used to handle paging through resposnes.
type PaginationResponse struct {
	Offset   int `json:"offset"`
	Total    int `json:"total"`
	PageSize int `json:"page_size"`
}

// service is a generic service encapsulating the client for talking to the GoCD server.
type service struct {
	client *Client
}

// Auth structure wrapping the Username and Password variables, which are used to get an Auth cookie header used for
// subsequent requests.
type Auth struct {
	Username string
	Password string
}

// Configuration object used to initialise a gocd lib client to interact with the GoCD server.
type Configuration struct {
	Server   string `yaml:"server"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	SslCheck bool   `yaml:"ssl_check,omitempty"`
}

// HasAuth checks whether or not we have the required Username/Password variables provided.
func (c *Configuration) HasAuth() bool {
	return (c.Username != "") && (c.Password != "")
}

// Client returns a client which allows us to interact with the GoCD Server.
func (c *Configuration) Client() *Client {
	return NewClient(c, nil)
}

// NewClient creates a new client based on the provided configuration payload, and optionally a custom httpClient to
// allow overriding of http client structures.
func NewClient(cfg *Configuration, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if strings.HasPrefix(cfg.Server, "https") {
		if !cfg.SslCheck {
			httpClient.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !cfg.SslCheck},
			}
		}
	}

	baseURL, _ := url.Parse(cfg.Server)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	c.common.client = c

	c.Username = cfg.Username
	c.Password = cfg.Password

	c.Agents = (*AgentsService)(&c.common)
	c.PipelineGroups = (*PipelineGroupsService)(&c.common)
	c.Stages = (*StagesService)(&c.common)
	c.Jobs = (*JobsService)(&c.common)
	c.PipelineTemplates = (*PipelineTemplatesService)(&c.common)
	c.Pipelines = (*PipelinesService)(&c.common)
	c.PipelineConfigs = (*PipelineConfigsService)(&c.common)
	c.Configuration = (*ConfigurationService)(&c.common)
	c.Encryption = (*EncryptionService)(&c.common)
	c.Plugins = (*PluginsService)(&c.common)

	return c
}

// NewRequest creates an HTTP requests to the GoCD API endpoints.
func (c *Client) NewRequest(method, urlStr string, body interface{}, apiVersion string) (*APIRequest, error) {
	request := &APIRequest{}

	// I'm not sure how to get this method to return an error intentionally for testing. For testing purposes, I've
	// added a switch so that the error handling in dependent methods can be tested.
	if os.Getenv("GOCD_RAISE_ERROR_NEW_REQUEST") == "yes" {
		return request, errors.New("Mock Testing Error")
	}

	rel, err := url.Parse("api/" + urlStr)
	if err != nil {
		return request, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)

		enc := json.NewEncoder(buf)
		enc.SetIndent("", "  ")
		err := enc.Encode(body)

		if err != nil {
			return nil, err
		}
		bdy, _ := ioutil.ReadAll(buf)
		request.Body = string(bdy)

		buf = new(bytes.Buffer)
		enc = json.NewEncoder(buf)
		enc.SetIndent("", "  ")
		enc.Encode(body)
	}

	req, err := http.NewRequest(method, u.String(), buf)
	request.HTTP = req
	if err != nil {
		return request, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if apiVersion != "" {
		req.Header.Set("Accept", apiVersion)
	}
	req.Header.Set("User-Agent", c.UserAgent)

	if c.cookie == "" {
		if c.Username != "" && c.Password != "" {
			req.SetBasicAuth(c.Username, c.Password)
		}
	} else {
		req.Header.Set("Cookie", c.cookie)
	}

	return request, nil
}

// Do takes an HTTP request and resposne the response from the GoCD API endpoint.
func (c *Client) Do(ctx context.Context, req *APIRequest, v interface{}, responseType string) (*APIResponse, error) {

	req.HTTP = req.HTTP.WithContext(ctx)

	response := &APIResponse{
		Request: req,
	}

	resp, err := c.client.Do(req.HTTP)
	if err != nil {
		return nil, err
	}

	response.HTTP = resp
	err = CheckResponse(response.HTTP)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			bdy, err := ioutil.ReadAll(resp.Body)
			if responseType == responseTypeXML {
				err = xml.Unmarshal(bdy, v)
			} else {
				err = json.Unmarshal(bdy, v)
			}
			response.Body = string(bdy)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

// CheckResponse asserts that the http response status code was 2xx.
func CheckResponse(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 400 {
		bdy, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf(
			"Received HTTP Status '%s': '%s'",
			response.Status,
			bdy,
		)
	}
	return nil
}
