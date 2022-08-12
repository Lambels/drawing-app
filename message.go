package main

type MessageType int

const (
	TypeUserJoined MessageType = iota + 1
	TypeUserLeft
	TypeUserFirst
	TypeClear
	TypeDraw
	TypeClose
	TypePoint
	TypePointStart
	TypePointEnd
)

type Message interface {
	Kind() MessageType
	SenderId() int64
}

type MessageUserJoined struct {
	Type MessageType `json:"type"`
	User *User       `json:"user"`
}

func (m *MessageUserJoined) Kind() MessageType {
	return m.Type
}

func (m *MessageUserJoined) SenderId() int64 {
	return -1
}

func newJoinMessage(user *User) Message {
	return &MessageUserJoined{
		Type: TypeUserJoined,
		User: user,
	}
}

type MessageUserLeft struct {
	Type MessageType `json:"type"`
	Id   int64       `json:"id"`
}

func (m *MessageUserLeft) Kind() MessageType {
	return m.Type
}

func (m *MessageUserLeft) SenderId() int64 {
	return -1
}

func newLeaveMessage(id int64) Message {
	return &MessageUserLeft{
		Type: TypeUserLeft,
		Id:   id,
	}
}

type MessageUserFirst struct {
	Type       MessageType `json:"type"`
	Id         int64       `json:"id"`
	Color      string      `json:"color"`
	PastPoints []Message   `json:"pastPoints"`
}

func (m *MessageUserFirst) Kind() MessageType {
	return m.Type
}

func (m *MessageUserFirst) SenderId() int64 {
	return m.Id
}

func newFirstMessage(id int64, color string, pastPoints []Message) Message {
	return &MessageUserFirst{
		Type:       TypeUserFirst,
		Id:         id,
		Color:      color,
		PastPoints: pastPoints,
	}
}

type MessagePoint struct {
	Type  MessageType `json:"type"`
	Id    int64       `json:"id"`
	Color string      `json:"color"`
	X     int         `json:"x"`
	Y     int         `json:"y"`
}

func (m *MessagePoint) Kind() MessageType {
	return m.Type
}

func (m *MessagePoint) SenderId() int64 {
	return m.Id
}
