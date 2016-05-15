package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joshheinrichs/go-scr/scr"
	"github.com/joshheinrichs/go-scr/scrr"
)

var router = scrr.NewRouter(
	scrr.NewRoute(scr.CmdAppend, "chat.messages", itou(messageHandler)),
	scrr.NewRoute(scr.CmdSet, "chat.users.#id", itou(userHandler)),
	scrr.NewRoute(scr.CmdSet, "chat.users.#id.name", itou(nameHandler)),
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var users = NewUsers()
var room = NewRoom()

type Users struct {
	rwMutex sync.RWMutex
	users   map[string]*User
	serial  int
}

func NewUsers() *Users {
	return &Users{
		users:  make(map[string]*User),
		serial: 0,
	}
}

func (users *Users) NewUser(ws *websocket.Conn) *User {
	users.rwMutex.Lock()
	defer users.rwMutex.Unlock()
	user := &User{
		ws:       ws,
		id:       strconv.Itoa(users.serial),
		requests: make(chan *scr.Request),
		room:     room,
	}
	users.serial++
	go user.listen()
	room.Join(user)
	return user
}

type User struct {
	rwMutex  sync.RWMutex
	id       string
	ws       *websocket.Conn
	requests chan *scr.Request
	room     *Room
}

func (user *User) listen() {
	go user.read()
	go user.write()
}

func (user *User) read() {
	for {
		request := new(scrr.Request)
		err := user.ws.ReadJSON(request)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Println("user closed connection:", err)
				break
			} else {
				log.Println("unexpected error:", err)
				// user.SendMessage(NewErrorMessage(err)
				break
			}
		}
		err = router.Handle(user, request)
		if err != nil {
			log.Println("unexpected error:", err)
			// user.SendMessage(NewErrorMessage(err))
			continue
		}
	}
	// user disconnected
	user.ws.Close()
	user.room.Leave(user)
	user.room = nil
	close(user.requests)
}

func (user *User) write() {
	for {
		request, more := <-user.requests
		if !more {
			break
		}
		err := user.ws.WriteJSON(request)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Println("user closed connection:", err)
				break
			} else {
				log.Println("unexpected error:", err)
				user.ws.Close()
				break
			}
		}
	}
	// user disconnected
}

func (user *User) SendRequest(r *scr.Request) {
	user.requests <- r
}

type Room struct {
	rwMutex sync.RWMutex
	users   map[string]*User
}

func NewRoom() *Room {
	return &Room{
		users: make(map[string]*User),
	}
}

func (room *Room) Join(u *User) {
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()
	room.users[u.id] = u
}

func (room *Room) Leave(u *User) {
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()
	delete(room.users, u.id)
}

func (room *Room) BroadcastMessage(m *Message) {
	room.rwMutex.RLock()
	defer room.rwMutex.RUnlock()
	r := scr.Append("chat.messages", m)
	for k := range room.users {
		room.users[k].SendRequest(r)
	}
}

func itou(f func(*User, *scrr.Request)) scrr.Handler {
	return func(s interface{}, r *scrr.Request) {
		f(s.(*User), r)
	}
}

func messageHandler(u *User, r *scrr.Request) {
	u.rwMutex.RLock()
	defer u.rwMutex.RUnlock()
	m := new(Message)
	r.UnmarshalValue(m)
	m.Time = time.Now()
	m.UserID = u.id
	u.room.BroadcastMessage(m)
}

func userHandler(u *User, r *scrr.Request) {
}

func nameHandler(u *User, r *scrr.Request) {
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	users.NewUser(ws)
}

func main() {
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
