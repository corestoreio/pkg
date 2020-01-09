package dml

import "context"

// EventFlag describes where and when an event might get dispatched.
type EventFlag uint8

// EventFlag constants define the concrete locations of dispatched events.
const (
	EventFlagUndefined EventFlag = iota
	EventFlagBeforeSelect
	EventFlagAfterSelect
	EventFlagBeforeInsert
	EventFlagAfterInsert
	EventFlagBeforeUpdate
	EventFlagAfterUpdate
	EventFlagBeforeUpsert
	EventFlagAfterUpsert
	EventFlagBeforeDelete
	EventFlagAfterDelete
	EventFlagMax // indicates maximum events available. Might change without notice.
)

type (
	ctxSkipEvents     struct{}
	ctxSkipTimestamps struct{}
)

// SkipEvents modifies a context to prevent events from running for any query it
// encounters.
func SkipEvents(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxSkipEvents{}, true)
}

// EventsAreSkipped returns true if the context skips events.
func EventsAreSkipped(ctx context.Context) bool {
	skip := ctx.Value(ctxSkipEvents{})
	return skip != nil && skip.(bool)
}

// SkipTimestamps modifies a context to prevent events from running for any
// query it encounters.
func SkipTimestamps(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxSkipTimestamps{}, true)
}

// TimestampsAreSkipped returns true if the context skips events.
func TimestampsAreSkipped(ctx context.Context) bool {
	skip := ctx.Value(ctxSkipTimestamps{})
	return skip != nil && skip.(bool)
}
