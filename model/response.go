package model

type ResponseCode byte

const (
	Error ResponseCode = iota
	Success
)

type Response struct {
	Code     ResponseCode `json:"code"`
	Msg      string       `json:"msg"`
	ErrorLog string       `json:"errorLog"`
	Data     interface{}  `json:"data"`
}
