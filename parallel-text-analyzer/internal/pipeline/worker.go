package pipeline

import (
	"bufio"
	"os"
	"strings"

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
