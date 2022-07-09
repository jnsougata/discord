package client

import (
	"encoding/json"
	"fmt"
	"github.com/disgo/core/objects"
	"github.com/disgo/core/router"
	"github.com/disgo/core/types"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"
)

type raw struct {
	SessionId   string                 `json:"session_id"`
	Application map[string]interface{} `json:"application"`
}

var events = map[string]interface{}{}
var queue []interface{}
var commands = map[string]interface{}{}
var bot *types.User
var execLocked = true

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
	seq := 0
	for {
		err := conn.WriteJSON(map[string]interface{}{"op": 1, "d": 251})
		seq++
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(dur) * time.Millisecond)
	}
}

func (c *Client) sendIdentification(conn *websocket.Conn, Token string, intent int) {
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
	events[name] = fn
}

func (c *Client) Queue(apc any, handler interface{}) {
	queue = append(queue, []interface{}{apc, handler})
}

func registerCommand(com any, token string, applicationId string, handler interface{}) {
	typePrefix := ""
	if reflect.TypeOf(com) == reflect.TypeOf(objects.SlashCommand{}) {
		typePrefix = "SLASH_"
	}
	c, _ := json.Marshal(com)
	cmd := map[string]interface{}{}
	err := json.Unmarshal(c, &cmd)
	if err != nil {
		panic(err)
	}
	guildId := cmd["guild_id"].(string)
	qualifiedName := ""
	path := ""
	if guildId != "" {
		qualifiedName = typePrefix + "GUILD_" + cmd["name"].(string) + "_" + guildId
		path = fmt.Sprintf("/applications/%s/guilds/%s/commands", applicationId, guildId)
		commands[qualifiedName] = handler
		delete(cmd, "guild_id")
	} else {
		qualifiedName = typePrefix + "GLOBAL_" + cmd["name"].(string)
		commands[qualifiedName] = handler
		path = fmt.Sprintf("/applications/%s/commands", applicationId)
	}
	r := router.New("POST", path, cmd, token)
	r.Request()
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
		err := conn.ReadJSON(&wsmsg)
		if err != nil {
			log.Fatal(err)
		}
		if wsmsg.Event == "READY" {
			var rc raw
			b, _ := json.Marshal(wsmsg.Data)
			err = json.Unmarshal(b, &rc)
			if err != nil {
				panic(err)
			}
			for _, cmd := range queue {
				go registerCommand(
					cmd.([]any)[0].(objects.SlashCommand),
					token,
					rc.Application["id"].(string),
					cmd.([]any)[1],
				)
			}
		}
		if wsmsg.Op == 10 {
			interval := wsmsg.Data["heartbeat_interval"].(float64)
			c.sendIdentification(conn, token, c.intent)
			go c.keepAlive(conn, int(interval))
		}
		eventHandler(wsmsg.Event, wsmsg.Data)
		if wsmsg.Event == "READY" {
			bot = types.BuildUser(wsmsg.Data["user"].(map[string]interface{}))
			execLocked = false
			go events[wsmsg.Event].(func(bot *types.User))(bot)
		}

	}
}

func eventHandler(event string, data map[string]interface{}) {
	if execLocked == true {
		return
	}
	if event == "MESSAGE_CREATE" {
		if _, ok := events[event]; ok {
			go events[event].(func(bot *types.User, message *types.Message))(bot, types.BuildMessage(data))
		}
	}
	if event == "INTERACTION_CREATE" {
		i := types.BuildInteraction(data)
		if i.Type == 1 {
			// interaction ping
		}
		if i.Type == 2 {
			qual := buildQualifiedName(i.GuildID, i.Data.Name, "SLASH")
			if _, ok := commands[qual]; ok {
				go commands[qual].(func(bot *types.User, interaction *types.Interaction))(bot, i)
			}
		}
		if i.Type == 3 {
			// handle component interaction
		}
		if i.Type == 4 {
			// handle auto-complete interaction
		}
		if i.Type == 5 {
			// handle modal submit interaction
		}
	}
}

func buildQualifiedName(guild string, cmdName string, cmdType string) string {
	if guild == "" {
		return cmdType + "_GLOBAL_" + cmdName
	}
	return cmdType + "_GUILD_" + cmdName + "_" + guild
}
