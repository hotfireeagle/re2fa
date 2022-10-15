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

type DrawFAResponse struct {
	Edges        []*Edge `json:"edges"`
	Nodes        []int   `json:"nodes"`
	StartState   int     `json:"startState"`
	AcceptStates []int   `json:"acceptStates"`
	DeadState    int     `json:"deadState"`
	Title        string  `json:"title"`
}

type FAItem struct {
	FA    FA     `json:"fa"`
	Title string `json:"title"`
}

type ApiListItem struct {
	Name string `json:"name"`
	Api  string `json:"api"`
}

type FA interface {
	Match(s string) bool
	ConvertToJSON() *DrawFAResponse
}

type MatchAnswer struct {
	Result bool    `json:"result"`
	Time   float64 `json:"time"`
}
