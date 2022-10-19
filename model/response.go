package model

type Response struct {
	//not full response schema
	Query Query `json:"query"`
}

type Query struct {
	Pages []Pages
}
type Pages struct {
	Title     string      `json:"title"`
	Missing   bool        `json:"missing"`
	Revisions []Revisions `json:"revisions"`
}

type Revisions struct {
	ContentModel  string `json:"contentmodel"`
	ContentFormat string `json:"contentformat"`
	Content       string `json:"content"`
}
