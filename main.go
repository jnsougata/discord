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
		gocmdHandler,
		objects.SlashCommand{
			Name:         "gocmd",
			Description:  "sample disgo command",
			DMPermission: true,
			TestGuildId:  877399405056102431,
		},
	)
	b.OnReady(OnReady)
	b.Run(os.Getenv("DISCORD_TOKEN"))
}

func OnReady(bot *types.User) {
	log.Println(fmt.Sprintf("Logged in as %s#%s (ID: %s)", bot.Username, bot.Discriminator, bot.ID))
	log.Println("---------")
}

func gocmdHandler(bot *types.User, interaction *types.Interaction) {
	interaction.SendResponse(
		&objects.InteractionMessage{
			Embeds: []objects.Embed{
				{
					Title:       "disgo",
					Description: "testing disgo interaction & slash commands",
					Color:       0x00FFFF,
				},
			},
		})
}
