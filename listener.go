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
	OnReady         func(bot BotUser)
	OnMessage       func(bot BotUser, message Message)
	OnGuildJoin     func(bot BotUser, guild Guild)
	OnGuildLeave    func(bot BotUser, guild Guild)
	OnInteraction   func(bot BotUser, interaction Interaction)
}
