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

type eventType uint8

const (
	// EventToSQLBefore gets dispatched before generating the SQL string.
	EventToSQLBefore eventType = iota
	maxEventTypes
)

// Functions which acts as event receivers to allow changes to the underlying
// unprocessed SQL query.
type (
	SelectEvent func(*Select)
	InsertEvent func(*Insert)
	UpdateEvent func(*Update)
	DeleteEvent func(*Delete)
)

// Events a type for embedding to define events for manipulating the SQL.
type Events struct {
	selectEvents [maxEventTypes][]SelectEvent
	insertEvents [maxEventTypes][]InsertEvent
	updateEvents [maxEventTypes][]UpdateEvent
	deleteEvents [maxEventTypes][]DeleteEvent
}

// NewEvents creates a new set of hooks for data manipulation language
func NewEvents() *Events {
	return new(Events)
}

// Merge merges one or more other hooks into the current hook.
func (h *Events) Merge(events ...*Events) *Events {
	if h == nil {
		h = NewEvents()
	}
	for _, et := range events {
		if et != nil {
			for idx, selEvs := range et.selectEvents {
				if eventType(idx) < maxEventTypes {
					for _, evt := range selEvs {
						h.AddSelect(eventType(idx), evt)
					}
				}
			}
		}
		//h.AddInsert(hs.insertEvents...)
		//h.AddUpdate(hs.updateEvents...)
		//h.AddDelete(hs.deleteEvents...)
	}
	return h
}

func (h *Events) AddSelect(et eventType, sh ...SelectEvent) *Events {
	if h == nil {
		h = NewEvents()
	}
	h.selectEvents[et] = append(h.selectEvents[et], sh...)
	return h
}

func (h *Events) dispatchSelect(et eventType, b *Select) {
	if h == nil {
		return
	}
	for _, e := range h.selectEvents[et] {
		e(b)
	}
}

func (h *Events) AddInsert(et eventType, sh ...InsertEvent) *Events {
	if h == nil {
		h = NewEvents()
	}
	h.insertEvents[et] = append(h.insertEvents[et], sh...)
	return h
}

func (h *Events) dispatchInsert(et eventType, b *Insert) {
	if h == nil {
		return
	}
	for _, e := range h.insertEvents[et] {
		e(b)
	}
}

func (h *Events) AddUpdate(et eventType, sh ...UpdateEvent) *Events {
	if h == nil {
		h = NewEvents()
	}
	h.updateEvents[et] = append(h.updateEvents[et], sh...)
	return h
}

func (h *Events) dispatchUpdate(et eventType, b *Update) {
	if h == nil {
		return
	}
	for _, e := range h.updateEvents[et] {
		e(b)
	}
}

func (h *Events) AddDelete(et eventType, sh ...DeleteEvent) *Events {
	if h == nil {
		h = NewEvents()
	}
	h.deleteEvents[et] = append(h.deleteEvents[et], sh...)
	return h
}

func (h *Events) dispatchDelete(et eventType, b *Delete) {
	if h == nil {
		return
	}
	for _, e := range h.deleteEvents[et] {
		e(b)
	}
}
