package main

import (
	"fmt"
	"github.com/jnsougata/disgo/core/bot"
	"github.com/jnsougata/disgo/core/intents"
	"github.com/jnsougata/disgo/core/presence"
	"github.com/jnsougata/disgo/core/user"
	"log"
	"os"
)

func main() {
	b := bot.New(intents.All(), false)
	b.AddCommands(toast)
	b.OnReady(onReady)
	b.SetPresence(pr)
	b.OnSocketReceive(onSocketReceive)
	b.Run(os.Getenv("DISCORD_TOKEN"))
}

func onReady(bot user.User) {
	log.Println(fmt.Sprintf(
		"Logged in as %s#%s (Id: %s)",
		bot.Username, bot.Discriminator, bot.Id,
	))
	log.Println("---------")
}

func onSocketReceive(_ map[string]interface{}) {
	//log.Println(payload)
}

var pr = presence.Presence{
	Status:       "online",
	AFK:          false,
	ClientStatus: "mobile",
	Activity: presence.Activity{
		Name: "LO:FI",
		Type: 3,
		URL:  "https://www.youtube.com/watch?v=e97w-GHsRMY",
	},
}
