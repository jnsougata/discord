package kind

import (
	"encoding/json"
	"log"
)

type User struct {
	ID            string `json:"id"`
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

func BuildUser(payload interface{}) *User {
	u := &User{}
	data, _ := json.Marshal(payload)
	err := json.Unmarshal(data, u)
	if err != nil {
		log.Fatal(err)
	}
	return u
}
