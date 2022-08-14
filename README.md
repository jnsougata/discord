# discord


# Quick Example

```go
package main

import (
	"fmt"
	"github.com/jnsougata/discord"
	"os"
)

func main() {

	bot := discord.New(discord.Intents())
	bot.Cached = true
	bot.Presence = discord.Presence{
          Since:  0,
          Status: discord.Online,
          Activity: discord.Activity{
                Name: "/ping",
                Type: discord.Listening,
          },
    }
	bot.Listeners = listeners
	bot.Commands(avatar)
	bot.Run(os.Getenv("DISCORD_TOKEN"))
}

var listeners = discord.Listeners{
    OnReady: func(bot discord.Bot) {
    	fmt.Println(fmt.Sprintf("Running %s#%s (Id: %s)", bot.Username, bot.Discriminator, bot.Id))
    	fmt.Println("-------")
    },
    // add more built-in listeners here
}

var avatar = discord.Command{
	Name:        "avatar",
	Description: "shows the avatar of the invoker",
	Execute: func(bot discord.Bot, ctx discord.Context, options discord.ResolvedOptions) {
		ctx.Send(discord.Response{Embed: discord.Embed{Image: discord.EmbedImage{Url: ctx.Author.Avatar.URL()}}})
	},
}

```
