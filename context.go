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

type Row struct {
	Components []component
}

type componentData struct {
	ComponentType int      `json:"component_type"`
	CustomId      string   `json:"custom_id"`
	Values        []string `json:"values"`
	Components    []Row    `json:"components"`
}

type Response struct {
	Content         string
	Embed           Embed
	Embeds          []Embed
	AllowedMentions []string
	TTS             bool
	Ephemeral       bool
	SuppressEmbeds  bool
	View            View
	File            File
	Files           []File
}

func (resp *Response) marshal() (map[string]interface{}, error) {
	flag := 0
	body := map[string]interface{}{}
	if resp.Content != "" {
		body["content"] = resp.Content
	}
	if checkTrueEmbed(resp.Embed) {
		resp.Embeds = append([]Embed{resp.Embed}, resp.Embeds...)
	}
	for i, em := range resp.Embeds {
		if !checkTrueEmbed(em) {
			resp.Embeds = append(resp.Embeds[:i], resp.Embeds[i+1:]...)
		}
	}
	if len(resp.Embeds) > 10 {
		resp.Embeds = resp.Embeds[:10]
	}
	body["embeds"] = []map[string]interface{}{}
	for _, em := range resp.Embeds {
		emd, err := em.marshal()
		if err != nil {
			return nil, err
		} else {
			body["embeds"] = append(body["embeds"].([]map[string]interface{}), emd)
		}
	}
	if resp.TTS {
		body["tts"] = true
	}
	if resp.Ephemeral {
		flag |= 1 << 6
	}
	if resp.SuppressEmbeds {
		flag |= 1 << 2
	}
	if resp.Ephemeral || resp.SuppressEmbeds {
		body["flags"] = flag
	}
	view, err := resp.View.marshal()
	if err == nil {
		body["components"] = view
	} else {
		return nil, err
	}
	if checkTrueFile(resp.File) {
		resp.Files = append([]File{resp.File}, resp.Files...)
	}
	body["attachments"] = []map[string]interface{}{}
	for i, f := range resp.Files {
		if checkTrueFile(f) {
			a := map[string]interface{}{
				"id":          i,
				"filename":    f.Name,
				"description": f.Description,
			}
			body["attachments"] = append(body["attachments"].([]map[string]interface{}), a)
		} else {
			resp.Files = append(resp.Files[:i], resp.Files[i+1:]...)
		}
	}
	return body, nil
}

type Followup struct {
	Id            string
	Content       string
	Embeds        []Embed
	ChannelId     string
	Flags         int
	ctx           Context
	token         string
	applicationId string
}

func (f *Followup) Delete() {
	if f.Flags|1<<6 == 1<<6|1<<2 || f.Flags == 0 {
		path := fmt.Sprintf("/webhooks/%s/%s/messages/%s", f.applicationId, f.token, f.Id)
		r := minimalReq("DELETE", path, nil, "")
		go r.fire()
	}
}

func (f *Followup) Edit(resp Response) (Followup, error) {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/%s", f.applicationId, f.token, f.Id)
	body, err := resp.marshal()
	body["wait"] = true
	r := multipartReq("PATCH", path, body, "", resp.Files...)
	fl := make(chan Followup, 1)
	go func() {
		bs, _ := io.ReadAll(r.fire().Body)
		var msg Message
		_ = json.Unmarshal(bs, &msg)
		fl <- Followup{
			Id:            msg.Id,
			Content:       msg.Content,
			Embeds:        msg.Embeds,
			ChannelId:     msg.ChannelId,
			Flags:         msg.Flags,
			token:         f.token,
			ctx:           f.ctx,
			applicationId: f.applicationId,
		}
	}()
	return <-fl, err
}

type Data struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	Type     int                    `json:"type"`
	Resolved map[string]interface{} `json:"resolved"`
	Options  []Option               `json:"options"`
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

func (ro *ResolvedOptions) String(name string) (bool, string) {
	val, ok := ro.strings[name]
	return ok, val
}

func (ro *ResolvedOptions) Integer(name string) (bool, int64) {
	val, ok := ro.integers[name]
	return ok, val
}

func (ro *ResolvedOptions) Boolean(name string) (bool, bool) {
	val, ok := ro.booleans[name]
	return ok, val
}

func (ro *ResolvedOptions) Number(name string) (bool, float64) {
	val, ok := ro.numbers[name]
	return ok, val
}

func (ro *ResolvedOptions) User(name string) (bool, User) {
	val, ok := ro.users[name]
	return ok, val
}

func (ro *ResolvedOptions) Role(name string) (bool, Role) {
	val, ok := ro.roles[name]
	return ok, val
}

func (ro *ResolvedOptions) Channel(name string) (bool, Channel) {
	val, ok := ro.channels[name]
	return ok, val
}

func (ro *ResolvedOptions) Mentionable(name string) (bool, interface{}) {
	val, ok := ro.mentionables[name]
	return ok, val
}

func (ro *ResolvedOptions) Attachment(name string) (bool, Attachment) {
	val, ok := ro.attachments[name]
	return ok, val
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
	token          string
	commandData    []Option
	componentData  componentData
	raw            map[string]interface{}
}

func (c *Context) OriginalResponse() Message {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := minimalReq("GET", path, nil, "")
	bs, _ := io.ReadAll(r.fire().Body)
	var m Message
	_ = json.Unmarshal(bs, &m)
	return m
}

func (c *Context) Send(resp Response) error {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	body, err := resp.marshal()
	r := multipartReq(
		"POST", path, map[string]interface{}{"type": 4, "data": body}, "", resp.Files...)
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

func (c *Context) SendFollowup(resp Response) (Followup, error) {
	path := fmt.Sprintf("/webhooks/%s/%s", c.ApplicationId, c.Token)
	body, err := resp.marshal()
	r := multipartReq("POST", path, body, "", resp.Files...)
	f := make(chan Followup, 1)
	go func() {
		bs, _ := io.ReadAll(r.fire().Body)
		var msg Message
		_ = json.Unmarshal(bs, &msg)
		f <- Followup{
			Id:            msg.Id,
			token:         c.Token,
			Content:       msg.Content,
			Embeds:        msg.Embeds,
			ChannelId:     c.ChannelId,
			Flags:         msg.Flags,
			ctx:           *c,
			applicationId: c.ApplicationId,
		}
	}()
	val, _ := <-f
	return val, err
}

func (c *Context) Edit(resp Response) error {
	body, err := resp.marshal()
	if c.Type == 2 {
		path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
		r := multipartReq("PATCH", path, body, "", resp.Files...)
		go r.fire()
	} else {
		path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
		pl := map[string]interface{}{"type": 7, "data": body}
		r := multipartReq("POST", path, pl, "", resp.Files...)
		go r.fire()
	}
	return err
}

func (c *Context) Delete() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := minimalReq("DELETE", path, nil, "")
	go r.fire()
}
