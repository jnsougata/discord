package attachment

import "encoding/json"

type Attachment struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Description string `json:"description"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Ephemeral   bool   `json:"ephemeral"`
}

type Partial struct {
	Id          string `json:"id"`
	Filename    string `json:"filename"`
	Description string `json:"description"`
}

func Unmarshal(payload interface{}) *Attachment {
	attachment := &Attachment{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, attachment)
	return attachment
}
