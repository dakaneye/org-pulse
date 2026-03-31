package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// CheckAuth verifies the gh CLI is authenticated by running gh auth status.
func CheckAuth(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "gh", "auth", "status")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh auth check: %s: %w", strings.TrimSpace(stderr.String()), err)
	}
	return nil
}

func graphQL(ctx context.Context, query string, variables map[string]any) (json.RawMessage, error) {
	args := []string{"api", "graphql"}
	for k, v := range variables {
		args = append(args, "-F", fmt.Sprintf("%s=%v", k, v))
	}
	args = append(args, "-f", fmt.Sprintf("query=%s", query))

	cmd := exec.CommandContext(ctx, "gh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gh api graphql: %s: %w", stderr.String(), err)
	}

	return stdout.Bytes(), nil
}
