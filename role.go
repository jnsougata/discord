package disgo

import "encoding/json"

type Role struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Color        int    `json:"color"`
	Hoist        bool   `json:"hoist"`
	Icon         string `json:"icon"`
	UnicodeEmoji bool   `json:"unicode_emoji"`
	Position     int    `json:"position"`
	Permissions  string `json:"permissions"`
	Managed      bool   `json:"managed"`
	Mentionable  bool   `json:"mentionable"`
	Tags         string `json:"tags"`
	GuildId      string `json:"guild_id"`
}

func unmarshalRole(payload interface{}) *Role {
	role := &Role{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, role)
	return role
}
