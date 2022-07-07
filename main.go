package main

import (
	"github.com/disgo/bot"
	"github.com/disgo/types"
	"log"
	"os"
)

func main() {
	b := bot.New(33283)
	b.OnMessage(OnMessage)
	b.OnReady(OnReady)
	b.Run(os.Getenv("DISCORD_TOKEN"))
}

func OnMessage(message *types.Message) {
	log.Println(message.Content)
}
func OnReady() {
	log.Println("[-------- READY --------]")
}
