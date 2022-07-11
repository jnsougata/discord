package kind

import (
	"encoding/json"
	"fmt"
	"github.com/disgo/core/models"
	"github.com/disgo/core/router"
	"log"
)

type InteractionData struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	Type     int                    `json:"type"`
	Resolved map[string]interface{} `json:"resolved"`
	Options  map[string]interface{} `json:"options"`

	// guild where interaction was created
	GuildID string `json:"guild_id"`

	// id of the targeted context
	TargetId string `json:"target_id"`
}

type Interaction struct {
	ID             string           `json:"id"`
	ApplicationID  string           `json:"application_id"`
	Type           int              `json:"type"`
	Data           *InteractionData `json:"data"`
	GuildID        string           `json:"guild_id"`
	ChannelID      string           `json:"channel_id"`
	Member         interface{}      `json:"member"`
	User           *User            `json:"user"`
	Token          string           `json:"token"`
	Version        int              `json:"version"`
	Message        interface{}      `json:"message"`
	AppPermissions string           `json:"app_permissions"`
	Locale         string           `json:"locale"`
	GuildLocale    string           `json:"guild_locale"`
}

func BuildInteraction(payload interface{}) *Interaction {
	i := &Interaction{}
	data, _ := json.Marshal(payload)
	err := json.Unmarshal(data, i)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func (i *Interaction) SendResponse(message *models.InteractionMessage) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, message.ToBody(), "")
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

func (i *Interaction) SendModal(modal *models.Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, modal.ToBody(), "")
	go r.Request()
}

func (i *Interaction) SendAutoComplete(choices ...*models.Choice) {
	payload := map[string]interface{}{
		"type": 8,
		"data": map[string]interface{}{"choices": choices},
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, payload, "")
	go r.Request()
}

func (i *Interaction) SendFollowup(choices ...*models.Choice) {
	payload := map[string]interface{}{
		"type": 8,
		"data": map[string]interface{}{"choices": choices},
	}
	path := fmt.Sprintf("/webhooks/%s/%s", i.ApplicationID, i.Token)
	r := router.New("POST", path, payload, "")
	go r.Request()
}
