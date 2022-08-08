package main

import (
	"fmt"
	"github.com/jnsougata/disgo"
	"os"
)

func main() {
	bot := disgo.Bot(disgo.Intents(), true, disgo.Presence{
		Since:  0,
		Status: "dnd",
		Activity: disgo.Activity{
			Name: "/ping",
			Type: 2,
		},
	})

	ping := disgo.ApplicationCommand{
		Name:        "ping",
		Description: "Pong!",
	}
	ping.SubCommands(disgo.SubCommand{
		Name:        "sub",
		Description: "Subcommand",
		Options: []disgo.Option{
			{
				Name:        "option",
				Description: "The option to use",
				Type:        disgo.StringOption,
				MaxLength:   2,
				MinLength:   1,
			},
		},
		Task: func(bot disgo.BotUser, ctx disgo.Context, options map[string]disgo.Option) {
			ctx.Channel.Send(disgo.Draft{Content: "Pong!....", DeleteAfter: 2})
		},
	})
	bot.AddCommands(ping)
	bot.OnReady(onReady)
	bot.Run(os.Getenv("DISCORD_TOKEN"))
}

func onReady(b disgo.BotUser) {
	fmt.Println(fmt.Sprintf("Running %s#%s (ID: %s)", b.Username, b.Discriminator, b.Id))
	fmt.Println("-------")
}
