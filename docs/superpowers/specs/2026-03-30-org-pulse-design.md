# org-pulse Design Spec

## Overview

Go CLI that generates weekly delivery health reports for a GitHub organization. Pulls PR velocity, review load balance, stale PRs, CI failure rate, and review burden distribution. Outputs markdown or JSON to stdout. No SaaS dependency.

Primary user: engineering leaders who want org-level visibility without commercial tooling.

## CLI Interface

```
org-pulse report --org <name> [--weeks 4] [--format md|json] [--repos repo1,repo2] [--exclude-repos repo1,repo2] [--concurrency 10]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--org` | (required) | GitHub organization name |
| `--weeks` | `4` | Lookback window in weeks |
| `--format` | `md` | Output format: `md` or `json` |
| `--repos` | (none) | Positive filter: only include these repos |
| `--exclude-repos` | (none) | Negative filter: exclude these repos |
| `--concurrency` | `10` | Max parallel `gh` subprocess calls |

- `--repos` and `--exclude-repos` are mutually exclusive.
- Archived and forked repos are excluded by default.
- Draft PRs are excluded from metrics.
- Output goes to stdout.

## Project Structure

```
org-pulse/
├── main.go                  # cobra setup, flag parsing, run
├── internal/
│   ├── github/              # gh CLI subprocess calls, GraphQL queries, data types
│   ├── metrics/             # pure computation functions
│   └── report/              # markdown and JSON rendering
├── go.mod
└── go.sum
```

Single binary. `cobra` for CLI framework. Distributed via goreleaser + Homebrew.

## Data Fetching (internal/github)

### Auth

Shells out to `gh api graphql`. Inherits user's existing `gh` auth. Validated on startup via `gh auth status`.

### Queries

1. **List repos** — paginated GraphQL query for all org repos. Client-side filter: remove archived, forks, then apply `--repos`/`--exclude-repos`.

2. **Per-repo PR data** — for each repo, fetch PRs updated within the `--weeks` window:
   - PR metadata: author, created/merged/closed timestamps, state, number, title
   - Reviews: reviewer, submitted timestamp, state (approved/changes_requested/commented)
   - Review requests: requested reviewers and when requested
   - Commits: last commit's check suites/runs (CI status)

### Concurrency

```
list repos (serial, single paginated query)
    → errgroup: fetch PR data per repo (bounded, default 10 workers)
        → collect into []RepoData
```

Bounded parallelism via `errgroup`. `gh` rate limiting respected implicitly by bounding workers.

### Data Types

Plain Go structs: `Repo`, `PullRequest`, `Review`, `ReviewRequest`, `CheckRun`. No interfaces, no generics. Data bags passed to metrics.

## Metrics Computation (internal/metrics)

All functions are pure: data in, number out. No side effects. `time.Now()` is never called — the current timestamp is passed as a parameter for deterministic testing.

### Org-Level

| Metric | Description |
|--------|-------------|
| PR velocity | Total opened, merged, closed in window |
| Median time-to-first-review | PR open → first review submission (business days) |
| Median time-to-merge | PR open → merge (business days) |
| Review rounds distribution | Histogram of changes_requested → re-review cycles before merge |
| CI failure rate | Failed GHA runs / total GHA runs (org aggregate) |

### Per-Repo

| Metric | Description |
|--------|-------------|
| PR velocity | Opened, merged, closed per repo |
| CI failure rate | Failed / total GHA runs per repo |

### Per-Person

| Metric | Description |
|--------|-------------|
| Review load | Number of reviews submitted per person |
| Review burden | Percentage of total org reviews per person |

### Stale PRs

A PR is stale if it matches any of these conditions (all thresholds: 7 business days):

| Category | Condition |
|----------|-----------|
| No review activity | Open with no review activity for N business days |
| No review response | Review requested, no response for N business days |
| Approved, not merged | Approved but not merged for N business days |

Each stale PR is tagged with which category it matched.

### Business Day Calculation

Dedicated, isolated pure function. Counts weekdays (Mon–Fri) between two `time.Time` values. No holiday awareness.

**Testing strategy for metrics:**
- Every metric is its own function with its own tests
- Business day function: exhaustive table-driven tests (same-day, cross-weekend, multi-week, Friday-to-Monday edge cases)
- Median: even count, odd count, single item, empty slice
- Stale PR detection: fixtures covering every category and boundary conditions (exactly N days, N-1, N+1)
- Review rounds: single approval, multiple rounds, comment-only reviews
- "Now" is always injected so tests are fully deterministic

## Report Rendering (internal/report)

### Markdown

```markdown
# Org Pulse: <org-name>
> 4-week report (Mar 3 – Mar 30, 2026) | 42 repos

## Summary
- PRs opened: 187 | merged: 162 | closed: 12
- Median time-to-first-review: 3.2 business days
- Median time-to-merge: 5.1 business days
- CI failure rate: 14.2%

## Review Load
| Reviewer       | Reviews | % of Total |
|----------------|---------|------------|
| @alice         | 47      | 29%        |
| @bob           | 31      | 19%        |

## Review Rounds Distribution
| Rounds | Count | % |
|--------|-------|---|
| 1      | 98    | 60% |
| 2      | 42    | 26% |
| 3+     | 22    | 14% |

## Stale PRs (14)
| PR | Repo | Age (bd) | Reason |
|----|------|----------|--------|
| #123 | repo-a | 12 | No review activity |
| #456 | repo-b | 8 | Approved, not merged |

## CI Failure Rate by Repo
| Repo   | Runs | Failures | Rate |
|--------|------|----------|------|
| repo-a | 320  | 58       | 18%  |
| repo-b | 210  | 12       | 6%   |
```

### JSON

Same data structure as the `Report` struct, marshaled directly. No special formatting logic.

### Rendering Principle

Rendering is dumb — takes a single `Report` struct and formats it. All logic lives in metrics, not in report.

## Error Handling

| Scenario | Behavior |
|----------|----------|
| `gh` not installed or not authenticated | Check on startup with `gh auth status`. Clear error message, exit 1. |
| Empty org / no repos match filter | Print "no repos found", exit 0. |
| Rate limiting | `gh` handles retry/backoff internally. |
| Single repo fetch failure | Surface error with context (which repo, what query). Fail the entire report — no partial results. |
| Repos with zero PRs in window | Included in repo count, omitted from per-repo tables. |

## PR Inclusion Rules

| Scenario | Included? |
|----------|-----------|
| Draft PRs | No |
| PRs opened before window but merged within it | Yes |
| Bot-authored PRs (dependabot, renovate) | Yes |
| PRs in archived repos | No |
| PRs in forked repos | No |

## Distribution

Single binary via goreleaser + Homebrew tap. Standard `go install` also supported.
