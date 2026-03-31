package metrics

import (
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
)

type StaleCategory string

const (
	StaleNoReviewActivity  StaleCategory = "No review activity"
	StaleNoReviewResponse  StaleCategory = "No review response"
	StaleApprovedNotMerged StaleCategory = "Approved, not merged"
)

type StalePR struct {
	PR       github.PullRequest
	Repo     string
	AgeDays  float64
	Category StaleCategory
}

func StalePRs(prs []github.PullRequest, now time.Time, thresholdDays float64) []StalePR {
	var stale []StalePR

	for _, pr := range prs {
		if pr.State != "OPEN" {
			continue
		}

		if cat, age, ok := classifyStale(pr, now, thresholdDays); ok {
			stale = append(stale, StalePR{
				PR:       pr,
				AgeDays:  age,
				Category: cat,
			})
		}
	}

	return stale
}

func classifyStale(pr github.PullRequest, now time.Time, threshold float64) (StaleCategory, float64, bool) {
	// Check approved but not merged
	for _, r := range pr.Reviews {
		if r.State == "APPROVED" {
			age := BusinessDaysBetween(r.SubmittedAt, now)
			if age >= threshold {
				return StaleApprovedNotMerged, age, true
			}
		}
	}

	// Check review requested but no response
	// Use pr.CreatedAt as reference since GitHub GraphQL doesn't expose
	// a timestamp on the ReviewRequest connection.
	if len(pr.ReviewRequests) > 0 && len(pr.Reviews) == 0 {
		age := BusinessDaysBetween(pr.CreatedAt, now)
		if age >= threshold {
			return StaleNoReviewResponse, age, true
		}
	}

	// Check no review activity at all
	if len(pr.Reviews) == 0 && len(pr.ReviewRequests) == 0 {
		age := BusinessDaysBetween(pr.CreatedAt, now)
		if age >= threshold {
			return StaleNoReviewActivity, age, true
		}
	}

	return "", 0, false
}
