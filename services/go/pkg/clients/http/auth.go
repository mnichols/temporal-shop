package http

import (
	"github.com/gojek/heimdall/v7"
	"net/http"
)

type basic struct {
	username string
	password string
}

func (b *basic) OnRequestStart(request *http.Request) {
	request.SetBasicAuth(b.username, b.password)
}

func (b *basic) OnRequestEnd(request *http.Request, response *http.Response) {
	// noop
}

func (b *basic) OnError(request *http.Request, err error) {
	// noop
}

func NewBasicAuth(username, password string) heimdall.Plugin {
	return &basic{
		username: username,
		password: password,
	}
}
