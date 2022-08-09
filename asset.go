package disgo

import (
	"fmt"
	"strings"
)

const base = "https://cdn.discordapp.com"

type Asset struct {
	Hash   string
	Size   int64
	Format string
	Extras string
}

func (a Asset) URL() string {
	if a.Format == "" {
		if strings.HasPrefix(a.Hash, "a_") {
			a.Format = "gif"
		} else {
			a.Format = "png"
		}
	}
	if a.Size == 0 {
		a.Size = 1024
	}
	if a.Hash != "" && a.Extras != "" {
		return fmt.Sprintf("%s/%s/%s.%s?size=%v", base, a.Extras, a.Hash, a.Format, a.Size)
	} else {
		return ""
	}
}

func (a Asset) CustomURL(size int64, format string) string {
	if format == "" {
		panic("Asset.Format is empty")
	}
	if size == 0 {
		panic("Asset.Size 0 is invalid")
	}
	if a.Hash != "" && a.Extras != "" {
		return fmt.Sprintf("%s/%s/%s.%s?size=%v", base, a.Extras, a.Hash, format, size)
	} else {
		return ""
	}
}
