package websocket

type Message struct {
	Method string      `json:"method"`
	FormId string      `json:"formId"`
	Data   interface{} `json:"data"`
}
