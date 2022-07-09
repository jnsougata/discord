package client

import (
	"encoding/json"
	"fmt"
	"github.com/disgo/core/router"
	"github.com/disgo/core/types"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

type raw struct {
	SessionId   string                 `json:"session_id"`
	Application map[string]interface{} `json:"application"`
}

var events = map[string]interface{}{}
var queue []map[string]interface{}
var commands = map[string]interface{}{}

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

func (c *Client) Queue(apc any) {
	com, _ := json.Marshal(apc)
	m := map[string]interface{}{}
	err := json.Unmarshal(com, &m)
	if err != nil {
		panic(err)
	}
	queue = append(queue, m)
}

func registerCommand(com map[string]interface{}, token string, applicationId string) {
	guildId := com["guild_id"].(string)
	if guildId != "" {
		qualifiedName := "GUILD_" + com["name"].(string) + "_" + guildId
		commands[qualifiedName] = nil
		delete(com, "guild_id")
		r := router.New(
			"POST",
			fmt.Sprintf("/applications/%s/guilds/%s/commands", applicationId, guildId),
			com,
			token,
		)
		r.Request()
	} else {
		qualifiedName := "GLOBAL_" + com["name"].(string)
		commands[qualifiedName] = nil
		r := router.New(
			"POST",
			fmt.Sprintf("/applications/%s/commands", applicationId),
			com,
			token,
		)
		r.Request()
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
				go registerCommand(cmd, token, rc.Application["id"].(string))
			}
		}
		if wsmsg.Op == 10 {
			interval := wsmsg.Data["heartbeat_interval"].(float64)
			c.sendIdentification(conn, token, c.intent)
			go c.keepAlive(conn, int(interval))
		}
		eventHandler(wsmsg.Event, wsmsg.Data)
	}
}

func eventHandler(event string, data map[string]interface{}) {
	if event == "MESSAGE_CREATE" {
		go events[event].(func(message *types.Message))(types.BuildMessage(data))
	}
	if event == "READY" {
		go events[event].(func())()
	}
	if event == "INTERACTION_CREATE" {
		i := types.BuildInteraction(data)
		go events[event].(func(interaction *types.Interaction))(i)
		if i.Type == 1 {
			// interaction ping
		}
		if i.Type == 2 {
			// handle application command interaction
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
