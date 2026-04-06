package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"

	"parallel-text-analyzer/internal/model"
	"parallel-text-analyzer/internal/pipeline"
)

func main() {
	path := flag.String("path", ".", "path to directory or file")
	ext := flag.String("ext", ".txt", "file extension to analyze")
	workers := flag.Int("workers", runtime.NumCPU(), "number of worker goroutines")
	flag.Parse()

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
	pipeline.StartWorkers(filePaths, results, &wg, *workers)

	for _, f := range files {
		filePaths <- f
	}
	close(filePaths)

	go func() {
		wg.Wait()
		close(results)
	}()

	for stats := range results {
		fmt.Printf("File: %s\n", stats.FileName)
		fmt.Printf("  Lines: %d  Words: %d  Chars: %d\n\n", stats.Lines, stats.Words, stats.Chars)
	}
}
