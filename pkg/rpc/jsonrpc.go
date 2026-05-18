package rpc

import "encoding/json"

type JSONRPCRequest struct {
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
	ID     json.RawMessage   `json:"id"`
}

type JSONRPCResponse struct {
	Result any             `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
	ID     json.RawMessage `json:"id"`
}
