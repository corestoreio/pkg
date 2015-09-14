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

// Subscriber allows you to listen to write actions. The order of calling
// each subscriber is totally random. If a subscriber panics, it gets securely
// removed without crashing the whole system.
type Subscriber interface {
	// MessageConfig when a configuration value will be written MessageConfig
	// gets called to allow you to listen to changes.
	MessageConfig(path string, sg ScopeGroup, s ScopeIDer)
}

// pubSub embedded pointer struct into the Manager
type pubSub struct {
	// subMap, subscribed writers are getting called when a write event
	// will happen. String is the path (aka topic) and int the Subscriber ID for later
	// removal.
	subMap     map[string]map[int]Subscriber
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
// Path allows you to filter to which path or part of a path you would like to listen.
// A path can be e.g. "system/smtp/host" to receive messages by single host changes or
// "system/smtp" to receive message from all smtp changes or "system" to receive changes
// for all paths beginning with "system". A path is equal to a topic in a PubSub system.
// Path cannot be empty means you cannot listen to all changes.
// Returns a unique identifier for the Subscriber for later removal, or an error.
func (ps *pubSub) Subscribe(path string, s Subscriber) (subscriptionID int, err error) {
	if path == "" {
		return 0, ErrPathEmpty
	}
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.subAutoInc++
	subscriptionID = ps.subAutoInc

	if _, ok := ps.subMap[path]; !ok {
		ps.subMap[path] = make(map[int]Subscriber)
	}
	ps.subMap[path][subscriptionID] = s

	return
}

// Unsubscribe removes a subscriber with a specific ID.
func (ps *pubSub) Unsubscribe(subscriptionID int) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for path, subs := range ps.subMap {
		if _, ok := subs[subscriptionID]; ok {
			delete(ps.subMap[path], subscriptionID) // mem leaks?
			if len(ps.subMap[path]) == 0 {
				delete(ps.subMap, path)
			}
			return nil
		}
	}
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

			if len(ps.subMap) == 0 {
				break
			}

			ps.mu.RLock()
			var evict []int

			if subs, ok := ps.subMap[a.pathLevel1()]; ok { // e.g.: system
				evict = append(evict, sendMessages(subs, a)...)
			}
			if subs, ok := ps.subMap[a.pathLevel2()]; ok { // e.g.: system/smtp
				evict = append(evict, sendMessages(subs, a)...)
			}
			if subs, ok := ps.subMap[a.pathLevel3()]; ok { // e.g.: system/smtp/host
				evict = append(evict, sendMessages(subs, a)...)
			}
			ps.mu.RUnlock()

			// remove all Subscribers which failed
			if len(evict) > 0 {
				for _, e := range evict {
					if err := ps.Unsubscribe(e); err != nil {
						log.Error("config.pubSub.publish.evict.Unsubscribe.err", "err", err, "subscriptionID", e)
					}
				}
			}
		}
	}
}

func sendMessages(subs map[int]Subscriber, a arg) (evict []int) {
	for id, s := range subs {
		if err := sendMsgRecoverable(id, s, a); err != nil {
			evict = append(evict, id) // mark Subscribers for removal which failed ...
		}
	}
	return
}

func sendMsgRecoverable(id int, sl Subscriber, a arg) (err error) {
	defer func() { // protect ... you'll never know
		if r := recover(); r != nil {
			if recErr, ok := r.(error); ok {
				err = log.Error("config.pubSub.publish.recover.err", "err", recErr)
			} else {
				err = log.Error("config.pubSub.publish.recover.r", "recover", r)
			}
			// the overall trick here is, that defer will assign a new error to err
			// and therefore will overwrite the returned nil value!
		}
	}()
	sl.MessageConfig(a.pa, a.sg, a.si)
	return
}

func newPubSub() *pubSub {
	return &pubSub{
		subMap:     make(map[string]map[int]Subscriber),
		publishArg: make(chan arg),
	}
}
