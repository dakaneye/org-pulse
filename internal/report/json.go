package report

import (
	"encoding/json"

	"github.com/dakaneye/org-pulse/internal/metrics"
)

func RenderJSON(r metrics.Report) (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
