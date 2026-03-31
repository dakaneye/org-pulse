package metrics

import "github.com/dakaneye/org-pulse/internal/github"

type CIRate struct {
	Total    int
	Failures int
	Rate     float64
}

func CIFailureRate(prs []github.PullRequest) CIRate {
	var total, failures int

	for _, pr := range prs {
		for _, run := range pr.CheckRuns {
			if run.Status != "COMPLETED" {
				continue
			}
			total++
			if run.Conclusion == "FAILURE" {
				failures++
			}
		}
	}

	var rate float64
	if total > 0 {
		rate = float64(failures) / float64(total) * 100
	}

	return CIRate{Total: total, Failures: failures, Rate: rate}
}
