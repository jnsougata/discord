package disgo

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

var islocked = false

type Socket struct {
	Intent       int
	Memoize      bool
	Presence     Presence
	interval     float64
	beatSent     int64
	beatAck      int64
	latency      int64
	self         *BotUser
	guilds       map[string]*Guild
	queue        []ApplicationCommand
	eventHooks   map[string]interface{}
	commandHooks map[string]interface{}
	sequence     int
	sessionId    string
}

func (sock *Socket) getGateway() string {
	data, _ := http.Get("https://discord.com/api/gateway")
	var payload map[string]string
	bytes, _ := io.ReadAll(data.Body)
	_ = json.Unmarshal(bytes, &payload)
	return fmt.Sprintf("%s?v=10&encoding=json", payload["url"])
}

func (sock *Socket) keepAlive(conn *websocket.Conn, dur int) {
	for {
		_ = conn.WriteJSON(map[string]interface{}{"op": 1, "d": 251})
		sock.beatSent = time.Now().UnixMilli()
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
	}
	if sock.Presence.Activity.Name != "" {
		d.Presence = sock.Presence.Marshal()
	}
	if sock.Presence.OnMobile {
		d.Properties.Browser = "Discord iOS"
	}
	payload := map[string]interface{}{"op": 2, "d": d}
	_ = conn.WriteJSON(payload)
}

func (sock *Socket) AddHandler(name string, handler interface{}) {
	if sock.eventHooks == nil {
		sock.eventHooks = make(map[string]interface{})
	}
	sock.eventHooks[name] = handler
}

func (sock *Socket) AddToQueue(commands ...ApplicationCommand) {
	for _, com := range commands {
		sock.queue = append(sock.queue, com)
	}
}

func (sock *Socket) registerCommand(com ApplicationCommand, token string, applicationId string) {
	var route string
	data, hook, guildId := com.Marshal()
	if guildId != 0 {
		route = fmt.Sprintf("/applications/%s/guilds/%v/commands", applicationId, guildId)
	} else {
		route = fmt.Sprintf("/applications/%s/commands", applicationId)
	}
	r := MinimalReq("POST", route, data, token)
	d := map[string]interface{}{}
	body, _ := io.ReadAll(r.Request().Body)
	_ = json.Unmarshal(body, &d)
	_, ok := d["id"]
	if ok {
		sock.commandHooks[d["id"].(string)] = hook
	} else {
		log.Fatal(
			fmt.Sprintf("Failed to register command {%s}. Reason: %s", com.Name, d["message"]))
	}
}

func (sock *Socket) Run(token string) {
	sock.guilds = make(map[string]*Guild)
	sock.commandHooks = make(map[string]interface{})
	wss := sock.getGateway()
	conn, _, err := websocket.DefaultDialer.Dial(wss, nil)
	if err != nil {
		log.Println(err)
	}
	for {
		var wsmsg struct {
			Op       int                    `json:"op"`
			Event    string                 `json:"t"`
			Sequence int                    `json:"s"`
			Data     map[string]interface{} `json:"d"`
		}
		_ = conn.ReadJSON(&wsmsg)
		var runtime struct {
			SessionId   string `json:"session_id"`
			Sequence    int
			Application struct {
				Id    string  `json:"id"`
				Flags float64 `json:"flags"`
			} `json:"application"`
		}
		sock.sequence = wsmsg.Sequence
		if wsmsg.Event == OnReady {
			ba, _ := json.Marshal(wsmsg.Data)
			_ = json.Unmarshal(ba, &runtime)
			for _, cmd := range sock.queue {
				go sock.registerCommand(cmd, token, runtime.Application.Id)
			}
			sock.self = Unmarshal(wsmsg.Data["user"].(map[string]interface{}))
			sock.self.Latency = sock.latency
			sock.self.IsReady = true
			islocked = false
			if hook, ok := sock.eventHooks[OnReady]; ok {
				go hook.(func(bot BotUser))(*sock.self)
			}
		}
		if wsmsg.Op == 10 {
			sock.interval = wsmsg.Data["heartbeat_interval"].(float64)
			sock.identify(conn, token, sock.Intent)
			go sock.keepAlive(conn, int(sock.interval))
		}
		if wsmsg.Op == 11 {
			sock.beatAck = time.Now().UnixMilli()
			sock.latency = sock.beatAck - sock.beatSent
			if sock.self != nil {
				sock.self.Latency = sock.latency
			}
		}
		if wsmsg.Op == 7 {
			_ = conn.WriteJSON(map[string]interface{}{
				"op": 6,
				"d": map[string]interface{}{
					"token":      token,
					"sequence":   sock.sequence,
					"session_id": sock.sessionId,
				},
			})
		}
		if wsmsg.Op == 9 {
			sock.identify(conn, token, sock.Intent)
		}
		sock.eventHandler(wsmsg.Event, wsmsg.Data)
		if h, ok := sock.eventHooks[OnSocketReceive]; ok {
			handler := h.(func(d map[string]interface{}))
			go handler(wsmsg.Data)
		}
		if wsmsg.Event == "GUILD_CREATE" {
			gld := UnmarshalGuild(wsmsg.Data)
			gld.UnmarshalRoles(wsmsg.Data["roles"].([]interface{}))
			gld.UnmarshalChannels(wsmsg.Data["channels"].([]interface{}))
			sock.guilds[gld.Id] = gld
			sock.self.Guilds = sock.guilds
			if sock.Memoize {
				requestMembers(conn, gld.Id)
			}
		}
		if wsmsg.Event == "GUILD_MEMBERS_CHUNK" {
			id := wsmsg.Data["guild_id"].(string)
			sock.guilds[id].UnmarshalMembers(wsmsg.Data["members"].([]interface{}))
		}
	}
}

