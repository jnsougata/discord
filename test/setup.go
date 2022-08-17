package main

import (
	"github.com/jnsougata/discord"
)

func setup() discord.Command {
	cmd := discord.Command{
		Name:              "setup",
		Description:       "Setup the bot",
		MemberPermissions: discord.Permissions.ManageGuild,
	}
	cmd.SubCommands(discord.SubCommand{
		Name:        "youtube",
		Description: "Sets up a youtube channel",
		Options: []discord.Option{
			{
				Name:        "channel",
				Description: "The channel ID/URL to set up",
				Type:        discord.StringOption,
				Required:    true,
			},
		},
		Execute: func(bot discord.Bot, ctx discord.Context, options discord.ResolvedOptions) {
			_ = ctx.Send(discord.Response{
				Content: "Setup youtube channel",
			})
		},
	})
	return cmd
}
