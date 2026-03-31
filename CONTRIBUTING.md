# Contributing

## Development

```bash
git clone https://github.com/dakaneye/org-pulse.git
cd org-pulse
go build ./...
go test ./...
```

## Before Submitting

1. `go build ./...` passes
2. `go test -race ./...` passes
3. `go vet ./...` clean
4. New functionality has tests

## Pull Requests

- Keep changes focused
- Update tests for new functionality
- Follow existing code style
