package event

type Handler func()

type handlers map[Event][]Handler

var eventHandlers handlers

func init() {
	eventHandlers = make(handlers)
}

func RegisterHandler(e Event, h Handler) {
	eventHandlers[e] = append(eventHandlers[e], h)
}
