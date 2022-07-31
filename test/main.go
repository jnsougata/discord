package main

import (
	"fmt"
	"github.com/jnsougata/disgo"
	"os"
)

func main() {
	b := disgo.Bot(disgo.Intents{}.All(), true, disgo.Presence{})
	b.OnReady(ready)
	b.AddCommands(cmd)
	b.Run(os.Getenv("DISCORD_TOKEN"))
}

func ready(b disgo.BotUser) {
	fmt.Println(fmt.Sprintf("%s is ready!", b.Username))
}

var cmd = disgo.ApplicationCommand{
	Name:        "ping",
	Description: "Ping!",
	Handler: func(bot disgo.BotUser, ctx disgo.Context, ops ...disgo.SlashCommandOption) {
		ctx.Send(disgo.Response{Content: fmt.Sprintf("Pong! %vms", bot.Latency)})
	},
}
