package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jnsougata/disgo/core/utils"
	"io"
	"log"
	"net/http"
)

const BASE = "https://discord.com/api/v10"

type Router struct {
	Token  string
	Path   string
	Data   map[string]interface{}
	Files  []utils.File
	Method string
}

func (obj *Router) Request() *http.Response {
	body, boundary := utils.MultiPartWriter(obj.Data, obj.Files)
	r, _ := http.NewRequest(obj.Method, BASE+obj.Path, io.NopCloser(bytes.NewBuffer(body)))
	r.Header.Set(`Content-Type`, fmt.Sprintf(`multipart/form-data; boundary=%s`, boundary))
	r.Header.Set(`Authorization`, fmt.Sprintf(`Bot %s`, obj.Token))
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}
	return resp
}

func New(method string, path string, data map[string]interface{}, token string, files []utils.File) *Router {
	return &Router{Token: token, Path: path, Data: data, Method: method, Files: files}
}

type MinimalRouter struct {
	Method string
	Token  string
	Path   string
	Data   map[string]interface{}
}

func (obj *MinimalRouter) Request() *http.Response {
	body, _ := json.Marshal(obj.Data)
	r, _ := http.NewRequest(obj.Method, BASE+obj.Path, io.NopCloser(bytes.NewBuffer(body)))
	r.Header.Set(`Content-Type`, `application/json`)
	r.Header.Set(`Authorization`, fmt.Sprintf(`Bot %s`, obj.Token))
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}
	return resp
}
func Minimal(method string, path string, data map[string]interface{}, token string) *MinimalRouter {
	return &MinimalRouter{Method: method, Path: path, Data: data, Token: token}
}
