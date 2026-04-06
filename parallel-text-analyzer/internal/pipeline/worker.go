package pipeline

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"parallel-text-analyzer/internal/model"
)

func ProcessFile(path string) (model.FileStats, error) {
	f, err := os.Open(path)
	if err != nil {
		return model.FileStats{}, err
	}
	defer f.Close()

	var lines, words, chars int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines++
		chars += len(line)
		words += len(strings.Fields(line))
	}

	if err := scanner.Err(); err != nil {
		return model.FileStats{}, err
	}

	return model.FileStats{
		FileName: path,
		Words:    words,
		Lines:    lines,
		Chars:    chars,
	}, nil
}

func StartWorkers(filePaths <-chan string, results chan<- model.FileStats, wg *sync.WaitGroup, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range filePaths {
				stats, err := ProcessFile(path)
				if err != nil {
					continue
				}
				results <- stats
			}
		}()
	}
}
