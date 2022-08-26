package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Message struct {
	Id                 string                   `json:"id"`
	ChannelId          string                   `json:"channel_id"`
	Author             User                     `json:"author"`
	Content            string                   `json:"content"`
	Timestamp          string                   `json:"timestamp"`
	EditedTimestamp    string                   `json:"edited_timestamp"`
	TTS                bool                     `json:"tts"`
	MentionEveryone    bool                     `json:"mention_everyone"`
	Mentions           []map[string]interface{} `json:"mentions"`
	RoleMentions       []string                 `json:"role_mentions"`
	ChannelMentions    []string                 `json:"channel_mentions"`
	Attachments        []map[string]interface{} `json:"attachments"`
	Embeds             []Embed                  `json:"embeds"`
	Reactions          []map[string]interface{} `json:"reactions"`
	Pinned             bool                     `json:"pinned"`
	WebhookId          string                   `json:"webhook_id"`
	Type               int                      `json:"types"`
	Activity           map[string]interface{}   `json:"activity"`
	Application        map[string]interface{}   `json:"application"`
	ApplicationId      string                   `json:"application_id"`
	MessageReference   map[string]interface{}   `json:"message_reference"`
	Flags              int                      `json:"flags"`
	ReferencedMessages map[string]interface{}   `json:"reference"`
	Interaction        map[string]interface{}   `json:"interaction"`
	Thread             map[string]interface{}   `json:"thread"`
	Components         []map[string]interface{} `json:"components"`
	Stickers           []map[string]interface{} `json:"sticker_items"`
	state              *state
}

func (m *Message) Reply(message Draft) (Message, error) {
	body, err := message.marshal()
	body["message_reference"] = map[string]interface{}{"message_id": m.Id}
	path := fmt.Sprintf("/channels/%s/messages", m.ChannelId)
	r := multipartReq("POST", path, body, m.state.Token, message.Files...)
	bs, _ := io.ReadAll(r.fire().Body)
	var msg Message
	_ = json.Unmarshal(bs, &msg)
	msg.state = m.state
	if message.DeleteAfter > 0 {
		go func() {
			time.Sleep(time.Second * time.Duration(message.DeleteAfter))
			msg.Delete()
		}()
	}
	return msg, err
}

func (m *Message) Delete() {
	path := fmt.Sprintf("/channels/%s/messages/%s", m.ChannelId, m.Id)
	go minimalReq("DELETE", path, nil, m.state.Token).fire()
}
