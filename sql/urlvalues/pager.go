package urlvalues

import (
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/dml"
)

type Pager struct {
	Limit  uint64
	Offset uint64

	// Default max limit is 1000.
	MaxLimit uint64
	// Default max offset is 1000000.
	MaxOffset uint64

	stickyErr error
}

func NewPager(values Values) *Pager {
	p := new(Pager)
	p.stickyErr = p.FromValues(values)
	return p
}

func (p *Pager) FromValues(values Values) error {
	limit, err := values.Int64("limit")
	if err != nil {
		return errors.WithStack(err)
	}
	p.Limit = uint64(limit)

	page, err := values.Int64("page")
	if err != nil {
		return errors.WithStack(err)
	}
	p.Offset = uint64(page)

	return nil
}

func (p *Pager) maxLimit() uint64 {
	if p.MaxLimit > 0 {
		return p.MaxLimit
	}
	return 1000
}

func (p *Pager) maxOffset() uint64 {
	if p.MaxOffset > 0 {
		return p.MaxOffset
	}
	return 1000000
}

func (p *Pager) GetLimit() uint64 {
	const defaultLimit = 100

	if p == nil {
		return defaultLimit
	}
	if p.Limit < 0 {
		return p.Limit
	}
	if p.Limit == 0 {
		return defaultLimit
	}
	if p.Limit > p.maxLimit() {
		return p.maxLimit()
	}
	return p.Limit
}

func (p *Pager) GetOffset() uint64 {
	if p == nil {
		return 0
	}
	if p.Offset > p.maxOffset() {
		return p.maxOffset()
	}
	return p.Offset
}

func (p *Pager) Pagination(a *dml.Artisan) (*dml.Artisan, error) {
	if p == nil {
		return a, nil
	}
	if p.stickyErr != nil {
		return nil, p.stickyErr
	}
	o := p.GetOffset()
	if o < 1 {
		o = 1
	}
	l := p.GetLimit()
	if l < 1 {
		l = 1
	}
	a = a.Paginate(o, l)
	return a, nil
}
