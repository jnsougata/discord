package disgo

func New(client Client) *connection {
	return &connection{
		Sock: &Socket{
			Intent:   client.Intent,
			Memoize:  client.Chunk,
			Presence: client.Presence,
		},
	}
}

type Client struct {
	Intent   int
	Chunk    bool
	Presence Presence
}

type connection struct {
	Sock *Socket
}

func (conn *connection) Run(token string) {
	conn.Sock.Run(token)
}

func (conn *connection) AddCommands(commands ...ApplicationCommand) {
	conn.Sock.RegistrationQueue(commands...)
}

func (conn *connection) OnSocketReceive(handler func(payload map[string]interface{})) {
	conn.Sock.AddHandler(OnSocketReceive, handler)
}

func (conn *connection) OnMessage(handler func(bot BotUser, message Message)) {
	conn.Sock.AddHandler(OnMessageCreate, handler)
}

func (conn *connection) OnReady(handler func(bot BotUser)) {
	conn.Sock.AddHandler(OnReady, handler)
}

func (conn *connection) OnInteraction(handler func(bot BotUser, ctx *Context)) {
	conn.Sock.AddHandler(OnInteractionCreate, handler)
}

func (conn *connection) OnGuildJoin(handler func(bot BotUser, guild Guild)) {
	conn.Sock.AddHandler(OnGuildCreate, handler)
}

func (conn *connection) OnGuildLeave(handler func(bot BotUser, guild Guild)) {
	conn.Sock.AddHandler(OnGuildDelete, handler)
}
