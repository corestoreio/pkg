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
	Select ListenersSelect
	Insert ListenersInsert
	Update ListenersUpdate
	Delete ListenersDelete
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
	ListenSelectFn
	ListenInsertFn
	ListenUpdateFn
	ListenDeleteFn
}

// <-------------------------COPY------------------------->

// ListenSelectFn receives the Select object pointer for modification when an event
// gets dispatched.
type ListenSelectFn func(*Select)

// selectListen wrapper struct because we might wrap the SelectReceiverFn from
// the SelectListen struct.
type selectListen struct {
	name string
	EventType
	ListenSelectFn
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

	nsl.ListenSelectFn = sl.ListenSelectFn
	return nsl
}

// ListenersSelect contains multiple select event listener
type ListenersSelect []selectListen

// Add adds multiple listener to the listener stack and transforms the listener
// functions according to the configuration.
func (se *ListenersSelect) Add(sls ...Listen) ListenersSelect {
	for idx, sl := range sls {
		if sl.ListenSelectFn != nil {
			*se = append(*se, makeSelectListen(idx, sl))
		}
	}
	return *se
}

// Merge merges other ListenersSelect into the current listeners.
func (se *ListenersSelect) Merge(sls ...ListenersSelect) ListenersSelect {
	for _, sl := range sls {
		*se = append(*se, sl...)
	}
	return *se
}

func (se ListenersSelect) dispatch(et EventType, b *Select) error {
	for i, s := range se {
		switch {
		case s.error != nil:
			return errors.Wrapf(s.error, "[dml] ListenersSelect.dispatch Index %d EventType: %s", i, et)
		case s.EventType == et && !(b.PropagationStopped && i > b.propagationStoppedAt):
			s.ListenSelectFn(b)
			if b.propagationStoppedAt == 0 && b.PropagationStopped {
				b.propagationStoppedAt = i
			}
		case s.EventType == et:
			if b.Log.IsDebug() {
				b.Log.Debug("dml.ListenersSelect.Dispatch.PropagationStopped",
					log.String("listener_name", s.name), log.Err(s.error), log.Stringer("event_type", s.EventType),
					log.Bool("propagation_stopped", b.PropagationStopped), log.Int("propagation_stopped_at", b.propagationStoppedAt),
				)
			}
		}
	}
	return nil
}

