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
	var finalEmbeds []embed.Embed
	if utils.CheckTrueEmbed(m.Embed) {
		finalEmbeds = append(finalEmbeds, m.Embed)
	}
	if len(m.Embeds) > 0 && len(m.Embeds) < 25 {
		for _, em := range m.Embeds {
			if utils.CheckTrueEmbed(em) {
				finalEmbeds = append(finalEmbeds, em)
			}
		}
	}
	body["embeds"] = finalEmbeds
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
	var finalFiles []file.File
	if utils.CheckTrueFile(m.File) {
		finalFiles = append(finalFiles, m.File)
	}
	for _, f := range m.Files {
		if utils.CheckTrueFile(f) {
			finalFiles = append(finalFiles, f)
		}
	}
	if len(finalFiles) > 0 {
		body["attachments"] = []map[string]interface{}{}
		for i, f := range finalFiles {
			a := map[string]interface{}{
				"id":          i,
				"filename":    f.Name,
				"description": f.Description,
			}
			body["attachments"] = append(body["attachments"].([]map[string]interface{}), a)
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

type Interaction struct {
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

func FromData(payload interface{}) *Interaction {
	i := &Interaction{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, i)
	return i
}

func (i *Interaction) SendResponse(message Response) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	body := map[string]interface{}{
		"type": 4,
		"data": message.ToBody(),
	}
	r := router.New("POST", path, body, "", message.Files)
	go r.Request()
}

func (i *Interaction) Ack() {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	body := map[string]interface{}{"type": 1}
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (i *Interaction) Defer(ephemeral bool) {
	body := map[string]interface{}{"type": 5}
	if ephemeral {
		body["data"] = map[string]interface{}{"flags": 1 << 6}
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (i *Interaction) SendModal(modal component.Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	r := router.New("POST", path, modal.ToBody(), "", nil)
	go r.Request()
}

func (i *Interaction) SendAutoComplete(choices ...command.Choice) {
	body := map[string]interface{}{
		"type": 8,
		"data": map[string]interface{}{"choices": choices},
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (i *Interaction) SendFollowup(resp Response) {
	path := fmt.Sprintf("/webhooks/%s/%s", i.ApplicationId, i.Token)
	r := router.New("POST", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}
