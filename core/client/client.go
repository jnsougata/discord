package client

import (
	"encoding/json"
	"fmt"
	"github.com/disgo/core/models"
	"github.com/disgo/core/router"
	"github.com/disgo/core/types"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	OnReady       = "READY"
	OnMessage     = "MESSAGE_CREATE"
	OnInteraction = "INTERACTION_CREATE"
)

type client struct {
	SessionId   string `json:"session_id"`
	Application struct {
		Id    string  `json:"id"`
		Flags float64 `json:"flags"`
	} `json:"application"`
	User struct {
		Bot           bool    `json:"bot"`
		Avatar        string  `json:"avatar"`
		Discriminator string  `json:"discriminator"`
		Flags         float64 `json:"flags"`
		MFAEnabled    bool    `json:"mfa_enabled"`
		Username      string  `json:"username"`
		Verified      bool    `json:"verified"`
	}
}

var bot *types.User
var execLocked = true
var queue []interface{}
var eventHooks = map[string]interface{}{}
var commandHooks = map[string]interface{}{}

func (c *Client) getGateway() string {
	data, _ := http.Get("https://discord.com/api/gateway")
	var payload map[string]interface{}
	bytes, _ := io.ReadAll(data.Body)
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		panic(err)
	}
	return payload["url"].(string) + "?v=10&encoding=json"
}

type Client struct{ intent int }

func New(intent int) *Client {
	return &Client{intent: intent}
}

func (c *Client) keepAlive(conn *websocket.Conn, dur int) {
	for {
		_ = conn.WriteJSON(map[string]interface{}{"op": 1, "d": 251})
		time.Sleep(time.Duration(dur) * time.Millisecond)
	}
}

func (c *Client) identify(conn *websocket.Conn, Token string, intent int) {
	payload := map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"token":   Token,
			"intents": intent,
			"properties": map[string]interface{}{
				"os":      "linux",
				"browser": "discord.go",
				"device":  "discord.go",
			},
		},
	}
	err := conn.WriteJSON(payload)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) AddHandler(name string, fn interface{}) {
	eventHooks[name] = fn
}

func (c *Client) Queue(apc any, hook interface{}) {
	queue = append(queue, []interface{}{apc, hook})
}

func registerCommand(command any, token string, applicationId string, hook interface{}) {
	var route string

	switch command.(type) {
	case models.SlashCommand:
		c, _ := json.Marshal(command)
		payload := map[string]interface{}{}
		_ = json.Unmarshal(c, &payload)
		guildId := payload["guild_id"].(string)
		switch guildId {
		case "":
			route = fmt.Sprintf("/applications/%s/guilds/%s/commands", applicationId, guildId)
			delete(payload, "guild_id")
		default:
			route = fmt.Sprintf("/applications/%s/commands", applicationId)
		}
		r := router.New("POST", route, payload, token)
		d := map[string]interface{}{}
		body, _ := io.ReadAll(r.Request().Body)
		_ = json.Unmarshal(body, &d)
		_, ok := d["id"]
		if ok {
			commandHooks[d["id"].(string)] = hook
		} else {
			errors, _ := json.Marshal(d["errors"])
			panic(fmt.Sprintf("Failed to register command (%s):\nErrors: %s", payload["name"], errors))
		}
	}
}

func (c *Client) Run(token string) {
	wss := c.getGateway()
	conn, _, err := websocket.DefaultDialer.Dial(wss, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		var wsmsg struct {
			Op       int
			Data     map[string]interface{} `json:"d"`
			Event    string                 `json:"t"`
			Sequence int                    `json:"s"`
		}
		_ = conn.ReadJSON(&wsmsg)

		if wsmsg.Event == OnReady {
			var c client
			b, _ := json.Marshal(wsmsg.Data)
			_ = json.Unmarshal(b, &c)
			for _, cmd := range queue {
				go registerCommand(cmd.([]any)[0].(models.SlashCommand), token, c.Application.Id, cmd.([]any)[1])
			}
			bot = types.BuildUser(wsmsg.Data["user"].(map[string]interface{}))
			execLocked = false

			if _, ok := eventHooks[OnReady]; ok {
				go eventHooks[OnReady].(func(bot *types.User))(bot)

			}
		}
		if wsmsg.Op == 10 {
			interval := wsmsg.Data["heartbeat_interval"].(float64)
			c.identify(conn, token, c.intent)
			go c.keepAlive(conn, int(interval))
		}
		eventHandler(wsmsg.Event, wsmsg.Data)

	}
}

func eventHandler(event string, data map[string]interface{}) {
	if execLocked == true {
		return
	}
	switch event {

	case OnMessage:
		if _, ok := eventHooks[event]; ok {
			eventHook := eventHooks[event].(func(bot *types.User, message *types.Message))
			go eventHook(bot, types.BuildMessage(data))
		}

	case OnInteraction:
		i := types.BuildInteraction(data)

		switch i.Type {

		case 1:
			// interaction ping
		case 2:
			if _, ok := commandHooks[i.Data.Id]; ok {
				commandHook := commandHooks[i.Data.Id].(func(bot *types.User, interaction *types.Interaction, options ...types.Option))
				go commandHook(bot, i, i.Data.Options...)
			}
		case 3:
			// handle component interaction
		case 4:
			// handle auto-complete interaction
		case 5:
			// handle modal submit interaction

		default:
			log.Println("Unknown interaction type: ", i.Type)
		}

	default:
	}
}
