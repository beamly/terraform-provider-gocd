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
	"sync"
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
	responseTypeText = "text"
)

//Logging Environment variables
const (
	gocdLogLevel = "GOCD_LOG"
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
	clientMu sync.Mutex // clientMu protects the client during multi-threaded calls
	client   *http.Client

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
	Environments      *EnvironmentsService
	Properties        *PropertiesService

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
		if cfg.SkipSslCheck {
			httpClient.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.SkipSslCheck},
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
	c.Environments = (*EnvironmentsService)(&c.common)
	c.Properties = (*PropertiesService)(&c.common)

	SetupLogging()

	return c
}

// Lock the client until release
func (c *Client) Lock() {
	c.clientMu.Lock()
}

// Unlock the client after a lock action
func (c *Client) Unlock() {
	c.clientMu.Unlock()
}

// NewRequest creates an HTTP requests to the GoCD API endpoints.
func (c *Client) NewRequest(method, urlStr string, body interface{}, apiVersion string) (req *APIRequest, err error) {
	var rel *url.URL
	var buf io.ReadWriter
	req = &APIRequest{}

	// I'm not sure how to get this method to return an error intentionally for testing. For testing purposes, I've
	// added a switch so that the error handling in dependent methods can be tested.
	if os.Getenv("GOCD_RAISE_ERROR_NEW_REQUEST") == "yes" {
		return req, errors.New("Mock Testing Error")
	}

	// Some calls
	if strings.HasPrefix(urlStr, "/") {
		urlStr = urlStr[1:]
	} else {
		urlStr = "api/" + urlStr
	}
	if rel, err = url.Parse(urlStr); err != nil {
		return req, err
	}

	u := c.BaseURL.ResolveReference(rel)
	if c.BaseURL.RawQuery != "" {
		u.RawQuery = c.BaseURL.RawQuery
	}

	if body != nil {
		buf = new(bytes.Buffer)

		enc := json.NewEncoder(buf)
		enc.SetIndent("", "  ")
		err := enc.Encode(body)

		if err != nil {
			return nil, err
		}
		bdy, _ := ioutil.ReadAll(buf)
		req.Body = string(bdy)

		buf = new(bytes.Buffer)
		enc = json.NewEncoder(buf)
		enc.SetIndent("", "  ")
		enc.Encode(body)
	}

	if req.HTTP, err = http.NewRequest(method, u.String(), buf); err != nil {
		return req, err
	}

	if body != nil {
		req.HTTP.Header.Set("Content-Type", "application/json")
	}
	if apiVersion != "" {
		req.HTTP.Header.Set("Accept", apiVersion)
	}
	req.HTTP.Header.Set("User-Agent", c.UserAgent)

	if c.cookie == "" {
		if c.Username != "" && c.Password != "" {
			req.HTTP.SetBasicAuth(c.Username, c.Password)
		}
	} else {
		req.HTTP.Header.Set("Cookie", c.cookie)
	}

	return
}

// Do takes an HTTP request and resposne the response from the GoCD API endpoint.
func (c *Client) Do(ctx context.Context, req *APIRequest, v interface{}, responseType string) (*APIResponse, error) {
	var err error
	var resp *http.Response

	req.HTTP = req.HTTP.WithContext(ctx)

	if resp, err = c.client.Do(req.HTTP); err != nil {
		return nil, err
	}

	r := &APIResponse{
		Request: req,
		HTTP:    resp,
	}

	if v != nil {
		if r.Body, err = readDoResponseBody(v, &r.HTTP.Body, responseType); err != nil {
			return nil, err
		}
	}

	if err = CheckResponse(r.HTTP); err != nil {
		return r, err
	}

	return r, err
}

func readDoResponseBody(v interface{}, body *io.ReadCloser, responseType string) (string, error) {

	if w, ok := v.(io.Writer); ok {
		_, err := io.Copy(w, *body)
		return "", err
	}

	bdy, err := ioutil.ReadAll(*body)
	if responseType == responseTypeText {
		strBody := string(bdy)
		v = &strBody
	} else if responseType == responseTypeXML {
		err = xml.Unmarshal(bdy, v)
	} else {
		err = json.Unmarshal(bdy, v)
	}
	if err == io.EOF {
		err = nil // ignore EOF errors caused by empty response body
	} else if err != nil {
		return "", nil
	}
	return string(bdy), nil
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
