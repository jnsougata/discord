package main

import (
	"fmt"
	"github.com/jnsougata/disgo/bot"
	"os"
)

func main() {
	c := New(Client{
		Intent: Intents{}.All(),
		Chunk:  true,
		Presence: Presence{
			Status:   "online",
			OnMobile: true,
			Activity: Activity{
				Type: 3,
				Name: "with Rick Astley",
				URL:  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			},
		},
	})
	c.OnReady(onReady)
	c.AddCommands(ping)
	c.Run(os.Getenv("DISCORD_TOKEN"))
}

func onReady(b bot.User) {
	fmt.Println(fmt.Sprintf("Running %s#%s (ID: %s)", b.Username, b.Discriminator, b.Id))
	fmt.Println("-------")
}

var ping = ApplicationCommand{
	Name:        "ping",
	Description: "shows the bot latency",
	Handler: func(b bot.User, ctx CommandContext, ops ...SlashCommandOption) {
		ctx.SendResponse(CommandResponse{Content: fmt.Sprintf("Pong! %dms", b.Latency)})
	},
}
