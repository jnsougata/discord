package discord

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"time"
)

var shared = state{}

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

type socket struct {
	own              *Bot
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
	secret           string
	queue            []Command
	commands         map[string]interface{}
	listeners        Listeners
	lazyGuilds       map[string]bool
	shardingMatrices sharding
	runtimeMatrices  runtime
}

func (ws *socket) gateway() string {
	version := 11
	encoding := "json"
	info := sharding{}
	r := minimalReq("GET", "/gateway/bot", nil, ws.secret)
	resp := r.fire()
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &info)
	ws.shardingMatrices = info
	return fmt.Sprintf("%s?v=%d&encoding=%s", info.URL, version, encoding)
}

func (ws *socket) pacemaker(conn *websocket.Conn, dur int) {
	for {
		_ = conn.WriteJSON(map[string]interface{}{"op": 1, "d": 251})
		ws.beatSent = time.Now().UnixMilli()
		time.Sleep(time.Duration(dur) * time.Millisecond)
	}
}

func (ws *socket) identify(conn *websocket.Conn, intent int) {
	identification := map[string]interface{}{
		"token":   ws.secret,
		"intents": intent,
	}
	identification["properties"] = map[string]string{
		"os":      "linux",
		"browser": "discord",
		"device":  "discord",
	}
	if ws.presence.Activity.Name != "" {
		identification["presence"] = ws.presence.Marshal()
	}
	if ws.presence.OnMobile {
		identification["properties"].(map[string]string)["browser"] = "Discord iOS"
	}
	payload := map[string]interface{}{"op": 2, "d": identification}
	_ = conn.WriteJSON(payload)
}

func (ws *socket) registerCommand(com Command, applicationId string) {
	var route string
	data, hook, guildId := com.marshal()
	if guildId != 0 {
		route = fmt.Sprintf("/applications/%s/guilds/%v/commands", applicationId, guildId)
	} else {
		route = fmt.Sprintf("/applications/%s/commands", applicationId)
	}
	r := minimalReq("POST", route, data, ws.secret)
	resp := map[string]interface{}{}
	body, _ := io.ReadAll(r.fire().Body)
	_ = json.Unmarshal(body, &resp)
	_, ok := resp["id"]
	if ok {
		commandId := resp["id"].(string)
		ws.commands[commandId] = hook
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
		message := resp["message"]
		panic(fmt.Sprintf("Failed to register command {%s}. Message: %s", com.Name, message))
	}
}

func (ws *socket) requestMembers(conn *websocket.Conn, guildId string) {
	_ = conn.WriteJSON(map[string]interface{}{
		"op": 8,
		"d": map[string]interface{}{
			"guild_id": guildId,
			"query":    "",
			"limit":    0,
		},
	})
}

func (ws *socket) run(token string) {
	shared.Guilds = make(map[string]*Guild)
	shared.Users = make(map[string]*User)
	shared.Token = token
	conn, _, err := websocket.DefaultDialer.Dial(ws.gateway(), nil)
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
		ws.sequence = wsmsg.Sequence
		if ws.listeners.OnSocketReceive != nil {
			go ws.listeners.OnSocketReceive(wsmsg)
		}
		r := runtime{}
		switch wsmsg.Event {
		case string(OnReady):
			ws.lazyGuilds = map[string]bool{}
			ba, _ := json.Marshal(wsmsg.Data)
			_ = json.Unmarshal(ba, &r)
			ws.runtimeMatrices = r
			for _, guild := range r.Guilds.([]interface{}) {
				id := guild.(map[string]interface{})["id"].(string)
				ws.lazyGuilds[id] = true
			}
			for _, cmd := range ws.queue {
				go ws.registerCommand(cmd, r.Application.Id)
			}
			ws.own = converter{
				state:   &shared,
				payload: wsmsg.Data["user"].(map[string]interface{}),
			}.Bot()
			ws.own.Latency = ws.latency
			ws.locked = false
			if ws.listeners.OnReady != nil {
				go ws.listeners.OnReady(*ws.own)
			} else {
				fmt.Println(fmt.Sprintf("Logged in as %s#%s (Id:%s)",
					ws.own.Username, ws.own.Discriminator, ws.own.Id))
				fmt.Println("---------")
			}
		case string(OnGuildJoin):
			g := converter{state: &shared, payload: wsmsg.Data}.Guild()
			g.clientId = ws.own.Id
			shared.Guilds[g.Id] = g
			ws.own.Guilds = shared.Guilds
			if ws.memoize {
				ws.requestMembers(conn, g.Id)
			}
			_, ok := ws.lazyGuilds[g.Id]
			if !ok && ws.listeners.OnGuildJoin != nil {
				go ws.listeners.OnGuildJoin(*ws.own, *g)
			}
		case string(OnGuildMembersChunk):
			id := wsmsg.Data["guild_id"].(string)
			members := wsmsg.Data["members"].([]interface{})
			shared.Guilds[id].fillMembers(members)
			ws.own.Users = shared.Users
		default:
			ws.handleDispatch(wsmsg.Event, wsmsg.Data)
		}

		switch wsmsg.OP {

		case 10:
			ws.interval = wsmsg.Data["heartbeat_interval"].(float64)
			ws.identify(conn, ws.intent)
			go ws.pacemaker(conn, int(ws.interval))
		case 11:
			ws.beatAck = time.Now().UnixMilli()
			ws.latency = ws.beatAck - ws.beatSent
			if ws.own != nil {
				ws.own.Latency = ws.latency
			}
		case 7:
			_ = conn.WriteJSON(map[string]interface{}{
				"op": 6,
				"d": map[string]interface{}{
					"token":      token,
					"sequence":   ws.sequence,
					"session_id": ws.runtimeMatrices.SessionId,
				},
			})
		case 9:
			ws.identify(conn, ws.intent)
		}
	}
}

