package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/multierr"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	server *http.Server
	hub    *DrawingHub
}

func NewServer(port string, bufCap int) *Server {
	return &Server{
		server: &http.Server{
			Addr:    port,
			Handler: nil,
		},
		hub: &DrawingHub{
			close:   make(chan error, 1),
			dataBuf: make([]Message, bufCap),
			colors:  make(map[string]int64),
			users:   make(map[int64]*User),
			read:    make(chan Message),
		},
	}
}

func (s *Server) Open() {
	s.hub.Open()

	http.HandleFunc("/ws", s.handleWebsocket)

	log.Printf("Starting server at %s\n", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Close(ctx context.Context) error {
	log.Println("Server and Hub shutting down gracefully")

	return multierr.Combine(
		s.hub.Close(ctx),
		s.server.Shutdown(ctx),
	)
}

func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "couldnt upgrade connection", http.StatusInternalServerError)
		return
	}

	user := &User{
		conn: socket,
	}

	s.hub.read <- newJoinMessage(user)
}
