package discord

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"time"
)

var s = state{}

type sharding struct {
	URL               string `json:"url"`
	Shards            int    `json:"shards"`
	SessionStartLimit struct {
		Total          int `json:"total"`
		Remaining      int `json:"remaining"`
		Reset          int `json:"reset_after"`
		MaxConcurrency int `json:"max_concurrency"`
	} `json:"session_start_limit"`
}

type runtime struct {
	SessionId        string      `json:"session_id"`
	Shard            []int       `json:"shard"`
	ResumeGatewayURL string      `json:"resume_gateway_url"`
	Guilds           interface{} `json:"guilds"`
	Application      struct {
		Id    string  `json:"id"`
		Flags float64 `json:"flags"`
	} `json:"application"`
}

type ws struct {
	locked           bool
	sequence         int
	intent           int
	memoize          bool
	beatSent         int64
	beatAck          int64
	latency          int64
	sessionId        string
	interval         float64
	presence         Presence
	self             *Bot
	secret           string
	queue            []Command
	commands         map[string]interface{}
	listeners        Listeners
	lazyGuilds       map[string]bool
	shardingMatrices sharding
	runtimeMatrices  runtime
}

func (sock *ws) gateway() string {
	version := 11
	encoding := "json"
	info := sharding{}
	r := minimalReq("GET", "/gateway/bot", nil, sock.secret)
	resp := r.fire()
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &info)
	sock.shardingMatrices = info
	return fmt.Sprintf("%s?v=%d&encoding=%s", info.URL, version, encoding)
}

func (sock *ws) keepAlive(conn *websocket.Conn, dur int) {
	for {
		_ = conn.WriteJSON(map[string]interface{}{"op": 1, "d": 251})
		sock.beatSent = time.Now().UnixMilli()
		time.Sleep(time.Duration(dur) * time.Millisecond)
	}
}

func (sock *ws) identify(conn *websocket.Conn, intent int) {
	identification := map[string]interface{}{
		"token":   sock.secret,
		"intents": intent,
	}
	identification["properties"] = map[string]string{
		"os":      "linux",
		"browser": "discord",
		"device":  "discord",
	}
	if sock.presence.Activity.Name != "" {
		identification["presence"] = sock.presence.Marshal()
	}
	if sock.presence.OnMobile {
		identification["properties"].(map[string]string)["browser"] = "Discord iOS"
	}
	payload := map[string]interface{}{"op": 2, "d": identification}
	_ = conn.WriteJSON(payload)
}

func (sock *ws) registerCommand(com Command, applicationId string) {
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
		sock.commands[commandId] = hook
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
		log.Fatal(fmt.Sprintf("Failed to register command {%s}. Message: %s", com.Name, d["message"]))
	}
}

func (sock *ws) requestMembers(conn *websocket.Conn, guildId string) {
	_ = conn.WriteJSON(map[string]interface{}{
		"op": 8,
		"d": map[string]interface{}{
			"guild_id": guildId,
			"query":    "",
			"limit":    0,
		},
	})
}

