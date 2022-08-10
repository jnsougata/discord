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
var cachedGuilds = map[string]*Guild{}

// ws is a Discord websocket connection,
// responsible for handling all ws events
type ws struct {
	sequence     int
	intent       int
	memoize      bool
	beatSent     int64
	beatAck      int64
	latency      int64
	sessionId    string
	interval     float64
	presence     Presence
	self         *BotUser
	secret       string
	queue        []ApplicationCommand
	eventHooks   map[string]interface{}
	commandHooks map[string]interface{}
}

func (sock *ws) getGateway() string {
	data, err := http.Get("https://discord.com/api/gateway")
	if err != nil {
		panic(err)
	}
	var payload map[string]string
	bytes, _ := io.ReadAll(data.Body)
	_ = json.Unmarshal(bytes, &payload)
	return fmt.Sprintf("%s?v=%s&encoding=json", payload["url"], "10")
}

func (sock *ws) keepAlive(conn *websocket.Conn, dur int) {
	for {
		_ = conn.WriteJSON(map[string]interface{}{"op": 1, "d": 251})
		sock.beatSent = time.Now().UnixMilli()
		time.Sleep(time.Duration(dur) * time.Millisecond)
	}
}

func (sock *ws) identify(conn *websocket.Conn, intent int) {

	d := map[string]interface{}{
		"token":   sock.secret,
		"intents": intent,
	}
	d["properties"] = map[string]string{
		"os":      "linux",
		"browser": "disgo",
		"device":  "disgo",
	}
	if sock.presence.Activity.Name != "" {
		d["presence"] = sock.presence.Marshal()
	}
	if sock.presence.OnMobile {
		d["properties"].(map[string]string)["browser"] = "Discord iOS"
	}
	payload := map[string]interface{}{"op": 2, "d": d}
	_ = conn.WriteJSON(payload)
}

func (sock *ws) AddHandler(name string, handler interface{}) {
	if sock.eventHooks == nil {
		sock.eventHooks = make(map[string]interface{})
	}
	sock.eventHooks[name] = handler
}

func (sock *ws) AddToQueue(commands ...ApplicationCommand) {
	for _, com := range commands {
		sock.queue = append(sock.queue, com)
	}
}

func (sock *ws) registerCommand(com ApplicationCommand, applicationId string) {
	var route string
	data, hook, guildId := com.marshal()
	if guildId != 0 {
		route = fmt.Sprintf("/applications/%s/guilds/%v/commands", applicationId, guildId)
	} else {
		route = fmt.Sprintf("/applications/%s/commands", applicationId)
	}
	r := minimalReq("POST", route, data, sock.secret)
	d := map[string]interface{}{}
	body, _ := io.ReadAll(r.fire().Body)
	_ = json.Unmarshal(body, &d)
	_, ok := d["id"]
	if ok {
		commandId := d["id"].(string)
		sock.commandHooks[commandId] = hook
		sc, there := subcommandBucket[com.uniqueId]
		if there {
			subcommandBucket[commandId] = sc
			delete(subcommandBucket, com.uniqueId)
		}
		scg, there := groupBucket[com.uniqueId]
		if there {
			groupBucket[commandId] = scg
			delete(groupBucket, com.uniqueId)
		}
	} else {
		log.Fatal(
			fmt.Sprintf("Failed to register command {%s}.\nMessage: %s", com.Name, d["message"]))
	}
}

func (sock *ws) Run(token string) {
	sock.secret = token
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
		if wsmsg.Event == onReady {
			ba, _ := json.Marshal(wsmsg.Data)
			_ = json.Unmarshal(ba, &runtime)
			for _, cmd := range sock.queue {
				go sock.registerCommand(cmd, runtime.Application.Id)
			}
			sock.self = Converter{
				token:   sock.secret,
				payload: wsmsg.Data["user"].(map[string]interface{}),
			}.Bot()
			sock.self.Latency = sock.latency
			sock.self.IsReady = true
			islocked = false
			if hook, ok := sock.eventHooks[onReady]; ok {
				go hook.(func(bot BotUser))(*sock.self)
			}
		}
		if wsmsg.Op == 10 {
			sock.interval = wsmsg.Data["heartbeat_interval"].(float64)
			sock.identify(conn, sock.intent)
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
			sock.identify(conn, sock.intent)
		}
		sock.eventHandler(wsmsg.Event, wsmsg.Data)
		if h, ok := sock.eventHooks[onSocketReceive]; ok {
			handler := h.(func(d map[string]interface{}))
			go handler(wsmsg.Data)
		}
		if wsmsg.Event == onGuildCreate {
			g := Converter{token: token, payload: wsmsg.Data}.Guild()
			g.clientId = sock.self.Id
			cachedGuilds[g.Id] = g
			sock.self.Guilds = cachedGuilds
			if sock.memoize {
				requestMembers(conn, g.Id)
			}
		}
		if wsmsg.Event == "GUILD_MEMBERS_CHUNK" {
			cachedGuilds[wsmsg.Data["guild_id"].(string)].unmarshalMembers(wsmsg.Data["members"].([]interface{}))
		}
	}
}