// String returns a list of all named event listeners.
func (se ListenersSelect) String() string {
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

// ListenInsertFn receives the Insert object pointer for modification when an event
// gets dispatched.
type ListenInsertFn func(*Insert)

// insertListen wrapper struct because we might wrap the InsertReceiverFn from
// the InsertListen struct.
type insertListen struct {
	name string
	EventType
	ListenInsertFn
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

	nsl.ListenInsertFn = sl.ListenInsertFn
	return nsl
}

// ListenersInsert contains multiple insert event listener
type ListenersInsert []insertListen

// Add adds multiple listener to the listener stack and transforms the listener
// functions according to the configuration.
func (se *ListenersInsert) Add(sls ...Listen) ListenersInsert {
	for idx, sl := range sls {
		if sl.ListenInsertFn != nil {
			*se = append(*se, makeInsertListen(idx, sl))
		}
	}
	return *se
}

// Merge merges other ListenersInsert into the current listeners.
func (se *ListenersInsert) Merge(sls ...ListenersInsert) ListenersInsert {
	for _, sl := range sls {
		*se = append(*se, sl...)
	}
	return *se
}

func (se ListenersInsert) dispatch(et EventType, b *Insert) error {
	for i, s := range se {
		switch {
		case s.error != nil:
			return errors.Wrapf(s.error, "[dml] ListenersInsert.dispatch Index %d EventType: %s", i, et)
		case s.EventType == et && !(b.PropagationStopped && i > b.propagationStoppedAt):
			s.ListenInsertFn(b)
			if b.propagationStoppedAt == 0 && b.PropagationStopped {
				b.propagationStoppedAt = i
			}
		case s.EventType == et:
			if b.Log.IsDebug() {
				b.Log.Debug("dml.ListenersInsert.Dispatch.PropagationStopped",
					log.String("listener_name", s.name), log.Err(s.error), log.Stringer("event_type", s.EventType),
					log.Bool("propagation_stopped", b.PropagationStopped), log.Int("propagation_stopped_at", b.propagationStoppedAt),
				)
			}
		}
	}
	return nil
}

// String returns a list of all named event listeners.
func (se ListenersInsert) String() string {
	var buf bytes.Buffer
	for i, li := range se {
		_, _ = buf.WriteString(li.name)
		if i < len(se)-1 {
			_, _ = buf.WriteString("; ")
		}
	}
	return buf.String()
}

// ListenUpdateFn receives the Update object pointer for modification when an event
// gets dispatched.
type ListenUpdateFn func(*Update)

// updateListen wrapper struct because we might wrap the UpdateReceiverFn from
// the UpdateListen struct.
type updateListen struct {
	name string
	EventType
	ListenUpdateFn
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

	nsl.ListenUpdateFn = sl.ListenUpdateFn
	return nsl
}

// ListenersUpdate contains multiple update event listener
type ListenersUpdate []updateListen

// Add adds multiple listener to the listener stack and transforms the listener
// functions according to the configuration.
func (se *ListenersUpdate) Add(sls ...Listen) ListenersUpdate {
	for idx, sl := range sls {
		if sl.ListenUpdateFn != nil {
			*se = append(*se, makeUpdateListen(idx, sl))
		}
	}
	return *se
}

// Merge merges other ListenersUpdate into the current listeners.
func (se *ListenersUpdate) Merge(sls ...ListenersUpdate) ListenersUpdate {
	for _, sl := range sls {
		*se = append(*se, sl...)
	}
	return *se
}

func (se ListenersUpdate) dispatch(et EventType, b *Update) error {
	for i, s := range se {
		switch {
		case s.error != nil:
			return errors.Wrapf(s.error, "[dml] ListenersUpdate.dispatch Index %d EventType: %s", i, et)
		case s.EventType == et && !(b.PropagationStopped && i > b.propagationStoppedAt):
			s.ListenUpdateFn(b)
			if b.propagationStoppedAt == 0 && b.PropagationStopped {
				b.propagationStoppedAt = i
			}
		case s.EventType == et:
			if b.Log.IsDebug() {
				b.Log.Debug("dml.ListenersUpdate.Dispatch.PropagationStopped",
					log.String("listener_name", s.name), log.Err(s.error), log.Stringer("event_type", s.EventType),
					log.Bool("propagation_stopped", b.PropagationStopped), log.Int("propagation_stopped_at", b.propagationStoppedAt),
				)
			}
		}
	}
	return nil
}

// String returns a list of all named event listeners.
func (se ListenersUpdate) String() string {
	var buf bytes.Buffer
	for i, li := range se {
		_, _ = buf.WriteString(li.name)
		if i < len(se)-1 {
			_, _ = buf.WriteString("; ")
		}
	}
	return buf.String()
}

// ListenDeleteFn receives the Delete object pointer for modification when an event
// gets dispatched.
type ListenDeleteFn func(*Delete)

// deleteListen wrapper struct because we might wrap the DeleteReceiverFn from
// the DeleteListen struct.
type deleteListen struct {
	name string
	EventType
	ListenDeleteFn
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

	nsl.ListenDeleteFn = sl.ListenDeleteFn
	return nsl
}

// ListenersDelete contains multiple delete event listener
type ListenersDelete []deleteListen

// Add adds multiple listener to the listener stack and transforms the listener
// functions according to the configuration.
func (se *ListenersDelete) Add(sls ...Listen) ListenersDelete {
	for idx, sl := range sls {
		if sl.ListenDeleteFn != nil {
			*se = append(*se, makeDeleteListen(idx, sl))
		}
	}
	return *se
}

// Merge merges other ListenersDelete into the current listeners.
func (se *ListenersDelete) Merge(sls ...ListenersDelete) ListenersDelete {
	for _, sl := range sls {
		*se = append(*se, sl...)
	}
	return *se
}

func (se ListenersDelete) dispatch(et EventType, b *Delete) error {
	for i, s := range se {
		switch {
		case s.error != nil:
			return errors.Wrapf(s.error, "[dml] ListenersDelete.dispatch Index %d EventType: %s", i, et)
		case s.EventType == et && !(b.PropagationStopped && i > b.propagationStoppedAt):
			s.ListenDeleteFn(b)
			if b.propagationStoppedAt == 0 && b.PropagationStopped {
				b.propagationStoppedAt = i
			}
		case s.EventType == et:
			if b.Log.IsDebug() {
				b.Log.Debug("dml.ListenersDelete.Dispatch.PropagationStopped",
					log.String("listener_name", s.name), log.Err(s.error), log.Stringer("event_type", s.EventType),
					log.Bool("propagation_stopped", b.PropagationStopped), log.Int("propagation_stopped_at", b.propagationStoppedAt),
				)
			}
		}
	}
	return nil
}

// String returns a list of all named event listeners.
func (se ListenersDelete) String() string {
	var buf bytes.Buffer
	for i, li := range se {
		_, _ = buf.WriteString(li.name)
		if i < len(se)-1 {
			_, _ = buf.WriteString("; ")
		}
	}
	return buf.String()
}
