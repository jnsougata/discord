package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/component"
	"github.com/jnsougata/disgo/core/consts"
	"github.com/jnsougata/disgo/core/guild"
	"github.com/jnsougata/disgo/core/interaction"
	"github.com/jnsougata/disgo/core/message"
	"github.com/jnsougata/disgo/core/router"
	"github.com/jnsougata/disgo/core/user"
	"io"
	"log"
	"net/http"
	"time"
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

var bot *user.User
var execLocked = true
var queue []command.ApplicationCommand
var eventHooks = map[string]interface{}{}
var commandHooks = map[string]interface{}{}

func (c *Client) getGateway() string {
	data, _ := http.Get("https://discord.com/api/gateway")
	var payload map[string]string
	bytes, _ := io.ReadAll(data.Body)
	_ = json.Unmarshal(bytes, &payload)
	return fmt.Sprintf("%s?v=10&encoding=json", payload["url"])
}

type Client struct {
	Intent  int
	Memoize bool
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

func (c *Client) Queue(apc command.ApplicationCommand) {
	queue = append(queue, apc)
}

func registerCommand(com command.ApplicationCommand, token string, applicationId string) {
	var route string
	data, hook, guildId := com.ToData()
	if guildId != 0 {
		route = fmt.Sprintf("/applications/%s/guilds/%v/commands", applicationId, guildId)
	} else {
		route = fmt.Sprintf("/applications/%s/commands", applicationId)
	}
	r := router.Minimal("POST", route, data, token)
	d := map[string]interface{}{}
	body, _ := io.ReadAll(r.Request().Body)
	_ = json.Unmarshal(body, &d)
	_, ok := d["id"]
	if ok {
		commandHooks[d["id"].(string)] = hook
	} else {
		log.Fatal(
			fmt.Sprintf("Failed to register command {%s}. Code %s", com.Name, d["message"]))
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
		if wsmsg.Event == consts.OnReady {
			var c client
			b, _ := json.Marshal(wsmsg.Data)
			_ = json.Unmarshal(b, &c)
			for _, cmd := range queue {
				go registerCommand(cmd, token, c.Application.Id)
			}
			bot = user.FromData(wsmsg.Data["user"].(map[string]interface{}))
			execLocked = false
			if _, ok := eventHooks[consts.OnReady]; ok {
				go eventHooks[consts.OnReady].(func(bot user.User))(*bot)
			}
		}
		if wsmsg.Op == 10 {
			interval := wsmsg.Data["heartbeat_interval"].(float64)
			c.identify(conn, token, c.Intent)
			go c.keepAlive(conn, int(interval))
		}
		eventHandler(wsmsg.Event, wsmsg.Data)
		if h, ok := eventHooks[consts.OnSocketReceive]; ok {
			handler := h.(func(map[string]interface{}))
			go handler(wsmsg.Data)
		}
	}
}

func eventHandler(event string, data map[string]interface{}) {
	if execLocked == true {
		return
	}
	switch event {

	case consts.OnMessageCreate:
		if _, ok := eventHooks[event]; ok {
			eventHook := eventHooks[event].(func(bot user.User, message message.Message))
			go eventHook(*bot, *message.FromData(data))
		}

	case consts.OnGuildCreate:
		if _, ok := eventHooks[event]; ok {
			eventHook := eventHooks[event].(func(bot user.User, guild guild.Guild))
			go eventHook(*bot, *guild.FromData(data))
		}

	case consts.OnInteractionCreate:

		i := interaction.FromData(data)

		if _, ok := eventHooks[event]; ok {
			eventHook := eventHooks[event].(func(bot user.User, ctx interaction.Interaction))
			go eventHook(*bot, *i)
		}
		switch i.Type {
		case 1:
			// interaction ping
		case 2:
			ctx := command.Context{
				Id:            i.Id,
				Token:         i.Token,
				ApplicationId: i.ApplicationId,
			}
			if _, ok := commandHooks[i.Data.Id]; ok {
				hook := commandHooks[i.Data.Id].(func(bot user.User, ctx command.Context, ops ...interaction.Option))
				go hook(*bot, ctx, i.Data.Options...)
			}
		case 3:
			factory := component.CallbackTasks
			ctx := component.FromData(data)
			switch ctx.Data.ComponentType {
			case 2:
				cb, ok := factory[ctx.Data.CustomId]
				if ok {
					callback := cb.(func(b user.User, ctx component.Context))
					go callback(*bot, *ctx)
				}
			case 3:
				cb, ok := factory[ctx.Data.CustomId]
				if ok {
					callback := cb.(func(b user.User, ctx component.Context, values ...string))
					go callback(*bot, *ctx, ctx.Data.Values...)
				}
			}
			tmp, ok := component.TimeoutTasks[ctx.Data.CustomId]
			if ok {
				onTimeoutHandler := tmp[1].(func(b user.User, i component.Context))
				duration := tmp[0].(float64)
				delete(component.TimeoutTasks, ctx.Data.CustomId)
				go scheduleTimeoutTask(duration, *bot, *ctx, onTimeoutHandler)
			}
		case 4:
			// handle auto-complete interaction
		case 5:
			cctx := component.FromData(data)
			callback, ok := component.CallbackTasks[cctx.Data.CustomId]
			if ok {
				go callback.(func(b user.User, cctx component.Context))(*bot, *cctx)
				delete(component.CallbackTasks, cctx.Data.CustomId)
			}
		default:
			log.Println("Unknown interaction type: ", i.Type)
		}
	default:
	}
}

func scheduleTimeoutTask(timeout float64, user user.User, cctx component.Context,
	handler func(bot user.User, cctx component.Context)) {
	time.Sleep(time.Duration(timeout) * time.Second)
	handler(user, cctx)
}
