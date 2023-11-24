package types

type SigValue struct {
	Scaddress string `json:"scaddress"`
	Sender    string `json:"sender"`
	Qty       uint32 `json:"qty"`
	ChainId   uint32 `json:"chainId"`
}
type ErrorResponse struct {
	Error   string      `json:"error"`
	Success bool        `json:"success" default:"false"`
	Data    interface{} `json:"data"`
}

type ResResponse struct {
	Error   string      `json:"error"`
	Success bool        `json:"success" default:"false"`
	Data    interface{} `json:"data"`
}
