package scr

const (
	CmdSet    = "set"
	CmdAppend = "append"
	CmdDelete = "delete"
	CmdCall   = "call"
)

type Request struct {
	Command string      `json:"command"`
	Path    string      `json:"path"`
	Value   interface{} `json:"value"`
}

func NewRequest(command, path string, value interface{}) *Request {
	return &Request{
		Command: command,
		Path:    path,
		Value:   value,
	}
}

func Set(path string, value interface{}) *Request {
	return NewRequest(CmdSet, path, value)
}

func Append(path string, value interface{}) *Request {
	return NewRequest(CmdAppend, path, value)
}

func Delete(path string, value interface{}) *Request {
	return NewRequest(CmdDelete, path, value)
}

func Call(path string, value interface{}) *Request {
	return NewRequest(CmdCall, path, value)
}
