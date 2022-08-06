package disgo

import (
	"encoding/json"
	"fmt"
	"io"
)

type Component struct {
	CustomId string   `json:"custom_id"`
	Type     int      `json:"type"`
	Value    string   `json:"value"`
	Values   []string `json:"values"`
}

type Row struct {
	Components []Component
}

type ComponentData struct {
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

func (resp *Response) Marshal() map[string]interface{} {
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
		body["embeds"] = append(body["embeds"].([]map[string]interface{}), em.Marshal())
	}
	if len(resp.AllowedMentions) > 0 && len(resp.AllowedMentions) <= 100 {
		body["allowed_mentions"] = resp.AllowedMentions
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
	if len(resp.View.ActionRows) > 0 {
		body["components"] = resp.View.marshal()
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
	return body
}

type Followup struct {
	Id            string
	Content       string
	Embeds        []Embed
	ChannelId     string
	Flags         int
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

func (f *Followup) Edit(resp Response) Followup {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/%s", f.applicationId, f.token, f.Id)
	body := resp.Marshal()
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
			applicationId: f.applicationId,
		}
	}()
	val, ok := <-fl
	if ok {
		return val
	}
	return Followup{}
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

type Context struct {
	Id             string `json:"id"`
	ApplicationId  string `json:"application_id"`
	Type           int    `json:"type"`
	Data           Data   `json:"data"`
	GuildId        string `json:"guild_id"`
	ChannelId      string `json:"channel_id"`
	Member         Member `json:"member"`
	User           User   `json:"user"`
	Token          string `json:"token"`
	Version        int    `json:"version"`
	AppPermissions string `json:"app_permissions"`
	Locale         string `json:"locale"`
	GuildLocale    string `json:"guild_locale"`
	TargetUser     User
	TargetMessage  Message
	componentData  ComponentData
	commandData    []Option
}

func unmarshalContext(payload interface{}) *Context {
	c := &Context{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, c)
	return c
}

func (c *Context) OriginalResponse() Message {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := minimalReq("GET", path, nil, "")
	bs, _ := io.ReadAll(r.fire().Body)
	var m Message
	_ = json.Unmarshal(bs, &m)
	return m
}

func (c *Context) Send(resp Response) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := multipartReq(
		"POST", path, map[string]interface{}{"type": 4, "data": resp.Marshal()}, "", resp.Files...)
	go r.fire()
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

func (c *Context) SendModal(modal Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := minimalReq("POST", path, modal.marshal(), "")
	go r.fire()
}

func (c *Context) SendFollowup(resp Response) Followup {
	path := fmt.Sprintf("/webhooks/%s/%s", c.ApplicationId, c.Token)
	r := multipartReq("POST", path, resp.Marshal(), "", resp.Files...)
	fl := make(chan Followup, 1)
	go func() {
		bs, _ := io.ReadAll(r.fire().Body)
		var msg Message
		_ = json.Unmarshal(bs, &msg)
		fl <- Followup{
			Id:            msg.Id,
			Content:       msg.Content,
			Embeds:        msg.Embeds,
			ChannelId:     c.ChannelId,
			Flags:         msg.Flags,
			token:         c.Token,
			applicationId: c.ApplicationId,
		}
	}()
	val, ok := <-fl
	if ok {
		return val
	}
	return Followup{}
}

func (c *Context) Edit(resp Response) {
	if c.Type == 2 {
		path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
		r := multipartReq("PATCH", path, resp.Marshal(), "", resp.Files...)
		go r.fire()
	} else {
		path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
		body := map[string]interface{}{"type": 7, "data": resp.Marshal()}
		r := multipartReq("POST", path, body, "", resp.Files...)
		go r.fire()
	}
}

func (c *Context) Delete() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := minimalReq("DELETE", path, nil, "")
	go r.fire()
}
