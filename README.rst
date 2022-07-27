disgo
-----

Quick Example
-------------
.. code:: go
package main

import (
	"fmt"
	"github.com/jnsougata/disgo/core/bot"
	"github.com/jnsougata/disgo/core/intents"
	"github.com/jnsougata/disgo/core/message"
	"github.com/jnsougata/disgo/core/user"
	"log"
	"os"
)

func main() {
	b := bot.New(intents.All(), false)
	b.AddCommand(cmd, handler)
	b.OnReady(onReady)
	b.OnMessage(onMessage)
	b.Run(os.Getenv("DISCORD_TOKEN"))
}

func onReady(bot user.User) {
	log.Println(fmt.Sprintf(
		"Logged in as %s#%s (Id: %s)",
		bot.Username, bot.Discriminator, bot.ID,
	))
	log.Println("---------")
}

func onMessage(_ user.User, msg message.Message) {
	log.Println(fmt.Sprintf(
		`[%s] %s#%s said: %s`, msg.Author.ID, msg.Author.Username, msg.Author.Discriminator, msg.Content))
}
