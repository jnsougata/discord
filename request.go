package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const root = "https://discord.com/api/v10"

type router struct {
	Token  string
	Path   string
	Data   map[string]interface{}
	Files  []File
	Method string
}

func (obj *router) fire() *http.Response {
	body, boundary := multipartWriter(obj.Data, obj.Files)
	r, _ := http.NewRequest(obj.Method, root+obj.Path, io.NopCloser(bytes.NewBuffer(body)))
	r.Header.Set(`Content-Type`, fmt.Sprintf(`multipart/form-data; boundary=%s`, boundary))
	if obj.Token != "" {
		r.Header.Set(`Authorization`, fmt.Sprintf(`Bot %s`, obj.Token))
	}
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode > 304 {
		em := make(map[string]interface{})
		b, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(b, &em)
		log.Println(em)
	}
	return resp
}

func multipartReq(method string, path string, data map[string]interface{}, token string, files ...File) *router {
	return &router{Token: token, Path: path, Data: data, Method: method, Files: files}
}

type minimalRouter struct {
	Method string
	Token  string
	Path   string
	Data   map[string]interface{}
}

func (obj *minimalRouter) fire() *http.Response {
	body, _ := json.Marshal(obj.Data)
	reader := io.NopCloser(bytes.NewBuffer(body))
	if obj.Method == "DELETE" || obj.Method == "GET" {
		reader = nil
	}
	r, _ := http.NewRequest(obj.Method, root+obj.Path, reader)
	r.Header.Set(`Content-Type`, `application/json`)
	if obj.Token != "" {
		r.Header.Set(`Authorization`, fmt.Sprintf(`Bot %s`, obj.Token))
	}
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
	}
	return resp
}

func minimalReq(method string, path string, data map[string]interface{}, token string) *minimalRouter {
	return &minimalRouter{Method: method, Path: path, Data: data, Token: token}
}
