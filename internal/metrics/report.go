package metrics

import (
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
)

type RepoCI struct {
	Repo string
	CIRate
}

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

func ComputeReport(org string, data []github.RepoData, now time.Time, weeks int, staleThreshold float64) Report {
	windowStart := now.AddDate(0, 0, -7*weeks)

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

	// Compute stale PRs per repo to avoid PR number collisions across repos
	for _, rd := range data {
		repoStale := StalePRs(rd.PullRequests, now, staleThreshold)
		for i := range repoStale {
			repoStale[i].Repo = rd.Repo.Name
		}
		r.StalePRs = append(r.StalePRs, repoStale...)
	}

	for _, rd := range data {
		if len(rd.PullRequests) == 0 {
			continue
		}
		rate := CIFailureRate(rd.PullRequests)
		if rate.Total == 0 {
			continue
		}
		r.RepoCIRates = append(r.RepoCIRates, RepoCI{
			Repo:   rd.Repo.Name,
			CIRate: rate,
		})
	}

	return r
}
