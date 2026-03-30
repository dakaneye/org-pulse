package github

import (
	"testing"
)

func TestCheckAuthParseError(t *testing.T) {
	err := CheckAuth()
	if err != nil {
		t.Skipf("gh not authenticated in test environment: %v", err)
	}
}
