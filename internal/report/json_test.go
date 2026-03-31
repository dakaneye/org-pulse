package report

import (
	"encoding/json"
	"testing"

	"github.com/dakaneye/org-pulse/internal/metrics"
)

func TestRenderJSON(t *testing.T) {
	r := metrics.Report{
		Org:       "test-org",
		RepoCount: 5,
		Velocity:  metrics.Velocity{Opened: 10, Merged: 8, Closed: 1},
	}

	out, err := RenderJSON(r)
	if err != nil {
		t.Fatalf("RenderJSON error: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if parsed["Org"] != "test-org" {
		t.Errorf("Org = %v, want test-org", parsed["Org"])
	}
}
