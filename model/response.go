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

type Edge struct {
	From  int    `json:"from"`
	To    int    `json:"to"`
	Label string `json:"label"`
}

type DrawNFAResponse struct {
	Nodes       []int  `json:"nodes"`
	Edges       []Edge `json:"edges"`
	StartState  int
	AcceptState int
}