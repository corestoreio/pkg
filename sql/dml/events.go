// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License at distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dml

import (
	"bytes"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// EventType defines where and when an event gets dispatched.
type EventType byte

func (et EventType) String() string {
	return string(et)
}

// List of possible dispatched events.
const (
	OnBeforeToSQL EventType = iota + 65
)

// ListenerBucket a type for embedding into other structs to define events for
// manipulating the SQL. Not an interface because interfaces are named with
// verbs ;-) Not yet thread safe.
type ListenerBucket struct {
	Select SelectListeners
	Insert InsertListeners
	Update UpdateListeners
	Delete DeleteListeners
}

// NewListenerBucket creates a new event container to which multiple listeners
// can be added to.
func NewListenerBucket(listeners ...Listen) (*ListenerBucket, error) {
	ec := new(ListenerBucket)
	ec.Select.Add(listeners...)
	ec.Insert.Add(listeners...)
	ec.Update.Add(listeners...)
	ec.Delete.Add(listeners...)

	for i, ls := range ec.Select {
		if ls.error != nil {
			return nil, errors.Wrapf(ls.error, "[dml] NewListenerBucket Select Index %d", i)
		}
	}
	for i, ls := range ec.Insert {
		if ls.error != nil {
			return nil, errors.Wrapf(ls.error, "[dml] NewListenerBucket Insert Index %d", i)
		}
	}
	for i, ls := range ec.Update {
		if ls.error != nil {
			return nil, errors.Wrapf(ls.error, "[dml] NewListenerBucket Update Index %d", i)
		}
	}
	for i, ls := range ec.Delete {
		if ls.error != nil {
			return nil, errors.Wrapf(ls.error, "[dml] NewListenerBucket Delete Index %d", i)
		}
	}
	return ec, nil
}

// MustNewListenerBucket same at NewListenerBucket but panics on error.
func MustNewListenerBucket(listeners ...Listen) *ListenerBucket {
	ec, err := NewListenerBucket(listeners...)
	if err != nil {
		panic(err)
	}
	return ec
}

// Merge merges other events into the current event container.
func (lb *ListenerBucket) Merge(buckets ...*ListenerBucket) *ListenerBucket {
	for _, b := range buckets {
		lb.Select = append(lb.Select, b.Select...)
		lb.Insert = append(lb.Insert, b.Insert...)
		lb.Update = append(lb.Update, b.Update...)
		lb.Delete = append(lb.Delete, b.Delete...)
	}
	return lb
}

// Listen an argument to create a new listener when an event gets dispatched by
// a "Select, Insert, Update, Delete" type. Implements Listener interface.
type Listen struct {
	// Name optionally set internal name to identify multiple different listeners.
	Name string
	// EventType defines when a listener gets called. Mandatory.
	EventType

	// Listeners. Set at least one listener.
	SelectFunc
	InsertFunc
	UpdateFunc
	DeleteFunc
}

// <-------------------------COPY------------------------->

// SelectFunc receives the Select object pointer for modification when an event
// gets dispatched.
type SelectFunc func(*Select)

// selectListen wrapper struct because we might wrap the SelectReceiverFn from
// the SelectListen struct.
type selectListen struct {
	name string
	EventType
	SelectFunc
	error
}

func makeSelectListen(idx int, sl Listen) selectListen {
	nsl := selectListen{
		name:      sl.Name,
		EventType: sl.EventType,
	}
	if nsl.EventType == 0 {
		nsl.error = errors.Empty.Newf("[dml] Eventype at empty for %q; index %d", nsl.name, idx)
	}

	nsl.SelectFunc = sl.SelectFunc
	return nsl
}

// SelectListeners contains multiple select event listener
type SelectListeners []selectListen

// Add adds multiple listener to the listener stack and transforms the listener
// functions according to the configuration.
func (se *SelectListeners) Add(sls ...Listen) SelectListeners {
	for idx, sl := range sls {
		if sl.SelectFunc != nil {
			*se = append(*se, makeSelectListen(idx, sl))
		}
	}
	return *se
}

// Merge merges other SelectListeners into the current listeners.
func (se *SelectListeners) Merge(sls ...SelectListeners) SelectListeners {
	for _, sl := range sls {
		*se = append(*se, sl...)
	}
	return *se
}

func (se SelectListeners) dispatch(et EventType, b *Select) error {
	for i, s := range se {
		switch {
		case s.error != nil:
			return errors.Wrapf(s.error, "[dml] SelectListeners.dispatch Index %d EventType: %s", i, et)
		case s.EventType == et && !(b.PropagationStopped && i > b.propagationStoppedAt):
			s.SelectFunc(b)
			if b.propagationStoppedAt == 0 && b.PropagationStopped {
				b.propagationStoppedAt = i
			}
		case s.EventType == et:
			if b.Log.IsDebug() {
				b.Log.Debug("dml.SelectListeners.Dispatch.PropagationStopped",
					log.String("listener_name", s.name), log.Err(s.error), log.Stringer("event_type", s.EventType),
					log.Bool("propagation_stopped", b.PropagationStopped), log.Int("propagation_stopped_at", b.propagationStoppedAt),
				)
			}
		}
	}
	return nil
}

// String returns a list of all named event listeners.
func (se SelectListeners) String() string {
	var buf bytes.Buffer
	for i, li := range se {
		_, _ = buf.WriteString(li.name)
		if i < len(se)-1 {
			_, _ = buf.WriteString("; ")
		}
	}
	return buf.String()
}

// <-------------------------/COPY------------------------->

// InsertFunc receives the Insert object pointer for modification when an event
// gets dispatched.
type InsertFunc func(*Insert)

// insertListen wrapper struct because we might wrap the InsertReceiverFn from
// the InsertListen struct.
type insertListen struct {
	name string
	EventType
	InsertFunc
	error
}

func makeInsertListen(idx int, sl Listen) insertListen {
	nsl := insertListen{
		name:      sl.Name,
		EventType: sl.EventType,
	}
	if nsl.EventType == 0 {
		nsl.error = errors.Empty.Newf("[dml] Eventype at empty for %q; index %d", nsl.name, idx)
	}

	nsl.InsertFunc = sl.InsertFunc
	return nsl
}

// InsertListeners contains multiple insert event listener
type InsertListeners []insertListen

// Add adds multiple listener to the listener stack and transforms the listener
// functions according to the configuration.
func (se *InsertListeners) Add(sls ...Listen) InsertListeners {
	for idx, sl := range sls {
		if sl.InsertFunc != nil {
			*se = append(*se, makeInsertListen(idx, sl))
		}
	}
	return *se
}

// Merge merges other InsertListeners into the current listeners.
func (se *InsertListeners) Merge(sls ...InsertListeners) InsertListeners {
	for _, sl := range sls {
		*se = append(*se, sl...)
	}
	return *se
}

func (se InsertListeners) dispatch(et EventType, b *Insert) error {
	for i, s := range se {
		switch {
		case s.error != nil:
			return errors.Wrapf(s.error, "[dml] InsertListeners.dispatch Index %d EventType: %s", i, et)
		case s.EventType == et && !(b.PropagationStopped && i > b.propagationStoppedAt):
			s.InsertFunc(b)
			if b.propagationStoppedAt == 0 && b.PropagationStopped {
				b.propagationStoppedAt = i
			}
		case s.EventType == et:
			if b.Log.IsDebug() {
				b.Log.Debug("dml.InsertListeners.Dispatch.PropagationStopped",
					log.String("listener_name", s.name), log.Err(s.error), log.Stringer("event_type", s.EventType),
					log.Bool("propagation_stopped", b.PropagationStopped), log.Int("propagation_stopped_at", b.propagationStoppedAt),
				)
			}
		}
	}
	return nil
}

// String returns a list of all named event listeners.
func (se InsertListeners) String() string {
	var buf bytes.Buffer
	for i, li := range se {
		_, _ = buf.WriteString(li.name)
		if i < len(se)-1 {
			_, _ = buf.WriteString("; ")
		}
	}
	return buf.String()
}

// UpdateFunc receives the Update object pointer for modification when an event
// gets dispatched.
type UpdateFunc func(*Update)

// updateListen wrapper struct because we might wrap the UpdateReceiverFn from
// the UpdateListen struct.
type updateListen struct {
	name string
	EventType
	UpdateFunc
	error
}

func makeUpdateListen(idx int, sl Listen) updateListen {
	nsl := updateListen{
		name:      sl.Name,
		EventType: sl.EventType,
	}
	if nsl.EventType == 0 {
		nsl.error = errors.Empty.Newf("[dml] Eventype at empty for %q; index %d", nsl.name, idx)
	}

	nsl.UpdateFunc = sl.UpdateFunc
	return nsl
}

// UpdateListeners contains multiple update event listener
type UpdateListeners []updateListen

// Add adds multiple listener to the listener stack and transforms the listener
// functions according to the configuration.
func (se *UpdateListeners) Add(sls ...Listen) UpdateListeners {
	for idx, sl := range sls {
		if sl.UpdateFunc != nil {
			*se = append(*se, makeUpdateListen(idx, sl))
		}
	}
	return *se
}

// Merge merges other UpdateListeners into the current listeners.
func (se *UpdateListeners) Merge(sls ...UpdateListeners) UpdateListeners {
	for _, sl := range sls {
		*se = append(*se, sl...)
	}
	return *se
}

func (se UpdateListeners) dispatch(et EventType, b *Update) error {
	for i, s := range se {
		switch {
		case s.error != nil:
			return errors.Wrapf(s.error, "[dml] UpdateListeners.dispatch Index %d EventType: %s", i, et)
		case s.EventType == et && !(b.PropagationStopped && i > b.propagationStoppedAt):
			s.UpdateFunc(b)
			if b.propagationStoppedAt == 0 && b.PropagationStopped {
				b.propagationStoppedAt = i
			}
		case s.EventType == et:
			if b.Log.IsDebug() {
				b.Log.Debug("dml.UpdateListeners.Dispatch.PropagationStopped",
					log.String("listener_name", s.name), log.Err(s.error), log.Stringer("event_type", s.EventType),
					log.Bool("propagation_stopped", b.PropagationStopped), log.Int("propagation_stopped_at", b.propagationStoppedAt),
				)
			}
		}
	}
	return nil
}

// String returns a list of all named event listeners.
func (se UpdateListeners) String() string {
	var buf bytes.Buffer
	for i, li := range se {
		_, _ = buf.WriteString(li.name)
		if i < len(se)-1 {
			_, _ = buf.WriteString("; ")
		}
	}
	return buf.String()
}

// DeleteFunc receives the Delete object pointer for modification when an event
// gets dispatched.
type DeleteFunc func(*Delete)

// deleteListen wrapper struct because we might wrap the DeleteReceiverFn from
// the DeleteListen struct.
type deleteListen struct {
	name string
	EventType
	DeleteFunc
	error
}

func makeDeleteListen(idx int, sl Listen) deleteListen {
	nsl := deleteListen{
		name:      sl.Name,
		EventType: sl.EventType,
	}
	if nsl.EventType == 0 {
		nsl.error = errors.Empty.Newf("[dml] Eventype at empty for %q; index %d", nsl.name, idx)
	}

	nsl.DeleteFunc = sl.DeleteFunc
	return nsl
}

// DeleteListeners contains multiple delete event listener
type DeleteListeners []deleteListen

// Add adds multiple listener to the listener stack and transforms the listener
// functions according to the configuration.
func (se *DeleteListeners) Add(sls ...Listen) DeleteListeners {
	for idx, sl := range sls {
		if sl.DeleteFunc != nil {
			*se = append(*se, makeDeleteListen(idx, sl))
		}
	}
	return *se
}

// Merge merges other DeleteListeners into the current listeners.
func (se *DeleteListeners) Merge(sls ...DeleteListeners) DeleteListeners {
	for _, sl := range sls {
		*se = append(*se, sl...)
	}
	return *se
}

func (se DeleteListeners) dispatch(et EventType, b *Delete) error {
	for i, s := range se {
		switch {
		case s.error != nil:
			return errors.Wrapf(s.error, "[dml] DeleteListeners.dispatch Index %d EventType: %s", i, et)
		case s.EventType == et && !(b.PropagationStopped && i > b.propagationStoppedAt):
			s.DeleteFunc(b)
			if b.propagationStoppedAt == 0 && b.PropagationStopped {
				b.propagationStoppedAt = i
			}
		case s.EventType == et:
			if b.Log.IsDebug() {
				b.Log.Debug("dml.DeleteListeners.Dispatch.PropagationStopped",
					log.String("listener_name", s.name), log.Err(s.error), log.Stringer("event_type", s.EventType),
					log.Bool("propagation_stopped", b.PropagationStopped), log.Int("propagation_stopped_at", b.propagationStoppedAt),
				)
			}
		}
	}
	return nil
}

// String returns a list of all named event listeners.
func (se DeleteListeners) String() string {
	var buf bytes.Buffer
	for i, li := range se {
		_, _ = buf.WriteString(li.name)
		if i < len(se)-1 {
			_, _ = buf.WriteString("; ")
		}
	}
	return buf.String()
}
