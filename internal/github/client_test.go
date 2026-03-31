package github

import (
	"context"
	"testing"
)

func TestCheckAuth(t *testing.T) {
	err := CheckAuth(context.Background())
	if err != nil {
		t.Skipf("gh not authenticated in test environment: %v", err)
	}
}
