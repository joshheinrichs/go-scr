// wscrc implements a websocket-based state change request client
package wscrc

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/joshheinrichs/go-scr/scr"
	"github.com/joshheinrichs/go-scr/scrr"
)

type Client struct {
	ws     *websocket.Conn
	router *scrr.Router
	state  interface{}
}

func New(ws *websocket.Conn, router *scrr.Router, state interface{}) *Client {
	client := &Client{
		ws:     ws,
		router: router,
		state:  state,
	}
	go client.read()
	return client
}

func (client *Client) read() {
	for {
		req := new(scrr.Request)
		err := client.ws.ReadJSON(req)
		if err != nil {
			log.Println("error:", err)
			break
		}
		client.router.Handle(client.state, req)
	}
}

func (client *Client) SendRequest(req *scr.Request) error {
	return client.ws.WriteJSON(req)
}
