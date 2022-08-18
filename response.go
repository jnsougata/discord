package discord

type Response struct {
	Content         string
	Embed           Embed
	Embeds          []Embed
	AllowedMentions []string
	TTS             bool
	Ephemeral       bool
	SuppressEmbeds  bool
	View            View
	File            File
	Files           []File
}

func (resp *Response) marshal() (map[string]interface{}, error) {
	flag := 0
	body := map[string]interface{}{}
	if resp.Content != "" {
		body["content"] = resp.Content
	}
	if checkTrueEmbed(resp.Embed) {
		resp.Embeds = append([]Embed{resp.Embed}, resp.Embeds...)
	}
	for i, em := range resp.Embeds {
		if !checkTrueEmbed(em) {
			resp.Embeds = append(resp.Embeds[:i], resp.Embeds[i+1:]...)
		}
	}
	if len(resp.Embeds) > 10 {
		resp.Embeds = resp.Embeds[:10]
	}
	body["embeds"] = []map[string]interface{}{}
	for _, em := range resp.Embeds {
		emd, err := em.marshal()
		if err != nil {
			return nil, err
		} else {
			body["embeds"] = append(body["embeds"].([]map[string]interface{}), emd)
		}
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
	view, err := resp.View.marshal()
	if err == nil {
		body["components"] = view
	} else {
		return nil, err
	}
	if checkTrueFile(resp.File) {
		resp.Files = append([]File{resp.File}, resp.Files...)
	}
	body["attachments"] = []map[string]interface{}{}
	for i, f := range resp.Files {
		if checkTrueFile(f) {
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
	return body, nil
}
