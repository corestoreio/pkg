package dbr

type alias struct {
	Expression string
	Alias      string
}

func newAlias(as ...string) alias {
	a := alias{
		Expression: as[0],
	}
	if len(as) > 1 {
		a.Alias = as[1]
	}
	return a
}

func (t alias) String() string {
	return Quoter.Alias(t.Expression, t.Alias)
}

func (t alias) QuoteAs() string {
	return Quoter.QuoteAs(t.Expression, t.Alias)
}
