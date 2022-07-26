package component

import (
	"encoding/json"
	"fmt"
	"github.com/jnsougata/disgo/core/attachment"
	"github.com/jnsougata/disgo/core/embed"
	"github.com/jnsougata/disgo/core/emoji"
	"github.com/jnsougata/disgo/core/modal"
	"github.com/jnsougata/disgo/core/router"
	"github.com/jnsougata/disgo/core/user"
	"log"
)

const (
	PrimaryButtonStyle   = 1
	SecondaryButtonStyle = 2
	SuccessButtonStyle   = 3
	DangerButtonStyle    = 4
	LinkButtonStyle      = 5
)

var CallbackFactory = map[string]interface{}{}

type Message struct {
	Content         string
	Embeds          []embed.Embed
	AllowedMentions []string
	Tts             bool
	Flags           int
	View            View
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

type SelectOption struct {
	Label       string        `json:"label"`
	Value       string        `json:"value"`
	Description string        `json:"description"`
	Emoji       emoji.Partial `json:"emoji"`
	Default     bool          `json:"default"`
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
	return map[string]interface{}{
		"type":      2,
		"style":     b.Style,
		"label":     b.Label,
		"emoji":     b.Emoji,
		"custom_id": b.CustomId,
		"url":       b.URL,
		"disabled":  b.Disabled,
	}
}

func (b *Button) AddCallback(handler func(bot user.User, interaction Interaction)) {
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

func (s *SelectMenu) AddCallback(handler func(bot user.User, interaction Interaction, values ...string)) {
	CallbackFactory[s.CustomId] = handler
}

func (s *SelectMenu) ToComponent() map[string]interface{} {
	return map[string]interface{}{
		"type":        3,
		"custom_id":   s.CustomId,
		"options":     s.Options,
		"placeholder": s.Placeholder,
		"min_values":  s.MinValues,
		"max_values":  s.MaxValues,
		"disabled":    s.Disabled,
	}
}

type ActionRow struct {
	Buttons    []Button
	SelectMenu SelectMenu
}

type View struct {
	ActionRows []ActionRow
	Callback   func(interaction Interaction, values ...string)
}

func (v *View) ToComponent() []interface{} {
	var c []interface{}
	ids := map[string]bool{}
	if len(v.ActionRows) > 0 && len(v.ActionRows) <= 5 {
		for _, row := range v.ActionRows {
			numButtons := 0
			tmp := map[string]interface{}{
				"type":       1,
				"components": []interface{}{},
			}
			for _, button := range row.Buttons {
				numButtons++
				if button.CustomId == "" && button.Style != LinkButtonStyle {
					log.Println(
						fmt.Sprintf("CustomId must be provided with non-link button `%s`", button.Label))
				} else if _, ok := ids[button.CustomId]; !ok {
					if numButtons <= 5 {
						ids[button.CustomId] = true
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
				if row.SelectMenu.CustomId == "" {
					log.Println("CustomId must be provided with Select Menu")
				} else if row.SelectMenu.MaxValues > 25 {
					log.Println("MaxValues must be less than or equals to 25")
				} else if row.SelectMenu.MinValues > row.SelectMenu.MaxValues {
					log.Println("MinValues must be less than or equals to MaxValues")
				} else if row.SelectMenu.MinValues < 0 {
					log.Println("MinValues must be greater than or equals to 0")
				} else if _, ok := ids[row.SelectMenu.CustomId]; !ok {
					if numButtons == 0 {
						ids[row.SelectMenu.CustomId] = true
						tmp["components"] = append(tmp["components"].([]interface{}), row.SelectMenu.ToComponent())
					} else {
						log.Println("An Action Row can contain one of these: (1x SelectMenu) or (max 5x Buttons)")
					}
				} else {
					log.Println(
						fmt.Sprintf("CustomId `%s` already used with a previous component", row.SelectMenu.CustomId))
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
	Member         interface{}            `json:"member"`
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

func (i *Interaction) Defer() {
	payload := map[string]interface{}{"type": 6}
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, payload, "")
	go r.Request()
}

func (i *Interaction) SendModal(modal modal.Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, modal.ToBody(), "")
	go r.Request()
}

func (i *Interaction) SendFollowup(message Message) {
	path := fmt.Sprintf("/webhooks/%s/%s", i.ApplicationID, i.Token)
	r := router.New("POST", path, message.ToBody(), "")
	go r.Request()
}
