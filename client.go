package main

import (
	"github.com/jnsougata/disgo/bot"
)

func New(client Client) *connection {
	return &connection{
		Sock: &Socket{
			Intent:   client.Intent,
			Memoize:  client.Chunk,
			Presence: client.Presence,
		},
	}
}

type Client struct {
	Intent   int
	Chunk    bool
	Presence Presence
}

type connection struct {
	Sock *Socket
}

func (conn *connection) Run(token string) {
	conn.Sock.Run(token)
}

func (conn *connection) AddCommands(commands ...ApplicationCommand) {
	conn.Sock.RegistrationQueue(commands...)
}

func (conn *connection) OnSocketReceive(handler func(payload map[string]interface{})) {
	conn.Sock.AddHandler(OnSocketReceive, handler)
}

func (conn *connection) OnMessage(handler func(bot bot.User, message Message)) {
	conn.Sock.AddHandler(OnMessageCreate, handler)
}

func (conn *connection) OnReady(handler func(bot bot.User)) {
	conn.Sock.AddHandler(OnReady, handler)
}

func (conn *connection) OnInteraction(handler func(bot bot.User, ctx CommandContext)) {
	conn.Sock.AddHandler(OnInteractionCreate, handler)
}

func (conn *connection) OnGuildJoin(handler func(bot bot.User, guild Guild)) {
	conn.Sock.AddHandler(OnGuildCreate, handler)
}

func (conn *connection) OnGuildLeave(handler func(bot bot.User, guild Guild)) {
	conn.Sock.AddHandler(OnGuildDelete, handler)
}
