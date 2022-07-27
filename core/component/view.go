package component

import (
	"encoding/json"
	"fmt"
	"github.com/jnsougata/disgo/core/embed"
	"github.com/jnsougata/disgo/core/emoji"
	"github.com/jnsougata/disgo/core/file"
	"github.com/jnsougata/disgo/core/member"
	"github.com/jnsougata/disgo/core/router"
	"github.com/jnsougata/disgo/core/user"
	"github.com/jnsougata/disgo/core/utils"
	"log"
)

const (
	BlueButton  = 1
	GreyButton  = 2
	GreenButton = 3
	RedButton   = 4
	LinkButton  = 5
)

var CallbackTasks = map[string]interface{}{}
var TimeoutTasks = map[string]interface{}{}

type Response struct {
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

type Button struct {
	Style    int
	Label    string
	Emoji    emoji.Partial
	URL      string
	Disabled bool
	CustomId string
}

func (b *Button) OnClick(handler func(bot user.User, interaction Interaction)) {
	b.CustomId = utils.AssignId("")
	CallbackTasks[b.CustomId] = handler
}

func (b *Button) ToComponent() map[string]interface{} {
	b.CustomId = utils.AssignId(b.CustomId)
	btn := map[string]interface{}{
		"type":      2,
		"custom_id": b.CustomId,
	}
	if b.Style != 0 {
		btn["style"] = b.Style
	} else {
		btn["style"] = BlueButton
	}
	if b.Label != "" {
		btn["label"] = b.Label
	} else {
		btn["label"] = "Button"
	}
	if b.Emoji.Id != "" {
		btn["emoji"] = b.Emoji
	}
	if b.URL != "" && b.Style == LinkButton {
		btn["url"] = b.URL
	}
	if b.Disabled {
		btn["disabled"] = true
	}
	return btn
}

type SelectOption struct {
	Label       string        `json:"label"`
	Value       string        `json:"value"`
	Description string        `json:"description,omitempty"`
	Emoji       emoji.Partial `json:"emoji,omitempty"`
	Default     bool          `json:"default,omitempty"`
}

type SelectMenu struct {
	Type        int
	CustomId    string
	Options     []SelectOption
	Placeholder string
	MinValues   int
	MaxValues   int
	Disabled    bool
}

func (s *SelectMenu) OnSelection(handler func(bot user.User, interaction Interaction, values ...string)) {
	s.CustomId = utils.AssignId("")
	CallbackTasks[s.CustomId] = handler
}

func (s *SelectMenu) ToComponent() map[string]interface{} {
	s.CustomId = utils.AssignId(s.CustomId)
	menu := map[string]interface{}{"type": 3, "custom_id": s.CustomId}
	if s.Placeholder != "" {
		menu["placeholder"] = s.Placeholder
	}
	if s.MinValues >= 0 {
		menu["min_values"] = s.MinValues
	} else {
		log.Println("MinValues must be greater than or equal to 0")
	}
	if s.MaxValues <= 25 {
		menu["max_values"] = s.MaxValues
	} else {
		log.Println("MaxValues must be less than or equal to 25")
	}
	if s.Disabled {
		menu["disabled"] = true
	}
	if len(s.Options) > 0 {
		menu["options"] = s.Options
	}
	return menu
}

type ActionRow struct {
	Buttons    []Button
	SelectMenu SelectMenu
}

type View struct {
	Timeout    float64
	ActionRows []ActionRow
	CustomId   string
}

func (v *View) OnTimeout(handler func(bot user.User, interaction Interaction)) {
	TimeoutTasks[".."] = handler
}

func (v *View) ToComponent() []interface{} {
	if v.Timeout == 0 {
		v.Timeout = 15 * 60
	}
	var undo = map[string]bool{}
	var c []interface{}
	if len(v.ActionRows) > 0 && len(v.ActionRows) <= 5 {
		for _, row := range v.ActionRows {
			num := 0
			tmp := map[string]interface{}{
				"type":       1,
				"components": []interface{}{},
			}
			for _, button := range row.Buttons {
				if num < 5 {
					undo[button.CustomId] = true
					tmp["components"] = append(tmp["components"].([]interface{}), button.ToComponent())
				}
				num++
			}
			if len(row.SelectMenu.Options) > 0 {
				if num == 0 {
					undo[row.SelectMenu.CustomId] = true
					tmp["components"] = append(tmp["components"].([]interface{}), row.SelectMenu.ToComponent())
				} else {
					log.Println("Single ActionRow can contain either 1x SelectMenu or max 5x Buttons")
				}
			}
			if len(tmp["components"].([]interface{})) > 0 {
				c = append(c, tmp)
				go utils.ScheduleDeletion(v.Timeout, CallbackTasks, undo)
			}
		}
	}
	return c
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

type Interaction struct {
	ID             string                 `json:"id"`
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

func FromData(payload interface{}) *Interaction {
	i := &Interaction{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, i)
	return i
}

func (i *Interaction) SendResponse(resp Response) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New(
		"POST", path, map[string]interface{}{"type": 4, "data": resp.ToBody()}, "", resp.Files)
	go r.Request()
}

func (i *Interaction) Defer() {
	body := map[string]interface{}{"type": 6}
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (i *Interaction) SendModal(modal Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, modal.ToBody(), "", nil)
	go r.Request()
}

func (i *Interaction) SendFollowup(resp Response) {
	path := fmt.Sprintf("/webhooks/%s/%s", i.ApplicationId, i.Token)
	r := router.New("POST", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}
