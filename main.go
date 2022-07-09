package main

import (
	"fmt"
	"github.com/disgo/core/bot"
	"github.com/disgo/core/objects"
	"github.com/disgo/core/types"
	"log"
	"os"
)

func main() {
	b := bot.New(33283)
	b.AddCommand(
		objects.SlashCommand{
			Name:         "gocmd",
			Description:  "sample disgo command",
			DMPermission: true,
			TestGuildId:  877399405056102431,
			Options: []objects.JSONMap{
				objects.Option{}.String(
					"string",
					"string type option",
					true,
					0,
					10,
					false,
				),
			},
		},
	)
	b.OnMessage(OnMessage)
	b.OnReady(OnReady)
	b.OnInteraction(OnInteraction)
	b.Run(os.Getenv("DISCORD_TOKEN"))
}

func OnMessage(bot *types.User, message *types.Message) {
	log.Println(message.Content)
}

func OnReady(bot *types.User) {
	log.Println(fmt.Sprintf("[-------- (%s#%s) --------]", bot.Username, bot.Discriminator))
}

func OnInteraction(bot *types.User, interaction *types.Interaction) {
	interaction.Respond(
		&objects.InteractionMessage{
			Content: "Hello GoLang!",
			Embeds: []objects.Embed{
				{
					Title:       "disgo",
					Description: "testing disgo interaction & slash commands",
					Color:       0x00FFFF,
				},
			},
		})
}
