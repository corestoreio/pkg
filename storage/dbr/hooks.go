package dbr

// These types are four callbacks to allow changes to the underlying SQL queries
// by a 3rd party package.
type (
	SelectHook func(*SelectBuilder)
	InsertHook func(*InsertBuilder)
	UpdateHook func(*UpdateBuilder)
	DeleteHook func(*DeleteBuilder)

	SelectHooks []SelectHook
	InsertHooks []InsertHook
	UpdateHooks []UpdateHook
	DeleteHooks []DeleteHook
)

func (sh SelectHooks) Apply(sb *SelectBuilder) {
	for _, h := range sh {
		h(sb)
	}
}

// Hook a type for embedding to define hooks for manipulating the SQL. DML
// stands for data manipulation language.
type Hook struct {
	SelectAfter SelectHooks
	InsertAfter InsertHooks
	UpdateAfter UpdateHooks
	DeleteAfter DeleteHooks
}

// NewHookDML creates a new set of hooks for data manipulation language
func NewHook() *Hook {
	return new(Hook)
}

func (h *Hook) Merge(hooks ...*Hook) *Hook {
	for _, hs := range hooks {
		h.AddSelectAfter(hs.SelectAfter...)
		h.AddInsertAfter(hs.InsertAfter...)
	}
	return h
}

func (h *Hook) AddSelectAfter(sh ...SelectHook) {
	h.SelectAfter = append(h.SelectAfter, sh...)
}

func (h *Hook) AddInsertAfter(sh ...InsertHook) {
	h.InsertAfter = append(h.InsertAfter, sh...)
}
