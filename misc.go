package main

import (
	"crypto/rand"
	"fmt"
	"time"
)

func AssignId(id string) string {
	if id == "" {
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		return fmt.Sprintf("%x", b)
	}
	return id
}

func ScheduleDeletion(timeout float64, loc map[string]interface{}, ids map[string]bool) {
	time.Sleep(time.Duration(timeout) * time.Second)
	for id := range ids {
		delete(loc, id)
	}
}
