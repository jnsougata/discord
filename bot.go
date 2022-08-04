package disgo

import (
	"encoding/json"
)

type BotUser struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	Banner        string `json:"banner"`
	Color         int    `json:"accent_color"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Flags         int    `json:"flags"`
	PublicFlags   int    `json:"public_flags"`
	Latency       int64  `json:"latency"`
	IsReady       bool
	Guilds        map[string]*Guild
}

func unmarshal(payload interface{}) *BotUser {
	bot := &BotUser{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, bot)
	return bot
}
