package pipeline

import "parallel-text-analyzer/internal/model"

func FilterResults(in <-chan model.FileStats, minWords int) <-chan model.FileStats {
	out := make(chan model.FileStats)

	go func() {
		defer close(out)
		for stats := range in {
			if stats.Words >= minWords {
				out <- stats
			}
		}
	}()

	return out
}
