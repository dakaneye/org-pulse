package metrics

import (
	"testing"
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
)

func TestComputeReport(t *testing.T) {
	now := d(2026, 3, 30)
	windowStart := d(2026, 3, 2)
	merged := d(2026, 3, 15)

	data := []github.RepoData{
		{
			Repo: github.Repo{Name: "repo-a", Owner: "org"},
			PullRequests: []github.PullRequest{
				{
					Number: 1, State: "MERGED", Author: "alice",
					CreatedAt: d(2026, 3, 5), MergedAt: &merged,
					Reviews: []github.Review{
						{Author: "bob", State: "APPROVED", SubmittedAt: d(2026, 3, 7)},
					},
					CheckRuns: []github.CheckRun{
						{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"},
					},
				},
				{
					Number: 2, State: "OPEN", Author: "charlie",
					CreatedAt: d(2026, 3, 10),
				},
			},
		},
		{
			Repo: github.Repo{Name: "repo-b", Owner: "org"},
			PullRequests: []github.PullRequest{
				{
					Number: 1, State: "OPEN", Author: "dave",
					CreatedAt: d(2026, 3, 3),
					ReviewRequests: []github.ReviewRequest{{Reviewer: "eve"}},
				},
			},
		},
	}

	r := ComputeReport("org", data, now, windowStart, 4, 7)

	if r.Org != "org" {
		t.Errorf("Org = %q, want %q", r.Org, "org")
	}
	if r.RepoCount != 2 {
		t.Errorf("RepoCount = %d, want 2", r.RepoCount)
	}
	if r.Velocity.Opened != 3 {
		t.Errorf("Opened = %d, want 3", r.Velocity.Opened)
	}
	if r.Velocity.Merged != 1 {
		t.Errorf("Merged = %d, want 1", r.Velocity.Merged)
	}

	// Stale PRs should be tagged with repo names
	for _, s := range r.StalePRs {
		if s.Repo == "" {
			t.Errorf("stale PR #%d has empty Repo", s.PR.Number)
		}
	}

	// Review load should include bob
	found := false
	for _, rl := range r.ReviewLoads {
		if rl.Reviewer == "bob" {
			found = true
		}
	}
	if !found {
		t.Error("ReviewLoads missing bob")
	}

	// CI rate should have repo-a
	if len(r.RepoCIRates) != 1 || r.RepoCIRates[0].Repo != "repo-a" {
		t.Errorf("RepoCIRates = %v, want [repo-a]", r.RepoCIRates)
	}
}

func TestComputeReportEmpty(t *testing.T) {
	now := d(2026, 3, 30)
	windowStart := d(2026, 3, 2)
	r := ComputeReport("org", nil, now, windowStart, 4, 7)
	if r.RepoCount != 0 {
		t.Errorf("RepoCount = %d, want 0", r.RepoCount)
	}
	if r.Velocity.Opened != 0 {
		t.Errorf("Opened = %d, want 0", r.Velocity.Opened)
	}
}

// Ensure window fields are populated correctly.
func TestComputeReportWindow(t *testing.T) {
	now := d(2026, 3, 30)
	windowStart := d(2026, 3, 2)
	r := ComputeReport("org", nil, now, windowStart, 4, 7)
	if !r.Window.Start.Equal(windowStart) {
		t.Errorf("Window.Start = %v, want %v", r.Window.Start, windowStart)
	}
	if !r.Window.End.Equal(now) {
		t.Errorf("Window.End = %v, want %v", r.Window.End, now)
	}
	if r.Window.Weeks != 4 {
		t.Errorf("Window.Weeks = %d, want 4", r.Window.Weeks)
	}
}

// Confirm MergedAt pointer is handled without panic on nil.
func TestComputeReportNilMergedAt(t *testing.T) {
	now := d(2026, 3, 30)
	windowStart := d(2026, 3, 2)
	data := []github.RepoData{
		{
			Repo: github.Repo{Name: "repo-a", Owner: "org"},
			PullRequests: []github.PullRequest{
				{Number: 1, State: "OPEN", Author: "alice", CreatedAt: d(2026, 3, 5)},
			},
		},
	}
	// Should not panic.
	_ = ComputeReport("org", data, now, windowStart, 4, 7)
}

var _ = time.Time{} // keep time import used
