# go-scr ![](https://godoc.org/github.com/joshheinrichs/go-scr?status.svg)

go-scr provides a model for communication via state change requests. This would ideally simplify cases where you have a shared state between clients and a server, and you want the server to dictate legal state changes and handle requests in a generalizable way. This also has advantages for minimizing the amount of traffic sent over a connection, as only changes to the state would be sent as opposed sending the entire state, which could contain a large amount of redundant information.

## Example

Take for example, a simple game of x's and o's. The clients and server may have a shared state as follows:

```json
{
	"turn": "o",
	"board": [
		["x", "", ""],
		["", "o", ""],
		["x", "", ""]
	]
}
```

If a client wanted to make a move, they could send a request to the server like so:
```json
{
	"command": "set",
	"path": "board.0.1",
	"value": "o"
}
```

The server would then verify that the proposed state change was legal, and then send the following scrs back to each of the clients:
```json
{
	"command": "set",
	"path": "board.0.1",
	"value": "o"
}
```
```json
{
	"command": "set",
	"path": "turn",
	"value": "x"
}
```

Handling these paths can be somewhat annoying, so a routing library, `scrr`, has been included, which matches a request's operation and path to a handler function.

In the example given above, the following router could be constructed for the server:
```go
	router := scrr.NewRouter(
		scrr.NewRoute(scr.CmdSet, "board.#x.#y", moveHandler),
	)
```

The values of `x` and `y` are dynamically set and passed into the handler so that paths can be handled in a generalizable way.

## Disclaimer

This package was specifically created for a personal project, and as such isn't as efficient or feature complete as I would ultimately like. If it works well I may improve upon this package in the future, but there are no current plans for continued development.
