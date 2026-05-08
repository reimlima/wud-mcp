package tools_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/reimlima/wud-mcp/client"
	"github.com/reimlima/wud-mcp/tools"
)

// setup creates a WUD test server with the given handler, wires a real client
// and MCP server with all tools registered, and returns the MCP server.
func setup(t *testing.T, handler http.Handler) *server.MCPServer {
	t.Helper()
	wud := httptest.NewServer(handler)
	t.Cleanup(wud.Close)
	t.Setenv("WUD_BASE_URL", wud.URL)
	c := client.New()
	s := server.NewMCPServer("test", "0")
	tools.RegisterApp(s, c)
	tools.RegisterContainers(s, c)
	tools.RegisterRegistries(s, c)
	tools.RegisterTriggers(s, c)
	tools.RegisterWatchers(s, c)
	tools.RegisterStore(s, c)
	return s
}

// call invokes a registered tool handler and returns its result.
func call(t *testing.T, s *server.MCPServer, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	st := s.GetTool(name)
	if st == nil {
		t.Fatalf("tool %q not registered", name)
	}
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	result, err := st.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("tool %q handler error: %v", name, err)
	}
	return result
}

func okHandler(path, body string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, path) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(body))
	})
}

func errHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})
}

// ── app ──────────────────────────────────────────────────────────────────────

