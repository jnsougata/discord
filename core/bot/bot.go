package bot

import (
	"github.com/disgo/core/client"
	"github.com/disgo/core/models"
	"github.com/disgo/core/types"
)

type Bot struct {
	intent int
	core   *client.Client
}

func New(intent int) *Bot {
	return &Bot{intent: intent, core: client.New(intent)}
}

func (bot *Bot) Run(token string) {
	bot.core.Run(token)
}

func (bot *Bot) OnMessage(handler func(bot *types.User, message *types.Message)) {
	bot.core.AddHandler("MESSAGE_CREATE", handler)
}

func (bot *Bot) OnReady(handler func(bot *types.User)) {
	bot.core.AddHandler("READY", handler)
}

func (bot *Bot) OnInteraction(handler func(bot *types.User, interaction *types.Interaction)) {
	bot.core.AddHandler("INTERACTION_CREATE", handler)
}

func (bot *Bot) AddCommand(
	hook func(bot *types.User, interaction *types.Interaction, options ...types.Option), command models.SlashCommand) {
	bot.core.Queue(command, hook)
}
