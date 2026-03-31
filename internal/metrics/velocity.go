package metrics

import (
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
)

type Velocity struct {
	Opened int
	Merged int
	Closed int
}

func PRVelocity(prs []github.PullRequest, windowStart, windowEnd time.Time) Velocity {
	var v Velocity
	for _, pr := range prs {
		if !pr.CreatedAt.Before(windowStart) && pr.CreatedAt.Before(windowEnd) {
			v.Opened++
		}
		if pr.MergedAt != nil && !pr.MergedAt.Before(windowStart) && pr.MergedAt.Before(windowEnd) {
			v.Merged++
		}
		if pr.State == "CLOSED" && pr.ClosedAt != nil && !pr.ClosedAt.Before(windowStart) && pr.ClosedAt.Before(windowEnd) {
			v.Closed++
		}
	}
	return v
}
