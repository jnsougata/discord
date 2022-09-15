package discord

import (
	"encoding/json"
	"fmt"
	"io"
)

type Interaction struct {
	Id             string                 `json:"id"`
	ApplicationId  string                 `json:"application_id"`
	Type           int                    `json:"type"`
	Data           map[string]interface{} `json:"data"`
	GuildId        string                 `json:"guild_id"`
	ChannelId      string                 `json:"channel_id"`
	User           User                   `json:"user"`
	Token          string                 `json:"token"`
	Version        int                    `json:"version"`
	AppPermissions string                 `json:"app_permissions"`
	Locale         string                 `json:"locale"`
	GuildLocale    string                 `json:"guild_locale"`
	TargetUser     User
	TargetMessage  Message
	Channel        Channel
	Guild          Guild
	Author         Member
	state          *state
}

func (i *Interaction) OriginalResponse() Message {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", i.ApplicationId, i.Token)
	r := minimalReq("GET", path, nil, "")
	bs, _ := io.ReadAll(r.fire().Body)
	var m Message
	_ = json.Unmarshal(bs, &m)
	m.state = i.state
	return m
}

func (i *Interaction) SendResponse(response *Response) error {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	body, err := response.marshal()
	payload := map[string]interface{}{"type": 4, "data": body}
	r := multipartReq(
		"POST", path, payload, "", response.Files...)
	go r.fire()
	return err
}

func (i *Interaction) DeferResponse(ephemeral bool) {
	body := map[string]interface{}{}
	if i.Type == 2 {
		body["type"] = 5
		if ephemeral {
			body["data"] = map[string]interface{}{"flags": 1 << 6}
		}
	} else {
		body["type"] = 6
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	r := minimalReq("POST", path, body, "")
	go r.fire()
}

func (i *Interaction) ModalResponse(modal *Modal) error {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	body, err := modal.marshal()
	r := minimalReq("POST", path, body, "")
	go r.fire()
	return err
}

func (i *Interaction) FollowupResponse(response *Response) (Message, error) {
	path := fmt.Sprintf("/webhooks/%s/%s", i.ApplicationId, i.Token)
	body, err := response.marshal()
	r := multipartReq("POST", path, body, "", response.Files...)
	m := make(chan Message, 1)
	go func() {
		bs, _ := io.ReadAll(r.fire().Body)
		var msg Message
		_ = json.Unmarshal(bs, &msg)
		msg.state = i.state
		m <- msg
	}()
	return <-m, err
}

func (i *Interaction) EditResponse(response *Response) error {
	body, err := response.marshal()
	if len(response.Embeds) == 0 {
		body["embeds"] = []map[string]interface{}{}
	}
	if len(response.Files) == 0 {
		body["attachments"] = []map[string]interface{}{}
	}
	if len(response.View.rows) == 0 {
		body["components"] = []map[string]interface{}{}
	}
	if i.Type == 2 {
		path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", i.ApplicationId, i.Token)
		r := multipartReq("PATCH", path, body, "", response.Files...)
		go r.fire()
	} else {
		path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
		nbody := map[string]interface{}{"type": 7, "data": body}
		r := multipartReq("POST", path, nbody, "", response.Files...)
		go r.fire()
	}
	return err
}

func (i *Interaction) DeleteResponse() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", i.ApplicationId, i.Token)
	r := minimalReq("DELETE", path, nil, "")
	go r.fire()
}
