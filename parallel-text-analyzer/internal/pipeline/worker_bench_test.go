package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"parallel-text-analyzer/internal/analyzer"
	"parallel-text-analyzer/internal/model"
)

func createBenchFiles(b *testing.B) string {
	dir := b.TempDir()
	content := []byte("The quick brown fox jumps over the lazy dog. " +
		"A journey of a thousand miles begins with a single step. " +
		"To be or not to be, that is the question. " +
		"All that glitters is not gold.\n")

	for i := 0; i < 50; i++ {
		path := filepath.Join(dir, "file"+string(rune('A'+i%26))+string(rune('0'+i/26))+".txt")
		data := make([]byte, 0, len(content)*20)
		for j := 0; j < 20; j++ {
			data = append(data, content...)
		}
		os.WriteFile(path, data, 0644)
	}
	return dir
}

func BenchmarkSequential(b *testing.B) {
	dir := createBenchFiles(b)
	files, _ := WalkDir(dir, ".txt", 0, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, f := range files {
			ProcessFile(f)
		}
	}
}

func BenchmarkParallel(b *testing.B) {
	dir := createBenchFiles(b)
	files, _ := WalkDir(dir, ".txt", 0, 0)
	analyzers := []analyzer.Analyzer{
		&analyzer.WordCountAnalyzer{},
		&analyzer.LineCountAnalyzer{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filePaths := make(chan string, len(files))
		results := make(chan model.FileStats, len(files))

		var wg sync.WaitGroup
		StartWorkers(context.Background(), filePaths, results, &wg, 4, analyzers)

		for _, f := range files {
			filePaths <- f
		}
		close(filePaths)

		go func() {
			wg.Wait()
			close(results)
		}()

		for range results {
		}
	}
}
