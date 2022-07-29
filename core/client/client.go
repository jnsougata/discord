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

type sockClient struct {
	socket *socket.Socket
}

func New(intent int, memoize bool) *sockClient {
	return &sockClient{socket: &socket.Socket{Intent: intent, Memoize: memoize}}
}

func (c *sockClient) Run(token string) {
	c.socket.Run(token)
}

func (c *sockClient) SetPresence(presence presence.Presence) {
	c.socket.StorePresenceData(presence)
}

func (c *sockClient) AddCommands(commands ...command.ApplicationCommand) {
	c.socket.Queue(commands...)
}

func (c *sockClient) OnSocketReceive(handler func(payload map[string]interface{})) {
	c.socket.AddHandler(consts.OnSocketReceive, handler)
}

func (c *sockClient) OnMessage(handler func(bot user.Bot, message message.Message)) {
	c.socket.AddHandler(consts.OnMessageCreate, handler)
}

func (c *sockClient) OnReady(handler func(bot user.Bot)) {
	c.socket.AddHandler(consts.OnReady, handler)
}

func (c *sockClient) OnInteraction(handler func(bot user.Bot, ctx command.Context)) {
	c.socket.AddHandler(consts.OnInteractionCreate, handler)
}

func (c *sockClient) OnGuildJoin(handler func(bot user.Bot, guild guild.Guild)) {
	c.socket.AddHandler(consts.OnGuildCreate, handler)
}

func (c *sockClient) OnGuildLeave(handler func(bot user.Bot, guild guild.Guild)) {
	c.socket.AddHandler(consts.OnGuildDelete, handler)
}
