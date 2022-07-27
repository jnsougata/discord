package utils

import (
	"github.com/jnsougata/disgo/core/embed"
	"github.com/jnsougata/disgo/core/file"
)

func CheckTrueEmbed(em embed.Embed) bool {
	return em.Title != "" || em.Description != "" || len(em.Fields) > 0 || em.Author.IconUrl != "" || em.Author.Name != "" || em.Footer.IconUrl != "" || em.Footer.Text != "" || em.Image.Url != "" || em.Thumbnail.Url != ""
}

func CheckTrueFile(f file.File) bool {
	return f.Name != "" && len(f.Content) > 0
}
