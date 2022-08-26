package discord

import (
	"encoding/json"
	"fmt"
	"io"
)

type component struct {
	CustomId string   `json:"custom_id"`
	Type     int      `json:"type"`
	Value    string   `json:"value"`
	Values   []string `json:"values"`
}

type row struct {
	Components []component
}

type componentData struct {
	ComponentType int      `json:"component_type"`
	CustomId      string   `json:"custom_id"`
	Values        []string `json:"values"`
	Components    []row    `json:"components"`
}

type Data struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	Type     int                    `json:"type"`
	Resolved map[string]interface{} `json:"resolved"`
	Options  []option               `json:"options"`
	GuildId  string                 `json:"guild_id"`
	TargetId string                 `json:"target_id"`
}

type ResolvedOptions struct {
	strings      map[string]string
	integers     map[string]int64
	booleans     map[string]bool
	numbers      map[string]float64
	users        map[string]User
	roles        map[string]Role
	channels     map[string]Channel
	mentionables map[string]interface{}
	attachments  map[string]Attachment
}

func (ro *ResolvedOptions) String(name string) (string, bool) {
	val, ok := ro.strings[name]
	return val, ok
}

func (ro *ResolvedOptions) Integer(name string) (int64, bool) {
	val, ok := ro.integers[name]
	return val, ok
}

func (ro *ResolvedOptions) Boolean(name string) (bool, bool) {
	val, ok := ro.booleans[name]
	return val, ok
}

func (ro *ResolvedOptions) Number(name string) (float64, bool) {
	val, ok := ro.numbers[name]
	return val, ok
}

func (ro *ResolvedOptions) User(name string) (User, bool) {
	val, ok := ro.users[name]
	return val, ok
}

func (ro *ResolvedOptions) Role(name string) (Role, bool) {
	val, ok := ro.roles[name]
	return val, ok
}

func (ro *ResolvedOptions) Channel(name string) (Channel, bool) {
	val, ok := ro.channels[name]
	return val, ok
}

func (ro *ResolvedOptions) Mentionable(name string) (interface{}, bool) {
	val, ok := ro.mentionables[name]
	return val, ok
}

func (ro *ResolvedOptions) Attachment(name string) (Attachment, bool) {
	val, ok := ro.attachments[name]
	return val, ok
}

type Context struct {
	Id             string `json:"id"`
	ApplicationId  string `json:"application_id"`
	Type           int    `json:"type"`
	Data           Data   `json:"data"`
	GuildId        string `json:"guild_id"`
	ChannelId      string `json:"channel_id"`
	User           User   `json:"user"`
	Token          string `json:"token"`
	Version        int    `json:"version"`
	AppPermissions string `json:"app_permissions"`
	Locale         string `json:"locale"`
	GuildLocale    string `json:"guild_locale"`
	TargetUser     User
	TargetMessage  Message
	Channel        Channel
	Guild          Guild
	Author         Member
	//token          string
	commandData   []option
	componentData componentData
	data          map[string]interface{}
	state         *state
}

func (c *Context) OriginalResponse() Message {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := minimalReq("GET", path, nil, "")
	bs, _ := io.ReadAll(r.fire().Body)
	var m Message
	_ = json.Unmarshal(bs, &m)
	return m
}

func (c *Context) Send(response Response) error {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	body, err := response.marshal()
	r := multipartReq("POST", path, map[string]interface{}{"type": 4, "data": body}, "", response.Files...)
	go r.fire()
	return err
}

func (c *Context) Defer(ephemeral bool) {
	body := map[string]interface{}{}
	if c.Type == 2 {
		body["type"] = 5
		if ephemeral {
			body["data"] = map[string]interface{}{"flags": 1 << 6}
		}
	} else {
		body["type"] = 6
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := minimalReq("POST", path, body, "")
	go r.fire()
}

func (c *Context) SendModal(modal Modal) error {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	body, err := modal.marshal()
	r := minimalReq("POST", path, body, "")
	go r.fire()
	return err
}

func (c *Context) SendFollowup(response Response) (Message, error) {
	path := fmt.Sprintf("/webhooks/%s/%s", c.ApplicationId, c.Token)
	body, err := response.marshal()
	r := multipartReq("POST", path, body, "", response.Files...)
	m := make(chan Message, 1)
	go func() {
		bs, _ := io.ReadAll(r.fire().Body)
		var msg Message
		_ = json.Unmarshal(bs, &msg)
		m <- msg
	}()
	return <-m, err
}

func (c *Context) Edit(response Response) error {
	body, err := response.marshal()
	if c.Type == 2 {
		path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
		r := multipartReq("PATCH", path, body, "", response.Files...)
		go r.fire()
	} else {
		path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
		p := map[string]interface{}{"type": 7, "data": body}
		r := multipartReq("POST", path, p, "", response.Files...)
		go r.fire()
	}
	return err
}

func (c *Context) Delete() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := minimalReq("DELETE", path, nil, "")
	go r.fire()
}
