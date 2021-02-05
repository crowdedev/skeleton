package events

type ModelEvent struct {
	data    interface{}
	id      string
	service string
}

func NewModelEvent(service string, data interface{}, id string) *ModelEvent {
	return &ModelEvent{
		service: service,
		data:    data,
		id:      id,
	}
}

func (e *ModelEvent) Data() interface{} {
	return e.data
}

func (e *ModelEvent) Id() string {
	return e.id
}

func (e *ModelEvent) Service() string {
	return e.service
}