func TestGetAppInfo_Success(t *testing.T) {
	s := setup(t, okHandler("/api/app", `{"name":"wud","version":"7.0.0"}`))
	result := call(t, s, "get_app_info", nil)
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetAppInfo_WUDError(t *testing.T) {
	s := setup(t, errHandler())
	result := call(t, s, "get_app_info", nil)
	if !result.IsError {
		t.Errorf("expected error result")
	}
}

// ── containers ───────────────────────────────────────────────────────────────

func TestListContainers_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/containers", `[]`))
	result := call(t, s, "list_containers", nil)
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestListContainers_Tool_WUDError(t *testing.T) {
	s := setup(t, errHandler())
	result := call(t, s, "list_containers", nil)
	if !result.IsError {
		t.Errorf("expected error result")
	}
}

func TestWatchAllContainers_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/containers/watch", `{"status":"ok"}`))
	result := call(t, s, "watch_all_containers", nil)
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetContainer_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/containers/", `{"id":"abc"}`))
	result := call(t, s, "get_container", map[string]any{"id": "abc"})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetContainer_Tool_MissingID(t *testing.T) {
	s := setup(t, okHandler("/api/containers/", `{}`))
	result := call(t, s, "get_container", map[string]any{})
	if !result.IsError {
		t.Errorf("expected error for missing id")
	}
}

func TestGetContainer_Tool_WUDError(t *testing.T) {
	s := setup(t, errHandler())
	result := call(t, s, "get_container", map[string]any{"id": "abc"})
	if !result.IsError {
		t.Errorf("expected error result")
	}
}

func TestGetContainerTriggers_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/containers/", `[]`))
	result := call(t, s, "get_container_triggers", map[string]any{"id": "abc"})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetContainerTriggers_Tool_MissingID(t *testing.T) {
	s := setup(t, okHandler("/api/containers/", `[]`))
	result := call(t, s, "get_container_triggers", map[string]any{})
	if !result.IsError {
		t.Errorf("expected error for missing id")
	}
}

func TestWatchContainer_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/containers/", `{"status":"ok"}`))
	result := call(t, s, "watch_container", map[string]any{"id": "abc"})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestWatchContainer_Tool_MissingID(t *testing.T) {
	s := setup(t, okHandler("/api/containers/", `{}`))
	result := call(t, s, "watch_container", map[string]any{})
	if !result.IsError {
		t.Errorf("expected error for missing id")
	}
}

func TestRunContainerTrigger_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/containers/", `{"status":"ok"}`))
	result := call(t, s, "run_container_trigger", map[string]any{
		"id":   "abc",
		"type": "slack",
		"name": "myslack",
	})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestRunContainerTrigger_Tool_MissingParams(t *testing.T) {
	s := setup(t, okHandler("/api/containers/", `{}`))
	result := call(t, s, "run_container_trigger", map[string]any{"id": "abc"})
	if !result.IsError {
		t.Errorf("expected error for missing params")
	}
}

func TestDeleteContainer_Tool_Success(t *testing.T) {
	s := setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	result := call(t, s, "delete_container", map[string]any{"id": "abc"})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestDeleteContainer_Tool_MissingID(t *testing.T) {
	s := setup(t, okHandler("/api/containers/", `{}`))
	result := call(t, s, "delete_container", map[string]any{})
	if !result.IsError {
		t.Errorf("expected error for missing id")
	}
}

// ── registries ───────────────────────────────────────────────────────────────

func TestListRegistries_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/registries", `[]`))
	result := call(t, s, "list_registries", nil)
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetRegistry_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/registries/", `{"type":"hub"}`))
	result := call(t, s, "get_registry", map[string]any{"type": "hub", "name": "dockerhub"})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetRegistry_Tool_MissingParams(t *testing.T) {
	s := setup(t, okHandler("/api/registries/", `{}`))
	result := call(t, s, "get_registry", map[string]any{"type": "hub"})
	if !result.IsError {
		t.Errorf("expected error for missing name")
	}
}

// ── triggers ─────────────────────────────────────────────────────────────────

func TestListTriggers_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/triggers", `[]`))
	result := call(t, s, "list_triggers", nil)
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetTrigger_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/triggers/", `{"type":"slack"}`))
	result := call(t, s, "get_trigger", map[string]any{"type": "slack", "name": "myslack"})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetTrigger_Tool_MissingParams(t *testing.T) {
	s := setup(t, okHandler("/api/triggers/", `{}`))
	result := call(t, s, "get_trigger", map[string]any{"type": "slack"})
	if !result.IsError {
		t.Errorf("expected error for missing name")
	}
}

func TestRunTrigger_Tool_Success_NoBody(t *testing.T) {
	s := setup(t, okHandler("/api/triggers/", `{"status":"ok"}`))
	result := call(t, s, "run_trigger", map[string]any{"type": "slack", "name": "myslack"})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestRunTrigger_Tool_Success_WithBody(t *testing.T) {
	s := setup(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "bad content-type", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	result := call(t, s, "run_trigger", map[string]any{
		"type":           "slack",
		"name":           "myslack",
		"container_json": `{"id":"abc","name":"nginx"}`,
	})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestRunTrigger_Tool_MissingParams(t *testing.T) {
	s := setup(t, okHandler("/api/triggers/", `{}`))
	result := call(t, s, "run_trigger", map[string]any{"type": "slack"})
	if !result.IsError {
		t.Errorf("expected error for missing name")
	}
}

// ── watchers ─────────────────────────────────────────────────────────────────

func TestListWatchers_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/watchers", `[]`))
	result := call(t, s, "list_watchers", nil)
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetWatcher_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/watchers/", `{"type":"docker"}`))
	result := call(t, s, "get_watcher", map[string]any{"type": "docker", "name": "local"})
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetWatcher_Tool_MissingParams(t *testing.T) {
	s := setup(t, okHandler("/api/watchers/", `{}`))
	result := call(t, s, "get_watcher", map[string]any{"type": "docker"})
	if !result.IsError {
		t.Errorf("expected error for missing name")
	}
}

// ── store ─────────────────────────────────────────────────────────────────────

func TestGetStoreConfig_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/store", `{"configuration":{"path":"/store","file":"store.json"}}`))
	result := call(t, s, "get_store_config", nil)
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetStoreConfig_Tool_WUDError(t *testing.T) {
	s := setup(t, errHandler())
	result := call(t, s, "get_store_config", nil)
	if !result.IsError {
		t.Errorf("expected error result")
	}
}

func TestGetLogConfig_Tool_Success(t *testing.T) {
	s := setup(t, okHandler("/api/log", `{"level":"info"}`))
	result := call(t, s, "get_log_config", nil)
	if result.IsError {
		t.Errorf("expected success, got error")
	}
}

func TestGetLogConfig_Tool_WUDError(t *testing.T) {
	s := setup(t, errHandler())
	result := call(t, s, "get_log_config", nil)
	if !result.IsError {
		t.Errorf("expected error result")
	}
}
