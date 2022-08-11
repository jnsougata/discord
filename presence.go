package discord

type ActivityType int
type Status string

const (
	Playing   ActivityType = 0
	Streaming ActivityType = 1
	Listening ActivityType = 2
	Watching  ActivityType = 3
	Competing ActivityType = 5
)

const (
	Online       Status = "online"
	Idle         Status = "idle"
	DoNotDisturb Status = "dnd"
	Invisible    Status = "invisible"
	Offline      Status = "offline"
)

type Presence struct {
	Since    int64
	Status   Status
	AFK      bool
	Activity Activity // base activity object
	OnMobile bool
}

func (p *Presence) Marshal() map[string]interface{} {
	presence := map[string]interface{}{}
	if p.Since != 0 {
		presence["since"] = p.Since
	}
	if p.Status != "" {
		presence["status"] = p.Status
	}
	if p.AFK {
		presence["afk"] = true
	}
	presence["activities"] = []map[string]interface{}{p.Activity.Marshal()}
	return presence
}

type Activity struct {
	Name string       `json:"name"`
	Type ActivityType `json:"type"`
	URL  string       `json:"url"` // "url" for ActivityType Streaming
}

func (a *Activity) Marshal() map[string]interface{} {
	body := map[string]interface{}{}
	body["type"] = a.Type
	if a.Name != "" {
		body["name"] = a.Name
	}
	if a.URL != "" && a.Type == Streaming {
		body["url"] = a.URL
	}
	return body
}