func (sock *ws) run(token string) {
	s.Guilds = make(map[string]*Guild)
	s.Users = make(map[string]*User)
	conn, _, err := websocket.DefaultDialer.Dial(sock.gateway(), nil)
	if err != nil {
		fmt.Println(err)
	}
	for {
		var wsmsg struct {
			OP       int                    `json:"op"`
			Event    string                 `json:"t"`
			Sequence int                    `json:"s"`
			Data     map[string]interface{} `json:"d"`
		}
		err = conn.ReadJSON(&wsmsg)
		if err != nil {
			fmt.Println(err)
		}
		sock.sequence = wsmsg.Sequence
		if sock.listeners.OnSocketReceive != nil {
			go sock.listeners.OnSocketReceive(wsmsg)
		}
		r := runtime{}
		switch wsmsg.Event {
		case string(OnReady):
			sock.lazyGuilds = map[string]bool{}
			ba, _ := json.Marshal(wsmsg.Data)
			_ = json.Unmarshal(ba, &r)
			sock.runtimeMatrices = r
			for _, guild := range r.Guilds.([]interface{}) {
				id := guild.(map[string]interface{})["id"].(string)
				sock.lazyGuilds[id] = true
			}
			for _, cmd := range sock.queue {
				go sock.registerCommand(cmd, r.Application.Id)
			}
			sock.self = Converter{
				token:   sock.secret,
				payload: wsmsg.Data["user"].(map[string]interface{}),
			}.Bot()
			sock.self.Latency = sock.latency
			sock.locked = false
			if sock.listeners.OnReady != nil {
				go sock.listeners.OnReady(*sock.self)
			} else {
				fmt.Println(fmt.Sprintf("Logged in as %s#%s (Id:%s)",
					sock.self.Username, sock.self.Discriminator, sock.self.Id))
				fmt.Println("---------")
			}
		case string(OnGuildJoin):
			g := Converter{token: token, payload: wsmsg.Data}.Guild()
			g.clientId = sock.self.Id
			s.Guilds[g.Id] = g
			sock.self.Guilds = s.Guilds
			if sock.memoize {
				sock.requestMembers(conn, g.Id)
			}
			_, ok := sock.lazyGuilds[g.Id]
			if !ok && sock.listeners.OnGuildJoin != nil {
				go sock.listeners.OnGuildJoin(*sock.self, *g)
			}
		case string(OnGuildMembersChunk):
			id := wsmsg.Data["guild_id"].(string)
			members := wsmsg.Data["members"].([]interface{})
			s.Guilds[id].fillMembers(members)
			sock.self.Users = s.Users
		default:
			sock.processEvents(wsmsg.Event, wsmsg.Data)
		}

		switch wsmsg.OP {

		case 10:
			sock.interval = wsmsg.Data["heartbeat_interval"].(float64)
			sock.identify(conn, sock.intent)
			go sock.keepAlive(conn, int(sock.interval))
		case 11:
			sock.beatAck = time.Now().UnixMilli()
			sock.latency = sock.beatAck - sock.beatSent
			if sock.self != nil {
				sock.self.Latency = sock.latency
			}
		case 7:
			_ = conn.WriteJSON(map[string]interface{}{
				"op": 6,
				"d": map[string]interface{}{
					"token":      token,
					"sequence":   sock.sequence,
					"session_id": sock.runtimeMatrices.SessionId,
				},
			})
		case 9:
			sock.identify(conn, sock.intent)
		}
	}
}

