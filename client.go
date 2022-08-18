package discord

// New is a function that represents a connection to discord.
func New(intent intent) *connection {
	return &connection{ws: &ws{intent: int(intent)}}
}

type connection struct {
	ws        *ws
	Cache     bool
	Presence  Presence
	Listeners Listeners
}

func (conn *connection) Run(token string) {
	if conn.Cache {
		conn.ws.memoize = true
	}
	if conn.Presence.Activity.Name != "" {
		conn.ws.presence = conn.Presence
	}
	conn.ws.locked = true
	conn.ws.secret = token
	conn.ws.listeners = conn.Listeners
	conn.ws.commands = make(map[string]interface{})
	conn.ws.run(token)
}

func (conn *connection) Commands(commands ...Command) {
	conn.ws.queue = append(conn.ws.queue, commands...)
}
