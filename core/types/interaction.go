package types

import (
	"encoding/json"
	"fmt"
	"github.com/disgo/core/objects"
	"github.com/disgo/core/router"
	"log"
	"time"
)

type Interaction struct {
	Type   int    `json:"type"`
	Token  string `json:"token"`
	Member struct {
		User struct {
			ID            string `json:"id"`
			Username      string `json:"username"`
			Avatar        string `json:"avatar"`
			Discriminator string `json:"discriminator"`
			PublicFlags   int    `json:"public_flags"`
		} `json:"user"`
		Roles        []string    `json:"roles"`
		PremiumSince interface{} `json:"premium_since"`
		Permissions  string      `json:"permissions"`
		Pending      bool        `json:"pending"`
		Nick         interface{} `json:"nick"`
		Mute         bool        `json:"mute"`
		JoinedAt     time.Time   `json:"joined_at"`
		IsPending    bool        `json:"is_pending"`
		Deaf         bool        `json:"deaf"`
	} `json:"member"`
	ID             string `json:"id"`
	GuildID        string `json:"guild_id"`
	AppPermissions string `json:"app_permissions"`
	GuildLocale    string `json:"guild_locale"`
	Locale         string `json:"locale"`
	Data           struct {
		Options []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"options"`
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"data"`
	ChannelID string `json:"channel_id"`
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

func (i *Interaction) SendResponse(message *objects.InteractionMessage) {
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

func (i *Interaction) SendModal(modal *objects.Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, modal.ToBody(), "")
	go r.Request()
}

func (i *Interaction) SendAutoComplete(choices ...*objects.Choice) {
	payload := map[string]interface{}{
		"type": 8,
		"data": map[string]interface{}{"choices": choices},
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.ID, i.Token)
	r := router.New("POST", path, payload, "")
	go r.Request()
}
