package metrics

import (
	"testing"
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
)

func TestStalePRs(t *testing.T) {
	now := d(2026, 3, 30) // Monday
	threshold := float64(7)

	prs := []github.PullRequest{
		{
			// No review activity for 10+ business days — stale
			Number: 1, State: "OPEN", CreatedAt: d(2026, 3, 10),
			Reviews: nil,
		},
		{
			// Review requested, no reviews submitted, created 17+ bd ago — stale
			Number: 2, State: "OPEN", CreatedAt: d(2026, 3, 5),
			ReviewRequests: []github.ReviewRequest{
				{Reviewer: "alice"},
			},
		},
		{
			// Approved 10+ business days ago but not merged — stale
			Number: 3, State: "OPEN", CreatedAt: d(2026, 3, 5),
			Reviews: []github.Review{
				{State: "APPROVED", SubmittedAt: d(2026, 3, 13)},
			},
		},
		{
			// Created 3 business days ago, no reviews — NOT stale (under threshold)
			Number: 4, State: "OPEN", CreatedAt: d(2026, 3, 25),
		},
		{
			// Merged — NOT stale
			Number: 5, State: "MERGED", CreatedAt: d(2026, 3, 5),
			MergedAt: func() *time.Time { t := d(2026, 3, 20); return &t }(),
		},
		{
			// Approved 5 business days ago — NOT stale (under threshold)
			Number: 6, State: "OPEN", CreatedAt: d(2026, 3, 5),
			Reviews: []github.Review{
				{State: "APPROVED", SubmittedAt: d(2026, 3, 23)},
			},
		},
	}

	stale := StalePRs(prs, now, threshold)

	if len(stale) != 3 {
		t.Fatalf("got %d stale PRs, want 3", len(stale))
	}

	staleNumbers := map[int]bool{}
	for _, s := range stale {
		staleNumbers[s.PR.Number] = true
	}

	for _, num := range []int{1, 2, 3} {
		if !staleNumbers[num] {
			t.Errorf("PR #%d should be stale", num)
		}
	}
}

func TestStalePRsBoundary(t *testing.T) {
	now := d(2026, 3, 30) // Monday
	threshold := float64(7)

	prs := []github.PullRequest{
		{
			// Exactly 7 business days — stale (at threshold)
			// Mar 30 (Mon) minus 7 business days = Mar 19 (Thu)
			Number: 1, State: "OPEN", CreatedAt: d(2026, 3, 19),
		},
		{
			// 6 business days — NOT stale (under threshold)
			// Mar 30 (Mon) minus 6 business days = Mar 20 (Fri)
			Number: 2, State: "OPEN", CreatedAt: d(2026, 3, 20),
		},
	}

	stale := StalePRs(prs, now, threshold)

	if len(stale) != 1 {
		t.Fatalf("got %d stale PRs, want 1", len(stale))
	}
	if stale[0].PR.Number != 1 {
		t.Errorf("expected PR #1 to be stale, got #%d", stale[0].PR.Number)
	}
}

func TestStalePRsEmpty(t *testing.T) {
	stale := StalePRs(nil, d(2026, 3, 30), 7)
	if len(stale) != 0 {
		t.Errorf("expected empty, got %d", len(stale))
	}
}
