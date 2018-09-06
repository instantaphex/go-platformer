package engine

type Subscriber interface {
	PostEvent(Event)
}

type Dispatcher struct {
	subscribers map[string][]Subscriber
}

type Subscription struct {
	Unsubscribe func()
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		subscribers: make(map[string][]Subscriber),
	}
}

func (d *Dispatcher) Subscribe(eventType string, sub Subscriber) Subscription {
	var idx int

	d.subscribers[eventType] = append(d.subscribers[eventType], sub)
	idx = len(d.subscribers) - 1

	return Subscription {
		Unsubscribe: func() {
			d.subscribers[eventType] = append(
				d.subscribers[eventType][:idx],
				d.subscribers[eventType][idx+1:]...
			)
		},
	}
}

func (d *Dispatcher) EmitEvent(e Event) {
	for _, sub := range d.subscribers[e.Type()] {
		sub.PostEvent(e)
	}
}
