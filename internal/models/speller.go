package models

type SpellError struct {
	Word string   `json:"word"`
	S    []string `json:"s"`
}