func (sock *ws) eventHandler(dispatch string, data map[string]interface{}) {
	if islocked {
		return
	}
	converter := Converter{token: sock.secret, payload: data}

	switch dispatch {

	case onMessageCreate:
		if event, ok := sock.eventHooks[dispatch]; ok {
			hook := event.(func(bot BotUser, message Message))
			go hook(*sock.self, *converter.Message())
		}

	case onGuildCreate:
		if event, ok := sock.eventHooks[dispatch]; ok {
			hook := event.(func(bot BotUser, guild Guild))
			go hook(*sock.self, *converter.Guild())
		}
	case onGuildDelete:
		if event, ok := sock.eventHooks[dispatch]; ok {
			guild := converter.Guild()
			hook := event.(func(bot BotUser, guild Guild))
			cached, exists := cachedGuilds[guild.Id]
			if exists {
				go hook(*sock.self, *cached)
				delete(cachedGuilds, guild.Id)
			} else {
				go hook(*sock.self, *guild)
			}
		}
	case onInteractionCreate:
		ctx := unmarshalContext(data)
		ctx.raw = data
		ctx.token = sock.secret
		if event, ok := sock.eventHooks[dispatch]; ok {
			hook := event.(func(bot BotUser, ctx Context))
			go hook(*sock.self, *ctx)
		}
		switch ctx.Type {
		case 1:
			// interaction ping
		case 2:
			if task, ok := sock.commandHooks[ctx.Data.Id]; ok {
				switch ctx.Data.Type {
				case 1:
					type subOption struct {
						Name    string   `json:"name"`
						Options []Option `json:"options"`
						Type    int      `json:"type"`
					}
					for _, option := range ctx.Data.Options {
						if int(option.Type) == 1 {
							subOptions := map[string]Option{}
							d := ctx.raw["data"].(map[string]interface{})["options"].([]interface{})
							ds, _ := json.Marshal(d)
							var so []subOption
							_ = json.Unmarshal(ds, &so)
							for _, opt := range so[0].Options {
								subOptions[opt.Name] = opt
							}
							if hook, ok := subcommandBucket[ctx.Data.Id]; ok {
								scTask, exists := hook.(map[string]interface{})[option.Name]
								if exists {
									hook := scTask.(func(bot BotUser, ctx Context, options map[string]Option))
									go hook(*sock.self, *ctx, subOptions)
								}
							}
						}
					}
					hook := task.(func(bot BotUser, ctx Context, ops map[string]Option))
					options := map[string]Option{}
					for _, option := range ctx.Data.Options {
						options[option.Name] = option
					}
					if hook != nil {
						go hook(*sock.self, *ctx, options)
					}
				case 2:
					target := ctx.Data.TargetId
					rud := ctx.Data.Resolved["users"].(map[string]interface{})[target]
					ctx.TargetUser = *Converter{token: sock.secret, payload: rud.(map[string]interface{})}.User()
					hook := task.(func(bot BotUser, ctx Context, _ map[string]Option))
					go hook(*sock.self, *ctx, nil)
				case 3:
					target := ctx.Data.TargetId
					rmd := ctx.Data.Resolved["messages"].(map[string]interface{})[target]
					ctx.TargetMessage = *Converter{token: sock.secret, payload: rmd.(map[string]interface{})}.Message()
					hook := task.(func(bot BotUser, ctx Context, _ map[string]Option))
					go hook(*sock.self, *ctx, nil)
				}
			} else {
				log.Printf("ApplicationCommand (%s) is not implemented.", ctx.Data.Id)
			}
		case 3:
			var compdata ComponentData
			da, _ := json.Marshal(data["data"].(map[string]interface{}))
			_ = json.Unmarshal(da, &compdata)
			ctx.componentData = compdata
			switch ctx.componentData.ComponentType {
			case 2:
				cb, ok := callbackTasks[ctx.componentData.CustomId]
				if ok {
					callback := cb.(func(b BotUser, ctx Context))
					go callback(*sock.self, *ctx)
				}
			case 3:
				cb, ok := callbackTasks[ctx.componentData.CustomId]
				if ok {
					callback := cb.(func(b BotUser, ctx Context, values ...string))
					go callback(*sock.self, *ctx, ctx.componentData.Values...)
				}
			}
			tmp, ok := timeoutTasks[ctx.componentData.CustomId]
			if ok {
				onTimeoutHandler := tmp[1].(func(b BotUser, ctx Context))
				duration := tmp[0].(float64)
				delete(timeoutTasks, ctx.componentData.CustomId)
				go scheduleTimeoutTask(duration, *sock.self, *ctx, onTimeoutHandler)
			}
		case 4:
			// handle auto-complete interaction
		case 5:
			callback, ok := callbackTasks[ctx.componentData.CustomId]
			if ok {
				go callback.(func(b BotUser, ctx *Context))(*sock.self, ctx)
				delete(callbackTasks, ctx.componentData.CustomId)
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
