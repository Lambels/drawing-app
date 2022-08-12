package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.uber.org/multierr"
)

type DrawingHub struct {
	idCount int64

	lastX   int
	dataBuf []Message
	colors  map[string]int64
	users   map[int64]*User
	read    chan Message
	close   chan error
}

func (h *DrawingHub) Open() {
	log.Println("Hub started listening")
	go h.listen()
}

func (h *DrawingHub) Close(ctx context.Context) error {
	close(h.read)

	select {
	case err := <-h.close:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *DrawingHub) listen() {
	defer func() {
		log.Println("Hub is closing")
		h.close <- h.closeAll()
	}()

	for msg := range h.read {
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
			msg := newFirstMessage(user.Id, user.Color, h.dataBuf)
			payload, _ := json.Marshal(msg)
			user.write(payload)

			// start user.
			user.Start()

		case TypeUserLeft:
			pMsg := msg.(*MessageUserLeft)
			userId := pMsg.Id
			user := h.users[userId]
			log.Printf("User: %d is leaving\n", userId)

			delete(h.users, userId)
			delete(h.colors, user.Color)

		case TypePointStart, TypePoint, TypePointEnd:
			log.Printf("Drawing: %d From user: %d", msg.Kind(), msg.SenderId())
			h.saveMessage(msg)
			h.broadcast(msg, msg.SenderId())
		}
	}
}

func (h *DrawingHub) saveMessage(msg Message) {
	if msg.Kind() != TypeDraw {
		return
	}

	if h.lastX == cap(h.dataBuf) {
		h.lastX = 0
	}
	h.dataBuf[h.lastX] = msg
	h.lastX++
}

func (h *DrawingHub) broadcast(msg Message, sender int64) {
	payload, _ := json.Marshal(msg)
	for _, user := range h.users {
		if user.Id != sender {
			user.write(payload)
		}
	}
}

func (h *DrawingHub) closeAll() error {
	var err error
	for _, user := range h.users {
		if err := user.conn.Close(); err != nil {
			multierr.AppendInto(&err, err)
		}
	}

	return err
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
