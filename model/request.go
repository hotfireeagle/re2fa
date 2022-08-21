package model

type GenerateFAPostData struct {
	RegExp string `json:"regexp" validate:"required"`
}

type FAMatchPostData struct {
	RegExp string `json:"regexp" validate:"required"`
	Text   string `json:"text"`
	Api    string `json:"api" validate:"required"`
}
