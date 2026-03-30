package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

	report := &cobra.Command{
		Use:   "report",
		Short: "Generate a delivery health report for a GitHub org",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(repos) > 0 && len(excludeRepos) > 0 {
				return fmt.Errorf("--repos and --exclude-repos are mutually exclusive")
			}
			// TODO: wire up in Task 16
			fmt.Fprintf(os.Stderr, "org=%s weeks=%d format=%s\n", org, weeks, format)
			return nil
		},
	}

	report.Flags().StringVar(&org, "org", "", "GitHub organization name (required)")
	report.Flags().IntVar(&weeks, "weeks", 4, "Lookback window in weeks")
	report.Flags().StringVar(&format, "format", "md", "Output format: md or json")
	report.Flags().StringSliceVar(&repos, "repos", nil, "Only include these repos")
	report.Flags().StringSliceVar(&excludeRepos, "exclude-repos", nil, "Exclude these repos")
	report.Flags().IntVar(&concurrency, "concurrency", 10, "Max parallel gh subprocess calls")
	_ = report.MarkFlagRequired("org")

	root.AddCommand(report)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
