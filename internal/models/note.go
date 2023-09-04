package models

type Note struct {
	ID      int    `json:"-"`
	Title   string `json:"title" validate:"required" format:"string" example:"заголовок заметки"`
	Content string `json:"content" validate:"required" format:"string" example:"содержание заметки"`
}
