package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type DrawingHub struct {
	idCount int64

	colors map[string]int64
	users  map[int64]*User
	read   chan *Message
}

func (h *DrawingHub) Open() {
	go h.listen()
}

func (h *DrawingHub) Close() {
	h.read <- &Message{Type: TypeClose}
}

func (h *DrawingHub) listen() {
	for {
		select {
		case msg := <-h.read:
			switch msg.Type {
			case TypeUserJoined:
				h.idCount++

				// generate unique color.
				color := genRandomColorHex()
				for _, ok := h.colors[color]; !ok; {
					color = genRandomColorHex()
				}

				// populate user.
				user := msg.Data[TypeUserJoined].(*User)
				user.Id = h.idCount
				user.Color = color
				user.commSend = h.read

				// add user to hub.
				h.users[user.Id] = user
				h.colors[color] = user.Id

				// start user.
				go user.Start()

			case TypeUserLeft:
				userId := msg.Data[TypeUserLeft].(int64)
				user := h.users[userId]

				delete(h.users, userId)
				delete(h.colors, user.Color)

			case TypeClose:
				h.closeAll()
				return

			default:
				h.broadcast(msg, msg.Data["sender"].(int64))
			}
		}
	}
}

func (h *DrawingHub) broadcast(msg *Message, sender int64) {
	payload, _ := json.Marshal(msg)
	for _, user := range h.users {
		if user.Id != sender {
			user.write(payload)
		}
	}
}

func (h *DrawingHub) closeAll() {
	for _, user := range h.users {
		user.conn.Close()
	}
}

func genRandomColorHex() string {
	rand.Seed(time.Now().Unix())
	r := rand.Intn(255)
	g := rand.Intn(255)
	b := rand.Intn(255)

	hexR := fmt.Sprintf("%x", r)
	if len(hexR) == 1 {
		hexR = "0" + hexR
	}

	hexG := fmt.Sprintf("%x", g)
	if len(hexG) == 1 {
		hexG = "0" + hexG
	}

	hexB := fmt.Sprintf("%x", b)
	if len(hexB) == 1 {
		hexB = "0" + hexB
	}

	return "#" + hexR + hexG + hexB
}
