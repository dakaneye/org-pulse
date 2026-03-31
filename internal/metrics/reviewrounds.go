package metrics

import "github.com/dakaneye/org-pulse/internal/github"

// RoundsBucket holds the count and percentage for a review rounds bucket.
type RoundsBucket struct {
	Label   string
	Count   int
	Percent float64
}

// ReviewRoundsDistribution returns the distribution of review rounds across merged PRs.
// Review rounds are counted as 1 (base) plus the number of CHANGES_REQUESTED reviews.
func ReviewRoundsDistribution(prs []github.PullRequest) []RoundsBucket {
	counts := map[int]int{}
	total := 0

	for _, pr := range prs {
		if pr.MergedAt == nil {
			continue
		}

		rounds := countRounds(pr.Reviews)
		counts[rounds]++
		total++
	}

	if total == 0 {
		return nil
	}

	var buckets []RoundsBucket
	for _, def := range []struct {
		label string
		match func(int) bool
	}{
		{"1", func(r int) bool { return r == 1 }},
		{"2", func(r int) bool { return r == 2 }},
		{"3+", func(r int) bool { return r >= 3 }},
	} {
		count := 0
		for rounds, c := range counts {
			if def.match(rounds) {
				count += c
			}
		}
		if count > 0 {
			buckets = append(buckets, RoundsBucket{
				Label:   def.label,
				Count:   count,
				Percent: float64(count) / float64(total) * 100,
			})
		}
	}

	return buckets
}

func countRounds(reviews []github.Review) int {
	rounds := 1
	for _, r := range reviews {
		if r.State == "CHANGES_REQUESTED" {
			rounds++
		}
	}
	return rounds
}
