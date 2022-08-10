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
		Task: func(bot disgo.BotUser, ctx disgo.Context, options map[string]disgo.Option) {
			//fmt.Println(bot.Users)
			ctx.Send(disgo.Response{Embed: disgo.Embed{Image: disgo.EmbedImage{Url: ctx.Author.Avatar.URL()}}})
		},
	}
	bot.AddCommands(ping)
	bot.OnReady(onReady)
	bot.Run(os.Getenv("DISCORD_TOKEN"))
}

func onReady(b disgo.BotUser) {
	fmt.Println(fmt.Sprintf("Running %s#%s (Id: %s)", b.Username, b.Discriminator, b.Id))
	fmt.Println("-------")
}
