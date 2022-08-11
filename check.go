package discord

func checkTrueEmbed(em Embed) bool {
	return em.Title != "" || em.Description != "" || len(em.Fields) > 0 || em.Author.IconUrl != "" || em.Author.Name != "" || em.Footer.IconUrl != "" || em.Footer.Text != "" || em.Image.Url != "" || em.Thumbnail.Url != ""
}

func checkTrueFile(f File) bool {
	return f.Name != "" && len(f.Content) > 0
}
