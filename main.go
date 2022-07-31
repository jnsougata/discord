package main

import (
	"fmt"
	"github.com/jnsougata/disgo/bot"
	"github.com/jnsougata/disgo/client"
	"github.com/jnsougata/disgo/intents"
	"github.com/jnsougata/disgo/presence"
	"os"
)

func main() {
	c := client.New(client.Client{
		Intent: intents.All(),
		Chunk:  true,
		Presence: presence.Presence{
			Status:   "online",
			OnMobile: true,
			Activity: presence.Activity{
				Type: 3,
				Name: "with Rick Astley",
				URL:  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			},
		},
	})
	c.OnReady(onReady)
	c.Run(os.Getenv("DISCORD_TOKEN"))
}

func onReady(b bot.User) {
	fmt.Println(fmt.Sprintf("Running %s#%s (ID: %s)", b.Username, b.Discriminator, b.Id))
	fmt.Println("-------")
}
