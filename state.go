package discord

type state struct {
	Token  string
	Users  map[string]*User
	Guilds map[string]*Guild
}
