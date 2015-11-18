package dbr

// These types are four callbacks to allow changes to the underlying SQL queries
// by a 3rd party package.
type (
	SelectCb func(*SelectBuilder) *SelectBuilder
	InsertCb func(*InsertBuilder) *InsertBuilder
	UpdateCb func(*UpdateBuilder) *UpdateBuilder
	DeleteCb func(*DeleteBuilder) *DeleteBuilder
)
