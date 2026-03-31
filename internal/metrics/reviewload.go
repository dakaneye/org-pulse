package metrics

import (
	"sort"

	"github.com/dakaneye/org-pulse/internal/github"
)

type ReviewerLoad struct {
	Reviewer string
	Count    int
	Percent  float64
}

func ReviewLoad(prs []github.PullRequest) []ReviewerLoad {
	counts := make(map[string]int)
	total := 0

	for _, pr := range prs {
		for _, r := range pr.Reviews {
			counts[r.Author]++
			total++
		}
	}

	if total == 0 {
		return nil
	}

	loads := make([]ReviewerLoad, 0, len(counts))
	for reviewer, count := range counts {
		loads = append(loads, ReviewerLoad{
			Reviewer: reviewer,
			Count:    count,
			Percent:  float64(count) / float64(total) * 100,
		})
	}

	sort.Slice(loads, func(i, j int) bool {
		return loads[i].Count > loads[j].Count
	})

	return loads
}
