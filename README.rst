disgo
-----

Quick Example
-------------

.. code:: go

    package main

    import (
        "fmt"
        "github.com/jnsougata/disgo/client"
        "github.com/jnsougata/disgo/intents"
        "github.com/jnsougata/disgo/presence"
        "log"
        "os"
    )

    func main() {
        d := Disgo(intents.All(), true)
        d.AddCommands(toast)
        d.OnReady(onReady)
        d.SetPresence(pr)
        d.OnSocketReceive(onSocketReceive)
        d.Run(os.Getenv("DISCORD_TOKEN"))
    }

    func onReady(bot client.User) {
        log.Println(fmt.Sprintf(
            "Logged in as %s#%s (Id: %s)",
            bot.Username, bot.Discriminator, bot.Id,
        ))
        log.Println("---------")
    }

    func onSocketReceive(_ map[string]interface{}) {
        //
    }

    var pr = presence.Presence{
        Status:       "idle",
        AFK:          false,
        ClientStatus: "mobile",
        Activity: presence.Activity{
            Name: "LO:FI",
            Type: 3,
            URL:  "https://www.youtube.com/watch?v=e97w-GHsRMY",
        },
    }

