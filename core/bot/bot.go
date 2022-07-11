package bot

import (
	"github.com/disgo/core/client"
	"github.com/disgo/core/kind"
	"github.com/disgo/core/models"
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

func (bot *Bot) OnMessage(handler func(bot *kind.User, message *kind.Message)) {
	bot.core.AddHandler("MESSAGE_CREATE", handler)
}

func (bot *Bot) OnReady(handler func(bot *kind.User)) {
	bot.core.AddHandler("READY", handler)
}

func (bot *Bot) OnInteraction(handler func(bot *kind.User, interaction *kind.Interaction)) {
	bot.core.AddHandler("INTERACTION_CREATE", handler)
}

func (bot *Bot) AddCommand(
	handler func(bot *kind.User, interaction *kind.Interaction),
	command models.SlashCommand,
) {
	bot.core.Queue(command, handler)
}
