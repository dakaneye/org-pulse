package metrics

import (
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
)

// RepoCI holds CI failure rate data for a single repository.
type RepoCI struct {
	Repo string
	CIRate
}

// Report holds all computed metrics for an org pulse report.
type Report struct {
	Org       string
	RepoCount int
	Window    struct {
		Start time.Time
		End   time.Time
		Weeks int
	}

	Velocity                Velocity
	MedianTimeToFirstReview float64
	MedianTimeToMerge       float64
	ReviewRounds            []RoundsBucket
	ReviewLoads             []ReviewerLoad
	StalePRs                []StalePR
	OrgCIRate               CIRate
	RepoCIRates             []RepoCI
}

// ComputeReport aggregates all metrics across the given repo data.
func ComputeReport(org string, data []github.RepoData, now time.Time, windowStart time.Time, weeks int, staleThreshold float64) Report {
	var allPRs []github.PullRequest
	for _, rd := range data {
		allPRs = append(allPRs, rd.PullRequests...)
	}

	r := Report{
		Org:       org,
		RepoCount: len(data),
	}
	r.Window.Start = windowStart
	r.Window.End = now
	r.Window.Weeks = weeks

	r.Velocity = PRVelocity(allPRs, windowStart, now)
	r.MedianTimeToFirstReview = MedianTimeToFirstReview(allPRs)
	r.MedianTimeToMerge = MedianTimeToMerge(allPRs)
	r.ReviewRounds = ReviewRoundsDistribution(allPRs)
	r.ReviewLoads = ReviewLoad(allPRs)
	r.OrgCIRate = CIFailureRate(allPRs)

	// Per-repo metrics: stale PRs and CI rates in a single pass
	for _, rd := range data {
		repoStale := StalePRs(rd.PullRequests, now, staleThreshold)
		for i := range repoStale {
			repoStale[i].Repo = rd.Repo.Name
		}
		r.StalePRs = append(r.StalePRs, repoStale...)

		if len(rd.PullRequests) > 0 {
			rate := CIFailureRate(rd.PullRequests)
			if rate.Total > 0 {
				r.RepoCIRates = append(r.RepoCIRates, RepoCI{
					Repo:   rd.Repo.Name,
					CIRate: rate,
				})
			}
		}
	}

	return r
}
