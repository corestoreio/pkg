// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import "sync"

type eventType uint8

const (
	eventToSQLBefore eventType = iota
	maxEventTypes
)

type (
	// SelectReceiverFn receives the Select object pointer for modification
	SelectReceiverFn func(*Select)
	// SelectEvents event object containing a list of different event types
	// which gets dispatched in different situations.
	SelectEvents struct {
		receivers [maxEventTypes][]SelectReceiverFn
	}
	// InsertReceiverFn receives the Select object pointer for modification
	InsertReceiverFn func(*Insert)
	// InsertEvents event object containing a list of different event types
	// which gets dispatched in different situations.
	InsertEvents struct {
		receivers [maxEventTypes][]InsertReceiverFn
	}
	// UpdateReceiverFn receives the Select object pointer for modification
	UpdateReceiverFn func(*Update)
	// UpdateEvents event object containing a list of different event types
	// which gets dispatched in different situations.
	UpdateEvents struct {
		receivers [maxEventTypes][]UpdateReceiverFn
	}
	// DeleteReceiverFn receives the Select object pointer for modification
	DeleteReceiverFn func(*Delete)
	// DeleteEvents event object containing a list of different event types
	// which gets dispatched in different situations.
	DeleteEvents struct {
		receivers [maxEventTypes][]DeleteReceiverFn
	}
)

// EventContainer a type for embedding into other structs to define events for
// manipulating the SQL. Not an interface because interfaces are named with
// verbs ;-) Not yet thread safe.
type EventContainer struct {
	Select SelectEvents
	Insert InsertEvents
	Update UpdateEvents
	Delete DeleteEvents
}

// Merge merges other events into the current event container.
func (e *EventContainer) Merge(events ...EventContainer) *EventContainer {
	for _, et := range events {
		e.Select.Merge(et.Select)
		e.Insert.Merge(et.Insert)
		e.Update.Merge(et.Update)
		e.Delete.Merge(et.Delete)
	}
	return e
}

// Merge merges other Select events into the current event object.
func (e *SelectEvents) Merge(ses ...SelectEvents) *SelectEvents {
	for _, ev := range ses {
		for idx, recs := range ev.receivers {
			e.receivers[idx] = append(e.receivers[idx], recs...)
		}
	}
	return e
}

// AddBeforeToSQL dispatches the events every time the ToSQL function gets
// called.
func (e *SelectEvents) AddBeforeToSQL(fns ...SelectReceiverFn) *SelectEvents {
	e.receivers[eventToSQLBefore] = append(e.receivers[eventToSQLBefore], fns...)
	return e
}

// AddBeforeToSQLOnce dispatches the events only once before ToSQL gets called.
// Subsequent calls to ToSQL do not trigger the ReceiverFn closures.
func (e *SelectEvents) AddBeforeToSQLOnce(fns ...SelectReceiverFn) *SelectEvents {
	newFns := make([]SelectReceiverFn, len(fns))
	for i, fn := range fns {
		fn := fn // catch variables because of the closure
		i := i
		var onesie sync.Once
		newFns[i] = func(b *Select) { onesie.Do(func() { fn(b) }) }
	}
	e.receivers[eventToSQLBefore] = append(e.receivers[eventToSQLBefore], newFns...)
	return e
}

func (e SelectEvents) dispatch(et eventType, b *Select) {
	for _, e := range e.receivers[et] {
		e(b)
	}
}

// Merge merges other Insert events into the current event object.
func (e *InsertEvents) Merge(events ...InsertEvents) *InsertEvents {
	for _, ev := range events {
		for idx, recs := range ev.receivers {
			e.receivers[idx] = append(e.receivers[idx], recs...)
		}
	}
	return e
}

// AddBeforeToSQL dispatches the events every time the ToSQL function gets
// called.
func (e *InsertEvents) AddBeforeToSQL(fns ...InsertReceiverFn) *InsertEvents {
	e.receivers[eventToSQLBefore] = append(e.receivers[eventToSQLBefore], fns...)
	return e
}

// AddBeforeToSQLOnce dispatches the events only once before ToSQL gets called.
// Subsequent calls to ToSQL do not trigger the ReceiverFn closures.
func (e *InsertEvents) AddBeforeToSQLOnce(fns ...InsertReceiverFn) *InsertEvents {
	newFns := make([]InsertReceiverFn, len(fns))
	for i, fn := range fns {
		fn := fn // catch variables because of the closure
		i := i
		var onesie sync.Once
		newFns[i] = func(b *Insert) { onesie.Do(func() { fn(b) }) }
	}
	e.receivers[eventToSQLBefore] = append(e.receivers[eventToSQLBefore], newFns...)
	return e
}

func (e InsertEvents) dispatch(et eventType, b *Insert) {
	for _, e := range e.receivers[et] {
		e(b)
	}
}

// Merge merges other Update events into the current event object.
func (e *UpdateEvents) Merge(events ...UpdateEvents) *UpdateEvents {
	for _, ev := range events {
		for idx, recs := range ev.receivers {
			e.receivers[idx] = append(e.receivers[idx], recs...)
		}
	}
	return e
}

// AddBeforeToSQL dispatches the events every time the ToSQL function gets
// called.
func (e *UpdateEvents) AddBeforeToSQL(fns ...UpdateReceiverFn) *UpdateEvents {
	e.receivers[eventToSQLBefore] = append(e.receivers[eventToSQLBefore], fns...)
	return e
}

// AddBeforeToSQLOnce dispatches the events only once before ToSQL gets called.
// Subsequent calls to ToSQL do not trigger the ReceiverFn closures.
func (e *UpdateEvents) AddBeforeToSQLOnce(fns ...UpdateReceiverFn) *UpdateEvents {
	newFns := make([]UpdateReceiverFn, len(fns))
	for i, fn := range fns {
		fn := fn // catch variables because of the closure
		i := i
		var onesie sync.Once
		newFns[i] = func(b *Update) { onesie.Do(func() { fn(b) }) }
	}
	e.receivers[eventToSQLBefore] = append(e.receivers[eventToSQLBefore], newFns...)
	return e
}

func (e UpdateEvents) dispatch(et eventType, b *Update) {
	for _, e := range e.receivers[et] {
		e(b)
	}
}

// Merge merges other Delete events into the current event object.
func (e *DeleteEvents) Merge(events ...DeleteEvents) *DeleteEvents {
	for _, ev := range events {
		for idx, recs := range ev.receivers {
			e.receivers[idx] = append(e.receivers[idx], recs...)
		}
	}
	return e
}

// AddBeforeToSQL dispatches the events every time the ToSQL function gets
// called.
func (e *DeleteEvents) AddBeforeToSQL(fns ...DeleteReceiverFn) *DeleteEvents {
	e.receivers[eventToSQLBefore] = append(e.receivers[eventToSQLBefore], fns...)
	return e
}

// AddBeforeToSQLOnce dispatches the events only once before ToSQL gets called.
// Subsequent calls to ToSQL do not trigger the ReceiverFn closures.
func (e *DeleteEvents) AddBeforeToSQLOnce(fns ...DeleteReceiverFn) *DeleteEvents {
	newFns := make([]DeleteReceiverFn, len(fns))
	for i, fn := range fns {
		fn := fn // catch variables because of the closure
		i := i
		var onesie sync.Once
		newFns[i] = func(b *Delete) { onesie.Do(func() { fn(b) }) }
	}
	e.receivers[eventToSQLBefore] = append(e.receivers[eventToSQLBefore], newFns...)
	return e
}

func (e DeleteEvents) dispatch(et eventType, b *Delete) {
	for _, e := range e.receivers[et] {
		e(b)
	}
}
