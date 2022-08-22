package discord

type intent int

type intents struct {
	Members   intent
	Contents  intent
	Presences intent
}

var Intents = intents{
	Members:   intent(1 << 1),
	Contents:  intent(1 << 15),
	Presences: intent(1 << 8),
}

func (i *intents) Defaults(extras ...intent) intent {
	nums := []int{0, 2, 3, 4, 5, 6, 7, 9, 10, 11, 12, 13, 14, 16, 17, 18}
	ini := 0
	for _, i := range nums {
		ini |= 1 << i
	}

	for _, e := range extras {
		ini |= int(e)
	}

	return intent(ini)
}

func (i *intents) All() intent {
	ini := 0
	for i := 0; i < 19; i++ {
		ini |= 1 << i
	}
	return intent(ini)
}
