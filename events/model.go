package events

type ModelEvent struct {
	data interface{}
}

func NewModelEvent(data interface{}) *ModelEvent {
	return &ModelEvent{
		data: data,
	}
}

func (e *ModelEvent) Data() interface{} {
	return e.data
}
