package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joshheinrichs/go-scr/scr"
	"github.com/joshheinrichs/go-scr/scrr"
	"github.com/joshheinrichs/go-scr/wscrc"
)

var wg sync.WaitGroup
var client *wscrc.Client

var dialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type State struct {
	messages []*Message
	users    map[string]string
}

func NewState() *State {
	return &State{
		messages: make([]*Message, 0),
		users:    make(map[string]string),
	}
}

func itos(f func(*State, *scrr.Request)) scrr.Handler {
	return func(state interface{}, r *scrr.Request) {
		f(state.(*State), r)
	}
}

func messageHandler(state *State, r *scrr.Request) {
	message := new(Message)
	err := r.UnmarshalValue(message)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(message)
}

func userHandler(state *State, r *scrr.Request) {
}

func nameHandler(state *State, r *scrr.Request) {
}

func errorHandler(state *State, r *scrr.Request) {
	var str string
	err := r.UnmarshalValue(&str)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Error: %s\n", str)
}

// input reads input from the stdin and sends state change requests to the
// server.
func input() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		err = client.SendRequest(scr.Append(
			"chat.messages",
			&Message{Text: string(text[:len(text)-1])},
		))
		if err != nil {
			break
		}
	}
	wg.Done()
}

func main() {
	router := scrr.NewRouter(
		scrr.NewRoute(scr.CmdAppend, "chat.messages", itos(messageHandler)),
		scrr.NewRoute(scr.CmdSet, "chat.users.#id", itos(userHandler)),
		scrr.NewRoute(scr.CmdSet, "chat.users.#id.name", itos(nameHandler)),
		scrr.NewRoute(scr.CmdSet, "error", itos(errorHandler)),
	)
	ws, res, err := dialer.Dial("ws://localhost:3000/ws", nil)
	if err != nil {
		log.Panic(err, res)
	}
	client = wscrc.New(ws, router, NewState())
	go input()
	wg.Add(1)
	wg.Wait()
}
