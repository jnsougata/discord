package disgo

type Draft struct {
	Content string
	Embed   Embed
	Embeds  []Embed
	//AllowedMentions []string
	TTS            bool
	View           View
	File           File
	Files          []File
	SuppressEmbeds bool
	Reference      any
}

func (d *Draft) marshal() map[string]interface{} {
	body := map[string]interface{}{}
	if d.Content != "" {
		body["content"] = d.Content
	}
	if checkTrueEmbed(d.Embed) {
		d.Embeds = append([]Embed{d.Embed}, d.Embeds...)
	}
	for i, em := range d.Embeds {
		if !checkTrueEmbed(em) {
			d.Embeds = append(d.Embeds[:i], d.Embeds[i+1:]...)
		}
	}
	if len(d.Embeds) > 10 {
		d.Embeds = d.Embeds[:10]
	}
	body["embeds"] = []map[string]interface{}{}
	for _, em := range d.Embeds {
		body["embeds"] = append(body["embeds"].([]map[string]interface{}), em.Marshal())
	}
	if checkTrueFile(d.File) {
		d.Files = append([]File{d.File}, d.Files...)
	}
	body["attachments"] = []map[string]interface{}{}
	for i, f := range d.Files {
		if checkTrueFile(f) {
			a := map[string]interface{}{
				"id":          i,
				"filename":    f.Name,
				"description": f.Description,
			}
			body["attachments"] = append(body["attachments"].([]map[string]interface{}), a)
		} else {
			d.Files = append(d.Files[:i], d.Files[i+1:]...)
		}
	}
	return body
}
