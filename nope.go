package main

import (
	"fmt"
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/interaction"
	"github.com/jnsougata/disgo/core/user"
)

var toast = command.ApplicationCommand{
	Name:        "ping",
	Description: "shows the latency of the client",
	Handler: func(b user.Bot, ctx command.Context, ops ...interaction.Option) {
		ctx.SendResponse(command.Message{Content: fmt.Sprintf("Pong! %vms", b.Latency)})
	},
}
