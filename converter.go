package disgo

import "encoding/json"

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
	return msg
}

func (c Converter) Guild() *Guild {
	guild := &Guild{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, guild)
	guild.token = c.token
	return guild
}

func (c Converter) Member() *Member {
	m := &Member{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, m)
	return m
}

func (c Converter) User() *User {
	u := &User{}
	data, _ := json.Marshal(c.payload)
	_ = json.Unmarshal(data, u)
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
