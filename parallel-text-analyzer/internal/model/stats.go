package model

type FileStats struct {
	FileName string
	Words    int
	Lines    int
	Chars    int
	WordFreq map[string]int
	Results  []AnalysisResult
}

type AnalysisResult struct {
	Name string
	Data map[string]interface{}
}

type WordCount struct {
	Word  string
	Count int
}
