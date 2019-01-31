// {{.Entity}} represents a single row for DB table `{{.TableName}}`.
// Auto generated.{{with .Comment}}
{{. -}}{{end}}{{- if .HasEasyJsonMarshaler }}
//easyjson:json{{end}}
type {{.Entity}} struct {
{{range .Columns}}{{GoCamelMaybePrivate .Field}} {{GoTypeNull .}}
		{{- if ne .StructTag "" -}}`{{.StructTag}}`{{- end}} {{.GoComment}}
{{end}} {{range .ReferencedCollections }} {{.}}
{{end}} }

// AssignLastInsertID updates the increment ID field with the last inserted ID
// from an INSERT operation. Implements dml.InsertIDAssigner. Auto generated.
func (e *{{.Entity}}) AssignLastInsertID(id int64) {
	{{range .Columns}}{{if and .IsPK .IsAutoIncrement}} e.{{GoCamelMaybePrivate .Field}} = {{GoType .}}(id)
	{{end}}{{end -}}
}

{{- range .Columns}}{{ if IsFieldPrivate .Field }}
// Set{{GoCamel .Field}} sets the data for a private and security sensitive
// field.
func (e *{{$.Entity}}) Set{{GoCamel .Field}}(d {{GoTypeNull .}}) *{{$.Entity}} {
	e.{{GoCamelMaybePrivate .Field}} = d
	return e
}

// Get{{GoCamel .Field}} returns the data from a private and security sensitive
// field.
func (e *{{$.Entity}}) Get{{GoCamel .Field}}() {{GoTypeNull .}} {
	return e.{{GoCamelMaybePrivate .Field}}
}
{{end }}{{end }}
// MapColumns implements interface ColumnMapper only partially. Auto generated.
func (e *{{.Entity}}) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm{{range .Columns}}.{{GoFuncNull .}}(&e.{{GoCamelMaybePrivate .Field}}){{end}}.Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c { {{range .Columns}}
			case "{{.Field}}"{{range .Aliases}},"{{.}}"{{end}}:
				cm.{{GoFuncNull .}}(&e.{{GoCamelMaybePrivate .Field}}){{end}}
			default:
				return errors.NotFound.Newf("[{{.Package}}] {{.Entity}} Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

// Empty empties all the fields of the current object. Also known as Reset.
func (e *{{.Entity}}) Empty() *{{.Entity}} { *e = {{.Entity}}{}; return e }

// {{.Collection}} represents a collection type for DB table {{.TableName}}
// Not thread safe. Auto generated.{{with .Comment}}
{{. -}}{{end}}{{- if .HasEasyJsonMarshaler }}
//easyjson:json{{end}}
type {{.Collection}} struct {
	Data           		[]*{{.Entity}} `json:"data,omitempty"`
	BeforeMapColumns	func(uint64, *{{.Entity}}) error `json:"-"`
	AfterMapColumns 	func(uint64, *{{.Entity}}) error `json:"-"`
}

// New{{.Collection}} creates a new initialized collection. Auto generated.
func New{{.Collection}}() *{{.Collection}} {
	{{/*
		TODO(idea): use a global pool which can register for each type the
		before/after mapcolumn function so that the dev does not need to assign
		each time. think if it's worth such a pattern.
	*/ -}}
	return &{{.Collection}}{
		Data: make([]*{{.Entity}}, 0, 5),
	}
}

// AssignLastInsertID traverses through the slice and sets a decrementing new
// ID to each entity.
func (cc *{{.Collection}}) AssignLastInsertID(id int64) {
	id++
	var j int64 = 1
	for i := len(cc.Data) - 1; i >= 0; i-- {
		cc.Data[i].AssignLastInsertID(id - j)
		j++
	}
}

func (cc *{{.Collection}}) scanColumns(cm *dml.ColumnMap,e *{{.Entity}}, idx uint64) error {
	if cc.BeforeMapColumns != nil {
		if err := cc.BeforeMapColumns(idx, e); err != nil {
			return errors.WithStack(err)
		}
	}
	if err := e.MapColumns(cm); err != nil {
		return errors.WithStack(err)
	}
	if cc.AfterMapColumns != nil {
		if err := cc.AfterMapColumns(idx, e); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// MapColumns implements dml.ColumnMapper interface. Auto generated.
func (cc *{{.Collection}}) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
		for i, e := range cc.Data {
			if err := cc.scanColumns(cm, e, uint64(i)); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapScan:
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		e := new({{.Entity}})
		if err := cc.scanColumns(cm, e, cm.Count); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, e)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next() {
			switch c := cm.Column(); c {
			{{- range .Columns.UniqueColumns -}}
			case "{{.Field}}"{{range .Aliases}},"{{.}}"{{end}}:
				cm = cm.{{GoFuncNull .}}s(cc.{{GoCamel .Field}}s()...)
			{{- end}}
			{{/* {{- range .Columns.UniquifiedColumns}}		// no idea if that is needed
			case "{{.Field}}"{{range .Aliases}},"{{.}}"{{end}}:
				cm = cm.{{GoFuncNull .}}s(cc.{{GoCamel .Field}}s()...)
			{{- end}} */}}
			default:
				return errors.NotFound.Newf("[{{.Package}}] {{.Collection}} Column %q not found", c)
			}
		}
	default:
		return errors.NotSupported.Newf("[{{.Package}}] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}
{{range .Columns.UniqueColumns}}
// {{GoCamel .Field}}s returns a slice or appends to a slice all values.
// Auto generated.
func (cc *{{$.Collection}}) {{GoCamel .Field}}s(ret ...{{GoTypeNull .}}) []{{GoTypeNull .}} {
	if ret == nil {
		ret = make([]{{GoTypeNull .}}, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.{{GoCamel .Field}})
	}
	return ret
} {{end}}

{{- range .Columns.UniquifiedColumns}}
// {{GoCamel .Field}}s belongs to the column `{{.Field}}`
// and returns a slice or appends to a slice only unique values of that column.
// The values will be filtered internally in a Go map. No DB query gets
// executed. Auto generated.
func (cc *{{$.Collection}}) Unique{{GoCamel .Field}}s(ret ...{{GoType .}}) []{{GoType .}} {
	if ret == nil {
		ret = make([]{{GoType .}}, 0, len(cc.Data))
	}
	{{/*
		TODO: a reusable map and use different algorithms depending on the size
		of the cc.Data slice. Sometimes a for/for loop runs faster than a map.
	*/}}
	dupCheck := make(map[{{GoType .}}]bool, len(cc.Data))
	for _, e := range cc.Data {
		if !dupCheck[e.{{GoPrimitiveNull .}}] {
			ret = append(ret, e.{{GoPrimitiveNull .}})
			dupCheck[e.{{GoPrimitiveNull .}}] = true
		}
	}
	return ret
} {{end}}
