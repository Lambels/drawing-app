package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	server *http.Server
	hub    *DrawingHub
}

func NewServer(port string) *Server {
	return &Server{
		server: &http.Server{
			Addr:    port,
			Handler: nil,
		},
		hub: &DrawingHub{
			colors: make(map[string]int64),
			users:  make(map[int64]*User),
			read:   make(chan Message),
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

func (s *Server) Close() error {
	s.hub.Close()

	cancelCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	log.Println("Server shutting down gracefully")
	return s.server.Shutdown(cancelCtx)
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
