package bot

import (
	"github.com/jnsougata/disgo/core/client"
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/consts"
	"github.com/jnsougata/disgo/core/guild"
	"github.com/jnsougata/disgo/core/message"
	"github.com/jnsougata/disgo/core/presence"
	"github.com/jnsougata/disgo/core/user"
)

type bot struct {
	intent int
	core   *client.Client
}

func New(intent int, memoize bool) *bot {
	return &bot{
		intent: intent,
		core: &client.Client{
			Intent:  intent,
			Memoize: memoize,
		},
	}
}

func (bot *bot) Run(token string) {
	bot.core.Run(token)
}

func (bot *bot) SetPresence(presence presence.Presence) {
	bot.core.StorePresenceData(presence)
}

func (bot *bot) AddCommands(commands ...command.ApplicationCommand) {
	bot.core.Queue(commands...)
}

func (bot *bot) OnSocketReceive(handler func(payload map[string]interface{})) {
	bot.core.AddHandler(consts.OnSocketReceive, handler)
}

func (bot *bot) OnMessage(handler func(bot user.Bot, message message.Message)) {
	bot.core.AddHandler(consts.OnMessageCreate, handler)
}

func (bot *bot) OnReady(handler func(bot user.Bot)) {
	bot.core.AddHandler(consts.OnReady, handler)
}

func (bot *bot) OnInteraction(handler func(bot user.Bot, ctx command.Context)) {
	bot.core.AddHandler(consts.OnInteractionCreate, handler)
}

func (bot *bot) OnGuildJoin(handler func(bot user.Bot, guild guild.Guild)) {
	bot.core.AddHandler(consts.OnGuildCreate, handler)
}

func (bot *bot) OnGuildLeave(handler func(bot user.Bot, guild guild.Guild)) {
	bot.core.AddHandler(consts.OnGuildDelete, handler)
}
