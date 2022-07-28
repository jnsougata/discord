package command

import (
	"fmt"
	"github.com/jnsougata/disgo/core/component"
	"github.com/jnsougata/disgo/core/embed"
	"github.com/jnsougata/disgo/core/file"
	"github.com/jnsougata/disgo/core/router"
	"github.com/jnsougata/disgo/core/utils"
)

type Message struct {
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

type Context struct {
	Id            string
	Token         string
	ApplicationId string
}

func (c *Context) SendResponse(resp Message) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	body := map[string]interface{}{
		"type": 4,
		"data": resp.ToBody(),
	}
	r := router.New("POST", path, body, "", resp.Files)
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

func (c *Context) SendAutoComplete(choices ...Choice) {
	body := map[string]interface{}{
		"type": 8,
		"data": map[string]interface{}{"choices": choices},
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (c *Context) SendFollowup(resp Message) {
	path := fmt.Sprintf("/webhooks/%s/%s", c.ApplicationId, c.Token)
	r := router.New("POST", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}

func (c *Context) DeleteOriginalResponse() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := router.Minimal("DELETE", path, nil, "")
	go r.Request()
}

func (c *Context) EditOriginalResponse(resp Message) {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := router.New("PATCH", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}
