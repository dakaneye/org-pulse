package metrics

import (
	"testing"
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
)

func TestPRVelocity(t *testing.T) {
	windowStart := d(2026, 3, 1)
	windowEnd := d(2026, 3, 30)
	merged := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	closed := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)

	prs := []github.PullRequest{
		{Number: 1, State: "OPEN", CreatedAt: d(2026, 3, 5)},
		{Number: 2, State: "MERGED", CreatedAt: d(2026, 3, 2), MergedAt: &merged},
		{Number: 3, State: "CLOSED", CreatedAt: d(2026, 3, 10), ClosedAt: &closed},
		{Number: 4, State: "MERGED", CreatedAt: d(2026, 2, 15), MergedAt: &merged},
		{Number: 5, State: "OPEN", CreatedAt: d(2026, 2, 10)},
	}

	v := PRVelocity(prs, windowStart, windowEnd)

	if v.Opened != 3 {
		t.Errorf("Opened = %d, want 3", v.Opened)
	}
	if v.Merged != 2 {
		t.Errorf("Merged = %d, want 2", v.Merged)
	}
	if v.Closed != 1 {
		t.Errorf("Closed = %d, want 1", v.Closed)
	}
}

func TestPRVelocityEmpty(t *testing.T) {
	v := PRVelocity(nil, d(2026, 3, 1), d(2026, 3, 30))
	if v.Opened != 0 || v.Merged != 0 || v.Closed != 0 {
		t.Errorf("expected all zeros for empty input, got %+v", v)
	}
}
