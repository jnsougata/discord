package main

import (
	"github.com/jnsougata/disgo/socket"
)

func Disgo(intent int, chunk bool) *Connection {
	return &Connection{Sock: &socket.Socket{Intent: intent, Memoize: chunk}}
}
