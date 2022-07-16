package models

type InteractionMessage struct {
	Content         string   `json:"content"`
	Embeds          []Embed  `json:"embeds"`
	AllowedMentions []string `json:"allowed_mentions"`
	Tts             bool     `json:"tts"`
	Flags           int      `json:"flags"`
	//Components      []types.Component  `json:"components"`
	//Attachments     []types.Attachment `json:"attachments"`
}

func (i *InteractionMessage) ToBody() map[string]interface{} {
	return map[string]interface{}{
		"type": 4,
		"data": map[string]interface{}{
			"content": i.Content,
			"embeds":  i.Embeds,
			"tts":     i.Tts,
			"flags":   i.Flags,
		},
	}
}
