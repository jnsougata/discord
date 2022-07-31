package main

type Presence struct {
	Since    int64    `json:"since"`
	Status   string   `json:"status"` // "online" or "idle" or "dnd" or "offline" or "invisible"
	AFK      bool     `json:"afk"`
	Activity Activity `json:"activities"` // base activity object
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
	presence["activities"] = []Activity{p.Activity}
	return presence
}

type Activity struct {
	Name string `json:"name"` // "name" of the activity
	Type int    `json:"type"` // 0 (game), 1 (streaming), 2 (listening), 3 (watching), 5 (competing)
	URL  string `json:"url"`  // "url" of type 3 (streaming) activity only
}

func (a *Activity) ToData() map[string]interface{} {
	activity := map[string]interface{}{}
	if a.Name != "" {
		activity["name"] = a.Name
	}
	activity["type"] = a.Type
	if a.URL != "" && a.Type == 1 {
		activity["url"] = a.URL
	}
	return activity
}
