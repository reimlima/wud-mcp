package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reimlima/wud-mcp/client"
)

func TestGetApp_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/app" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"wud","version":"7.0.0"}`))
	}))
	defer srv.Close()

	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()

	data, err := c.GetApp(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"name":"wud","version":"7.0.0"}` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestGetApp_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()

	_, err := c.GetApp(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListContainers_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/containers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()

	data, err := c.ListContainers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `[]` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		_, _ = w.Write([]byte(`{"name":"wud","version":"7.0.0"}`))
	}))
	defer srv.Close()

	t.Setenv("WUD_BASE_URL", srv.URL)
	t.Setenv("WUD_API_USER", "admin")
	t.Setenv("WUD_API_PASSWORD", "secret")
	c := client.New()

	data, err := c.GetApp(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestEmptyBody_ReturnsOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()

	data, err := c.WatchAllContainers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"status":"ok"}` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestGetContainer_URLEscaping(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/containers/my-container" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"my-container"}`))
	}))
	defer srv.Close()

	t.Setenv("WUD_BASE_URL", srv.URL)
	c := client.New()

	data, err := c.GetContainer(context.Background(), "my-container")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}
