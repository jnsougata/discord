package utils

import "fmt"

func MakeHash(data []byte) string {
	var hash uint64
	for _, b := range data {
		hash = uint64(b) + (hash << 6) + (hash << 16) - hash
	}
	return fmt.Sprintf("%X", hash)
}
