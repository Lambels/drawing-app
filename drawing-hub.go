package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type DrawingHub struct {
	idCount int64

	colors map[string]int64
	users  map[int64]*User
	read   chan Message
}

func (h *DrawingHub) Open() {
	log.Println("Hub started listening")
	go h.listen()
}

func (h *DrawingHub) Close() {
	h.read <- newCloseMessage()
}

func (h *DrawingHub) listen() {
	for {
		select {
		case msg := <-h.read:
			switch msg.Kind() {
			case TypeUserJoined:
				h.idCount++
				log.Printf("User: %d joined", h.idCount)

				// generate unique color.
				color := genRandomColorHex()
				_, ok := h.colors[color]
				for ok {
					color = genRandomColorHex()
					_, ok = h.colors[color]
				}

				// populate user.
				pMsg := msg.(*MessageUserJoined)
				user := pMsg.User
				user.Id = h.idCount
				user.Color = color
				user.commSend = h.read

				// add user to hub.
				h.users[user.Id] = user
				h.colors[color] = user.Id

				// notify client side of their color and id.
				msg := newFirstMessage(user.Id, user.Color)
				payload, _ := json.Marshal(msg)
				user.write(payload)

				// start user.
				go user.Start()

			case TypeUserLeft:
				pMsg := msg.(*MessageUserLeft)
				userId := pMsg.Id
				user := h.users[userId]
				log.Printf("User: %d is leaving\n", userId)

				delete(h.users, userId)
				delete(h.colors, user.Color)

			case TypeClose:
				log.Println("Hub is closing")
				h.closeAll()
				return

			default:
				h.broadcast(msg, msg.SenderId())
			}
		}
	}
}

func (h *DrawingHub) broadcast(msg Message, sender int64) {
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
