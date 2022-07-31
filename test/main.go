package main

import (
	"github.com/jnsougata/disgo"
	"os"
)

func main() {
	b := disgo.Bot(0, false, disgo.Presence{})
	b.Run(os.Getenv("DISCORD_TOKEN"))
}
