package discord

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type converter struct {
	state   *state
	payload interface{}
}

func (c converter) Bot() *Bot {
	bot := &Bot{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, bot)
	return bot
}

func (c converter) Message() *Message {
	msg := &Message{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, msg)
	msg.state = c.state
	return msg
}

func (c converter) Guild() *Guild {
	raw := c.payload.(map[string]interface{})
	iconHash := raw["icon"]
	bannerHash := raw["banner"]
	splashHash := raw["splash"]
	dSplashHash := raw["discovery_splash"]
	guild := &Guild{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, guild)
	guild.state = c.state
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

func (c converter) Member() *Member {
	m := &Member{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, m)
	m.state = c.state
	u := converter{payload: c.payload.(map[string]interface{})["user"], state: c.state}.User()
	shared.Users[u.Id] = u
	m.fillUser(u)
	avatarHash := c.payload.(map[string]interface{})["avatar"]
	if reflect.TypeOf(avatarHash) != nil {
		m.Avatar = Asset{Hash: avatarHash.(string), Format: "png", Size: 1024, Extras: "avatars/" + m.Id}
	} else {
		m.Avatar = u.Avatar
	}
	return m
}

func (c converter) User() *User {
	u := &User{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, u)
	u.state = c.state
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

func (c converter) Role() *Role {
	role := &Role{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, role)
	return role
}

func (c converter) Channel() *Channel {
	ch := &Channel{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, ch)
	ch.state = c.state
	return ch
}

func (c converter) Emoji() *Emoji {
	emoji := &Emoji{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, emoji)
	return emoji
}

func (c converter) Embed() *Embed {
	embed := &Embed{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, embed)
	return embed
}

func (c converter) Attachment() *Attachment {
	attachment := &Attachment{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, attachment)
	return attachment
}

func (c converter) Interaction() *Interaction {
	i := &Interaction{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, i)
	i.state = c.state
	guild, okg := shared.Guilds[i.GuildId]
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
	conv := converter{state: i.state}
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

func (c converter) Context() *Context {
	ctx := &Context{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, ctx)
	ctx.state = c.state
	ctx.data = c.payload.(map[string]interface{})
	guild, okg := shared.Guilds[ctx.GuildId]
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
	conv := converter{state: ctx.state}
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

func buildResolved(options []Option, resolved map[string]interface{}, state *state) *ResolvedOptions {
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
		conv := converter{state: state}
		if option.Type == OptionTypes.String {
			ro.strings[option.Name] = option.Value.(string)
		}
		if option.Type == OptionTypes.Integer {
			ro.integers[option.Name] = option.Value.(int64)
		}
		if option.Type == OptionTypes.Boolean {
			ro.booleans[option.Name] = option.Value.(bool)
		}
		if option.Type == OptionTypes.Number {
			ro.numbers[option.Name] = option.Value.(float64)
		}
		if option.Type == OptionTypes.Channel {
			channelId := option.Value.(string)
			conv.payload = resolved["channels"].(map[string]interface{})[channelId]
			ro.channels[option.Name] = *conv.Channel()
		}
		if option.Type == OptionTypes.Role {
			roleId := option.Value.(string)
			conv.payload = resolved["roles"].(map[string]interface{})[roleId]
			ro.roles[option.Name] = *conv.Role()
		}
		if option.Type == OptionTypes.Mentionable {
			ro.mentionables[option.Name] = option.Value
		}
		if option.Type == OptionTypes.Attachment {
			attachmentId := option.Value.(string)
			conv.payload = resolved["attachments"].(map[string]interface{})[attachmentId]
			ro.attachments[option.Name] = *conv.Attachment()
		}
		if option.Type == OptionTypes.User {
			userId := option.Value.(string)
			conv.payload = resolved["users"].(map[string]interface{})[userId]
			ro.users[option.Name] = *conv.User()
		}
	}
	return &ro
}
