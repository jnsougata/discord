package discord

type ListenerType string

const (
	OnReady             ListenerType = "READY"
	OnMessage           ListenerType = "MESSAGE_CREATE"
	OnGuildJoin         ListenerType = "GUILD_CREATE"
	OnGuildLeave        ListenerType = "GUILD_DELETE"
	OnInteraction       ListenerType = "INTERACTION_CREATE"
	OnGuildMembersChunk ListenerType = "GUILD_MEMBERS_CHUNK"
)

type Listeners struct {
	OnSocketReceive func(payload interface{})
	OnReady         func(bot Bot)
	OnMessage       func(bot Bot, message Message)
	OnGuildJoin     func(bot Bot, guild Guild)
	OnGuildLeave    func(bot Bot, guild Guild)
	OnInteraction   func(bot Bot, interaction Interaction)
}
