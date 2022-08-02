package disgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const BASE = "https://discord.com/api/v10"

type Router struct {
	Token  string
	Path   string
	Data   map[string]interface{}
	Files  []File
	Method string
}

func (obj *Router) Fire() *http.Response {
	body, boundary := MultiPartWriter(obj.Data, obj.Files)
	r, _ := http.NewRequest(obj.Method, BASE+obj.Path, io.NopCloser(bytes.NewBuffer(body)))
	r.Header.Set(`Content-Type`, fmt.Sprintf(`multipart/form-data; boundary=%s`, boundary))
	r.Header.Set(`Authorization`, fmt.Sprintf(`Bot %s`, obj.Token))
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

func MultipartReq(method string, path string, data map[string]interface{}, token string, files ...File) *Router {
	return &Router{Token: token, Path: path, Data: data, Method: method, Files: files}
}

type MinimalRouter struct {
	Method string
	Token  string
	Path   string
	Data   map[string]interface{}
}

func (obj *MinimalRouter) Fire() *http.Response {
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

func MinimalReq(method string, path string, data map[string]interface{}, token string) *MinimalRouter {
	return &MinimalRouter{Method: method, Path: path, Data: data, Token: token}
}
