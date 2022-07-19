package main

import (
	"fmt"
	"github.com/jnsougata/disgo/core/bot"
	"github.com/jnsougata/disgo/core/models"
	"github.com/jnsougata/disgo/core/types"
	"log"
	"os"
)

func main() {

	b := bot.New(33283)
	cmd := models.NewSlashCommand(
		"go", "invokes the gocmd handler",
		0, true, 877399405056102431,
		models.Option{
			Name:        "process",
			Description: "the user to invoke the command on",
			Type:        models.SubCommandType,
			Options: []models.Option{
				{
					Name:        "id",
					Type:        models.StringType,
					Required:    true,
					Description: "the pid to invoke the command on",
				},
				{
					Name:        "name",
					Type:        models.StringType,
					Required:    true,
					Description: "the name to invoke the command on",
				},
			},
		})
	b.AddCommand(goCommand, cmd)
	b.OnReady(OnReady)
	b.Run(os.Getenv("DISCORD_TOKEN"))
}

func OnReady(bot *types.User) {
	log.Println(fmt.Sprintf("Logged in as %s#%s (ID: %s)", bot.Username, bot.Discriminator, bot.ID))
	log.Println("---------")
}

func goCommand(_ *types.User, interaction *types.Interaction, options ...types.Option) {
	interaction.SendResponse(
		&models.InteractionMessage{
			Embeds: []models.Embed{
				{
					Title:       fmt.Sprintf("PID: %s", options[0].Options[0].Value.(string)),
					Description: fmt.Sprintf("Process: %s", options[0].Options[1].Value.(string)),
					Color:       0x00FFFF,
				},
			},
		})
}
