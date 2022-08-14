package discord

import (
	"crypto/rand"
	"fmt"
	"time"
)

func assignId() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func scheduleDeletion(timeout float64, loc map[string]interface{}, ids map[string]bool) {
	time.Sleep(time.Duration(timeout) * time.Second)
	for id := range ids {
		delete(loc, id)
	}
}
