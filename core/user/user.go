package user

import (
	"encoding/json"
)

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	System        bool   `json:"system"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	Banner        string `json:"banner"`
	Color         int    `json:"accent_color"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
}

type Bot struct {
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
}

func MakeBot(payload interface{}) *Bot {
	bot := &Bot{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, bot)
	return bot
}

func FromData(payload interface{}) *Bot {
	user := &Bot{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, user)
	return user
}
