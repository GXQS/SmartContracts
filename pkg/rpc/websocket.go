package rpc

type WSMessage struct {
	Topic string `json:"topic"`
	Data  []byte `json:"data"`
}
