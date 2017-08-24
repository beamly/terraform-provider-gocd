package gocd

import (
	"context"
	"net/url"
)

// EncryptionService describes the HAL _link resource for the api response object for a pipelineconfig
type EncryptionService service

// CipherText sescribes the response from the api with an encrypted value.
type CipherText struct {
	EncryptedValue string       `json:"encrypted_value"`
	Links          EncryptLinks `json:"_links"`
}

// EncryptLinks describes the HAL _link resource for the api response object for a collection of agent objects.
//go:generate gocd-response-links-generator -type=EncryptLinks
type EncryptLinks struct {
	Self *url.URL `json:"self"`
	Doc  *url.URL `json:"doc"`
}

// Encrypt takes a plaintext value and returns a cipher text.
func (es *EncryptionService) Encrypt(ctx context.Context, value string) (*CipherText, *APIResponse, error) {

	c := CipherText{}
	_, resp, err := es.client.postAction(ctx, &APIClientRequest{
		Path:         "admin/encrypt",
		ResponseBody: &c,
		RequestBody: &struct {
			Value string `json:"value"`
		}{Value: value},
		APIVersion: apiV1,
	})

	return &c, resp, err
}
