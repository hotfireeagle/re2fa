package model

type GenerateFAPostData struct {
	RegExp string `json:"regexp" validate:"required"`
}
