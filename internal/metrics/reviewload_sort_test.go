package metrics

import (
	"testing"

	"github.com/dakaneye/org-pulse/internal/github"
)

func TestReviewLoadSortOrder(t *testing.T) {
	prs := []github.PullRequest{
		{
			Reviews: []github.Review{
				{Author: "alice"},
				{Author: "alice"},
				{Author: "alice"},
				{Author: "bob"},
			},
		},
	}

	loads := ReviewLoad(prs)

	if len(loads) != 2 {
		t.Fatalf("got %d reviewers, want 2", len(loads))
	}
	if loads[0].Reviewer != "alice" {
		t.Errorf("first reviewer = %s, want alice (highest count)", loads[0].Reviewer)
	}
	if loads[1].Reviewer != "bob" {
		t.Errorf("second reviewer = %s, want bob", loads[1].Reviewer)
	}
}
