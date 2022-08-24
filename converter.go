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

func (c Converter) Bot() *Bot {
	bot := &Bot{}
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

func (c Converter) Interaction() *Interaction {
	i := &Interaction{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, i)
	i.token = c.token
	guild, okg := s.Guilds[i.GuildId]
	if okg {
		i.Guild = *guild
		channel, okc := guild.Channels[i.ChannelId]
		if okc {
			i.Channel = *channel
		} else {
			i.Channel = Channel{}
		}
	} else {
		i.Guild = Guild{}
		i.Channel = Channel{}
	}
	conv := Converter{token: i.token}
	userData := c.payload.(map[string]interface{})["user"]
	memberData := c.payload.(map[string]interface{})["member"]
	if reflect.TypeOf(userData) != nil {
		conv.payload = userData.(map[string]interface{})
		i.User = *conv.User()
	}
	if reflect.TypeOf(memberData) != nil {
		id := memberData.(map[string]interface{})["user"].(map[string]interface{})["id"].(string)
		member, ok := i.Guild.Members[id]
		if ok {
			i.Author = *member
		} else {
			conv.payload = memberData.(map[string]interface{})
			i.Author = *conv.Member()
		}
	}
	return i
}

func (c Converter) Context() *Context {
	ctx := &Context{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, ctx)
	ctx.token = c.token
	ctx.raw = c.payload.(map[string]interface{})
	guild, okg := s.Guilds[ctx.GuildId]
	if okg {
		ctx.Guild = *guild
		channel, okc := guild.Channels[ctx.ChannelId]
		if okc {
			ctx.Channel = *channel
		} else {
			ctx.Channel = Channel{}
		}
	} else {
		ctx.Guild = Guild{}
		ctx.Channel = Channel{}
	}
	conv := Converter{token: ctx.token}
	userData := c.payload.(map[string]interface{})["user"]
	memberData := c.payload.(map[string]interface{})["member"]
	if reflect.TypeOf(userData) != nil {
		conv.payload = userData.(map[string]interface{})
		ctx.User = *conv.User()
	}
	if reflect.TypeOf(memberData) != nil {
		id := memberData.(map[string]interface{})["user"].(map[string]interface{})["id"].(string)
		member, ok := guild.Members[id]
		if ok {
			ctx.Author = *member
		} else {
			conv.payload = memberData.(map[string]interface{})
			ctx.Author = *conv.Member()
		}
	}
	return ctx
}

func buildRO(options []option, resolved map[string]interface{}, secret string) *ResolvedOptions {
	ro := ResolvedOptions{}
	ro.strings = map[string]string{}
	ro.integers = map[string]int64{}
	ro.booleans = map[string]bool{}
	ro.numbers = map[string]float64{}
	ro.channels = map[string]Channel{}
	ro.roles = map[string]Role{}
	ro.mentionables = map[string]interface{}{}
	ro.attachments = map[string]Attachment{}
	ro.users = map[string]User{}
	for _, option := range options {
		newConv := Converter{token: secret}
		if option.Type == stringOption {
			ro.strings[option.Name] = option.Value.(string)
		}
		if option.Type == integerOption {
			ro.integers[option.Name] = option.Value.(int64)
		}
		if option.Type == booleanOption {
			ro.booleans[option.Name] = option.Value.(bool)
		}
		if option.Type == numberOption {
			ro.numbers[option.Name] = option.Value.(float64)
		}
		if option.Type == channelOption {
			channelId := option.Value.(string)
			newConv.payload = resolved["channels"].(map[string]interface{})[channelId]
			ro.channels[option.Name] = *newConv.Channel()
		}
		if option.Type == roleOption {
			roleId := option.Value.(string)
			newConv.payload = resolved["roles"].(map[string]interface{})[roleId]
			ro.roles[option.Name] = *newConv.Role()
		}
		if option.Type == mentionableOption {
			ro.mentionables[option.Name] = option.Value
		}
		if option.Type == attachmentOption {
			attachmentId := option.Value.(string)
			newConv.payload = resolved["attachments"].(map[string]interface{})[attachmentId]
			ro.attachments[option.Name] = *newConv.Attachment()
		}
		if option.Type == userOption {
			userId := option.Value.(string)
			newConv.payload = resolved["users"].(map[string]interface{})[userId]
			ro.users[option.Name] = *newConv.User()
		}
	}
	return &ro
}
