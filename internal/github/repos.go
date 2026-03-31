package github

import (
	"context"
	"encoding/json"
	"fmt"
)

type repoListResponse struct {
	Data struct {
		Organization struct {
			Repositories struct {
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
				Nodes []struct {
					Name       string `json:"name"`
					Owner      struct{ Login string } `json:"owner"`
					IsArchived bool `json:"isArchived"`
					IsFork     bool `json:"isFork"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"organization"`
	} `json:"data"`
}

const repoListQuery = `
query($org: String!, $cursor: String) {
  organization(login: $org) {
    repositories(first: 100, after: $cursor) {
      pageInfo { hasNextPage endCursor }
      nodes { name owner { login } isArchived isFork }
    }
  }
}`

// ListRepos returns all repositories in the given GitHub organization, paginating as needed.
func ListRepos(ctx context.Context, org string) ([]Repo, error) {
	var all []Repo
	var cursor string

	for {
		vars := map[string]any{"org": org}
		if cursor != "" {
			vars["cursor"] = cursor
		}

		raw, err := graphQL(ctx, repoListQuery, vars)
		if err != nil {
			return nil, fmt.Errorf("list repos for %s: %w", org, err)
		}

		var resp repoListResponse
		if err := json.Unmarshal(raw, &resp); err != nil {
			return nil, fmt.Errorf("parse repo list: %w", err)
		}

		repos := resp.Data.Organization.Repositories
		for _, n := range repos.Nodes {
			all = append(all, Repo{
				Name:       n.Name,
				Owner:      n.Owner.Login,
				IsArchived: n.IsArchived,
				IsFork:     n.IsFork,
			})
		}

		if !repos.PageInfo.HasNextPage {
			break
		}
		cursor = repos.PageInfo.EndCursor
	}

	return all, nil
}

// FilterRepos removes archived and forked repos and applies optional include/exclude name lists.
func FilterRepos(repos []Repo, include, exclude []string) []Repo {
	includeSet := toSet(include)
	excludeSet := toSet(exclude)

	var filtered []Repo
	for _, r := range repos {
		if r.IsArchived || r.IsFork {
			continue
		}
		if len(includeSet) > 0 && !includeSet[r.Name] {
			continue
		}
		if excludeSet[r.Name] {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

func toSet(items []string) map[string]bool {
	if len(items) == 0 {
		return nil
	}
	s := make(map[string]bool, len(items))
	for _, item := range items {
		s[item] = true
	}
	return s
}
