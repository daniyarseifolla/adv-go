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

	"parallel-text-analyzer/internal/aggregator"
	"parallel-text-analyzer/internal/model"
	"parallel-text-analyzer/internal/pipeline"
)

func main() {
	path := flag.String("path", ".", "path to directory or file")
	ext := flag.String("ext", ".txt", "file extension to analyze")
	workers := flag.Int("workers", runtime.NumCPU(), "number of worker goroutines")
	topWords := flag.Int("top-words", 0, "show top N most frequent words")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	files, err := pipeline.WalkDir(*path, *ext)
	if err != nil {
		log.Fatalf("error walking path: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("no files found")
		return
	}

	filePaths := make(chan string, len(files))
	results := make(chan model.FileStats, len(files))

	var wg sync.WaitGroup
	pipeline.StartWorkers(ctx, filePaths, results, &wg, *workers)

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
		fmt.Printf("File: %s\n", stats.FileName)
		fmt.Printf("  Lines: %d  Words: %d  Chars: %d\n\n", stats.Lines, stats.Words, stats.Chars)
		agg.Merge(stats)
	}

	if *topWords > 0 {
		top := agg.TopWords(*topWords)
		fmt.Printf("Top %d words:\n", len(top))
		for i, wc := range top {
			fmt.Printf("  %d. %-15s — %d\n", i+1, wc.Word, wc.Count)
		}
	}
}
