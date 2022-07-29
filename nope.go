package main

import (
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/component"
	"github.com/jnsougata/disgo/core/embed"
	"github.com/jnsougata/disgo/core/emoji"
	"github.com/jnsougata/disgo/core/interaction"
	"github.com/jnsougata/disgo/core/user"
)

var toast = command.ApplicationCommand{
	Name:        "toast",
	Description: "sends a toast",
	Handler: func(b user.User, ctx command.Context, ops ...interaction.Option) {
		emo := emoji.Partial{Name: "toast", Id: "885467894979375124"}
		btn := component.Button{
			Label: "Delete",
			Style: 3,
			OnClick: func(bot user.User, cc component.Context) {
				ctx.DeleteOriginalResponse()
			},
		}
		menu := component.SelectMenu{
			Placeholder: "Select multiple options",
			Options: []component.SelectOption{
				{Label: "one", Value: "1", Emoji: emo},
				{Label: "two", Value: "2", Emoji: emo},
				{Label: "three", Value: "3", Emoji: emo},
			},
			MinValues: 1,
			MaxValues: 1,
		}
		view := component.View{}
		view.AddButtons(btn)
		view.AddSelectMenu(menu)
		ctx.SendResponse(command.Message{
			Embed: embed.Embed{Description: "Here's a toast for you!"},
			View:  view,
		})
	},
}
