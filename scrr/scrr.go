// Package scrr is for a state change request router
package scrr

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/joshheinrichs/go-scr/scr"
)

const (
	pathDelimiter   = "."
	pathParamPrefix = "#"
)

var ErrInavlidCommand = errors.New("Invalid command")
var ErrInvalidPath = errors.New("Invalid path")
var ErrNoMatchingPath = errors.New("No matching path")

var regexpHandlerPath = regexp.MustCompile("^(#?[a-zA-Z0-9_]+(\\.#?[a-zA-Z0-9_]+)*)?$")
var regexpRequestPath = regexp.MustCompile("^([a-zA-Z0-9_]+(\\.[a-zA-Z0-9_]+)*)?$")

type Request struct {
	Command    string            `json:"command"`
	Path       string            `json:"path"`
	pathParts  []string          `json:"-"`
	PathParams map[string]string `json:"-"`
	Value      json.RawMessage   `json:"value"`
}

func NewRequest(command, path string, value interface{}) *Request {
	jsonValue, _ := json.Marshal(value)
	return &Request{
		Command: command,
		Path:    path,
		Value:   jsonValue,
	}
}

func (req *Request) UnmarshalValue(value interface{}) error {
	return json.Unmarshal(req.Value, value)
}

type Route struct {
	Command   string
	Path      string
	pathParts []string
	Handler   Handler
}

type Handler func(interface{}, *Request)

func NewRouter(routes ...*Route) *Router {
	return &Router{
		routes: routes,
	}
}

// NewRoute returns a new route, or nil if the route is invalid.
func NewRoute(command, path string, handler Handler) *Route {
	if command != scr.CmdSet &&
		command != scr.CmdAppend &&
		command != scr.CmdDelete &&
		command != scr.CmdCall {
		return nil
	}
	if !regexpHandlerPath.MatchString(path) {
		log.Println(ErrInvalidPath, path)
		return nil
	}
	return &Route{
		Command:   command,
		Path:      path,
		pathParts: strings.Split(path, pathDelimiter),
		Handler:   handler,
	}
}

func (route *Route) Match(req *Request) bool {
	if req.Command != route.Command {
		return false
	}
	if len(route.pathParts) != len(req.pathParts) {
		return false
	}
	pathParams := make(map[string]string)
	for i, _ := range route.pathParts {
		if route.pathParts[i][0] == '#' {
			pathParams[strings.Trim(route.pathParts[i], pathParamPrefix)] = req.pathParts[i]
		} else if route.pathParts[i] != req.pathParts[i] {
			return false
		}
	}
	req.PathParams = pathParams
	return true
}

type Router struct {
	routes []*Route
}

func (router *Router) Handle(state interface{}, req *Request) error {
	if !regexpRequestPath.MatchString(req.Path) {
		return ErrInvalidPath
	}
	req.pathParts = strings.Split(req.Path, pathDelimiter)
	for _, route := range router.routes {
		if route.Match(req) {
			route.Handler(state, req)
			return nil
		}
	}
	return ErrNoMatchingPath
}
