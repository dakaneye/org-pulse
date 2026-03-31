// Package github provides GitHub data fetching via the gh CLI.
package github

import "time"

// Repo represents a GitHub repository with ownership and state metadata.
type Repo struct {
	Name       string
	Owner      string
	IsArchived bool
	IsFork     bool
}

// PullRequest holds metadata, reviews, review requests, and check runs for a single PR.
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

// Review represents a single code review submission on a pull request.
type Review struct {
	Author      string
	State       string // APPROVED, CHANGES_REQUESTED, COMMENTED
	SubmittedAt time.Time
}

// ReviewRequest represents a pending review request on a pull request.
type ReviewRequest struct {
	Reviewer string
}

// CheckRun represents a single CI check run associated with a pull request commit.
type CheckRun struct {
	Name       string
	Status     string // COMPLETED, IN_PROGRESS, QUEUED
	Conclusion string // SUCCESS, FAILURE, NEUTRAL, CANCELLED, TIMED_OUT, ACTION_REQUIRED, SKIPPED
}

// RepoData pairs a repository with all pull requests fetched for it.
type RepoData struct {
	Repo         Repo
	PullRequests []PullRequest
}
