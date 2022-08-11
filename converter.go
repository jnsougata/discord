package discord

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type Converter struct {
	token   string
	payload interface{}
}

func (c Converter) Bot() *BotUser {
	bot := &BotUser{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, bot)
	return bot
}

func (c Converter) Message() *Message {
	msg := &Message{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, msg)
	msg.token = c.token
	return msg
}

func (c Converter) Guild() *Guild {
	raw := c.payload.(map[string]interface{})
	iconHash := raw["icon"]
	bannerHash := raw["banner"]
	splashHash := raw["splash"]
	dSplashHash := raw["discovery_splash"]
	guild := &Guild{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, guild)
	guild.token = c.token
	guild.fillRoles(raw["roles"].([]interface{}))
	guild.fillChannels(raw["channels"].([]interface{}))
	asset := &Asset{Format: "png", Size: 1024}
	if reflect.TypeOf(iconHash) != nil {
		asset.Extras = "icons/" + guild.Id
		asset.Hash = iconHash.(string)
		guild.Icon = *asset
	}
	if reflect.TypeOf(splashHash) != nil {
		asset.Extras = "splashes/" + guild.Id
		asset.Hash = splashHash.(string)
		guild.Splash = *asset
	}
	if reflect.TypeOf(bannerHash) != nil {
		asset.Extras = "banners/" + guild.Id
		asset.Hash = bannerHash.(string)
		guild.Banner = *asset
	}
	if reflect.TypeOf(dSplashHash) != nil {
		asset.Extras = "discovery_splashes/" + guild.Id
		asset.Hash = dSplashHash.(string)
		guild.DiscoverySplash = *asset
	}
	return guild
}

func (c Converter) Member() *Member {
	m := &Member{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, m)
	m.token = c.token
	u := Converter{payload: c.payload.(map[string]interface{})["user"], token: c.token}.User()
	s.Users[u.Id] = u
	m.fillUser(u)
	avatarHash := c.payload.(map[string]interface{})["avatar"]
	if reflect.TypeOf(avatarHash) != nil {
		m.Avatar = Asset{Hash: avatarHash.(string), Format: "png", Size: 1024, Extras: "avatars/" + m.Id}
	} else {
		m.Avatar = u.Avatar
	}
	return m
}

func (c Converter) User() *User {
	u := &User{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, u)
	u.token = c.token
	avatarHash := c.payload.(map[string]interface{})["avatar"]
	bannerHash := c.payload.(map[string]interface{})["banner"]
	if reflect.TypeOf(avatarHash) != nil {
		u.Avatar = Asset{Hash: avatarHash.(string), Format: "png", Size: 1024, Extras: "avatars/" + u.Id}
	} else {
		discmInt, _ := strconv.Atoi(u.Discriminator)
		hash := strconv.Itoa(discmInt % 5)
		u.Avatar = Asset{Hash: hash, Format: "png", Size: 1024, Extras: "embed/avatars"}
	}
	if reflect.TypeOf(bannerHash) != nil {
		u.Banner = Asset{Hash: bannerHash.(string), Format: "png", Size: 1024, Extras: "banners/" + u.Id}
	}
	return u
}

func (c Converter) Role() *Role {
	role := &Role{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, role)
	return role
}

func (c Converter) Channel() *Channel {
	ch := &Channel{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, ch)
	ch.token = c.token
	return ch
}

func (c Converter) Emoji() *Emoji {
	emoji := &Emoji{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, emoji)
	return emoji
}

func (c Converter) Embed() *Embed {
	embed := &Embed{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, embed)
	return embed
}

func (c Converter) Attachment() *Attachment {
	attachment := &Attachment{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, attachment)
	return attachment
}
