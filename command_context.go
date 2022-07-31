package main

import (
	"fmt"
)

type CommandResponse struct {
	Content         string
	Embed           Embed // gets the priority over embeds
	Embeds          []Embed
	AllowedMentions []string
	TTS             bool
	Ephemeral       bool
	SuppressEmbeds  bool
	View            View
	File            File // gets the priority over files
	Files           []File
}

func (resp *CommandResponse) Marshal() map[string]interface{} {
	flag := 0
	body := map[string]interface{}{}
	if resp.Content != "" {
		body["content"] = resp.Content
	}
	if CheckTrueEmbed(resp.Embed) {
		resp.Embeds = append([]Embed{resp.Embed}, resp.Embeds...)
	}
	for i, em := range resp.Embeds {
		if !CheckTrueEmbed(em) {
			resp.Embeds = append(resp.Embeds[:i], resp.Embeds[i+1:]...)
		}
	}
	if len(resp.Embeds) > 10 {
		resp.Embeds = resp.Embeds[:10]
	}
	body["embeds"] = []map[string]interface{}{}
	for _, em := range resp.Embeds {
		body["embeds"] = append(body["embeds"].([]map[string]interface{}), em.Marshal())
	}
	if len(resp.AllowedMentions) > 0 && len(resp.AllowedMentions) <= 100 {
		body["allowed_mentions"] = resp.AllowedMentions
	}
	if resp.TTS {
		body["tts"] = true
	}
	if resp.Ephemeral {
		flag |= 1 << 6
	}
	if resp.SuppressEmbeds {
		flag |= 1 << 2
	}
	if resp.Ephemeral || resp.SuppressEmbeds {
		body["flags"] = flag
	}
	if len(resp.View.ActionRows) > 0 {
		body["components"] = resp.View.ToComponent()
	}
	if CheckTrueFile(resp.File) {
		resp.Files = append([]File{resp.File}, resp.Files...)
	}
	body["attachments"] = []map[string]interface{}{}
	for i, f := range resp.Files {
		if CheckTrueFile(f) {
			a := map[string]interface{}{
				"id":          i,
				"filename":    f.Name,
				"description": f.Description,
			}
			body["attachments"] = append(body["attachments"].([]map[string]interface{}), a)
		} else {
			resp.Files = append(resp.Files[:i], resp.Files[i+1:]...)
		}
	}
	return body
}

type CommandContext struct {
	Id            string
	Token         string
	ApplicationId string
}

func (c *CommandContext) SendResponse(resp CommandResponse) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	body := map[string]interface{}{
		"type": 4,
		"data": resp.Marshal(),
	}
	r := MultipartReq("POST", path, body, "", resp.Files)
	go r.Request()
}

func (c *CommandContext) Ack() {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	body := map[string]interface{}{"type": 1}
	r := MinimalReq("POST", path, body, "")
	go r.Request()
}

func (c *CommandContext) Defer(ephemeral bool) {
	body := map[string]interface{}{"type": 5}
	if ephemeral {
		body["data"] = map[string]interface{}{"flags": 1 << 6}
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := MultipartReq("POST", path, body, "", nil)
	go r.Request()
}

func (c *CommandContext) SendModal(modal Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := MultipartReq("POST", path, modal.ToBody(), "", nil)
	go r.Request()
}

func (c *CommandContext) SendAutoComplete(choices ...Choice) {
	body := map[string]interface{}{
		"type": 8,
		"data": map[string]interface{}{"choices": choices},
	}
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.Id, c.Token)
	r := MultipartReq("POST", path, body, "", nil)
	go r.Request()
}

func (c *CommandContext) SendFollowup(resp CommandResponse) {
	path := fmt.Sprintf("/webhooks/%s/%s", c.ApplicationId, c.Token)
	r := MultipartReq("POST", path, resp.Marshal(), "", resp.Files)
	go r.Request()
}

func (c *CommandContext) DeleteOriginalResponse() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := MinimalReq("DELETE", path, nil, "")
	go r.Request()
}

func (c *CommandContext) EditOriginalResponse(resp CommandResponse) {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := MultipartReq("PATCH", path, resp.Marshal(), "", resp.Files)
	go r.Request()
}
