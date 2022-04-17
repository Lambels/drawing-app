package main

const (
	TypeUserJoined = iota + 1
	TypeUserLeft
	TypeClear
	TypeDraw
	TypeClose
)

type Message struct {
	Type int                         `json:"type"`
	Data map[interface{}]interface{} `json:"data"`
}

func newJoinMessage(user *User) *Message {
	msj := Message{
		Type: TypeUserJoined,
		Data: make(map[interface{}]interface{}),
	}
	msj.Data[TypeUserJoined] = user
	return &msj
}

func newLeaveMessage(id int64) *Message {
	msj := Message{
		Type: TypeUserLeft,
		Data: make(map[interface{}]interface{}),
	}
	msj.Data[TypeUserLeft] = id
	return &msj
}
