package router

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const BASE = "https://discord.com/api/v10"

type Router struct {
	Token  string
	Path   string
	Body   map[string]interface{}
	Method string
}

func (obj *Router) Request() *http.Response {
	bodyByte, _ := json.Marshal(obj.Body)
	r, _ := http.NewRequest(obj.Method, BASE+obj.Path, io.NopCloser(bytes.NewBuffer(bodyByte)))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bot "+obj.Token)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}
	return resp
}

func New(method string, path string, body map[string]interface{}, token string) *Router {
	return &Router{Token: token, Path: path, Body: body, Method: method}
}
