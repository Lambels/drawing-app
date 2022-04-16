package main

type DrawingHub struct {
	idCount int

	users map[int]*User
	read  <-chan *Message
}

func (h *DrawingHub) Start() {

}

func (h *DrawingHub) listen() {

}

func (h *DrawingHub) broadcast(msg *Message) {

}
