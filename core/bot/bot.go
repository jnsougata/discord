package bot

import (
	"github.com/disgo/core/client"
	"github.com/disgo/core/types"
)

type Bot struct {
	intent int
	core   *client.Client
}

func New(intent int) *Bot {
	c := client.New(intent)
	return &Bot{intent: intent, core: c}
}

func (bot *Bot) Run(token string) {
	bot.core.Run(token)
}

func (bot *Bot) OnMessage(handler func(message *types.Message)) {
	bot.core.AddHandler("MESSAGE_CREATE", handler)
}

func (bot *Bot) OnReady(handler func()) {
	bot.core.AddHandler("READY", handler)
}

func (bot *Bot) OnInteraction(handler func(interaction *types.Interaction)) {
	bot.core.AddHandler("INTERACTION_CREATE", handler)
}

func (bot *Bot) CommandQueue(commands ...any) {
	bot.core.RegistrationQueue(commands)
}
