package main

import (
	"fmt"

	"parallel-text-analyzer/internal/model"
)

func printFileStats(stats model.FileStats) {
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
}

func printTopWords(top []model.WordCount) {
	fmt.Printf("Top %d words (all files):\n", len(top))
	for i, wc := range top {
		fmt.Printf("  %d. %-15s — %d\n", i+1, wc.Word, wc.Count)
	}
}
