package dsl

// ErrWrap is a wrapper for struct for storing syntax errors
type ErrWrap struct {
	Incorrect []Incorrect `json:"errors"`
}

// Incorrect is a struct for storing incorrect data
type Incorrect struct {
	LineNumber     int    `json:"lineNumber"`
	ErrDescription string `json:"errDescription"`
}

// QuizWrap is a wrapper for struct for storing quiz data
type QuizWrap struct {
	Quiz Quiz `json:"quiz"`
}

// Quiz is a struct for storing quiz data
type Quiz struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

// Block is a struct for storing block of lines
type Block struct {
	start int
	end   int
}
