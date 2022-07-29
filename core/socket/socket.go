package socket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/component"
	"github.com/jnsougata/disgo/core/consts"
	"github.com/jnsougata/disgo/core/guild"
	"github.com/jnsougata/disgo/core/interaction"
	"github.com/jnsougata/disgo/core/member"
	"github.com/jnsougata/disgo/core/message"
	"github.com/jnsougata/disgo/core/presence"
	"github.com/jnsougata/disgo/core/router"
	"github.com/jnsougata/disgo/core/user"
	"io"
	"log"
	"net/http"
	"time"
)

type runtime struct {
	SessionId   string `json:"session_id"`
	Sequence    int
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

var bot *user.Bot
var latency int64
var beatSent int64
var interval float64
var execLocked = true
var beatReceived int64
var presenceStruct presence.Presence
var queue []command.ApplicationCommand
var eventHooks = map[string]interface{}{}
var commandHooks = map[string]interface{}{}
var guilds = map[string]guild.Guild{}

func (sock *Socket) getGateway() string {
	data, _ := http.Get("https://discord.com/api/gateway")
	var payload map[string]string
	bytes, _ := io.ReadAll(data.Body)
	_ = json.Unmarshal(bytes, &payload)
	return fmt.Sprintf("%s?v=10&encoding=json", payload["url"])
}

type Socket struct {
	Intent  int
	Memoize bool
}

func (sock *Socket) keepAlive(conn *websocket.Conn, dur int) {
	for {
		_ = conn.WriteJSON(map[string]interface{}{"op": 1, "d": 251})
		beatSent = time.Now().UnixMilli()
		time.Sleep(time.Duration(dur) * time.Millisecond)
	}
}

func (sock *Socket) identify(conn *websocket.Conn, Token string, intent int) {
	type properties struct {
		Os      string `json:"os"`
		Browser string `json:"browser"`
		Device  string `json:"device"`
	}
	type data struct {
		Token      string                 `json:"token"`
		Intents    int                    `json:"intents"`
		Properties properties             `json:"properties"`
		Presence   map[string]interface{} `json:"presence"`
	}
	d := data{
		Token:   Token,
		Intents: intent,
		Properties: properties{
			Os:      "linux",
			Browser: "disgo",
			Device:  "disgo",
		},
		Presence: presenceStruct.ToData(),
	}
	if presenceStruct.ClientStatus == "mobile" {
		d.Properties.Browser = "Discord iOS"
	}
	payload := map[string]interface{}{"op": 2, "d": d}
	_ = conn.WriteJSON(payload)
}

func (sock *Socket) AddHandler(name string, fn interface{}) {
	eventHooks[name] = fn
}

func (sock *Socket) Queue(commands ...command.ApplicationCommand) {
	for _, com := range commands {
		queue = append(queue, com)
	}
}

func (sock *Socket) StorePresenceData(pr presence.Presence) {
	presenceStruct = pr
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

func (sock *Socket) Run(token string) {
	wss := sock.getGateway()
	conn, _, err := websocket.DefaultDialer.Dial(wss, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		var wsmsg struct {
			Op       int                    `json:"op"`
			Event    string                 `json:"t"`
			Sequence int                    `json:"s"`
			Data     map[string]interface{} `json:"d"`
		}
		_ = conn.ReadJSON(&wsmsg)
		var r runtime
		r.Sequence = wsmsg.Sequence
		if wsmsg.Event == consts.OnReady {
			b, _ := json.Marshal(wsmsg.Data)
			_ = json.Unmarshal(b, &r)
			for _, cmd := range queue {
				go registerCommand(cmd, token, r.Application.Id)
			}
			bot = user.MakeBot(wsmsg.Data["user"].(map[string]interface{}))
			bot.Latency = latency
			execLocked = false
			if _, ok := eventHooks[consts.OnReady]; ok {
				go eventHooks[consts.OnReady].(func(bot user.Bot))(*bot)
			}
		}
		if wsmsg.Op == 10 {
			interval = wsmsg.Data["heartbeat_interval"].(float64)
			sock.identify(conn, token, sock.Intent)
			go sock.keepAlive(conn, int(interval))
		}
		if wsmsg.Op == 11 {
			beatReceived = time.Now().UnixMilli()
			latency = beatReceived - beatSent
			if bot != nil {
				bot.Latency = latency
			}
		}
		if wsmsg.Op == 7 {
			_ = conn.WriteJSON(map[string]interface{}{
				"op": 6,
				"d": map[string]interface{}{
					"token":      token,
					"session_id": r.SessionId,
					"sequence":   r.Sequence,
				},
			})
		}
		if wsmsg.Op == 9 {
			sock.identify(conn, token, sock.Intent)
		}
		eventHandler(wsmsg.Event, wsmsg.Data)
		if h, ok := eventHooks[consts.OnSocketReceive]; ok {
			handler := h.(func(d map[string]interface{}))
			go handler(wsmsg.Data)
		}
		if wsmsg.Event == "GUILD_CREATE" {
			g := guild.FromData(wsmsg.Data)
			guilds[g.Id] = *g
			if sock.Memoize {
				requestMembers(conn, g.Id)
			}
		}
		if wsmsg.Event == "GUILD_MEMBERS_CHUNK" {
			execLocked = true
			var ms []member.Member
			ma, _ := json.Marshal(wsmsg.Data["members"])
			_ = json.Unmarshal(ma, &ms)
			id := wsmsg.Data["guild_id"].(string)
			g := guilds[id]
			g.Members = ms
			execLocked = false
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
			eventHook := eventHooks[event].(func(bot user.Bot, message message.Message))
			go eventHook(*bot, *message.FromData(data))
		}

	case consts.OnGuildCreate:
		if _, ok := eventHooks[event]; ok {
			eventHook := eventHooks[event].(func(bot user.Bot, guild guild.Guild))
			go eventHook(*bot, *guild.FromData(data))
		}

	case consts.OnInteractionCreate:

		i := interaction.FromData(data)

		if _, ok := eventHooks[event]; ok {
			eventHook := eventHooks[event].(func(bot user.Bot, ctx interaction.Interaction))
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
				hook := commandHooks[i.Data.Id].(func(bot user.Bot, ctx command.Context, ops ...interaction.Option))
				go hook(*bot, ctx, i.Data.Options...)
			}
		case 3:
			factory := component.CallbackTasks
			cc := component.FromData(data)
			switch cc.Data.ComponentType {
			case 2:
				cb, ok := factory[cc.Data.CustomId]
				if ok {
					callback := cb.(func(b user.Bot, ctx component.Context))
					go callback(*bot, *cc)
				}
			case 3:
				cb, ok := factory[cc.Data.CustomId]
				if ok {
					callback := cb.(func(b user.Bot, ctx component.Context, values ...string))
					go callback(*bot, *cc, cc.Data.Values...)
				}
			}
			tmp, ok := component.TimeoutTasks[cc.Data.CustomId]
			if ok {
				onTimeoutHandler := tmp[1].(func(b user.Bot, i component.Context))
				duration := tmp[0].(float64)
				delete(component.TimeoutTasks, cc.Data.CustomId)
				go scheduleTimeoutTask(duration, *bot, *cc, onTimeoutHandler)
			}
		case 4:
			// handle auto-complete interaction
		case 5:
			cc := component.FromData(data)
			callback, ok := component.CallbackTasks[cc.Data.CustomId]
			if ok {
				go callback.(func(b user.Bot, cc component.Context))(*bot, *cc)
				delete(component.CallbackTasks, cc.Data.CustomId)
			}
		default:
			log.Println("Unknown interaction type: ", i.Type)
		}
	default:
	}
}

func scheduleTimeoutTask(timeout float64, user user.Bot, cc component.Context,
	handler func(bot user.Bot, cc component.Context)) {
	time.Sleep(time.Duration(timeout) * time.Second)
	handler(user, cc)
}

func requestMembers(conn *websocket.Conn, guildId string) {
	_ = conn.WriteJSON(map[string]interface{}{
		"op": 8,
		"d": map[string]interface{}{
			"guild_id": guildId,
			"query":    "",
			"limit":    0,
		},
	})
}
