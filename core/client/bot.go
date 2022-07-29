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

type client struct {
	socket *socket.Socket
}

func New(intent int, memoize bool) *client {
	return &client{socket: &socket.Socket{Intent: intent, Memoize: memoize}}
}

func (c *client) Run(token string) {
	c.socket.Run(token)
}

func (c *client) SetPresence(presence presence.Presence) {
	c.socket.StorePresenceData(presence)
}

func (c *client) AddCommands(commands ...command.ApplicationCommand) {
	c.socket.Queue(commands...)
}

func (c *client) OnSocketReceive(handler func(payload map[string]interface{})) {
	c.socket.AddHandler(consts.OnSocketReceive, handler)
}

func (c *client) OnMessage(handler func(bot user.Bot, message message.Message)) {
	c.socket.AddHandler(consts.OnMessageCreate, handler)
}

func (c *client) OnReady(handler func(bot user.Bot)) {
	c.socket.AddHandler(consts.OnReady, handler)
}

func (c *client) OnInteraction(handler func(bot user.Bot, ctx command.Context)) {
	c.socket.AddHandler(consts.OnInteractionCreate, handler)
}

func (c *client) OnGuildJoin(handler func(bot user.Bot, guild guild.Guild)) {
	c.socket.AddHandler(consts.OnGuildCreate, handler)
}

func (c *client) OnGuildLeave(handler func(bot user.Bot, guild guild.Guild)) {
	c.socket.AddHandler(consts.OnGuildDelete, handler)
}
