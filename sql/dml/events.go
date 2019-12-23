package dml

// EventFlag describes where and when an event might get dispatched.
type EventFlag int

// EventFlag constants define the concrete locations of dispatched events.
const (
	EventFlagUndefined EventFlag = iota
	EventFlagBeforeSelect
	EventFlagAfterSelect
	EventFlagBeforeInsert
	EventFlagAfterInsert
	EventFlagBeforeUpdate
	EventFlagAfterUpdate
	EventFlagBeforeDelete
	EventFlagAfterDelete
)
