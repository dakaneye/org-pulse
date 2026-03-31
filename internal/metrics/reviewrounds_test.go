package metrics

import (
	"testing"

	"github.com/dakaneye/org-pulse/internal/github"
)

func TestReviewRounds(t *testing.T) {
	merged := d(2026, 3, 15)
	prs := []github.PullRequest{
		{
			Number: 1, State: "MERGED", MergedAt: &merged,
			Reviews: []github.Review{
				{State: "APPROVED", SubmittedAt: d(2026, 3, 5)},
			},
		},
		{
			Number: 2, State: "MERGED", MergedAt: &merged,
			Reviews: []github.Review{
				{State: "CHANGES_REQUESTED", SubmittedAt: d(2026, 3, 3)},
				{State: "APPROVED", SubmittedAt: d(2026, 3, 5)},
			},
		},
		{
			Number: 3, State: "MERGED", MergedAt: &merged,
			Reviews: []github.Review{
				{State: "CHANGES_REQUESTED", SubmittedAt: d(2026, 3, 2)},
				{State: "CHANGES_REQUESTED", SubmittedAt: d(2026, 3, 4)},
				{State: "APPROVED", SubmittedAt: d(2026, 3, 6)},
			},
		},
		{
			Number: 4, State: "OPEN",
			Reviews: []github.Review{
				{State: "CHANGES_REQUESTED", SubmittedAt: d(2026, 3, 2)},
			},
		},
		{
			Number: 5, State: "MERGED", MergedAt: &merged,
			Reviews: []github.Review{
				{State: "COMMENTED", SubmittedAt: d(2026, 3, 3)},
				{State: "APPROVED", SubmittedAt: d(2026, 3, 5)},
			},
		},
	}

	dist := ReviewRoundsDistribution(prs)

	expected := map[string]int{"1": 2, "2": 1, "3+": 1}
	for _, bucket := range dist {
		want, ok := expected[bucket.Label]
		if !ok {
			t.Errorf("unexpected bucket %s", bucket.Label)
			continue
		}
		if bucket.Count != want {
			t.Errorf("bucket %s count = %d, want %d", bucket.Label, bucket.Count, want)
		}
	}
}

func TestReviewRoundsEmpty(t *testing.T) {
	dist := ReviewRoundsDistribution(nil)
	if len(dist) != 0 {
		t.Errorf("expected empty, got %d buckets", len(dist))
	}
}
