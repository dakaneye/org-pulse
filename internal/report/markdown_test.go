package report

import (
	"strings"
	"testing"
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
	"github.com/dakaneye/org-pulse/internal/metrics"
)

func TestRenderMarkdown(t *testing.T) {
	merged := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	r := metrics.Report{
		Org:       "test-org",
		RepoCount: 2,
		Velocity:  metrics.Velocity{Opened: 10, Merged: 8, Closed: 1},
		MedianTimeToFirstReview: 2.5,
		MedianTimeToMerge:       4.0,
		ReviewRounds: []metrics.RoundsBucket{
			{Label: "1", Count: 5, Percent: 62.5},
			{Label: "2", Count: 3, Percent: 37.5},
		},
		ReviewLoads: []metrics.ReviewerLoad{
			{Reviewer: "alice", Count: 10, Percent: 66.7},
			{Reviewer: "bob", Count: 5, Percent: 33.3},
		},
		StalePRs: []metrics.StalePR{
			{
				PR:       github.PullRequest{Number: 42, MergedAt: &merged},
				Repo:     "repo-a",
				AgeDays:  12,
				Category: metrics.StaleNoReviewActivity,
			},
		},
		OrgCIRate: metrics.CIRate{Total: 100, Failures: 14, Rate: 14.0},
		RepoCIRates: []metrics.RepoCI{
			{Repo: "repo-a", CIRate: metrics.CIRate{Total: 60, Failures: 10, Rate: 16.7}},
			{Repo: "repo-b", CIRate: metrics.CIRate{Total: 40, Failures: 4, Rate: 10.0}},
		},
	}
	r.Window.Start = time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	r.Window.End = time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC)
	r.Window.Weeks = 4

	out := RenderMarkdown(r)

	checks := []string{
		"# Org Pulse: test-org",
		"4-week report",
		"2 repos",
		"PRs opened: 10",
		"merged: 8",
		"closed: 1",
		"2.5 business days",
		"4.0 business days",
		"14.0%",
		"alice",
		"bob",
		"#42",
		"repo-a",
		"No review activity",
	}

	for _, check := range checks {
		if !strings.Contains(out, check) {
			t.Errorf("output missing %q", check)
		}
	}
}
