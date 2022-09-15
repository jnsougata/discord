package discord

type ChannelMessage struct {
	Content        string
	Embed          Embed
	Embeds         []Embed
	TTS            bool
	View           View
	File           File
	Files          []File
	SuppressEmbeds bool
	Reference      any
	DeleteAfter    float64
	//AllowedMentions []string
}

func (msg *ChannelMessage) marshal() (map[string]interface{}, error) {
	body := map[string]interface{}{}
	if msg.Content != "" {
		body["content"] = msg.Content
	}
	if checkTrueEmbed(msg.Embed) {
		msg.Embeds = append([]Embed{msg.Embed}, msg.Embeds...)
	}
	for i, em := range msg.Embeds {
		if !checkTrueEmbed(em) {
			msg.Embeds = append(msg.Embeds[:i], msg.Embeds[i+1:]...)
		}
	}
	if len(msg.Embeds) > 10 {
		msg.Embeds = msg.Embeds[:10]
	}
	body["embeds"] = []map[string]interface{}{}
	for _, em := range msg.Embeds {
		emd, err := em.marshal()
		if err != nil {
			return nil, err
		} else {
			body["embeds"] = append(body["embeds"].([]map[string]interface{}), emd)
		}
	}
	if checkTrueFile(msg.File) {
		msg.Files = append([]File{msg.File}, msg.Files...)
	}
	body["attachments"] = []map[string]interface{}{}
	for i, f := range msg.Files {
		if checkTrueFile(f) {
			a := map[string]interface{}{
				"id":          i,
				"filename":    f.Name,
				"description": f.Description,
			}
			body["attachments"] = append(body["attachments"].([]map[string]interface{}), a)
		} else {
			msg.Files = append(msg.Files[:i], msg.Files[i+1:]...)
		}
	}
	return body, nil
}
