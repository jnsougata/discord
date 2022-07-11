package models

type InteractionMessage struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
	//	AllowedMentions []string                 `json:"allowed_mentions"`
	Tts   bool `json:"tts"`
	Flags int  `json:"flags"`
	//	Components      []map[string]interface{} `json:"components"`
	//	Attachments     []map[string]interface{} `json:"attachments"`
}

func (i *InteractionMessage) ToBody() map[string]interface{} {
	return map[string]interface{}{
		"kind": 4,
		"data": map[string]interface{}{
			"content": i.Content,
			"embeds":  i.Embeds,
			"tts":     i.Tts,
			"flags":   i.Flags,
		},
	}
}
