package disgo

type PartialEmoji struct {
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
