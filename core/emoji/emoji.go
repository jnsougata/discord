package emoji

import "encoding/json"

type Partial struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Animated bool   `json:"animated"`
}

type Emoji struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	RequireColons bool     `json:"require_colons"`
	Managed       bool     `json:"managed"`
	Animated      bool     `json:"animated"`
	Available     bool     `json:"available"`
}

func FromData(payload interface{}) *Emoji {
	emoji := &Emoji{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, emoji)
	return emoji
}
