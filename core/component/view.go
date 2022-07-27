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

var CallbackFactory = map[string]interface{}{}
var Ids = map[string]bool{}

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

type SelectOption struct {
	Label       string        `json:"label"`
	Value       string        `json:"value"`
	Description string        `json:"description,omitempty"`
	Emoji       emoji.Partial `json:"emoji,omitempty"`
	Default     bool          `json:"default,omitempty"`
}

type Button struct {
	Style    int
	Label    string
	Emoji    emoji.Partial
	CustomId string
	URL      string
	Disabled bool
}

func (b *Button) ToComponent() map[string]interface{} {
	if b.CustomId == "" && b.Style != LinkButton {
		panic("CustomId is required for each Button")
	}
	btn := map[string]interface{}{"type": 2, "custom_id": b.CustomId}
	if b.Style != 0 {
		btn["style"] = b.Style
	}
	if b.Label != "" {
		btn["label"] = b.Label
	}
	if b.Emoji.ID != "" {
		btn["emoji"] = b.Emoji
	}
	if b.URL != "" {
		btn["url"] = b.URL
	}
	if b.Disabled {
		btn["disabled"] = true
	}
	return btn
}

func (b *Button) Callback(handler func(bot user.User, interaction Interaction)) {
	CallbackFactory[b.CustomId] = handler
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

func (s *SelectMenu) Callback(handler func(bot user.User, interaction Interaction, values ...string)) {
	CallbackFactory[s.CustomId] = handler
}

func (s *SelectMenu) ToComponent() map[string]interface{} {
	if s.CustomId == "" {
		log.Println("CustomId is required for each SelectMenu")
	}
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
	ActionRows []ActionRow
}

func (v *View) ToComponent() []interface{} {
	var c []interface{}
	if len(v.ActionRows) > 0 && len(v.ActionRows) <= 5 {
		for _, row := range v.ActionRows {
			numButtons := 0
			tmp := map[string]interface{}{
				"type":       1,
				"components": []interface{}{},
			}
			for _, button := range row.Buttons {
				numButtons++
				if _, ok := Ids[button.CustomId]; !ok {
					if numButtons <= 5 {
						Ids[button.CustomId] = true
						tmp["components"] = append(tmp["components"].([]interface{}), button.ToComponent())
					} else {
						log.Println("An Action Row can either contain max 5x Buttons")
					}
				} else {
					log.Println(
						fmt.Sprintf("CustomId `%s` already used with a previous component", button.CustomId))
				}
			}
			if len(row.SelectMenu.Options) > 0 {
				if _, ok := Ids[row.SelectMenu.CustomId]; !ok {
					if numButtons == 0 {
						Ids[row.SelectMenu.CustomId] = true
						tmp["components"] = append(tmp["components"].([]interface{}), row.SelectMenu.ToComponent())
					} else {
						log.Println("An Action Row can contain one of these: (1x SelectMenu) or (max 5x Buttons)")
					}
				} else {
					log.Println(
						fmt.Sprintf(
							"CustomId `%s` already used with a previous component", row.SelectMenu.CustomId))
				}
			}
			if len(tmp["components"].([]interface{})) > 0 {
				c = append(c, tmp)
			}
		}
	}
	return c
}

type Data struct {
	Type   int      `json:"component_type"`
	Id     string   `json:"custom_id"`
	Values []string `json:"values"`
}

type Interaction struct {
	ID             string                 `json:"id"`
	ApplicationID  string                 `json:"application_id"`
	Type           int                    `json:"type"`
	Data           Data                   `json:"data"`
	GuildID        string                 `json:"guild_id"`
	ChannelID      string                 `json:"channel_id"`
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
	path := fmt.Sprintf("/webhooks/%s/%s", i.ApplicationID, i.Token)
	r := router.New("POST", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}
