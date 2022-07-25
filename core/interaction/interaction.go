package interaction

import (
	"encoding/json"
	"fmt"
	"github.com/jnsougata/disgo/core/attachment"
	"github.com/jnsougata/disgo/core/command"
	"github.com/jnsougata/disgo/core/component"
	"github.com/jnsougata/disgo/core/embed"
	"github.com/jnsougata/disgo/core/modal"
	"github.com/jnsougata/disgo/core/router"
	"github.com/jnsougata/disgo/core/user"
)

type Message struct {
	Content         string
	Embeds          []embed.Embed
	AllowedMentions []string
	Tts             bool
	Flags           int
	View            component.View
	Attachments     []attachment.Partial
}

func (m *Message) ToBody() map[string]interface{} {
	return map[string]interface{}{
		"content":    m.Content,
		"embeds":     m.Embeds,
		"tts":        m.Tts,
		"flags":      m.Flags,
		"components": m.View.ToComponent(),
	}
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
	GuildID  string                 `json:"guild_id"`
	TargetId string                 `json:"target_id"`
}

type Interaction struct {
	ID             string      `json:"id"`
	ApplicationID  string      `json:"application_id"`
	Type           int         `json:"type"`
	Data           Data        `json:"data"`
	GuildID        string      `json:"guild_id"`
	ChannelID      string      `json:"channel_id"`
	Member         interface{} `json:"member"`
	User           *user.User  `json:"user"`
	Token          string      `json:"token"`
	Version        int         `json:"version"`
	Message        interface{} `json:"message"`
	AppPermissions string      `json:"app_permissions"`
	Locale         string      `json:"locale"`
	GuildLocale    string      `json:"guild_locale"`
}

func FromData(payload interface{}) *Interaction {
	i := &Interaction{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, i)
	return i
}

func (i *Interaction) SendResponse(message Message) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, map[string]interface{}{"type": 4, "data": message.ToBody()}, "")
	go r.Request()
}

func (i *Interaction) Ack() {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, map[string]interface{}{"type": 1}, "")
	go r.Request()
}

func (i *Interaction) Defer(ephemeral bool) {
	payload := map[string]interface{}{"type": 5}
	if ephemeral {
		payload["data"] = map[string]interface{}{"flags": 1 << 6}
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, payload, "")
	go r.Request()
}

func (i *Interaction) SendModal(modal modal.Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, modal.ToBody(), "")
	go r.Request()
}

func (i *Interaction) SendAutoComplete(choices ...command.Choice) {
	payload := map[string]interface{}{
		"type": 8,
		"data": map[string]interface{}{"choices": choices},
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, payload, "")
	go r.Request()
}

func (i *Interaction) SendFollowup(message Message) {
	path := fmt.Sprintf("/webhooks/%s/%s", i.ApplicationID, i.Token)
	r := router.New("POST", path, message.ToBody(), "")
	go r.Request()
}
