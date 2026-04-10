package pipeline

import (
	"bufio"
	"context"
	"os"
	"strings"
	"sync"

	"parallel-text-analyzer/internal/analyzer"
	"parallel-text-analyzer/internal/model"
)

func ProcessFile(path string) (model.FileStats, error) {
	f, err := os.Open(path)
	if err != nil {
		return model.FileStats{}, err
	}
	defer f.Close()

	var lines, words, chars int
	freq := make(map[string]int)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines++
		chars += len(line)
		for _, w := range strings.Fields(line) {
			words++
			freq[strings.ToLower(w)]++
		}
	}

	if err := scanner.Err(); err != nil {
		return model.FileStats{}, err
	}

	return model.FileStats{
		FileName: path,
		Words:    words,
		Lines:    lines,
		Chars:    chars,
		WordFreq: freq,
	}, nil
}

func runAnalyzers(content string, analyzers []analyzer.Analyzer) []model.AnalysisResult {
	results := make([]model.AnalysisResult, len(analyzers))

	var wg sync.WaitGroup
	for i, a := range analyzers {
		wg.Add(1)
		go func(idx int, an analyzer.Analyzer) {
			defer wg.Done()
			results[idx] = an.Analyze(content)
		}(i, a)
	}
	wg.Wait()

	return results
}

func StartWorkers(ctx context.Context, filePaths <-chan string, results chan<- model.FileStats, wg *sync.WaitGroup, numWorkers int, analyzers []analyzer.Analyzer) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range filePaths {
				select {
				case <-ctx.Done():
					return
				default:
				}

				stats, err := ProcessFile(path)
				if err != nil {
					continue
				}

				if len(analyzers) > 0 {
					content, err := os.ReadFile(path)
					if err == nil {
						stats.Results = runAnalyzers(string(content), analyzers)
					}
				}

				results <- stats
			}
		}()
	}
}
