package objects

type Modal struct {
	//TODO
}

func (i *Modal) ToBody() map[string]interface{} {
	return map[string]interface{}{
		"type": 9,
		"data": map[string]interface{}{},
	}
}
