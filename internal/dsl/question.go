package dsl

// Question is a struct for storing question data
type Question struct {
	Text           string   `json:"text"`
	Multiple       bool     `json:"multiple"`
	Options        []Answer `json:"options"`
	Recommendation string   `json:"recommendation"`
}
