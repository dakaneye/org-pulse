package metrics

import (
	"sort"

	"github.com/dakaneye/org-pulse/internal/github"
)

func Median(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sorted := make([]float64, len(vals))
	copy(sorted, vals)
	sort.Float64s(sorted)

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

func MedianTimeToFirstReview(prs []github.PullRequest) float64 {
	var days []float64
	for _, pr := range prs {
		if len(pr.Reviews) == 0 {
			continue
		}
		earliest := pr.Reviews[0].SubmittedAt
		for _, r := range pr.Reviews[1:] {
			if r.SubmittedAt.Before(earliest) {
				earliest = r.SubmittedAt
			}
		}
		days = append(days, BusinessDaysBetween(pr.CreatedAt, earliest))
	}
	return Median(days)
}

func MedianTimeToMerge(prs []github.PullRequest) float64 {
	var days []float64
	for _, pr := range prs {
		if pr.MergedAt == nil {
			continue
		}
		days = append(days, BusinessDaysBetween(pr.CreatedAt, *pr.MergedAt))
	}
	return Median(days)
}
