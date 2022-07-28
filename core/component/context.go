package component

import (
	"encoding/json"
	"fmt"
	"github.com/jnsougata/disgo/core/member"
	"github.com/jnsougata/disgo/core/router"
	"github.com/jnsougata/disgo/core/user"
)

type Component struct {
	CustomId string   `json:"custom_id"`
	Type     int      `json:"type"`
	Value    string   `json:"value"`
	Values   []string `json:"values"`
}

type Row struct {
	Components []Component
}

type Data struct {
	ComponentType int      `json:"component_type"`
	CustomId      string   `json:"custom_id"`
	Values        []string `json:"values"`
	Components    []Row    `json:"components"`
}

type Context struct {
	ID             string                 `json:"id"`
	ApplicationId  string                 `json:"application_id"`
	Type           int                    `json:"type"`
	Data           Data                   `json:"data"`
	GuildId        string                 `json:"guild_id"`
	ChannelId      string                 `json:"channel_id"`
	Member         member.Member          `json:"member"`
	User           user.User              `json:"user"`
	Token          string                 `json:"token"`
	Version        int                    `json:"version"`
	Message        map[string]interface{} `json:"message"`
	AppPermissions string                 `json:"app_permissions"`
	Locale         string                 `json:"locale"`
	GuildLocale    string                 `json:"guild_locale"`
}

func FromData(payload interface{}) *Context {
	i := &Context{}
	data, _ := json.Marshal(payload)
	_ = json.Unmarshal(data, i)
	return i
}

func (c *Context) SendResponse(resp Message) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.ID, c.Token)
	r := router.New(
		"POST", path, map[string]interface{}{"type": 4, "data": resp.ToBody()}, "", resp.Files)
	go r.Request()
}

func (c *Context) Ack() {
	body := map[string]interface{}{"type": 6}
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.ID, c.Token)
	r := router.New("POST", path, body, "", nil)
	go r.Request()
}

func (c *Context) SendModal(modal Modal) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.ID, c.Token)
	r := router.New("POST", path, modal.ToBody(), "", nil)
	go r.Request()
}

func (c *Context) SendFollowup(resp Message) {
	path := fmt.Sprintf("/webhooks/%s/%s", c.ApplicationId, c.Token)
	r := router.New("POST", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}

func (c *Context) DeleteOriginalResponse() {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := router.Minimal("DELETE", path, nil, "")
	go r.Request()
}

func (c *Context) EditOriginalResponse(resp Message) {
	path := fmt.Sprintf("/webhooks/%s/%s/messages/@original", c.ApplicationId, c.Token)
	r := router.New("PATCH", path, resp.ToBody(), "", resp.Files)
	go r.Request()
}

func (c *Context) SendUpdate(resp Message) {
	path := fmt.Sprintf("/interactions/%s/%s/callback", c.ID, c.Token)
	body := map[string]interface{}{"type": 7, "data": resp.ToBody()}
	r := router.New("POST", path, body, "", resp.Files)
	go r.Request()
}
