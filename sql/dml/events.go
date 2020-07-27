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

type ctxKeyQueryOption uint8

// QueryOptions provides different options while executing code for SQL queries.
type QueryOptions struct {
	SkipEvents     bool // skips above defined EventFlag
	SkipTimestamps bool // skips generating timestamps (TODO)
	SkipRelations  bool // skips executing relation based SQL code
}

// WithContextQueryOptions adds options for executing queries, mostly in generated code.
func WithContextQueryOptions(ctx context.Context, qo QueryOptions) context.Context {
	return context.WithValue(ctx, ctxKeyQueryOption(0), qo)
}

// FromContextQueryOptions returns the options from the context.
func FromContextQueryOptions(ctx context.Context) QueryOptions {
	v := ctx.Value(ctxKeyQueryOption(0))
	if qo, ok := v.(QueryOptions); ok {
		return qo
	}
	return QueryOptions{}
}
