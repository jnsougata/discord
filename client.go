package discord

// Bot is a function that represents a connection to discord.
func Bot(intent Intent, cache bool, presence Presence) *connection {
	return &connection{
		sock: &ws{intent: int(intent), memoize: cache, presence: presence}}
}

type connection struct {
	sock *ws
}

func (con *connection) Run(token string) {
	con.sock.Run(token)
}

func (con *connection) AddCommands(commands ...ApplicationCommand) {
	con.sock.AddToQueue(commands...)
}

func (con *connection) OnSocketReceive(handler func(payload map[string]interface{})) {
	con.sock.AddHandler(onSocketReceive, handler)
}

func (con *connection) OnMessage(handler func(bot BotUser, message Message)) {
	con.sock.AddHandler(onMessageCreate, handler)
}

func (con *connection) OnReady(handler func(bot BotUser)) {
	con.sock.AddHandler(onReady, handler)
}

func (con *connection) OnInteraction(handler func(bot BotUser, ctx *Context)) {
	con.sock.AddHandler(onInteractionCreate, handler)
}

func (con *connection) OnGuildJoin(handler func(bot BotUser, guild Guild)) {
	con.sock.AddHandler(onGuildCreate, handler)
}

func (con *connection) OnGuildLeave(handler func(bot BotUser, guild Guild)) {
	con.sock.AddHandler(onGuildDelete, handler)
}
