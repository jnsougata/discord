package interaction

import (
	"encoding/json"
	"fmt"
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/component"
	"github.com/jnsougata/disgo/core/embed"
	"github.com/jnsougata/disgo/core/file"
	"github.com/jnsougata/disgo/core/member"
	"github.com/jnsougata/disgo/core/router"
	"github.com/jnsougata/disgo/core/user"
	"github.com/jnsougata/disgo/core/utils"
)

type Response struct {
	Content         string
	Embed           embed.Embed
	Embeds          []embed.Embed
	AllowedMentions []string
	TTS             bool
	Ephemeral       bool
	SuppressEmbeds  bool
	View            component.View
	File            file.File
	Files           []file.File
}

func (m *Response) ToBody() map[string]interface{} {
	flag := 0
	body := map[string]interface{}{}
	if m.Content != "" {
		body["content"] = m.Content
	}
	if utils.CheckTrueEmbed(m.Embed) {
		m.Embeds = append([]embed.Embed{m.Embed}, m.Embeds...)
	}
	if len(m.Embeds) > 0 && len(m.Embeds) < 10 {
		for i, em := range m.Embeds {
			if !utils.CheckTrueEmbed(em) {
				m.Embeds = append(m.Embeds[:i], m.Embeds[i+1:]...)
			}
			if i > 10 {
				m.Embeds = append(m.Embeds[:i], m.Embeds[i+1:]...)
			}
		}
	}
	body["embeds"] = m.Embeds
	if len(m.AllowedMentions) > 0 && len(m.AllowedMentions) <= 100 {
		body["allowed_mentions"] = m.AllowedMentions
	}
	if m.TTS {
		body["tts"] = true
	}
	if m.Ephemeral {
		flag |= 1 << 6
	}
	if m.SuppressEmbeds {
		flag |= 1 << 2
	}
	if m.Ephemeral || m.SuppressEmbeds {
		body["flags"] = flag
	}
	if len(m.View.ActionRows) > 0 {
		body["components"] = m.View.ToComponent()
	}
	if utils.CheckTrueFile(m.File) {
		m.Files = append([]file.File{m.File}, m.Files...)
	}
	body["attachments"] = []map[string]interface{}{}
	for i, f := range m.Files {
		if utils.CheckTrueFile(f) {
			a := map[string]interface{}{
				"id":          i,
				"filename":    f.Name,
				"description": f.Description,
			}
			body["attachments"] = append(body["attachments"].([]map[string]interface{}), a)
		} else {
			m.Files = append(m.Files[:i], m.Files[i+1:]...)
		}
	}
	return body
}

type Option struct {
	Name    string      `json:"name"`
	Type    int         `json:"type"`
	Value   interface{} `json:"value"`
	Options []Option    `json:"options"`
	Focused bool        `json:"focused"`
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
	Id             string        `json:"id"`
	ApplicationId  string        `json:"application_id"`
	Type           int           `json:"type"`
	Data           Data          `json:"data"`
	GuildId        string        `json:"guild_id"`
	ChannelId      string        `json:"channel_id"`
	Member         member.Member `json:"member"`
	User           user.User     `json:"user"`
	Token          string        `json:"token"`
	Version        int           `json:"version"`
	Message        interface{}   `json:"message"`
	AppPermissions string        `json:"app_permissions"`
	Locale         string        `json:"locale"`
	GuildLocale    string        `json:"guild_locale"`
}

func FromData(payload interface{}) *Context {
	i := &Context{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, i)
	return i
}

func (c *Context) SendResponse(message Response) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	body := map[string]interface{}{
		"type": 4,
		"data": message.ToBody(),
	}
	r := router.New("POST", path, body, "", message.Files)
	go r.Request()
}

func (c *Context) Ack() {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	body := map[string]interface{}{"type": 1}
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (c *Context) Defer(ephemeral bool) {
	body := map[string]interface{}{"type": 5}
	if ephemeral {
		body["data"] = map[string]interface{}{"flags": 1 << 6}
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (c *Context) SendModal(modal component.Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := router.New("POST", path, modal.ToBody(), "", nil)
	go r.Request()
}

func (c *Context) SendAutoComplete(choices ...command.Choice) {
	body := map[string]interface{}{
		"type": 8,
		"data": map[string]interface{}{"choices": choices},
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (c *Context) SendFollowup(resp Response) {
	path := fmt.Sprintf("/webhooks/%s/%s", c.ApplicationId, c.Token)
	r := router.New("POST", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}

func (c *Context) DeleteOriginalResponse() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := router.Minimal("DELETE", path, nil, "")
	go r.Request()
}

func (c *Context) EditOriginalResponse(resp Response) {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := router.New("PATCH", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}
