package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// User represents a connection.
type User struct {
	Id    int64  `json:"id"`
	Color string `json:"color"` // hex value.

	commSend chan<- Message  // communicate with hub.
	conn     *websocket.Conn // communicate with client.
}

// Start starts the reading pump in its own go-routine.
func (u *User) Start() {
	go u.readPump()
}

// readPump reads messages from the conn and pumps them to the hub to be broadcasted.
func (u *User) readPump() {
	defer func() {
		u.commSend <- newLeaveMessage(u.Id)
		u.conn.Close()
	}()

	for {
		var msg MessagePoint
		if err := u.conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("readPump: unexpected error: %v", err)
			}
			break
		}

		u.commSend <- &msg
	}
}

// write attempts to write msg to u.conn, if an error occurs the connection
// will get closed.
func (u *User) write(msg []byte) {
	if err := u.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		u.conn.Close()
	}
}
