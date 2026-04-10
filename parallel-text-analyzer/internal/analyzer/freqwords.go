package analyzer

import (
	"sort"
	"strings"

	"parallel-text-analyzer/internal/model"
)

type MostFrequentWordsAnalyzer struct {
	TopN int
}

func (a *MostFrequentWordsAnalyzer) Name() string {
	return "FreqWords"
}

func (a *MostFrequentWordsAnalyzer) Analyze(content string) model.AnalysisResult {
	freq := make(map[string]int)
	for _, w := range strings.Fields(content) {
		freq[strings.ToLower(w)]++
	}

	type pair struct {
		word  string
		count int
	}

	pairs := make([]pair, 0, len(freq))
	for w, c := range freq {
		pairs = append(pairs, pair{w, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	n := a.TopN
	if n > len(pairs) {
		n = len(pairs)
	}

	top := make([]model.WordCount, n)
	for i := 0; i < n; i++ {
		top[i] = model.WordCount{Word: pairs[i].word, Count: pairs[i].count}
	}

	return model.AnalysisResult{
		Name: a.Name(),
		Data: map[string]interface{}{
			"top_words": top,
		},
	}
}
