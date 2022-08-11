package discord

// Bot is a function that represents a connection to discord.
func Bot(intent Intent, cache bool, presence Presence) *connection {
	return &connection{
		ws: &ws{intent: int(intent), memoize: cache, presence: presence}}
}

type connection struct {
	ws        *ws
	Listeners Listeners
}

func (conn *connection) Run(token string) {
	conn.ws.locked = true
	conn.ws.secret = token
	conn.ws.listeners = conn.Listeners
	conn.ws.commands = make(map[string]interface{})
	conn.ws.run(token)
}

func (conn *connection) Commands(commands ...ApplicationCommand) {
	conn.ws.queue = append(conn.ws.queue, commands...)
}
