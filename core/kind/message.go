package kind

import (
	"encoding/json"
	"log"
)

type Message struct {
	Id                 string                   `json:"id"`
	ChannelId          string                   `json:"channel_id"`
	Author             map[string]interface{}   `json:"author"`
	Content            string                   `json:"content"`
	Timestamp          string                   `json:"timestamp"`
	EditedTimestamp    string                   `json:"edited_timestamp"`
	TTS                bool                     `json:"tts"`
	MentionEveryone    bool                     `json:"mention_everyone"`
	Mentions           []map[string]interface{} `json:"mentions"`
	RoleMentions       []string                 `json:"role_mentions"`
	ChannelMentions    []string                 `json:"channel_mentions"`
	Attachments        []map[string]interface{} `json:"attachments"`
	Embeds             []map[string]interface{} `json:"embeds"`
	Reactions          []map[string]interface{} `json:"reactions"`
	Pinned             bool                     `json:"pinned"`
	WebhookId          string                   `json:"webhook_id"`
	Type               int                      `json:"kind"`
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
}

func BuildMessage(payload interface{}) *Message {
	msg := &Message{}
	data, _ := json.Marshal(payload)
	err := json.Unmarshal(data, msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}
