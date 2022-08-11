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
		Status: "dnd",
		Activity: discord.Activity{
			Name: "/ping",
			Type: 2,
		},
	})

	ping := discord.ApplicationCommand{
		Name:        "ping",
		Description: "Pong!",
		Task: func(bot discord.BotUser, ctx discord.Context, options map[string]discord.Option) {
			for _, role := range ctx.Author.Roles {
				fmt.Println(role.Permissions)
			}
			ctx.Send(discord.Response{Embed: discord.Embed{Image: discord.EmbedImage{Url: ctx.Author.Avatar.URL()}}})
		},
	}
	bot.AddCommands(ping)
	bot.OnReady(onReady)
	bot.Run(os.Getenv("DISCORD_TOKEN"))
}

func onReady(b discord.BotUser) {
	fmt.Println(fmt.Sprintf("Running %s#%s (Id: %s)", b.Username, b.Discriminator, b.Id))
	fmt.Println("-------")
}

```
