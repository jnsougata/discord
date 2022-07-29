package main

import (
	"github.com/jnsougata/disgo/client"
	"github.com/jnsougata/disgo/command"
	"github.com/jnsougata/disgo/consts"
	"github.com/jnsougata/disgo/guild"
	"github.com/jnsougata/disgo/message"
	"github.com/jnsougata/disgo/presence"
	"github.com/jnsougata/disgo/socket"
)

type Disgo struct {
	Intent   int
	Chunk    bool
	Presence presence.Presence
}

func (d Disgo) New() *connection {
	return &connection{
		Sock: &socket.Socket{
			Intent:   d.Intent,
			Memoize:  d.Chunk,
			Presence: d.Presence,
		},
	}
}

type connection struct {
	Sock *socket.Socket
}

func (conn *connection) Run(token string) {
	conn.Sock.Run(token)
}

func (conn *connection) AddCommands(commands ...command.ApplicationCommand) {
	conn.Sock.RegistrationQueue(commands...)
}

func (conn *connection) OnSocketReceive(handler func(payload map[string]interface{})) {
	conn.Sock.AddHandler(consts.OnSocketReceive, handler)
}

func (conn *connection) OnMessage(handler func(bot client.User, message message.Message)) {
	conn.Sock.AddHandler(consts.OnMessageCreate, handler)
}

func (conn *connection) OnReady(handler func(bot client.User)) {
	conn.Sock.AddHandler(consts.OnReady, handler)
}

func (conn *connection) OnInteraction(handler func(bot client.User, ctx command.Context)) {
	conn.Sock.AddHandler(consts.OnInteractionCreate, handler)
}

func (conn *connection) OnGuildJoin(handler func(bot client.User, guild guild.Guild)) {
	conn.Sock.AddHandler(consts.OnGuildCreate, handler)
}

func (conn *connection) OnGuildLeave(handler func(bot client.User, guild guild.Guild)) {
	conn.Sock.AddHandler(consts.OnGuildDelete, handler)
}
