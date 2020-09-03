package main

import "sync"

type Broker struct {
	Subscribers map[string][]func(e string, p interface{})
}

var brokerM = sync.Mutex{} // Global

func (b *Broker) SubscribeToEvent(e string, s func(e string, p interface{})) {
	if b.Subscribers == nil {
		b.Subscribers = map[string][]func(e string, p interface{}){}
	}
	if b.Subscribers[e] == nil {
		b.Subscribers[e] = []func(e string, p interface{}){}
	}
	b.Subscribers[e] = append(b.Subscribers[e], s)
}

func (b *Broker) EmitEvent(e string, payload interface{}) {
	// TODO: make this async
	brokerM.Lock()
	defer brokerM.Unlock()
	for i := range b.Subscribers[e] {
		go b.Subscribers[e][i](
			e,
			payload,
		)
	}
}
