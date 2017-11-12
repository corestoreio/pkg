// {{.Collection}} represents a collection type for DB table {{ .TableName }}
// Not thread safe. Auto generated.
{{.Comment -}}
type {{.Collection}} struct {
	Data           		[]*{{.Entity}}
	BeforeMapColumns	func(uint64, *{{.Entity}}) error
	AfterMapColumns 	func(uint64, *{{.Entity}}) error
}

// Make{{.Collection}} creates a new initialized collection. Auto generated.
func Make{{.Collection}}() {{.Collection}} {
	return {{.Collection}}{
		Data: make([]*{{.Entity}}, 0, 5),
	}
}

func (cc {{.Collection}}) scanColumns(cm *dml.ColumnMap,e *{{.Entity}}, idx uint64) error {
	if err := cc.BeforeMapColumns(idx, e); err != nil {
		return errors.WithStack(err)
	}
	if err := e.MapColumns(cm); err != nil {
		return errors.WithStack(err)
	}
	if err := cc.AfterMapColumns(idx, e); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// MapColumns implements dml.ColumnMapper interface. Auto generated.
func (cc {{.Collection}}) MapColumns(cm *dml.ColumnMap) error {
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
		e := New{{.Entity}}()
		if err := cc.scanColumns(cm, e, cm.Count); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, e)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next() {
			switch c := cm.Column(); c {
			{{- range .Columns.UniqueColumns -}}
			case "{{.Field }}"{{ range .Aliases }},"{{.}}"{{end}}:
				cm.Args = cm.Args.{{GoFuncNull .}}s(cc.{{ToGoCamelCase .Field}}s()...)
			{{- end}}
			{{- range .Columns.UniquifiedColumns }}
			case "{{.Field }}"{{ range .Aliases }},"{{.}}"{{end}}:
				cm.Args = cm.Args.{{GoFunc .}}s(cc.{{ToGoCamelCase .Field}}s()...){{end}}
			default:
				return errors.NewNotFoundf("[{{.Package}}] {{.Collection}} Column %q not found", c)
			}
		}
	default:
		return errors.NewNotSupportedf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}
{{ range .Columns.UniqueColumns }}
// {{ToGoCamelCase .Field}}s returns a slice or appends to a slice all values.
// Auto generated.
func (cc {{$.Collection}}) {{ToGoCamelCase .Field}}s(ret ...{{GoTypeNull .}}) []{{GoTypeNull .}} {
	if ret == nil {
		ret = make([]{{GoTypeNull .}}, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.{{ToGoCamelCase .Field}})
	}
	return ret
} {{end}}

{{- range .Columns.UniquifiedColumns }}
// {{ToGoCamelCase .Field}}s belongs to the column `{{.Field}}`
// and returns a slice or appends to a slice only unique values of that column.
// The values will be filtered internally in a Go map. No DB query gets
// executed. Auto generated.
func (cc {{$.Collection}}) {{ToGoCamelCase .Field}}s(ret ...{{GoType .}}) []{{GoType .}} {
	if ret == nil {
		ret = make([]{{GoType .}}, 0, len(cc.Data))
	}
	{{/*
		TODO: a reusable map and use different algorithms depending on the size
		of the cc.Data slice. Sometimes a for/for loop runs faster than a map.
	*/}}
	dupCheck := make(map[{{GoType .}}]struct{}, len(cc.Data))
	for _, e := range cc.Data {
		if _, ok := dupCheck[e.{{GoPrimitive .}}]; !ok {
			ret = append(ret, e.{{GoPrimitive .}})
			dupCheck[e.{{GoPrimitive .}}] = struct{}{}
		}
	}
	return ret
} {{end}}
