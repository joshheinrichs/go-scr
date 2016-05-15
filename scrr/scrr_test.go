package scrr

import (
	"log"
	"testing"

	"github.com/joshheinrichs/go-scr/scr"
	"github.com/stretchr/testify/assert"
)

func FooHandler(i interface{}, req *Request) {
	log.Println(req.PathParams["b"])
	log.Println(req.PathParams["c"])
}

func TestMatch(t *testing.T) {
	router := NewRouter(
		NewRoute(scr.CmdSet, "a.#b.#c.d", FooHandler),
	)
	assert.NoError(t, router.Handle(nil, &Request{
		Command: scr.CmdSet,
		Path:    "a.ayy.lmao.d",
	}))
	assert.Error(t, router.Handle(nil, &Request{
		Command: scr.CmdAppend,
		Path:    "foo.bar",
	}))
}
