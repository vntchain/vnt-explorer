package common

type Test struct {
	Name string `json:"name"`
}

type Rpc struct {
	Jsonrpc	string			`json:"jsonrpc"`
	Method	string			`json:"method"`
	Params	[]interface{}	`json:"params"`
	Id		int				`json:"id"`
}

type Error struct {
	Code 	int		`json:"code"`
	Message	string	`json:"message"`
}

type Response struct {
	Jsonrpc	string			`json:"jsonrpc"`
	Id		int				`json:"id"`
	Result	interface{}		`json:result`
	Error	*Error			`json: error`
}

func NewRpc() *Rpc {
	return &Rpc{
		Jsonrpc: "2.0",
		Id: 1,
		Params: make([]interface{}, 0),
	}
}