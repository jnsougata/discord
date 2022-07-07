package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const BASE = "https://discord.com/api/v10"

type Router struct {
	Token  string
	Path   string
	Body   map[string]interface{}
	Method string
}

func (obj *Router) Request() (*http.Response, error) {
	bodyByte, _ := json.Marshal(obj.Body)
	r, _ := http.NewRequest(obj.Method, BASE+obj.Path, io.NopCloser(bytes.NewBuffer(bodyByte)))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bot "+obj.Token)
	client := &http.Client{}
	return client.Do(r)
}
