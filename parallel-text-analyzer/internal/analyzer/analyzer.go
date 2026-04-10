package analyzer

import "parallel-text-analyzer/internal/model"

type Analyzer interface {
	Analyze(content string) model.AnalysisResult
	Name() string
}
