package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"parallel-text-analyzer/internal/aggregator"
	"parallel-text-analyzer/internal/analyzer"
	"parallel-text-analyzer/internal/model"
	"parallel-text-analyzer/internal/pipeline"
)

func main() {
	path := flag.String("path", ".", "path to directory or file")
	ext := flag.String("ext", ".txt", "file extension to analyze")
	workers := flag.Int("workers", runtime.NumCPU(), "number of worker goroutines")
	topWords := flag.Int("top-words", 0, "show top N most frequent words")
	minSize := flag.Int64("min-size", 0, "minimum file size in bytes")
	maxSize := flag.Int64("max-size", 0, "maximum file size in bytes")
	cpuProf := flag.String("cpuprofile", "", "write cpu profile to file")
	memProf := flag.String("memprofile", "", "write memory profile to file")
	flag.Parse()

	stopProfile := startCPUProfile(*cpuProf)
	defer stopProfile()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	files, err := pipeline.WalkDir(*path, *ext, *minSize, *maxSize)
	if err != nil {
		log.Fatalf("error walking path: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("no files found")
		return
	}

	n := *topWords
	if n == 0 {
		n = 5
	}

	analyzers := []analyzer.Analyzer{
		&analyzer.WordCountAnalyzer{},
		&analyzer.LineCountAnalyzer{},
		&analyzer.MostFrequentWordsAnalyzer{TopN: n},
	}

	fmt.Printf("Starting analysis: %d files, %d workers\n\n", len(files), *workers)
	start := time.Now()

	filePaths := make(chan string, len(files))
	results := make(chan model.FileStats, len(files))

	var wg sync.WaitGroup
	pipeline.StartWorkers(ctx, filePaths, results, &wg, *workers, analyzers)

	for _, f := range files {
		filePaths <- f
	}
	close(filePaths)

	go func() {
		wg.Wait()
		close(results)
	}()

	agg := aggregator.New()

	for stats := range results {
		printFileStats(stats)
		agg.Merge(stats)
	}

	if *topWords > 0 {
		printTopWords(agg.TopWords(*topWords))
	}

	fmt.Printf("\nDone in %s\n", time.Since(start))

	writeMemProfile(*memProf)
}
