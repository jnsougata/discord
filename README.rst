disgo
-----

Quick Example
-------------

.. code:: go

    package main

    import (
        "fmt"
        "github.com/jnsougata/disgo/client"
        "github.com/jnsougata/disgo/command"
        "github.com/jnsougata/disgo/intents"
        "github.com/jnsougata/disgo/interaction"
        "github.com/jnsougata/disgo/presence"
        "log"
        "os"
    )

    func main() {
        c := client.New(client.Client{Intent: intents.All(), Chunk: true, Presence: p})
        c.AddCommands(ping)
        c.OnReady(onReady)
        c.Run(os.Getenv("DISCORD_TOKEN"))
    }

    func onReady(bot client.User) {
        log.Println(fmt.Sprintf(
            "Logged in as %s#%s (Id: %s)",
            bot.Username, bot.Discriminator, bot.Id,
        ))
        log.Println("---------")
    }

    var p = presence.Presence{
        Status:       "idle",
        AFK:          false,
        ClientStatus: "mobile",
        Activity: presence.Activity{
            Name: "LO:FI",
            Type: 3,
            URL:  "https://www.youtube.com/....",
        },
    }

    var ping = command.ApplicationCommand{
        Name:        "ping",
        Description: "shows the latency of the client",
        Handler: func(b client.User, ctx command.Context, ops ...interaction.Option) {
            ctx.SendResponse(command.Message{Content: fmt.Sprintf("Pong! %vms", b.Latency)})
        },
    }