func (sock *Socket) eventHandler(event string, data map[string]interface{}) {
	if islocked {
		return
	}
	switch event {

	case OnMessageCreate:
		if event, ok := sock.eventHooks[event]; ok {
			hook := event.(func(bot BotUser, message Message))
			go hook(*sock.self, *UnmarshalMessage(data))
		}

	case OnGuildCreate:
		if event, ok := sock.eventHooks[event]; ok {
			hook := event.(func(bot BotUser, guild Guild))
			go hook(*sock.self, *UnmarshalGuild(data))
		}

	case OnInteractionCreate:
		ctx := UnmarshalContext(data)
		if event, ok := sock.eventHooks[event]; ok {
			hook := event.(func(bot BotUser, ctx Context))
			go hook(*sock.self, *ctx)
		}
		switch ctx.Type {
		case 1:
			// interaction ping
		case 2:
			if ev, ok := sock.commandHooks[ctx.Data.Id]; ok {
				hook := ev.(func(bot BotUser, ctx Context, ops ...SlashCommandOption))
				go hook(*sock.self, *ctx, ctx.Data.Options...)
			}
		case 3:
			var compdata ComponentData
			da, _ := json.Marshal(data["data"].(map[string]interface{}))
			_ = json.Unmarshal(da, &compdata)
			ctx.ComponentData = compdata
			switch ctx.ComponentData.ComponentType {
			case 2:
				cb, ok := callbackTasks[ctx.ComponentData.CustomId]
				if ok {
					callback := cb.(func(b BotUser, ctx Context))
					go callback(*sock.self, *ctx)
				}
			case 3:
				cb, ok := callbackTasks[ctx.ComponentData.CustomId]
				if ok {
					callback := cb.(func(b BotUser, ctx Context, values ...string))
					go callback(*sock.self, *ctx, ctx.ComponentData.Values...)
				}
			}
			tmp, ok := timeoutTasks[ctx.ComponentData.CustomId]
			if ok {
				onTimeoutHandler := tmp[1].(func(b BotUser, ctx Context))
				duration := tmp[0].(float64)
				delete(timeoutTasks, ctx.ComponentData.CustomId)
				go scheduleTimeoutTask(duration, *sock.self, *ctx, onTimeoutHandler)
			}
		case 4:
			// handle auto-complete interaction
		case 5:
			callback, ok := callbackTasks[ctx.ComponentData.CustomId]
			if ok {
				go callback.(func(b BotUser, ctx *Context))(*sock.self, ctx)
				delete(callbackTasks, ctx.ComponentData.CustomId)
			}
		default:
			log.Println("Unknown interaction type: ", ctx.Type)
		}
	default:
	}
}

func scheduleTimeoutTask(timeout float64, user BotUser, ctx Context,
	handler func(bot BotUser, ctx Context)) {
	time.Sleep(time.Duration(timeout) * time.Second)
	handler(user, ctx)
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
