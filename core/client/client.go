package client

import (
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/consts"
	"github.com/jnsougata/disgo/core/guild"
	"github.com/jnsougata/disgo/core/message"
	"github.com/jnsougata/disgo/core/presence"
	"github.com/jnsougata/disgo/core/socket"
	"github.com/jnsougata/disgo/core/user"
)

type connection struct {
	socket *socket.Socket
}

func New(intent int, chunk bool) *connection {
	return &connection{socket: &socket.Socket{Intent: intent, Memoize: chunk}}
}

func (c *connection) Run(token string) {
	c.socket.Run(token)
}

func (c *connection) SetPresence(presence presence.Presence) {
	c.socket.StorePresenceData(presence)
}

func (c *connection) AddCommands(commands ...command.ApplicationCommand) {
	c.socket.RegistrationQueue(commands...)
}

func (c *connection) OnSocketReceive(handler func(payload map[string]interface{})) {
	c.socket.AddHandler(consts.OnSocketReceive, handler)
}

func (c *connection) OnMessage(handler func(bot user.Bot, message message.Message)) {
	c.socket.AddHandler(consts.OnMessageCreate, handler)
}

func (c *connection) OnReady(handler func(bot user.Bot)) {
	c.socket.AddHandler(consts.OnReady, handler)
}

func (c *connection) OnInteraction(handler func(bot user.Bot, ctx command.Context)) {
	c.socket.AddHandler(consts.OnInteractionCreate, handler)
}

func (c *connection) OnGuildJoin(handler func(bot user.Bot, guild guild.Guild)) {
	c.socket.AddHandler(consts.OnGuildCreate, handler)
}

func (c *connection) OnGuildLeave(handler func(bot user.Bot, guild guild.Guild)) {
	c.socket.AddHandler(consts.OnGuildDelete, handler)
}
