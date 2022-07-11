package models

type Modal struct {
	//TODO
}

func (i *Modal) ToBody() map[string]interface{} {
	return map[string]interface{}{
		"kind": 9,
		"data": map[string]interface{}{},
	}
}
