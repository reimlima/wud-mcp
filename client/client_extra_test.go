package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/reimlima/wud-mcp/client"
)

func newServer(t *testing.T, path, response string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Errorf("unexpected path: got %s, want %s", r.URL.Path, path)
		}
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestNew_EmptyBaseURL_UsesDefault(t *testing.T) {
	t.Setenv("WUD_BASE_URL", "")
	c := client.New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNew_ValidTimeout(t *testing.T) {
	srv := newServer(t, "/api/app", `{"name":"wud","version":"7.0.0"}`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	t.Setenv("WUD_TIMEOUT", "30")
	c := client.New()
	_, err := c.GetApp(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidTimeout_FallsBackToDefault(t *testing.T) {
	srv := newServer(t, "/api/app", `{"name":"wud","version":"7.0.0"}`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	t.Setenv("WUD_TIMEOUT", "notanumber")
	c := client.New()
	_, err := c.GetApp(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_ZeroTimeout_FallsBackToDefault(t *testing.T) {
	srv := newServer(t, "/api/app", `{"name":"wud","version":"7.0.0"}`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	t.Setenv("WUD_TIMEOUT", "0")
	c := client.New()
	_, err := c.GetApp(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDo_ConnectionRefused(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	t.Setenv("WUD_BASE_URL", srv.URL)
	srv.Close()
	c := client.New()
	_, err := c.GetApp(context.Background())
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}

func TestDo_ContentTypeSetWithBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(srv.Close)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	_, err := c.RunTrigger(context.Background(), "smtp", "mysmtp", strings.NewReader(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetContainerTriggers_Success(t *testing.T) {
	srv := newServer(t, "/api/containers/abc/triggers", `[{"type":"slack"}]`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.GetContainerTriggers(context.Background(), "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestWatchContainer_Success(t *testing.T) {
	srv := newServer(t, "/api/containers/abc/watch", `{"status":"ok"}`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.WatchContainer(context.Background(), "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestRunContainerTrigger_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "/api/containers/abc/triggers/slack/myslack"
		if r.URL.Path != want {
			t.Errorf("unexpected path: got %s, want %s", r.URL.Path, want)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	t.Cleanup(srv.Close)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.RunContainerTrigger(context.Background(), "abc", "slack", "myslack")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestDeleteContainer_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/containers/abc" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.DeleteContainer(context.Background(), "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"status":"ok"}` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestListRegistries_Success(t *testing.T) {
	srv := newServer(t, "/api/registries", `[]`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.ListRegistries(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `[]` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestGetRegistry_Success(t *testing.T) {
	srv := newServer(t, "/api/registries/hub/dockerhub", `{"type":"hub"}`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.GetRegistry(context.Background(), "hub", "dockerhub")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestListTriggers_Success(t *testing.T) {
	srv := newServer(t, "/api/triggers", `[]`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.ListTriggers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `[]` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestGetTrigger_Success(t *testing.T) {
	srv := newServer(t, "/api/triggers/slack/myslack", `{"type":"slack"}`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.GetTrigger(context.Background(), "slack", "myslack")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestRunTrigger_Success_NoBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/triggers/slack/myslack" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	t.Cleanup(srv.Close)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.RunTrigger(context.Background(), "slack", "myslack", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestListWatchers_Success(t *testing.T) {
	srv := newServer(t, "/api/watchers", `[]`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.ListWatchers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `[]` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestGetWatcher_Success(t *testing.T) {
	srv := newServer(t, "/api/watchers/docker/local", `{"type":"docker"}`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.GetWatcher(context.Background(), "docker", "local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestGetStore_Success(t *testing.T) {
	srv := newServer(t, "/api/store", `{"configuration":{"path":"/store","file":"store.json"}}`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.GetStore(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestGetLog_Success(t *testing.T) {
	srv := newServer(t, "/api/log", `{"level":"info"}`)
	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()
	data, err := c.GetLog(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}
