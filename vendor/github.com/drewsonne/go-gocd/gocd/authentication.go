package gocd

import "context"

// Login sends basic auth to the GoCD Server and sets an auth cookie in the client to enable cookie based auth
// for future requests.
func (c *Client) Login(ctx context.Context) error {
	req, err := c.NewRequest("GET", "api/agents", nil, apiV2)
	if err != nil {
		return err
	}
	req.HTTP.SetBasicAuth(c.Username, c.Password)

	resp, err := c.Do(ctx, req, nil, responseTypeJSON)
	if err == nil {
		c.cookie = resp.HTTP.Header["Set-Cookie"][0]
	}
	return err
}
