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
}

type FAItem struct {
	FA    *DrawFAResponse `json:"fa"`
	Title string          `json:"title"`
}

func NewDoubleFA(fa1 *DrawFAResponse, t1 string, fa2 *DrawFAResponse, t2 string) interface{} {
	result := make([]*FAItem, 0)
	result = append(result, &FAItem{fa1, t1})
	result = append(result, &FAItem{fa2, t2})
	return result
}

type ApiListItem struct {
	Name string `json:"name"`
	Api  string `json:"api"`
}
