// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"sync"

	"errors"

	"github.com/corestoreio/csfw/utils/log"
)

var ErrPublisherClosed = errors.New("config Manager Publisher already closed")

// Subscriber allows you listen to write actions. The order of calling
// each subscriber is totally random.
type Subscriber interface {
	// Message when a configuration value will be written Message gets
	// called to allow you to listen to changes.
	Message(path string, sg ScopeGroup, s ScopeIDer)
}

// pubSub embedded pointer struct into the Manager
type pubSub struct {
	// subWriters subscribe writers are getting called when a write even
	// will happen.
	subWriters map[int]Subscriber
	subAutoInc int // subAutoInc increased whenever a Subscriber has been added
	mu         sync.RWMutex
	publishArg chan arg
	closed     bool
}

// Close closes the internal channel for the pubsub Goroutine. Prevents a leaking
// Goroutine.
func (ps *pubSub) Close() error {
	if ps.closed {
		return ErrPublisherClosed
	}
	ps.closed = true
	close(ps.publishArg)
	return nil
}

// Subscribe adds a Subscriber to be called when a write event happens.
// Returns a unique identifier for the Subscriber for later removal.
func (ps *pubSub) Subscribe(s Subscriber) (subscriptionID int, err error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.subAutoInc++
	ps.subWriters[ps.subAutoInc] = s
	subscriptionID = ps.subAutoInc
	return
}

// Unsubscribe removes a subscriber with a specific ID.
func (ps *pubSub) Unsubscribe(subscriptionID int) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, ok := ps.subWriters[subscriptionID]; ok {
		ps.subWriters[subscriptionID] = nil // avoid mem leaks
	}
	delete(ps.subWriters, subscriptionID)
	return nil
}

// sendMsg sends the arg into the channel
func (ps *pubSub) sendMsg(a arg) {
	if false == ps.closed {
		ps.publishArg <- a
	}
}

// publish runs in a Goroutine and listens on the channel publishArg. Every time
// a message is coming in, it calls all subscribers. We must run asynchronously
// because we don't know how long each subscriber needs.
func (ps *pubSub) publish() {

	for {
		select {
		case a, ok := <-ps.publishArg:
			if !ok {
				// channel closed
				return
			}

			if len(ps.subWriters) > 0 {
				ps.mu.RLock()
				for _, s := range ps.subWriters {
					sendMsgRecoverable(s, a)
				}
				ps.mu.RUnlock()
			}
		}
	}
}

func sendMsgRecoverable(sl Subscriber, a arg) {
	defer func() { // protect ... you'll never know
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Error("config.pubSub.publish.recover.err", "err", err)
			} else {
				log.Error("config.pubSub.publish.recover.wtf", "recover", r)
			}
		}
	}()
	sl.Message(a.pa, a.sg, a.si)
}

func newPubSub() *pubSub {
	return &pubSub{
		subWriters: make(map[int]Subscriber),
		publishArg: make(chan arg),
	}
}
