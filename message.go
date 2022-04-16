package main

const (
	TypeUserJoined = iota + 1
	TypeUserLeft
	TypeClear
	TypeDraw
)

type Message struct {
	Type int                         `json:"type"`
	Data map[interface{}]interface{} `json:"data"`
}

func newLeaveMessage(id int64) *Message {
	msj := Message{
		Type: TypeUserLeft,
		Data: make(map[interface{}]interface{}),
	}
	msj.Data[TypeUserLeft] = struct{}{}
	return &msj
}
