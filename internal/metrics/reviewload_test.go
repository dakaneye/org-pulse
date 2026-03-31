package metrics

import (
	"testing"

	"github.com/dakaneye/org-pulse/internal/github"
)

func TestReviewLoad(t *testing.T) {
	prs := []github.PullRequest{
		{
			Number: 1,
			Reviews: []github.Review{
				{Author: "alice", State: "APPROVED"},
				{Author: "bob", State: "CHANGES_REQUESTED"},
			},
		},
		{
			Number: 2,
			Reviews: []github.Review{
				{Author: "alice", State: "APPROVED"},
				{Author: "alice", State: "COMMENTED"},
				{Author: "charlie", State: "APPROVED"},
			},
		},
	}

	loads := ReviewLoad(prs)

	expected := map[string]int{"alice": 3, "bob": 1, "charlie": 1}
	for _, rl := range loads {
		want, ok := expected[rl.Reviewer]
		if !ok {
			t.Errorf("unexpected reviewer %s", rl.Reviewer)
			continue
		}
		if rl.Count != want {
			t.Errorf("%s count = %d, want %d", rl.Reviewer, rl.Count, want)
		}
	}

	var totalPct float64
	for _, rl := range loads {
		totalPct += rl.Percent
	}
	if totalPct < 99.9 || totalPct > 100.1 {
		t.Errorf("percentages sum to %v, want ~100", totalPct)
	}
}

func TestReviewLoadEmpty(t *testing.T) {
	loads := ReviewLoad(nil)
	if len(loads) != 0 {
		t.Errorf("expected empty, got %d entries", len(loads))
	}
}
