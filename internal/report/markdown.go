// Package report provides rendering of metrics reports as markdown or JSON.
package report

import (
	"fmt"
	"strings"

	"github.com/dakaneye/org-pulse/internal/metrics"
)

// RenderMarkdown formats a Report as a human-readable markdown document.
func RenderMarkdown(r metrics.Report) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# Org Pulse: %s\n", r.Org)
	fmt.Fprintf(&b, "> %d-week report (%s – %s) | %d repos\n\n",
		r.Window.Weeks,
		r.Window.Start.Format("Jan 2"),
		r.Window.End.Format("Jan 2, 2006"),
		r.RepoCount)

	// Summary
	b.WriteString("## Summary\n")
	fmt.Fprintf(&b, "- PRs opened: %d | merged: %d | closed: %d\n",
		r.Velocity.Opened, r.Velocity.Merged, r.Velocity.Closed)
	fmt.Fprintf(&b, "- Median time-to-first-review: %s\n", formatDuration(r.MedianTimeToFirstReview))
	fmt.Fprintf(&b, "- Median time-to-merge: %s\n", formatDuration(r.MedianTimeToMerge))
	fmt.Fprintf(&b, "- CI failure rate: %.1f%%\n\n", r.OrgCIRate.Rate)

	// Review Load
	if len(r.ReviewLoads) > 0 {
		b.WriteString("## Review Load\n")
		b.WriteString("| Reviewer | Reviews | % of Total |\n")
		b.WriteString("|----------|---------|------------|\n")
		for _, rl := range r.ReviewLoads {
			fmt.Fprintf(&b, "| @%s | %d | %.0f%% |\n", rl.Reviewer, rl.Count, rl.Percent)
		}
		b.WriteString("\n")
	}

	// Review Rounds
	if len(r.ReviewRounds) > 0 {
		b.WriteString("## Review Rounds Distribution\n")
		b.WriteString("| Rounds | Count | % |\n")
		b.WriteString("|--------|-------|---|\n")
		for _, bucket := range r.ReviewRounds {
			fmt.Fprintf(&b, "| %s | %d | %.0f%% |\n", bucket.Label, bucket.Count, bucket.Percent)
		}
		b.WriteString("\n")
	}

	// Stale PRs
	if len(r.StalePRs) > 0 {
		fmt.Fprintf(&b, "## Stale PRs (%d)\n", len(r.StalePRs))
		b.WriteString("| PR | Repo | Age (bd) | Reason |\n")
		b.WriteString("|----|------|----------|--------|\n")
		for _, s := range r.StalePRs {
			fmt.Fprintf(&b, "| #%d | %s | %.0f | %s |\n",
				s.PR.Number, s.Repo, s.AgeDays, s.Category)
		}
		b.WriteString("\n")
	}

	// CI by Repo
	if len(r.RepoCIRates) > 0 {
		b.WriteString("## CI Failure Rate by Repo\n")
		b.WriteString("| Repo | Runs | Failures | Rate |\n")
		b.WriteString("|------|------|----------|------|\n")
		for _, rc := range r.RepoCIRates {
			fmt.Fprintf(&b, "| %s | %d | %d | %.1f%% |\n",
				rc.Repo, rc.Total, rc.Failures, rc.Rate)
		}
		b.WriteString("\n")
	}

	return b.String()
}

func formatDuration(businessDays float64) string {
	if businessDays == 0 {
		return "0 hours"
	}
	if businessDays < 1 {
		hours := businessDays * 24
		return fmt.Sprintf("%.1f hours", hours)
	}
	return fmt.Sprintf("%.1f business days", businessDays)
}
