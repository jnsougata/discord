package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const root = "https://discord.com/api/v10"

type multipartRouter struct {
	Token  string
	Path   string
	Data   map[string]interface{}
	Files  []File
	Method string
}

func (obj *multipartRouter) fire() *http.Response {
	body, boundary := multipartWriter(obj.Data, obj.Files)
	r, _ := http.NewRequest(obj.Method, root+obj.Path, io.NopCloser(bytes.NewBuffer(body)))
	r.Header.Set(`Content-Type`, fmt.Sprintf(`multipart/form-data; boundary=%s`, boundary))
	if obj.Token != "" {
		r.Header.Set(`Authorization`, fmt.Sprintf(`Bot %s`, obj.Token))
	}
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println("Client ran into an error while sending the request:", err)
	}
	if resp.StatusCode > 304 {
		e := make(map[string]interface{})
		b, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(b, &e)
		errcode := e["code"].(float64)
		msg := e["message"].(string)
		fmt.Println(fmt.Sprintf(
			"HTTP error %v occured (%v: %s) Method: %s Path: %s",
			resp.StatusCode, errcode, msg, obj.Method, root+obj.Path))
	}
	return resp
}

func multipartReq(
	method string, path string, data map[string]interface{}, token string, files ...File) *multipartRouter {
	return &multipartRouter{Token: token, Path: path, Data: data, Method: method, Files: files}
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
		fmt.Println("Client ran into an error while sending the request:", err)
	}
	if resp.StatusCode > 304 {
		e := make(map[string]interface{})
		b, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(b, &e)
		errcode := e["code"].(float64)
		msg := e["message"].(string)
		fmt.Println(fmt.Sprintf(
			"HTTP error %v occured (%v: %s) Method: %s Path: %s",
			resp.StatusCode, errcode, msg, obj.Method, root+obj.Path))
	}
	return resp
}

func minimalReq(method string, path string, data map[string]interface{}, token string) *minimalRouter {
	return &minimalRouter{Method: method, Path: path, Data: data, Token: token}
}
