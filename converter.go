package disgo

import "encoding/json"

func unmarshalBot(payload interface{}) *BotUser {
	bot := &BotUser{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, bot)
	return bot
}

func unmarshalMessage(payload interface{}) *Message {
	msg := &Message{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, msg)
	return msg
}

func unmarshalGuild(payload interface{}) *Guild {
	guild := &Guild{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, guild)
	return guild
}

func unmarshalMember(payload interface{}) *Member {
	m := &Member{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, m)
	return m
}

func unmarshalUser(payload interface{}) *User {
	u := &User{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, u)
	return u
}

func unmarshalRole(payload interface{}) *Role {
	role := &Role{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, role)
	return role
}

func unmarshalChannel(payload interface{}) *Channel {
	ch := &Channel{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, ch)
	return ch
}

func unmarshalEmoji(payload interface{}) *Emoji {
	emoji := &Emoji{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, emoji)
	return emoji
}

func unmarshalInteraction(payload interface{}) *Interaction {
	i := &Interaction{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, i)
	return i
}

func unmarshalEmbed(payload interface{}) *Embed {
	embed := &Embed{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, embed)
	return embed
}

func unmarshalAttachment(payload interface{}) *Attachment {
	attachment := &Attachment{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, attachment)
	return attachment
}
