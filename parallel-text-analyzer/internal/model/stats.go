package model

type FileStats struct {
	FileName string
	Words    int
	Lines    int
	Chars    int
	WordFreq map[string]int
}

type WordCount struct {
	Word  string
	Count int
}
