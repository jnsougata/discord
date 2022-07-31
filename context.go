package disgo

import (
	"encoding/json"
	"fmt"
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
	if CheckTrueEmbed(resp.Embed) {
		resp.Embeds = append([]Embed{resp.Embed}, resp.Embeds...)
	}
	for i, em := range resp.Embeds {
		if !CheckTrueEmbed(em) {
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
		body["components"] = resp.View.ToComponent()
	}
	if CheckTrueFile(resp.File) {
		resp.Files = append([]File{resp.File}, resp.Files...)
	}
	body["attachments"] = []map[string]interface{}{}
	for i, f := range resp.Files {
		if CheckTrueFile(f) {
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

type Context struct {
	Id             string                 `json:"id"`
	ApplicationId  string                 `json:"application_id"`
	Type           int                    `json:"type"`
	Data           InteractionData        `json:"data"`
	GuildId        string                 `json:"guild_id"`
	ChannelId      string                 `json:"channel_id"`
	Member         Member                 `json:"member"`
	User           User                   `json:"user"`
	Token          string                 `json:"token"`
	Version        int                    `json:"version"`
	Message        map[string]interface{} `json:"message"`
	AppPermissions string                 `json:"app_permissions"`
	Locale         string                 `json:"locale"`
	GuildLocale    string                 `json:"guild_locale"`
	ComponentData  ComponentData          `json:"x_component"`
	CommandData    []SlashCommandOption   `json:"x_command"`
}

func UnmarshalContext(payload interface{}) *Context {
	c := &Context{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, c)
	return c
}

func (c *Context) Send(resp Response) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := MultipartReq(
		"POST", path, map[string]interface{}{"type": 4, "data": resp.Marshal()}, "", resp.Files)
	go r.Request()
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
	r := MinimalReq("POST", path, body, "")
	go r.Request()
}

func (c *Context) SendModal(modal Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := MinimalReq("POST", path, modal.Marshal(), "")
	go r.Request()
}

func (c *Context) SendFollowup(resp Response) {
	path := fmt.Sprintf("/webhooks/%s/%s", c.ApplicationId, c.Token)
	r := MultipartReq("POST", path, resp.Marshal(), "", resp.Files)
	go r.Request()
	// TODO: handle followup
}

func (c *Context) Edit(resp Response) {
	if c.Type == 2 {
		path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
		r := MultipartReq("PATCH", path, resp.Marshal(), "", resp.Files)
		go r.Request()
	} else {
		path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
		body := map[string]interface{}{"type": 7, "data": resp.Marshal()}
		r := MultipartReq("POST", path, body, "", resp.Files)
		go r.Request()
	}
}

func (c *Context) Delete() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := MinimalReq("DELETE", path, nil, "")
	go r.Request()
}
