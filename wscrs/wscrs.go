// wscrc implements a websocket-based state change request server
package wscrs

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/joshheinrichs/go-scr/scrr"
)

type Server struct {
	ws     *websocket.Conn
	router *scrr.Router
	state  interface{}
}

func New(ws *websocket.Conn, router *scrr.Router, state interface{}) *Server {
	server := &Server{
		ws:     ws,
		router: router,
		state:  state,
	}
	go server.read()
	return server
}

func (server *Server) read() {
	for {
		req := new(scrr.Request)
		err := server.ws.ReadJSON(req)
		if err != nil {
			log.Println("error:", err)
			break
		}
		server.router.Handle(server.state, req)
	}
}

func (server *Server) SendRequest(req *scrr.Request) error {
	return server.ws.WriteJSON(req)
}
