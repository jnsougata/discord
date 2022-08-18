package discord

func checkTrueEmbed(em Embed) bool {
	return em.Title != "" || em.Description != "" || len(em.Fields) > 0 || em.Author.IconURL != "" || em.Author.Name != "" || em.Footer.IconURL != "" || em.Footer.Text != "" || em.Image.URL != "" || em.Thumbnail.URL != ""
}

func checkTrueFile(f File) bool {
	return f.Name != "" && len(f.Content) > 0
}
