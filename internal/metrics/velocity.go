package metrics

import (
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
)

// Velocity holds PR opened/merged/closed counts within a time window.
type Velocity struct {
	Opened int
	Merged int
	Closed int
}

// PRVelocity counts PRs opened, merged, and closed within [windowStart, windowEnd).
func PRVelocity(prs []github.PullRequest, windowStart, windowEnd time.Time) Velocity {
	var v Velocity
	for _, pr := range prs {
		if inWindow(pr.CreatedAt, windowStart, windowEnd) {
			v.Opened++
		}
		if pr.MergedAt != nil && inWindow(*pr.MergedAt, windowStart, windowEnd) {
			v.Merged++
		}
		if pr.State == "CLOSED" && pr.ClosedAt != nil && inWindow(*pr.ClosedAt, windowStart, windowEnd) {
			v.Closed++
		}
	}
	return v
}

// inWindow returns true if t is in the half-open interval [start, end).
func inWindow(t, start, end time.Time) bool {
	return !t.Before(start) && t.Before(end)
}
