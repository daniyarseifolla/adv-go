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
	"parallel-text-analyzer/internal/analyzer"
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

	n := *topWords
	if n == 0 {
		n = 5
	}

	analyzers := []analyzer.Analyzer{
		&analyzer.WordCountAnalyzer{},
		&analyzer.LineCountAnalyzer{},
		&analyzer.MostFrequentWordsAnalyzer{TopN: n},
	}

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
		fmt.Printf("File: %s\n", stats.FileName)
		for _, r := range stats.Results {
			switch r.Name {
			case "WordCount":
				fmt.Printf("  [%s] Words: %v\n", r.Name, r.Data["words"])
			case "LineCount":
				fmt.Printf("  [%s] Lines: %v\n", r.Name, r.Data["lines"])
			case "FreqWords":
				if top, ok := r.Data["top_words"].([]model.WordCount); ok {
					fmt.Printf("  [%s] Top: ", r.Name)
					for i, wc := range top {
						if i > 0 {
							fmt.Print(", ")
						}
						fmt.Printf("%s(%d)", wc.Word, wc.Count)
					}
					fmt.Println()
				}
			}
		}
		fmt.Println()
		agg.Merge(stats)
	}

	if *topWords > 0 {
		top := agg.TopWords(*topWords)
		fmt.Printf("Top %d words (all files):\n", len(top))
		for i, wc := range top {
			fmt.Printf("  %d. %-15s — %d\n", i+1, wc.Word, wc.Count)
		}
	}
}
