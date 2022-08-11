# disgo


# Quick Example

```go
package main

import (
	"fmt"
	"github.com/jnsougata/discord"
	"os"
)

func main() {

	bot := discord.Bot(discord.Intents(), true, discord.Presence{
		Since:  0,
		Status: discord.Online,
		Activity: discord.Activity{
			Name: "/ping",
			Type: discord.Listening,
		},
	})
	bot.Listeners = discord.Listeners{
		OnReady: func(bot discord.BotUser) {
			fmt.Println(fmt.Sprintf("Running %s#%s (Id: %s)", bot.Username, bot.Discriminator, bot.Id))
			fmt.Println("-------")
		},
	}
	bot.Commands(ping)
	bot.Run(os.Getenv("DISCORD_TOKEN"))
}

var ping = discord.ApplicationCommand{
	Name:        "ping",
	Description: "Pong!",
	Task: func(bot discord.BotUser, ctx discord.Context, options map[string]discord.Option) {
		ctx.Send(
			discord.Response{Embed: discord.Embed{Image: discord.EmbedImage{Url: ctx.Author.Avatar.URL()}}})
	},
}

```
