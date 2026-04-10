package analyzer

import (
	"strings"

	"parallel-text-analyzer/internal/model"
)

type WordCountAnalyzer struct{}

func (a *WordCountAnalyzer) Name() string {
	return "WordCount"
}

func (a *WordCountAnalyzer) Analyze(content string) model.AnalysisResult {
	words := len(strings.Fields(content))
	return model.AnalysisResult{
		Name: a.Name(),
		Data: map[string]interface{}{
			"words": words,
		},
	}
}
