package event

func Trigger(e Event) {
	for _, h := range eventHandlers[e] {
		h()
	}
}
