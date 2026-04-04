package main

import (
	"flag"
	"fmt"
	"log"

	"parallel-text-analyzer/internal/pipeline"
)

func main() {
	path := flag.String("path", ".", "path to directory or file")
	ext := flag.String("ext", ".txt", "file extension to analyze")
	_ = flag.Int("workers", 1, "number of worker goroutines")
	flag.Parse()

	files, err := pipeline.WalkDir(*path, *ext)
	if err != nil {
		log.Fatalf("error walking path: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("no files found")
		return
	}

	for _, f := range files {
		stats, err := pipeline.ProcessFile(f)
		if err != nil {
			log.Printf("error processing %s: %v", f, err)
			continue
		}
		fmt.Printf("File: %s\n", stats.FileName)
		fmt.Printf("  Lines: %d  Words: %d  Chars: %d\n\n", stats.Lines, stats.Words, stats.Chars)
	}
}