func (sock *ws) processEvents(dispatch string, data map[string]interface{}) {
	if sock.locked {
		return
	}
	conv := Converter{token: sock.secret, payload: data}

	switch dispatch {

	case string(OnMessage):
		if sock.listeners.OnMessage != nil {
			go sock.listeners.OnMessage(*sock.self, *conv.Message())
		}

	case string(OnGuildLeave):
		id := data["id"].(string)
		_, ok := sock.lazyGuilds[id]
		if ok {
			delete(sock.lazyGuilds, id)
		}
		if sock.listeners.OnGuildLeave != nil {
			cached, exists := s.Guilds[id]
			if exists {
				go sock.listeners.OnGuildLeave(*sock.self, *cached)
				delete(s.Guilds, id)
			} else {
				var guild Guild
				md, _ := json.Marshal(data)
				_ = json.Unmarshal(md, &guild)
				go sock.listeners.OnGuildLeave(*sock.self, guild)
			}
		}

	case string(OnInteraction):
		if sock.listeners.OnInteraction != nil {
			go sock.listeners.OnInteraction(*sock.self, *conv.Interaction())
		}
		ctx := conv.Context()
		switch ctx.Type {
		case 1:
			// interaction ping
		case 2:
			if task, ok := sock.commands[ctx.Data.Id]; ok {
				switch ctx.Data.Type {
				case int(SlashCommand):
					if len(ctx.Data.Options) > 0 && int(ctx.Data.Options[0].Type) == 1 {
						type subcommand struct {
							Name    string   `json:"name"`
							Options []Option `json:"options"`
							Type    int      `json:"type"`
						}
						d := ctx.raw["data"].(map[string]interface{})["options"].([]interface{})
						ds, _ := json.Marshal(d)
						var subcommands []subcommand
						_ = json.Unmarshal(ds, &subcommands)
						scos := buildRO(subcommands[0].Options, ctx.Data.Resolved, sock.secret)
						if sc, okz := subcommandBucket[ctx.Data.Id]; okz {
							scTask, exists := sc.(map[string]interface{})[ctx.Data.Options[0].Name]
							if exists {
								hook := scTask.(func(bot Bot, ctx Context, options ResolvedOptions))
								go hook(*sock.self, *ctx, *scos)
							}
						}
					} else {
						hook := task.(func(bot Bot, ctx Context, options ResolvedOptions))
						if hook != nil {
							ro := buildRO(ctx.Data.Options, ctx.Data.Resolved, sock.secret)
							go hook(*sock.self, *ctx, *ro)
						}
					}
				case int(UserCommand):
					target := ctx.Data.TargetId
					rud := ctx.Data.Resolved["users"].(map[string]interface{})[target]
					ctx.TargetUser = *Converter{token: sock.secret, payload: rud.(map[string]interface{})}.User()
					hook := task.(func(bot Bot, ctx Context, _ map[string]Option))
					go hook(*sock.self, *ctx, nil)
				case int(MessageCommand):
					target := ctx.Data.TargetId
					rmd := ctx.Data.Resolved["messages"].(map[string]interface{})[target]
					ctx.TargetMessage = *Converter{token: sock.secret, payload: rmd.(map[string]interface{})}.Message()
					hook := task.(func(bot Bot, ctx Context, _ map[string]Option))
					go hook(*sock.self, *ctx, nil)
				}
			} else {
				log.Printf("Command (%s) is not implemented.", ctx.Data.Id)
			}

		case 3:
			var compdata componentData
			da, _ := json.Marshal(data["data"].(map[string]interface{}))
			_ = json.Unmarshal(da, &compdata)
			ctx.componentData = compdata
			switch ctx.componentData.ComponentType {
			case 2:
				cb, ok := callbackTasks[ctx.componentData.CustomId]
				if ok {
					callback := cb.(func(b Bot, ctx Context))
					go callback(*sock.self, *ctx)
				}
			case 3:
				cb, ok := callbackTasks[ctx.componentData.CustomId]
				if ok {
					callback := cb.(func(b Bot, ctx Context, values ...string))
					go callback(*sock.self, *ctx, ctx.componentData.Values...)
				}
			}
			tmp, ok := timeoutTasks[ctx.componentData.CustomId]
			if ok {
				onTimeoutHandler := tmp[1].(func(b Bot, ctx Context))
				duration := tmp[0].(float64)
				delete(timeoutTasks, ctx.componentData.CustomId)
				go scheduleTimeoutTask(duration, *sock.self, *ctx, onTimeoutHandler)
			}
		case 4:
			// handle auto-complete interaction
		case 5:
			callback, ok := callbackTasks[ctx.componentData.CustomId]
			if ok {
				go callback.(func(b Bot, ctx *Context))(*sock.self, ctx)
				delete(callbackTasks, ctx.componentData.CustomId)
			}
		default:
			log.Println("Unknown interaction type: ", ctx.Type)
		}
	default:
	}
}

func scheduleTimeoutTask(timeout float64, user Bot, ctx Context, task func(bot Bot, ctx Context)) {
	time.Sleep(time.Duration(timeout) * time.Second)
	task(user, ctx)
}
