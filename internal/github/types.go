package github

import "time"

type Repo struct {
	Name       string
	Owner      string
	IsArchived bool
	IsFork     bool
}

type PullRequest struct {
	Number    int
	Title     string
	Author    string
	State     string // OPEN, MERGED, CLOSED
	IsDraft   bool
	CreatedAt time.Time
	MergedAt  *time.Time
	ClosedAt  *time.Time
	UpdatedAt time.Time
	Reviews   []Review
	ReviewRequests []ReviewRequest
	CheckRuns []CheckRun
}

type Review struct {
	Author      string
	State       string // APPROVED, CHANGES_REQUESTED, COMMENTED
	SubmittedAt time.Time
}

type ReviewRequest struct {
	Reviewer    string
	RequestedAt time.Time
}

type CheckRun struct {
	Name       string
	Status     string // COMPLETED, IN_PROGRESS, QUEUED
	Conclusion string // SUCCESS, FAILURE, NEUTRAL, CANCELLED, TIMED_OUT, ACTION_REQUIRED, SKIPPED
}

type RepoData struct {
	Repo         Repo
	PullRequests []PullRequest
}
