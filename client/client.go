package client

import (
	"encoding/json"
	"github.com/disgo/types"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

var events = map[string]interface{}{}

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

type Client struct {
	intent int
}

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
		if wsmsg.Op == 10 {
			interval := wsmsg.Data["heartbeat_interval"].(float64)
			c.sendIdentification(conn, token, c.intent)
			go c.keepAlive(conn, int(interval))
		}
		if wsmsg.Event == "MESSAGE_CREATE" {
			go events["MESSAGE_CREATE"].(func(message *types.Message))(types.ToMessage(wsmsg.Data))
		}
		if wsmsg.Event == "READY" {
			go events["READY"].(func())()
		}
	}
}