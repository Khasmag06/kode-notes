package models

type Note struct {
	ID      int    `json:"-"`
	Title   string `json:"title" format:"string" example:"заголовок заметки"`
	Content string `json:"content" format:"string" example:"содержание заметки"`
}
