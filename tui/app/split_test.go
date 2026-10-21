package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/dinoDanic/diny/commit"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
)

func TestLoadSplitPlanRecovery(t *testing.T) {
	valid := []commit.SplitGroup{
		{Order: 1, Type: "refactor", Message: "refactor(base): reorganize settings", Files: []string{"modules/base.nix"}},
		{Order: 2, Type: "fix", Message: "fix(hardware): configure audio", Files: []string{"modules/hardware.nix"}},
	}
	duplicate := []commit.SplitGroup{
		{Order: 1, Type: "refactor", Message: "refactor(base): reorganize settings", Files: []string{"modules/base.nix"}},
		{Order: 2, Type: "fix", Message: "fix(hardware): configure audio", Files: []string{"modules/base.nix", "modules/hardware.nix"}},
	}
	missing := valid[:1]
	unknown := []commit.SplitGroup{{Order: 1, Files: []string{"not-staged.nix"}}}
	for _, tc := range []struct {
		name      string
		plans     [][]commit.SplitGroup
		status    int
		wantCalls int
		wantError bool
	}{
		{name: "valid immediately", plans: [][]commit.SplitGroup{valid}, wantCalls: 1},
		{name: "duplicate corrected", plans: [][]commit.SplitGroup{duplicate, valid}, wantCalls: 2},
		{name: "missing corrected", plans: [][]commit.SplitGroup{missing, valid}, wantCalls: 2},
		{name: "unknown corrected", plans: [][]commit.SplitGroup{unknown, valid}, wantCalls: 2},
		{name: "persistent conflict rejected", plans: [][]commit.SplitGroup{duplicate}, wantCalls: 3, wantError: true},
		{name: "provider failure returned without retry", status: http.StatusUnauthorized, wantCalls: 1, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if tc.status != 0 {
					w.WriteHeader(tc.status)
					return
				}
				var request struct {
					Messages []struct {
						Content string `json:"content"`
					} `json:"messages"`
				}
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					t.Error(err)
					w.WriteHeader(400)
					return
				}
				if calls > 1 {
					user := request.Messages[1].Content
					if !strings.Contains(user, "Keep hardware separate") || !strings.Contains(user, "modules/base.nix") {
						t.Error("repair request lost the user feedback or conflicting staged path")
					}
				}
				index := calls - 1
				if index >= len(tc.plans) {
					index = len(tc.plans) - 1
				}
				plan, _ := json.Marshal(map[string]any{"groups": tc.plans[index]})
				json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]string{"content": string(plan)}}}})
			}))
			defer server.Close()
			cfg := &config.Config{Provider: config.ProviderConfig{BaseURL: server.URL, Model: "test", TimeoutSeconds: 5}}
			staged := []git.StagedFile{{Status: "M", Path: "modules/base.nix"}, {Status: "M", Path: "modules/hardware.nix"}}
			result := loadSplitPlan("diff --git a/modules/base.nix b/modules/base.nix\ndiff --git a/modules/hardware.nix b/modules/hardware.nix", cfg, staged, nil, "Keep hardware separate")()
			if calls != tc.wantCalls {
				t.Errorf("provider calls = %d, want %d", calls, tc.wantCalls)
			}
			if tc.wantError {
				if msg, ok := result.(errMsg); !ok || msg.err == nil {
					t.Fatalf("invalid plan escaped as %T: %+v", result, result)
				}
				return
			}
			ready, ok := result.(splitPlanReadyMsg)
			if !ok {
				t.Fatalf("split plan did not recover: %+v", result)
			}
			if !reflect.DeepEqual(ready.plan, valid) {
				t.Fatalf("plan = %+v, want %+v", ready.plan, valid)
			}
		})
	}
}
