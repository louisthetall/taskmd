package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/driangle/taskmd/apps/cli/internal/sync"
)

func TestJiraSource_Name(t *testing.T) {
	src := &JiraSource{}
	if src.Name() != "jira" {
		t.Fatalf("expected name %q, got %q", "jira", src.Name())
	}
}

func TestJiraSource_ValidateConfig_Valid(t *testing.T) {
	src := &JiraSource{}
	cfg := sync.SourceConfig{
		Project:  "PROJ",
		TokenEnv: "JIRA_TOKEN",
		UserEnv:  "JIRA_USER",
	}

	if err := src.ValidateConfig(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJiraSource_ValidateConfig_MissingProject(t *testing.T) {
	src := &JiraSource{}
	cfg := sync.SourceConfig{
		TokenEnv: "JIRA_TOKEN",
		UserEnv:  "JIRA_USER",
	}
	err := src.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestJiraSource_ValidateConfig_MissingTokenEnv(t *testing.T) {
	src := &JiraSource{}
	cfg := sync.SourceConfig{
		Project: "PROJ",
		UserEnv: "JIRA_USER",
	}
	err := src.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing token_env")
	}
}

func TestJiraSource_ValidateConfig_MissingUserEnv(t *testing.T) {
	src := &JiraSource{}
	cfg := sync.SourceConfig{
		Project:  "PROJ",
		TokenEnv: "JIRA_TOKEN",
	}
	err := src.ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error for missing user_env")
	}
}

func newTestIssues() []jiraIssue {
	return []jiraIssue{
		{
			Key: "PROJ-1",
			Fields: jiraFields{
				Summary: "Fix login bug",
				Description: map[string]any{
					"type": "doc",
					"content": []any{
						map[string]any{
							"type": "paragraph",
							"content": []any{
								map[string]any{"type": "text", "text": "Login is broken"},
							},
						},
					},
				},
				Status:   &jiraStatus{Name: "In Progress"},
				Priority: &jiraPriority{Name: "High"},
				Assignee: &jiraUser{DisplayName: "Alice Smith"},
				Labels:   []string{"bug", "urgent"},
				Created:  "2025-01-15T10:30:00.000+0000",
				Updated:  "2025-01-16T14:00:00.000+0000",
			},
		},
		{
			Key: "PROJ-2",
			Fields: jiraFields{
				Summary: "Add search feature",
				Status:  &jiraStatus{Name: "To Do"},
				Labels:  []string{"enhancement"},
				Created: "2025-01-17T09:00:00.000+0000",
				Updated: "2025-01-17T09:00:00.000+0000",
			},
		},
	}
}

func fetchTestTasks(t *testing.T, issues []jiraIssue) ([]sync.ExternalTask, string) {
	t.Helper()

	resp := searchResponse{Total: len(issues), Issues: issues}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") {
			t.Errorf("expected Basic auth, got: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("unexpected accept header: %s", r.Header.Get("Accept"))
		}
		if !strings.Contains(r.URL.Path, "/rest/api/3/search") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	t.Setenv("TEST_JIRA_TOKEN", "test-token")
	t.Setenv("TEST_JIRA_USER", "alice@example.com")

	src := &JiraSource{HTTPClient: srv.Client()}
	cfg := sync.SourceConfig{
		Project:  "PROJ",
		BaseURL:  srv.URL,
		TokenEnv: "TEST_JIRA_TOKEN",
		UserEnv:  "TEST_JIRA_USER",
	}

	tasks, err := src.FetchTasks(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return tasks, srv.URL
}

func TestJiraSource_FetchTasks_Basic(t *testing.T) {
	tasks, srvURL := fetchTestTasks(t, newTestIssues())

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	task := tasks[0]
	if task.ExternalID != "PROJ-1" {
		t.Errorf("expected ExternalID %q, got %q", "PROJ-1", task.ExternalID)
	}
	if task.Title != "Fix login bug" {
		t.Errorf("expected Title %q, got %q", "Fix login bug", task.Title)
	}
	if task.Description != "Login is broken" {
		t.Errorf("expected Description %q, got %q", "Login is broken", task.Description)
	}
	if task.Status != "In Progress" {
		t.Errorf("expected Status %q, got %q", "In Progress", task.Status)
	}
	if task.Priority != "High" {
		t.Errorf("expected Priority %q, got %q", "High", task.Priority)
	}
	if task.Assignee != "Alice Smith" {
		t.Errorf("expected Assignee %q, got %q", "Alice Smith", task.Assignee)
	}
	expectedURL := srvURL + "/browse/PROJ-1"
	if task.URL != expectedURL {
		t.Errorf("expected URL %q, got %q", expectedURL, task.URL)
	}
	if len(task.Labels) != 2 || task.Labels[0] != "bug" || task.Labels[1] != "urgent" {
		t.Errorf("expected labels [bug urgent], got %v", task.Labels)
	}
	if task.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if task.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}

	task2 := tasks[1]
	if task2.Assignee != "" {
		t.Errorf("expected empty Assignee, got %q", task2.Assignee)
	}
	if task2.Priority != "" {
		t.Errorf("expected empty Priority, got %q", task2.Priority)
	}
}

func TestJiraSource_FetchTasks_Pagination(t *testing.T) {
	callCount := 0

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		startAt := r.URL.Query().Get("startAt")

		switch startAt {
		case "0", "":
			issues := make([]jiraIssue, maxResults)
			for i := range issues {
				issues[i] = jiraIssue{
					Key: fmt.Sprintf("PROJ-%d", i+1),
					Fields: jiraFields{
						Summary: fmt.Sprintf("Issue %d", i+1),
						Status:  &jiraStatus{Name: "To Do"},
					},
				}
			}
			resp := searchResponse{
				StartAt:    0,
				MaxResults: maxResults,
				Total:      maxResults + 1,
				Issues:     issues,
			}
			json.NewEncoder(w).Encode(resp)
		default:
			resp := searchResponse{
				StartAt:    maxResults,
				MaxResults: maxResults,
				Total:      maxResults + 1,
				Issues: []jiraIssue{
					{
						Key: fmt.Sprintf("PROJ-%d", maxResults+1),
						Fields: jiraFields{
							Summary: fmt.Sprintf("Issue %d", maxResults+1),
							Status:  &jiraStatus{Name: "To Do"},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer srv.Close()

	t.Setenv("TEST_PAGE_TOKEN", "tok")
	t.Setenv("TEST_PAGE_USER", "user@example.com")

	src := &JiraSource{HTTPClient: srv.Client()}
	cfg := sync.SourceConfig{
		Project:  "PROJ",
		BaseURL:  srv.URL,
		TokenEnv: "TEST_PAGE_TOKEN",
		UserEnv:  "TEST_PAGE_USER",
	}

	tasks, err := src.FetchTasks(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != maxResults+1 {
		t.Fatalf("expected %d tasks, got %d", maxResults+1, len(tasks))
	}
	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
}

func TestJiraSource_FetchTasks_JQLFilter(t *testing.T) {
	var receivedJQL string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedJQL = r.URL.Query().Get("jql")
		json.NewEncoder(w).Encode(searchResponse{Total: 0, Issues: []jiraIssue{}})
	}))
	defer srv.Close()

	t.Setenv("TEST_JQL_TOKEN", "tok")
	t.Setenv("TEST_JQL_USER", "user@example.com")

	src := &JiraSource{HTTPClient: srv.Client()}
	cfg := sync.SourceConfig{
		Project:  "MYPROJ",
		BaseURL:  srv.URL,
		TokenEnv: "TEST_JQL_TOKEN",
		UserEnv:  "TEST_JQL_USER",
		Filters: map[string]any{
			"jql": "status = \"To Do\" AND priority = High",
		},
	}

	_, err := src.FetchTasks(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(receivedJQL, "MYPROJ") {
		t.Errorf("expected JQL to contain project, got: %s", receivedJQL)
	}
	if !strings.Contains(receivedJQL, "status = \"To Do\" AND priority = High") {
		t.Errorf("expected JQL to contain custom filter, got: %s", receivedJQL)
	}
}

func TestJiraSource_FetchTasks_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"errorMessages":["Client must be authenticated"]}`)
	}))
	defer srv.Close()

	t.Setenv("TEST_ERR_TOKEN", "bad-tok")
	t.Setenv("TEST_ERR_USER", "user@example.com")

	src := &JiraSource{HTTPClient: srv.Client()}
	cfg := sync.SourceConfig{
		Project:  "PROJ",
		BaseURL:  srv.URL,
		TokenEnv: "TEST_ERR_TOKEN",
		UserEnv:  "TEST_ERR_USER",
	}

	_, err := src.FetchTasks(cfg)
	if err == nil {
		t.Fatal("expected error for API error response")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("expected error to contain status code 401, got: %v", err)
	}
}

func TestJiraSource_FetchTasks_NilFields(t *testing.T) {
	resp := searchResponse{
		Total: 1,
		Issues: []jiraIssue{
			{
				Key: "PROJ-10",
				Fields: jiraFields{
					Summary:  "Minimal issue",
					Status:   nil,
					Priority: nil,
					Assignee: nil,
					Labels:   nil,
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	t.Setenv("TEST_NIL_TOKEN", "tok")
	t.Setenv("TEST_NIL_USER", "user@example.com")

	src := &JiraSource{HTTPClient: srv.Client()}
	cfg := sync.SourceConfig{
		Project:  "PROJ",
		BaseURL:  srv.URL,
		TokenEnv: "TEST_NIL_TOKEN",
		UserEnv:  "TEST_NIL_USER",
	}

	tasks, err := src.FetchTasks(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]
	if task.Status != "" {
		t.Errorf("expected empty status, got %q", task.Status)
	}
	if task.Priority != "" {
		t.Errorf("expected empty priority, got %q", task.Priority)
	}
	if task.Assignee != "" {
		t.Errorf("expected empty assignee, got %q", task.Assignee)
	}
	if task.Labels != nil {
		t.Errorf("expected nil labels, got %v", task.Labels)
	}
}

func TestJiraSource_FetchTasks_EmptyResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(searchResponse{Total: 0, Issues: []jiraIssue{}})
	}))
	defer srv.Close()

	t.Setenv("TEST_EMPTY_TOKEN", "tok")
	t.Setenv("TEST_EMPTY_USER", "user@example.com")

	src := &JiraSource{HTTPClient: srv.Client()}
	cfg := sync.SourceConfig{
		Project:  "PROJ",
		BaseURL:  srv.URL,
		TokenEnv: "TEST_EMPTY_TOKEN",
		UserEnv:  "TEST_EMPTY_USER",
	}

	tasks, err := src.FetchTasks(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestJiraSource_FetchTasks_UnsetToken(t *testing.T) {
	os.Unsetenv("UNSET_JIRA_TOKEN")
	t.Setenv("TEST_JIRA_USER_SET", "user@example.com")

	src := &JiraSource{}
	cfg := sync.SourceConfig{
		Project:  "PROJ",
		BaseURL:  "https://example.atlassian.net",
		TokenEnv: "UNSET_JIRA_TOKEN",
		UserEnv:  "TEST_JIRA_USER_SET",
	}

	_, err := src.FetchTasks(cfg)
	if err == nil {
		t.Fatal("expected error for unset token env var")
	}
	if !strings.Contains(err.Error(), "UNSET_JIRA_TOKEN") {
		t.Errorf("expected error to mention env var name, got: %v", err)
	}
}

func TestJiraSource_FetchTasks_MissingBaseURL(t *testing.T) {
	t.Setenv("TEST_URL_TOKEN", "tok")
	t.Setenv("TEST_URL_USER", "user@example.com")

	src := &JiraSource{}
	cfg := sync.SourceConfig{
		Project:  "PROJ",
		TokenEnv: "TEST_URL_TOKEN",
		UserEnv:  "TEST_URL_USER",
	}

	_, err := src.FetchTasks(cfg)
	if err == nil {
		t.Fatal("expected error for missing base_url")
	}
	if !strings.Contains(err.Error(), "base_url") {
		t.Errorf("expected error to mention base_url, got: %v", err)
	}
}

func TestADFToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name: "simple paragraph",
			input: map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type": "paragraph",
						"content": []any{
							map[string]any{"type": "text", "text": "Hello world"},
						},
					},
				},
			},
			expected: "Hello world",
		},
		{
			name: "heading",
			input: map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type":  "heading",
						"attrs": map[string]any{"level": float64(2)},
						"content": []any{
							map[string]any{"type": "text", "text": "Title"},
						},
					},
				},
			},
			expected: "## Title",
		},
		{
			name: "bold and italic",
			input: map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type": "paragraph",
						"content": []any{
							map[string]any{
								"type":  "text",
								"text":  "bold",
								"marks": []any{map[string]any{"type": "strong"}},
							},
							map[string]any{"type": "text", "text": " and "},
							map[string]any{
								"type":  "text",
								"text":  "italic",
								"marks": []any{map[string]any{"type": "em"}},
							},
						},
					},
				},
			},
			expected: "**bold** and *italic*",
		},
		{
			name: "inline code",
			input: map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type": "paragraph",
						"content": []any{
							map[string]any{"type": "text", "text": "Use "},
							map[string]any{
								"type":  "text",
								"text":  "fmt.Println",
								"marks": []any{map[string]any{"type": "code"}},
							},
						},
					},
				},
			},
			expected: "Use `fmt.Println`",
		},
		{
			name: "link",
			input: map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type": "paragraph",
						"content": []any{
							map[string]any{
								"type": "text",
								"text": "Click here",
								"marks": []any{
									map[string]any{
										"type":  "link",
										"attrs": map[string]any{"href": "https://example.com"},
									},
								},
							},
						},
					},
				},
			},
			expected: "[Click here](https://example.com)",
		},
		{
			name: "bullet list",
			input: map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type": "bulletList",
						"content": []any{
							map[string]any{
								"type": "listItem",
								"content": []any{
									map[string]any{"type": "text", "text": "Item one"},
								},
							},
							map[string]any{
								"type": "listItem",
								"content": []any{
									map[string]any{"type": "text", "text": "Item two"},
								},
							},
						},
					},
				},
			},
			expected: "- Item one\n- Item two",
		},
		{
			name: "ordered list",
			input: map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type": "orderedList",
						"content": []any{
							map[string]any{
								"type": "listItem",
								"content": []any{
									map[string]any{"type": "text", "text": "First"},
								},
							},
							map[string]any{
								"type": "listItem",
								"content": []any{
									map[string]any{"type": "text", "text": "Second"},
								},
							},
						},
					},
				},
			},
			expected: "1. First\n2. Second",
		},
		{
			name: "code block",
			input: map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type":  "codeBlock",
						"attrs": map[string]any{"language": "go"},
						"content": []any{
							map[string]any{"type": "text", "text": "fmt.Println(\"hi\")"},
						},
					},
				},
			},
			expected: "```go\nfmt.Println(\"hi\")\n```",
		},
		{
			name: "blockquote",
			input: map[string]any{
				"type": "doc",
				"content": []any{
					map[string]any{
						"type": "blockquote",
						"content": []any{
							map[string]any{
								"type": "paragraph",
								"content": []any{
									map[string]any{"type": "text", "text": "A quote"},
								},
							},
						},
					},
				},
			},
			expected: "> A quote",
		},
		{
			name:     "nil description",
			input:    nil,
			expected: "<nil>",
		},
		{
			name:     "plain string fallback",
			input:    "just a string",
			expected: "just a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adfToMarkdown(tt.input)
			if got != tt.expected {
				t.Errorf("adfToMarkdown:\n  got:      %q\n  expected: %q", got, tt.expected)
			}
		})
	}
}

func TestBuildJQL(t *testing.T) {
	tests := []struct {
		name     string
		project  string
		filters  map[string]any
		contains []string
	}{
		{
			name:     "project only",
			project:  "PROJ",
			filters:  nil,
			contains: []string{`project = "PROJ"`},
		},
		{
			name:    "with JQL filter",
			project: "MYPROJ",
			filters: map[string]any{
				"jql": "status = \"Done\"",
			},
			contains: []string{`project = "MYPROJ"`, "AND", `status = "Done"`},
		},
		{
			name:    "empty JQL filter ignored",
			project: "PROJ",
			filters: map[string]any{
				"jql": "",
			},
			contains: []string{`project = "PROJ"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildJQL(tt.project, tt.filters)
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("buildJQL(%q, %v) = %q, want to contain %q", tt.project, tt.filters, got, want)
				}
			}
		})
	}
}
