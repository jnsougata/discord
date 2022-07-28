package bot

import (
	"github.com/jnsougata/disgo/core/client"
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/consts"
	"github.com/jnsougata/disgo/core/guild"
	"github.com/jnsougata/disgo/core/message"
	"github.com/jnsougata/disgo/core/user"
)

type Bot struct {
	intent int
	core   *client.Client
}

func New(intent int, memoize bool) *Bot {
	return &Bot{
		intent: intent,
		core: &client.Client{
			Intent:  intent,
			Memoize: memoize,
		},
	}
}

func (bot *Bot) Run(token string) {
	bot.core.Run(token)
}

func (bot *Bot) AddCommand(command command.ApplicationCommand) {
	bot.core.Queue(command)
}

func (bot *Bot) OnSocketReceive(handler func(payload map[string]interface{})) {
	bot.core.AddHandler(consts.OnSocketReceive, handler)
}

func (bot *Bot) OnMessage(handler func(bot user.User, message message.Message)) {
	bot.core.AddHandler(consts.OnMessageCreate, handler)
}

func (bot *Bot) OnReady(handler func(bot user.User)) {
	bot.core.AddHandler(consts.OnReady, handler)
}

func (bot *Bot) OnInteraction(handler func(bot user.User, ctx command.Context)) {
	bot.core.AddHandler(consts.OnInteractionCreate, handler)
}

func (bot *Bot) OnGuildJoin(handler func(bot user.User, guild guild.Guild)) {
	bot.core.AddHandler(consts.OnGuildCreate, handler)
}

func (bot *Bot) OnGuildLeave(handler func(bot user.User, guild guild.Guild)) {
	bot.core.AddHandler(consts.OnGuildDelete, handler)
}
