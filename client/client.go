package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
}

func New() *Client {
	baseURL := os.Getenv("WUD_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000" //nolint:revive // DevSkim: ignore DS162092
	}

	timeout := 10
	if t := os.Getenv("WUD_TIMEOUT"); t != "" {
		if v, err := strconv.Atoi(t); err == nil && v > 0 {
			timeout = v
		}
	}

	return &Client{
		baseURL:  baseURL,
		username: os.Getenv("WUD_API_USER"),
		password: os.Getenv("WUD_API_PASSWORD"),
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(b))
	}
	if len(b) == 0 {
		return []byte(`{"status":"ok"}`), nil
	}
	return b, nil
}

func (c *Client) GetApp(ctx context.Context) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/app", nil)
}

func (c *Client) ListContainers(ctx context.Context) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/containers", nil)
}

func (c *Client) WatchAllContainers(ctx context.Context) ([]byte, error) {
	return c.do(ctx, http.MethodPost, "/api/containers/watch", nil)
}

func (c *Client) GetContainer(ctx context.Context, id string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/containers/"+url.PathEscape(id), nil)
}

func (c *Client) GetContainerTriggers(ctx context.Context, id string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/containers/"+url.PathEscape(id)+"/triggers", nil)
}

func (c *Client) WatchContainer(ctx context.Context, id string) ([]byte, error) {
	return c.do(ctx, http.MethodPost, "/api/containers/"+url.PathEscape(id)+"/watch", nil)
}

func (c *Client) RunContainerTrigger(ctx context.Context, id, trigType, name string) ([]byte, error) {
	path := "/api/containers/" + url.PathEscape(id) + "/triggers/" + url.PathEscape(trigType) + "/" + url.PathEscape(name)
	return c.do(ctx, http.MethodPost, path, nil)
}

func (c *Client) DeleteContainer(ctx context.Context, id string) ([]byte, error) {
	return c.do(ctx, http.MethodDelete, "/api/containers/"+url.PathEscape(id), nil)
}

func (c *Client) ListRegistries(ctx context.Context) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/registries", nil)
}

func (c *Client) GetRegistry(ctx context.Context, regType, name string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/registries/"+url.PathEscape(regType)+"/"+url.PathEscape(name), nil)
}

func (c *Client) ListTriggers(ctx context.Context) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/triggers", nil)
}

func (c *Client) GetTrigger(ctx context.Context, trigType, name string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/triggers/"+url.PathEscape(trigType)+"/"+url.PathEscape(name), nil)
}

func (c *Client) RunTrigger(ctx context.Context, trigType, name string, body io.Reader) ([]byte, error) {
	return c.do(ctx, http.MethodPost, "/api/triggers/"+url.PathEscape(trigType)+"/"+url.PathEscape(name), body)
}

func (c *Client) ListWatchers(ctx context.Context) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/watchers", nil)
}

func (c *Client) GetWatcher(ctx context.Context, watchType, name string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/watchers/"+url.PathEscape(watchType)+"/"+url.PathEscape(name), nil)
}

func (c *Client) GetStore(ctx context.Context) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/store", nil)
}

func (c *Client) GetLog(ctx context.Context) ([]byte, error) {
	return c.do(ctx, http.MethodGet, "/api/log", nil)
}
