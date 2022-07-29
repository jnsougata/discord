package component

import (
	"encoding/json"
	"fmt"
	"github.com/jnsougata/disgo/embed"
	"github.com/jnsougata/disgo/file"
	"github.com/jnsougata/disgo/member"
	"github.com/jnsougata/disgo/router"
	"github.com/jnsougata/disgo/user"
	"github.com/jnsougata/disgo/utils"
)

type Message struct {
	Content         string
	Embed           embed.Embed
	Embeds          []embed.Embed
	AllowedMentions []string
	TTS             bool
	Ephemeral       bool
	SuppressEmbeds  bool
	View            View
	File            file.File
	Files           []file.File
}

func (m *Message) ToBody() map[string]interface{} {
	flag := 0
	body := map[string]interface{}{}
	if m.Content != "" {
		body["content"] = m.Content
	}
	if utils.CheckTrueEmbed(m.Embed) {
		m.Embeds = append([]embed.Embed{m.Embed}, m.Embeds...)
	}
	for i, em := range m.Embeds {
		if !utils.CheckTrueEmbed(em) {
			m.Embeds = append(m.Embeds[:i], m.Embeds[i+1:]...)
		}
	}
	if len(m.Embeds) > 10 {
		m.Embeds = m.Embeds[:10]
	}
	body["embeds"] = []map[string]interface{}{}
	for _, em := range m.Embeds {
		body["embeds"] = append(body["embeds"].([]map[string]interface{}), em.ToBody())
	}
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

type Component struct {
	CustomId string   `json:"custom_id"`
	Type     int      `json:"type"`
	Value    string   `json:"value"`
	Values   []string `json:"values"`
}

type Row struct {
	Components []Component
}

type Data struct {
	ComponentType int      `json:"component_type"`
	CustomId      string   `json:"custom_id"`
	Values        []string `json:"values"`
	Components    []Row    `json:"components"`
}

type Context struct {
	Id             string                 `json:"id"`
	ApplicationId  string                 `json:"application_id"`
	Type           int                    `json:"type"`
	Data           Data                   `json:"data"`
	GuildId        string                 `json:"guild_id"`
	ChannelId      string                 `json:"channel_id"`
	Member         member.Member          `json:"member"`
	User           user.User              `json:"user"`
	Token          string                 `json:"token"`
	Version        int                    `json:"version"`
	Message        map[string]interface{} `json:"message"`
	AppPermissions string                 `json:"app_permissions"`
	Locale         string                 `json:"locale"`
	GuildLocale    string                 `json:"guild_locale"`
}

func FromData(payload interface{}) *Context {
	i := &Context{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, i)
	return i
}

func (c *Context) SendResponse(resp Message) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := router.New(
		"POST", path, map[string]interface{}{"type": 4, "data": resp.ToBody()}, "", resp.Files)
	go r.Request()
}

func (c *Context) Ack() {
	body := map[string]interface{}{"type": 6}
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (c *Context) SendModal(modal Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := router.New("POST", path, modal.ToBody(), "", nil)
	go r.Request()
}

func (c *Context) SendFollowup(resp Message) {
	path := fmt.Sprintf("/webhooks/%s/%s", c.ApplicationId, c.Token)
	r := router.New("POST", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}

func (c *Context) EditOriginalMessage(resp Message) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	body := map[string]interface{}{"type": 7, "data": resp.ToBody()}
	r := router.New("POST", path, body, "", resp.Files)
	go r.Request()
}
