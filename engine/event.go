package engine

import (
	"bytes"
	"encoding/gob"
)


type (
	// eventBus : Pub-sub for communicating between game systems
	eventBus struct {
		subscribers map[Topic][]Subscriber
	}

	// Topic : The subject of the event being pub/sub'd
	Topic string

	// Event : Base class for events
	Event struct {
		Topic   Topic
		Payload []byte
	}

	// Subscriber : Receives events
	Subscriber interface {
		Call(e Event)
	}
)

// Subscribe : Add a Subscriber to a Topic
func (e *eventBus) Subscribe(topic Topic, sub Subscriber) error {
	if subs, ok := e.subscribers[topic]; ok {
		subs = append(subs, sub)
		return nil
	}
	e.subscribers[topic] = []Subscriber{sub}
	return nil
}

// Publish : publish an event to subscribers of Topic
func (e *eventBus) Publish(topic Topic, payload interface{}) error {
	var err error
	b := bytes.Buffer{}
	enc := gob.NewEncoder(&b)
	if err = enc.Encode(payload); err != nil {
		return err
	}
	// TODO: Is there a risk of collision by re-using this event? At this scale, does it matter?
	event := Event{
		Topic: topic,
		Payload: b.Bytes(),
	}
	for _, s := range e.subscribers[topic] {
		// TODO: How are publishing errors handled? Do we accept a handler on sub?
		go func(e Event) {
			s.Call(e)
		}(event)
	}
	return nil
}

// GetPayload : Deserialize event payload. Ideally, Event is extended to streamline typing
func (e *Event) GetPayload(output interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(e.Payload))
	return dec.Decode(e)
}
