package metrics

import (
	"testing"

	"github.com/dakaneye/org-pulse/internal/github"
)

func TestCIFailureRate(t *testing.T) {
	prs := []github.PullRequest{
		{
			Number: 1,
			CheckRuns: []github.CheckRun{
				{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"},
				{Name: "lint", Status: "COMPLETED", Conclusion: "FAILURE"},
			},
		},
		{
			Number: 2,
			CheckRuns: []github.CheckRun{
				{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"},
				{Name: "lint", Status: "COMPLETED", Conclusion: "SUCCESS"},
			},
		},
	}

	rate := CIFailureRate(prs)

	if rate.Total != 4 {
		t.Errorf("Total = %d, want 4", rate.Total)
	}
	if rate.Failures != 1 {
		t.Errorf("Failures = %d, want 1", rate.Failures)
	}
	if rate.Rate != 25.0 {
		t.Errorf("Rate = %v, want 25.0", rate.Rate)
	}
}

func TestCIFailureRateSkipsNonCompleted(t *testing.T) {
	prs := []github.PullRequest{
		{
			Number: 1,
			CheckRuns: []github.CheckRun{
				{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"},
				{Name: "build", Status: "IN_PROGRESS", Conclusion: ""},
				{Name: "lint", Status: "COMPLETED", Conclusion: "SKIPPED"},
			},
		},
	}

	rate := CIFailureRate(prs)

	if rate.Total != 2 {
		t.Errorf("Total = %d, want 2", rate.Total)
	}
	if rate.Failures != 0 {
		t.Errorf("Failures = %d, want 0", rate.Failures)
	}
}

func TestCIFailureRateEmpty(t *testing.T) {
	rate := CIFailureRate(nil)
	if rate.Total != 0 || rate.Failures != 0 || rate.Rate != 0 {
		t.Errorf("expected all zeros, got %+v", rate)
	}
}
