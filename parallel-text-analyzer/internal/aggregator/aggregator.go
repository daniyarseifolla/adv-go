package aggregator

import (
	"sort"
	"sync"

	"parallel-text-analyzer/internal/model"
)

type Aggregator struct {
	mu        sync.Mutex
	globalFreq map[string]int
}

func New() *Aggregator {
	return &Aggregator{
		globalFreq: make(map[string]int),
	}
}

func (a *Aggregator) Merge(fs model.FileStats) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for word, count := range fs.WordFreq {
		a.globalFreq[word] += count
	}
}

func (a *Aggregator) TopWords(n int) []model.WordCount {
	a.mu.Lock()
	defer a.mu.Unlock()

	all := make([]model.WordCount, 0, len(a.globalFreq))
	for word, count := range a.globalFreq {
		all = append(all, model.WordCount{Word: word, Count: count})
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Count > all[j].Count
	})

	if n > len(all) {
		n = len(all)
	}
	return all[:n]
}