func (ws *socket) handleDispatch(dispatch string, data map[string]interface{}) {
	if ws.locked {
		return
	}
	conv := converter{state: &shared, payload: data}

	switch dispatch {

	case string(OnMessage):
		if ws.listeners.OnMessage != nil {
			go ws.listeners.OnMessage(*ws.own, *conv.Message())
		}

	case string(OnGuildLeave):
		id := data["id"].(string)
		_, ok := ws.lazyGuilds[id]
		if ok {
			delete(ws.lazyGuilds, id)
		}
		if ws.listeners.OnGuildLeave != nil {
			cached, exists := shared.Guilds[id]
			if exists {
				go ws.listeners.OnGuildLeave(*ws.own, *cached)
				delete(shared.Guilds, id)
			} else {
				var guild Guild
				md, _ := json.Marshal(data)
				_ = json.Unmarshal(md, &guild)
				go ws.listeners.OnGuildLeave(*ws.own, guild)
			}
		}

	case string(OnInteraction):
		if ws.listeners.OnInteraction != nil {
			go ws.listeners.OnInteraction(*ws.own, *conv.Interaction())
		}
		ctx := conv.Context()
		switch ctx.Type {
		case 1:
			// interaction ping
		case 2:
			if task, ok := ws.commands[ctx.Data.Id]; ok {
				switch ctx.Data.Type {
				case CommandTypes.Slash:
					if len(ctx.Data.Options) > 0 && int(ctx.Data.Options[0].Type) == 1 {
						type subcommand struct {
							Name    string   `json:"name"`
							Options []Option `json:"Options"`
							Type    int      `json:"type"`
						}
						d := ctx.data["data"].(map[string]interface{})["Options"].([]interface{})
						ds, _ := json.Marshal(d)
						var subcommands []subcommand
						_ = json.Unmarshal(ds, &subcommands)
						scos := buildResolved(subcommands[0].Options, ctx.Data.Resolved, &shared)
						if sc, okz := subcommandBucket[ctx.Data.Id]; okz {
							scTask, exists := sc.(map[string]interface{})[ctx.Data.Options[0].Name]
							if exists {
								hook := scTask.(func(bot Bot, ctx Context, options ResolvedOptions))
								go hook(*ws.own, *ctx, *scos)
							}
						}
					} else {
						hook := task.(func(bot Bot, ctx Context, options ResolvedOptions))
						if hook != nil {
							go hook(*ws.own, *ctx, ResolvedOptions{})
						}
					}
				case CommandTypes.User:
					target := ctx.Data.TargetId
					rud := ctx.Data.Resolved["users"].(map[string]interface{})[target]
					ctx.TargetUser = *converter{state: &shared, payload: rud.(map[string]interface{})}.User()
					hook := task.(func(bot Bot, ctx Context, _ map[string]Option))
					go hook(*ws.own, *ctx, nil)
				case CommandTypes.Message:
					target := ctx.Data.TargetId
					rmd := ctx.Data.Resolved["messages"].(map[string]interface{})[target]
					ctx.TargetMessage = *converter{state: &shared, payload: rmd.(map[string]interface{})}.Message()
					hook := task.(func(bot Bot, ctx Context, _ map[string]Option))
					go hook(*ws.own, *ctx, nil)
				}
			} else {
				fmt.Printf("Command (%s) is not implemented.", ctx.Data.Id)
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
					go callback(*ws.own, *ctx)
				}
			case 3:
				cb, ok := callbackTasks[ctx.componentData.CustomId]
				if ok {
					callback := cb.(func(b Bot, ctx Context, values ...string))
					go callback(*ws.own, *ctx, ctx.componentData.Values...)
				}
			}
			tmp, ok := timeoutTasks[ctx.componentData.CustomId]
			if ok {
				onTimeoutHandler := tmp[1].(func(b Bot, ctx Context))
				duration := tmp[0].(float64)
				delete(timeoutTasks, ctx.componentData.CustomId)
				go scheduleTimeoutTask(duration, *ws.own, *ctx, onTimeoutHandler)
			}
		case 4:
			// handle auto-complete interaction
		case 5:
			callback, ok := callbackTasks[ctx.componentData.CustomId]
			if ok {
				go callback.(func(b Bot, ctx *Context))(*ws.own, ctx)
				delete(callbackTasks, ctx.componentData.CustomId)
			}
		default:
			fmt.Println("Unknown interaction type: ", ctx.Type)
		}
	default:
	}
}
