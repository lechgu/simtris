package blocks

import (
	"log"

	"github.com/gorilla/websocket"
)

// Session ...
type Session struct {
	conn *websocket.Conn
}

// NewSession ...
func NewSession(conn *websocket.Conn) *Session {
	return &Session{
		conn: conn,
	}
}

func commandPump(conn *websocket.Conn, commands chan string, done chan bool) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			done <- true
		}
		commands <- string(msg)
	}
}

// Run ...
func (s *Session) Run() {
	model := NewModel()
	commands := make(chan string)
	done := make(chan bool)
	go model.Run()
	go commandPump(s.conn, commands, done)
	for {
		select {
		case <-done:
			model.done <- true
			return
		case state := <-model.updates:
			err := s.conn.WriteMessage(1, state)
			if err != nil {
				log.Println(err)
				return
			}
		case cmd := <-commands:
			model.commands <- cmd
		}
	}
}
