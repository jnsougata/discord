package discord

import (
	"encoding/json"
	"errors"
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
	token          string
}

func (i *Interaction) OriginalResponse() Message {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", i.ApplicationId, i.Token)
	r := minimalReq("GET", path, nil, "")
	bs, _ := io.ReadAll(r.fire().Body)
	var m Message
	_ = json.Unmarshal(bs, &m)
	return m
}

func (i *Interaction) Send(resp Response) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	body, err := resp.marshal()
	if err != nil {
		interactionalExceptionWrapper(err, i.Id)
	} else {
		r := multipartReq(
			"POST", path, map[string]interface{}{"type": 4, "data": body}, "", resp.Files...)
		go r.fire()
	}
}

func (i *Interaction) Defer(ephemeral bool) {
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

func (i *Interaction) SendModal(modal Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
	r := minimalReq("POST", path, modal.marshal(), "")
	go r.fire()
}

func (i *Interaction) SendFollowup(resp Response) (Followup, error) {
	path := fmt.Sprintf("/webhooks/%s/%s", i.ApplicationId, i.Token)
	body, err := resp.marshal()
	if err != nil {
		interactionalExceptionWrapper(err, i.Id)
		return Followup{}, err
	} else {
		r := multipartReq("POST", path, body, "", resp.Files...)
		fl := make(chan Followup, 1)
		go func() {
			bs, _ := io.ReadAll(r.fire().Body)
			var msg Message
			_ = json.Unmarshal(bs, &msg)
			fl <- Followup{
				Id:            msg.Id,
				Content:       msg.Content,
				Embeds:        msg.Embeds,
				ChannelId:     i.ChannelId,
				Flags:         msg.Flags,
				token:         i.Token,
				applicationId: i.ApplicationId,
			}
		}()
		val, ok := <-fl
		if ok {
			return val, nil
		}
		return Followup{}, errors.New("no followup received")
	}
}

func (i *Interaction) Edit(resp Response) {
	body, err := resp.marshal()
	if err != nil {
		interactionalExceptionWrapper(err, i.Id)
	} else {
		if i.Type == 2 {
			path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", i.ApplicationId, i.Token)
			r := multipartReq("PATCH", path, body, "", resp.Files...)
			go r.fire()
		} else {
			path := fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token)
			nbody := map[string]interface{}{"type": 7, "data": body}
			r := multipartReq("POST", path, nbody, "", resp.Files...)
			go r.fire()
		}
	}
}

func (i *Interaction) Delete() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", i.ApplicationId, i.Token)
	r := minimalReq("DELETE", path, nil, "")
	go r.fire()
}

func interactionalExceptionWrapper(err error, id string) {
	fmt.Println(fmt.Sprintf("Interaction `%s` panicked!", id))
	fmt.Println("Error:", err)
}
