package models

type Note struct {
	ID      int    `json:"id,omitempty"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
