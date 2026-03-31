package metrics

import (
	"testing"

	"github.com/dakaneye/org-pulse/internal/github"
)

func TestMedian(t *testing.T) {
	tests := []struct {
		name string
		vals []float64
		want float64
	}{
		{name: "empty", vals: nil, want: 0},
		{name: "single", vals: []float64{5}, want: 5},
		{name: "odd count", vals: []float64{1, 3, 5}, want: 3},
		{name: "even count", vals: []float64{1, 3, 5, 7}, want: 4},
		{name: "two items", vals: []float64{2, 4}, want: 3},
		{name: "already sorted", vals: []float64{10, 20, 30}, want: 20},
		{name: "unsorted", vals: []float64{30, 10, 20}, want: 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := median(tt.vals)
			if got != tt.want {
				t.Errorf("Median(%v) = %v, want %v", tt.vals, got, tt.want)
			}
		})
	}
}

func TestMedianTimeToFirstReview(t *testing.T) {
	// PR created Monday Mar 2, first review Wednesday Mar 4 = 2 business days
	// PR created Wednesday Mar 4, first review Monday Mar 9 = 3 business days
	prs := []github.PullRequest{
		{
			Number: 1, CreatedAt: d(2026, 3, 2),
			Reviews: []github.Review{
				{SubmittedAt: d(2026, 3, 4), State: "APPROVED"},
			},
		},
		{
			Number: 2, CreatedAt: d(2026, 3, 4),
			Reviews: []github.Review{
				{SubmittedAt: d(2026, 3, 9), State: "CHANGES_REQUESTED"},
			},
		},
	}

	got := MedianTimeToFirstReview(prs)
	// median of [2, 3] = 2.5
	if got != 2.5 {
		t.Errorf("MedianTimeToFirstReview = %v, want 2.5", got)
	}
}

func TestMedianTimeToFirstReviewNoReviews(t *testing.T) {
	prs := []github.PullRequest{
		{Number: 1, CreatedAt: d(2026, 3, 2)},
	}
	got := MedianTimeToFirstReview(prs)
	if got != 0 {
		t.Errorf("MedianTimeToFirstReview with no reviews = %v, want 0", got)
	}
}

func TestMedianTimeToMerge(t *testing.T) {
	// PR created Monday Mar 2, merged Friday Mar 6 = 4 business days
	// PR created Wednesday Mar 4, merged Wednesday Mar 11 = 5 business days
	merged1 := d(2026, 3, 6)
	merged2 := d(2026, 3, 11)
	prs := []github.PullRequest{
		{Number: 1, State: "MERGED", CreatedAt: d(2026, 3, 2), MergedAt: &merged1},
		{Number: 2, State: "MERGED", CreatedAt: d(2026, 3, 4), MergedAt: &merged2},
		{Number: 3, State: "OPEN", CreatedAt: d(2026, 3, 5)}, // not merged, excluded
	}

	got := MedianTimeToMerge(prs)
	// median of [4, 5] = 4.5
	if got != 4.5 {
		t.Errorf("MedianTimeToMerge = %v, want 4.5", got)
	}
}
