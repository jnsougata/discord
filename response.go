package discord

type Response struct {
	Content        string
	Embed          Embed
	Embeds         []Embed
	TTS            bool
	Ephemeral      bool
	SuppressEmbeds bool
	View           View
	File           File
	Files          []File
	//AllowedMentions []string
}

func (r *Response) marshal() (map[string]interface{}, error) {
	flag := 0
	body := map[string]interface{}{}
	if r.Content != "" {
		body["content"] = r.Content
	}
	if checkTrueEmbed(r.Embed) {
		r.Embeds = append([]Embed{r.Embed}, r.Embeds...)
	}
	for i, em := range r.Embeds {
		if !checkTrueEmbed(em) {
			r.Embeds = append(r.Embeds[:i], r.Embeds[i+1:]...)
		}
	}
	if len(r.Embeds) > 10 {
		r.Embeds = r.Embeds[:10]
	}
	body["embeds"] = []map[string]interface{}{}
	for _, em := range r.Embeds {
		emd, err := em.marshal()
		if err != nil {
			return nil, err
		} else {
			body["embeds"] = append(body["embeds"].([]map[string]interface{}), emd)
		}
	}
	if r.TTS {
		body["tts"] = true
	}
	if r.Ephemeral {
		flag |= 1 << 6
	}
	if r.SuppressEmbeds {
		flag |= 1 << 2
	}
	if r.Ephemeral || r.SuppressEmbeds {
		body["flags"] = flag
	}
	view, err := r.View.marshal()
	if err == nil {
		body["components"] = view
	} else {
		return nil, err
	}
	if checkTrueFile(r.File) {
		r.Files = append([]File{r.File}, r.Files...)
	}
	body["attachments"] = []map[string]interface{}{}
	for i, f := range r.Files {
		if checkTrueFile(f) {
			a := map[string]interface{}{
				"id":          i,
				"filename":    f.Name,
				"description": f.Description,
			}
			body["attachments"] = append(body["attachments"].([]map[string]interface{}), a)
		} else {
			r.Files = append(r.Files[:i], r.Files[i+1:]...)
		}
	}
	return body, nil
}
