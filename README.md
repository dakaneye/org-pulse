# org-pulse

[![CI](https://github.com/dakaneye/org-pulse/actions/workflows/ci.yml/badge.svg)](https://github.com/dakaneye/org-pulse/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dakaneye/org-pulse)](https://goreportcard.com/report/github.com/dakaneye/org-pulse)
[![Go Reference](https://pkg.go.dev/badge/github.com/dakaneye/org-pulse.svg)](https://pkg.go.dev/github.com/dakaneye/org-pulse)
[![Release](https://img.shields.io/github/v/release/dakaneye/org-pulse)](https://github.com/dakaneye/org-pulse/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

CLI tool that generates weekly delivery health reports for GitHub organizations. No SaaS required.

## Metrics

- **PR velocity** -- opened, merged, closed counts
- **Median time-to-first-review** -- in business days
- **Median time-to-merge** -- in business days
- **Review rounds distribution** -- how many review cycles before merge
- **Review load** -- who's carrying review burden
- **Stale PRs** -- no activity, no review response, approved but not merged
- **CI failure rate** -- per-repo and org-wide (GitHub Actions)

## Install

```bash
# Go install
go install github.com/dakaneye/org-pulse@latest

# Or download a release binary
# https://github.com/dakaneye/org-pulse/releases
```

## Prerequisites

- [gh CLI](https://cli.github.com/) installed and authenticated (`gh auth login`)

## Usage

```bash
# Full org report (last 4 weeks, markdown)
org-pulse report --org myorg

# JSON output
org-pulse report --org myorg --format json

# Custom time window
org-pulse report --org myorg --weeks 2

# Filter to specific repos
org-pulse report --org myorg --repos repo-a,repo-b

# Exclude repos
org-pulse report --org myorg --exclude-repos archived-thing

# Tune concurrency
org-pulse report --org myorg --concurrency 20
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--org` | (required) | GitHub organization name |
| `--weeks` | `4` | Lookback window in weeks |
| `--format` | `md` | Output format: `md` or `json` |
| `--repos` | (none) | Only include these repos (comma-separated) |
| `--exclude-repos` | (none) | Exclude these repos (comma-separated) |
| `--concurrency` | `10` | Max parallel `gh` API calls |

### Output

Markdown output goes to stdout. Pipe it wherever you want:

```bash
# Save to file
org-pulse report --org myorg > report.md

# Copy to clipboard (macOS)
org-pulse report --org myorg | pbcopy
```

## How It Works

- Uses `gh api graphql` to fetch data (inherits your existing `gh` auth)
- Fetches repos in the org, then PRs per repo in parallel
- Skips archived repos, forks, and draft PRs by default
- All time metrics use business days (Mon-Fri, no holidays)
- Fails fast if any repo fetch fails (no partial reports)

## License

[MIT](LICENSE)
