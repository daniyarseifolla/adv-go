package analyzer

import (
	"strings"

	"parallel-text-analyzer/internal/model"
)

type LineCountAnalyzer struct{}

func (a *LineCountAnalyzer) Name() string {
	return "LineCount"
}

func (a *LineCountAnalyzer) Analyze(content string) model.AnalysisResult {
	lines := strings.Count(content, "\n")
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		lines++
	}
	return model.AnalysisResult{
		Name: a.Name(),
		Data: map[string]interface{}{
			"lines": lines,
		},
	}
}
