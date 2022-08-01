# disgo


# Quick Example

```go
package main

import (
	"fmt"
	"github.com/jnsougata/disgo"
	"os"
)

func main() {
	bot := disgo.Bot(disgo.Intents{}.All(), true, disgo.Presence{
		Since:  0,
		Status: "dnd",
		Activity: disgo.Activity{
			Name: "/ping",
			Type: 2,
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

var ping = disgo.ApplicationCommand{
	Name:        "ping",
	Description: "shows the bot latency",
	Handler: func(b disgo.BotUser, ctx disgo.Context, ops ...disgo.SlashCommandOption) {
		v := disgo.View{}
		v.AddButtons(disgo.Button{
			Style: 4,
			Label: "del",
			OnClick: func(bot disgo.BotUser, cc disgo.Context) {
				ctx.Delete()
			},
		})
		ctx.Send(disgo.Response{Content: fmt.Sprintf("Pong! %dms", b.Latency), View: v})
	},
}

```
