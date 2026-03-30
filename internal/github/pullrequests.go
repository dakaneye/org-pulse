package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type prListResponse struct {
	Data struct {
		Repository struct {
			PullRequests struct {
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
				Nodes []prNode `json:"nodes"`
			} `json:"pullRequests"`
		} `json:"repository"`
	} `json:"data"`
}

type prNode struct {
	Number    int                  `json:"number"`
	Title     string               `json:"title"`
	Author    struct{ Login string } `json:"author"`
	State     string               `json:"state"`
	IsDraft   bool                 `json:"isDraft"`
	CreatedAt time.Time            `json:"createdAt"`
	MergedAt  *time.Time           `json:"mergedAt"`
	ClosedAt  *time.Time           `json:"closedAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	Reviews   struct {
		Nodes []struct {
			Author      struct{ Login string } `json:"author"`
			State       string                 `json:"state"`
			SubmittedAt time.Time              `json:"submittedAt"`
		} `json:"nodes"`
	} `json:"reviews"`
	ReviewRequests struct {
		Nodes []struct {
			RequestedReviewer struct{ Login string } `json:"requestedReviewer"`
		} `json:"nodes"`
	} `json:"reviewRequests"`
	Commits struct {
		Nodes []struct {
			Commit struct {
				CheckSuites struct {
					Nodes []struct {
						CheckRuns struct {
							Nodes []struct {
								Name       string `json:"name"`
								Status     string `json:"status"`
								Conclusion string `json:"conclusion"`
							} `json:"nodes"`
						} `json:"checkRuns"`
					} `json:"nodes"`
				} `json:"checkSuites"`
			} `json:"commit"`
		} `json:"nodes"`
	} `json:"commits"`
}

const prListQuery = `
query($owner: String!, $repo: String!, $cursor: String) {
  repository(owner: $owner, name: $repo) {
    pullRequests(first: 50, after: $cursor, orderBy: {field: UPDATED_AT, direction: DESC}) {
      pageInfo { hasNextPage endCursor }
      nodes {
        number title
        author { login }
        state isDraft
        createdAt mergedAt closedAt updatedAt
        reviews(first: 50) {
          nodes { author { login } state submittedAt }
        }
        reviewRequests(first: 20) {
          nodes { requestedReviewer { ... on User { login } } }
        }
        commits(last: 1) {
          nodes {
            commit {
              checkSuites(first: 10) {
                nodes {
                  checkRuns(first: 50) {
                    nodes { name status conclusion }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`

func FetchPullRequests(ctx context.Context, owner, repo string, since time.Time) ([]PullRequest, error) {
	var all []PullRequest
	var cursor string

	for {
		vars := map[string]any{"owner": owner, "repo": repo}
		if cursor != "" {
			vars["cursor"] = cursor
		}

		raw, err := graphQL(ctx, prListQuery, vars)
		if err != nil {
			return nil, fmt.Errorf("fetch PRs for %s/%s: %w", owner, repo, err)
		}

		var resp prListResponse
		if err := json.Unmarshal(raw, &resp); err != nil {
			return nil, fmt.Errorf("parse PRs for %s/%s: %w", owner, repo, err)
		}

		prs := resp.Data.Repository.PullRequests
		pastWindow := false

		for _, n := range prs.Nodes {
			if n.UpdatedAt.Before(since) {
				pastWindow = true
				break
			}
			if n.IsDraft {
				continue
			}

			pr := PullRequest{
				Number:    n.Number,
				Title:     n.Title,
				Author:    n.Author.Login,
				State:     n.State,
				IsDraft:   n.IsDraft,
				CreatedAt: n.CreatedAt,
				MergedAt:  n.MergedAt,
				ClosedAt:  n.ClosedAt,
				UpdatedAt: n.UpdatedAt,
			}

			for _, r := range n.Reviews.Nodes {
				pr.Reviews = append(pr.Reviews, Review{
					Author:      r.Author.Login,
					State:       r.State,
					SubmittedAt: r.SubmittedAt,
				})
			}

			for _, rr := range n.ReviewRequests.Nodes {
				pr.ReviewRequests = append(pr.ReviewRequests, ReviewRequest{
					Reviewer: rr.RequestedReviewer.Login,
				})
			}

			if len(n.Commits.Nodes) > 0 {
				commit := n.Commits.Nodes[0].Commit
				for _, suite := range commit.CheckSuites.Nodes {
					for _, run := range suite.CheckRuns.Nodes {
						pr.CheckRuns = append(pr.CheckRuns, CheckRun{
							Name:       run.Name,
							Status:     run.Status,
							Conclusion: run.Conclusion,
						})
					}
				}
			}

			all = append(all, pr)
		}

		if pastWindow || !prs.PageInfo.HasNextPage {
			break
		}
		cursor = prs.PageInfo.EndCursor
	}

	return all, nil
}
