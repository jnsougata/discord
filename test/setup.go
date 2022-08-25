package main

import (
	"github.com/jnsougata/discord"
)

func setup() discord.Command {
	cmd := discord.Command{
		Name:        "setup",
		Description: "Setup the bot",
		Permissions: []discord.Permission{
			discord.Permissions.ManageGuild,
			discord.Permissions.SendTTSMessages,
		},
	}
	//cmd.OptionBOOLEAN("cache", "Enable the cache", true)
	//cmd.OptionUSER("presence", "Set the presence of the bot", true)
	cmd.Execute = func(bot discord.Bot, ctx discord.Context, options discord.ResolvedOptions) {
		ctx.Send(discord.Response{Content: "Setup complete!"})
		//ctx.Channel.Send(discord.Draft{Content: "Hello, I'm a bot!"})
	}
	return cmd
}
