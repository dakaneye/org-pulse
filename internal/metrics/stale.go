package metrics

import (
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
)

// StaleCategory classifies why an open PR is considered stale.
type StaleCategory string

const (
	// StaleNoReviewActivity marks PRs with no review requests and no reviews.
	StaleNoReviewActivity StaleCategory = "No review activity"
	// StaleNoReviewResponse marks PRs where a reviewer was requested but has not responded.
	StaleNoReviewResponse StaleCategory = "No review response"
	// StaleApprovedNotMerged marks PRs that have been approved but remain unmerged.
	StaleApprovedNotMerged StaleCategory = "Approved, not merged"
)

// StalePR pairs an open pull request with its staleness classification and age.
type StalePR struct {
	PR       github.PullRequest
	Repo     string
	AgeDays  float64
	Category StaleCategory
}

// StalePRs returns open PRs that exceed thresholdDays of inactivity, classified by category.
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

	// For remaining checks, age is based on PR creation since GitHub GraphQL
	// doesn't expose a timestamp on ReviewRequest objects.
	age := BusinessDaysBetween(pr.CreatedAt, now)
	if age < threshold {
		return "", 0, false
	}

	if len(pr.ReviewRequests) > 0 && len(pr.Reviews) == 0 {
		return StaleNoReviewResponse, age, true
	}
	if len(pr.Reviews) == 0 && len(pr.ReviewRequests) == 0 {
		return StaleNoReviewActivity, age, true
	}

	return "", 0, false
}
