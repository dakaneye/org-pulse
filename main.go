package main

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/dakaneye/org-pulse/internal/github"
	"github.com/dakaneye/org-pulse/internal/metrics"
	"github.com/dakaneye/org-pulse/internal/report"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func main() {
	root := &cobra.Command{
		Use:   "org-pulse",
		Short: "GitHub org delivery health reports",
	}

	var (
		org          string
		weeks        int
		format       string
		repos        []string
		excludeRepos []string
		concurrency  int
	)

	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Generate a delivery health report for a GitHub org",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(repos) > 0 && len(excludeRepos) > 0 {
				return fmt.Errorf("--repos and --exclude-repos are mutually exclusive")
			}

			ctx := cmd.Context()

			if err := github.CheckAuth(ctx); err != nil {
				return err
			}
			now := time.Now()
			since := now.AddDate(0, 0, -7*weeks)

			fmt.Fprintf(os.Stderr, "Fetching repos for %s...\n", org)
			allRepos, err := github.ListRepos(ctx, org)
			if err != nil {
				return err
			}

			filtered := github.FilterRepos(allRepos, repos, excludeRepos)
			if len(filtered) == 0 {
				fmt.Fprintln(os.Stderr, "No repos found matching filters.")
				return nil
			}
			fmt.Fprintf(os.Stderr, "Found %d repos. Fetching PR data...\n", len(filtered))

			var (
				mu      sync.Mutex
				results []github.RepoData
			)

			g, gctx := errgroup.WithContext(ctx)
			g.SetLimit(concurrency)

			for _, repo := range filtered {
				g.Go(func() error {
					prs, err := github.FetchPullRequests(gctx, repo.Owner, repo.Name, since)
					if err != nil {
						return fmt.Errorf("repo %s: %w", repo.Name, err)
					}
					mu.Lock()
					results = append(results, github.RepoData{Repo: repo, PullRequests: prs})
					mu.Unlock()
					return nil
				})
			}

			if err := g.Wait(); err != nil {
				return err
			}

			sort.Slice(results, func(i, j int) bool {
				return results[i].Repo.Name < results[j].Repo.Name
			})

			r := metrics.ComputeReport(org, results, now, since, weeks, 7)

			switch format {
			case "json":
				out, err := report.RenderJSON(r)
				if err != nil {
					return err
				}
				fmt.Println(out)
			default:
				fmt.Print(report.RenderMarkdown(r))
			}

			return nil
		},
	}

	reportCmd.Flags().StringVar(&org, "org", "", "GitHub organization name (required)")
	reportCmd.Flags().IntVar(&weeks, "weeks", 4, "Lookback window in weeks")
	reportCmd.Flags().StringVar(&format, "format", "md", "Output format: md or json")
	reportCmd.Flags().StringSliceVar(&repos, "repos", nil, "Only include these repos")
	reportCmd.Flags().StringSliceVar(&excludeRepos, "exclude-repos", nil, "Exclude these repos")
	reportCmd.Flags().IntVar(&concurrency, "concurrency", 10, "Max parallel gh subprocess calls")
	_ = reportCmd.MarkFlagRequired("org")

	root.AddCommand(reportCmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
