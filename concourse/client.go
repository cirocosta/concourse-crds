package concourse

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type Client struct {
	url   string
	token string
}

// ps.: this IS NOT right - the token will eventually expire (sooner than you
// think).
//
func New(ctx context.Context, url, user, pass string) (client *Client, err error) {
	client = &Client{url: url}

	client.token, err = passwordGrant(ctx, url, user, pass)
	if err != nil {
		err = fmt.Errorf("password grant: %w", err)
		return
	}

	return
}

func (c *Client) DestroyPipeline(ctx context.Context, team, name string) (err error) {
	url := c.url + "/api/v1/teams/" + team + "/pipelines/" + name

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		err = fmt.Errorf("new delete request: %w", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("req: %w", err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = fmt.Errorf("non-ok status:  %s", resp.Status)
		return
	}

	return
}

func (c *Client) UnpausePipeline(ctx context.Context, team, name string) (err error) {
	err = c.pauseOrUnpausePipeline(ctx, team, name, "unpause")
	if err != nil {
		err = fmt.Errorf("unpause pipeline: %w", err)
		return
	}

	return
}

func (c *Client) PausePipeline(ctx context.Context, team, name string) (err error) {
	err = c.pauseOrUnpausePipeline(ctx, team, name, "pause")
	if err != nil {
		err = fmt.Errorf("unpause pipeline: %w", err)
		return
	}

	return
}

func (c *Client) SetPipeline(ctx context.Context, team, name string, config []byte) (err error) {
	pipelineVersion, err := c.getPipelineVersion(ctx, team, name)
	if err != nil {
		err = fmt.Errorf("get pipeline version: %w", err)
		return
	}

	url := c.url + "/api/v1/teams/" + team + "/pipelines/" + name + "/config"
	body := bytes.NewBuffer(config)

	req, err := http.NewRequestWithContext(ctx, "PUT", url, body)
	if err != nil {
		err = fmt.Errorf("new put request: %w", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", "application/x-yaml")

	if pipelineVersion != "" {
		req.Header.Add("X-Concourse-Config-Version", pipelineVersion)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("req: %w", err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = fmt.Errorf("non-ok status:  %s", resp.Status)
		return
	}

	return
}

func (c *Client) getPipelineVersion(ctx context.Context, team, name string) (version string, err error) {
	url := c.url + "/api/v1/teams/" + team + "/pipelines/" + name + "/config"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		err = fmt.Errorf("new put request: %w", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("req: %w", err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode == 404 { // pipeline not found
			err = nil
			return
		}

		err = fmt.Errorf("non-ok status:  %s", resp.Status)
		return
	}

	version = resp.Header.Get("X-Concourse-Config-Version")

	return
}

func (c *Client) pauseOrUnpausePipeline(ctx context.Context, team, name, action string) (err error) {
	url := c.url + "/api/v1/teams/" + team + "/pipelines/" + name + "/" + action

	req, err := http.NewRequestWithContext(ctx, "PUT", url, nil)
	if err != nil {
		err = fmt.Errorf("new put request: %w", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("req: %w", err)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = fmt.Errorf("non-ok status:  %s", resp.Status)
		return
	}

	return
}

func passwordGrant(ctx context.Context, url, user, pass string) (accessToken string, err error) {
	oauth2Config := oauth2.Config{
		ClientID:     "fly",
		ClientSecret: "Zmx5",
		Endpoint: oauth2.Endpoint{
			TokenURL: url + "/sky/token",
		},
		Scopes: []string{
			"openid", "profile", "email", "federated:id", "groups",
		},
	}

	token, err := oauth2Config.PasswordCredentialsToken(ctx, user, pass)
	if err != nil {
		err = fmt.Errorf("password creds token: %w", err)
		return
	}

	accessToken = token.AccessToken
	return
}
