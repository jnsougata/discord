package main

import (
	"encoding/json"
	"fmt"
	"github.com/disgo/core/bot"
	"github.com/disgo/core/objects"
	"github.com/disgo/core/types"
	"log"
	"os"
)

func main() {
	b := bot.New(33283)
	b.OnMessage(OnMessage)
	b.OnReady(OnReady)
	b.OnInteraction(OnInteraction)
	b.Run(os.Getenv("DISCORD_TOKEN"))
}

func OnMessage(message *types.Message) {
	log.Println(message.Content)
}
func OnReady() {
	log.Println("[-------- READY --------]")
	newCom := &objects.SlashCommand{
		Name:         "ping",
		Description:  "shows bot latency",
		DMPermission: true,
		TestGuildId:  123456,
		Options: []objects.JSONMap{
			objects.Option{}.String(
				"name",
				"type of the component to show ping for",
				false,
				0,
				100,
				false,
			),
		},
	}
	v, _ := json.MarshalIndent(newCom, "", "  ")
	fmt.Println(string(v))
}

func OnInteraction(interaction *types.Interaction) {
	interaction.Respond(
		&objects.InteractionMessage{
			Content: "Hello GoLang!",
			Embeds: []objects.Embed{
				{
					Title:       "disgo",
					Description: "testing disgo interaction",
					Color:       0x00FFFF,
				},
				{
					Title:       interaction.Member.User.Username,
					Description: "maintainer disgo interaction",
					Color:       0x00FFFF,
				},
			},
			Flags: 1 << 6,
		})
}
