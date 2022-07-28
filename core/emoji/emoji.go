package emoji

import "encoding/json"

type Partial struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Animated bool   `json:"animated,omitempty"`
}

type Emoji struct {
	Id            int64    `json:"id,string"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	Managed       bool     `json:"managed"`
	Animated      bool     `json:"animated"`
	Available     bool     `json:"available"`
	RequireColons bool     `json:"require_colons"`
}

func FromData(payload interface{}) *Emoji {
	emoji := &Emoji{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, emoji)
	return emoji
}
