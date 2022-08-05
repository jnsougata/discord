# disgo


# Quick Example

```go
package main

import (
	"fmt"
	"github.com/jnsougata/disgo"
	"os"
)

func main() {
	bot := disgo.Bot(disgo.Intents(), true, disgo.Presence{
		Since:  0,
		Status: "dnd",
		Activity: disgo.Activity{
			Name: "/ping",
			Type: 2,
		},
	})
	bot.OnReady(onReady)
	bot.Run(os.Getenv("DISCORD_TOKEN"))
}

func onReady(b disgo.BotUser) {
	fmt.Println(fmt.Sprintf("Running %s#%s (ID: %s)", b.Username, b.Discriminator, b.Id))
	fmt.Println("-------")
}
```
