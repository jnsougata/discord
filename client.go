package discord

// Bot is a function that represents a connection to discord.
func Bot(intent Intent, cache bool, presence Presence) *connection {
	return &connection{
		sock: &ws{intent: int(intent), memoize: cache, presence: presence}}
}

type connection struct {
	sock      *ws
	Listeners Listeners
}

func (con *connection) Run(token string) {
	con.sock.listeners = con.Listeners
	con.sock.secret = token
	con.sock.run(token)
}

func (con *connection) Commands(commands ...ApplicationCommand) {
	con.sock.queue = append(con.sock.queue, commands...)
}
