package main

import (
	"github.com/jnsougata/disgo/client"
	command2 "github.com/jnsougata/disgo/command"
	"github.com/jnsougata/disgo/consts"
	"github.com/jnsougata/disgo/guild"
	"github.com/jnsougata/disgo/message"
	"github.com/jnsougata/disgo/presence"
	"github.com/jnsougata/disgo/socket"
)

type Connection struct {
	Sock *socket.Socket
}

func (conn *Connection) Run(token string) {
	conn.Sock.Run(token)
}

func (conn *Connection) SetPresence(presence presence.Presence) {
	conn.Sock.StorePresenceData(presence)
}

func (conn *Connection) AddCommands(commands ...command2.ApplicationCommand) {
	conn.Sock.RegistrationQueue(commands...)
}

func (conn *Connection) OnSocketReceive(handler func(payload map[string]interface{})) {
	conn.Sock.AddHandler(consts.OnSocketReceive, handler)
}

func (conn *Connection) OnMessage(handler func(bot client.User, message message.Message)) {
	conn.Sock.AddHandler(consts.OnMessageCreate, handler)
}

func (conn *Connection) OnReady(handler func(bot client.User)) {
	conn.Sock.AddHandler(consts.OnReady, handler)
}

func (conn *Connection) OnInteraction(handler func(bot client.User, ctx command2.Context)) {
	conn.Sock.AddHandler(consts.OnInteractionCreate, handler)
}

func (conn *Connection) OnGuildJoin(handler func(bot client.User, guild guild.Guild)) {
	conn.Sock.AddHandler(consts.OnGuildCreate, handler)
}

func (conn *Connection) OnGuildLeave(handler func(bot client.User, guild guild.Guild)) {
	conn.Sock.AddHandler(consts.OnGuildDelete, handler)
}
