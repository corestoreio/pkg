// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dmlgen

// TplDBAC contains the template code = DataBaseAccessCode
const TplDBAC = `// {{.Entity}} represents a type for DB table {{.Tick}}{{.TableName}}{{.Tick}}
// Generated via dmlgen.
type {{.Entity}} struct {
	{{ range .Columns }}{{ToGoCamelCase .Field}} {{GoTypeNull .}} {{ $.Tick }}json:"{{.Field}},omitempty"{{ $.Tick }} {{.GoComment}}
{{ end }} }

// New{{.Entity}} creates a new pointer with pre-initialized fields.
func New{{.Entity}}() *{{.Entity}} {
	return &{{.Entity}}{}
}

// AssignLastInsertID updates the increment ID field with the last inserted ID
// from an INSERT operation. Implements dml.InsertIDAssigner
func (e *{{.Entity}}) AssignLastInsertID(id int64) {
	{{ range .Columns }}{{if .IsPK}} e.{{ToGoCamelCase .Field}} = {{GoTypeNull .}}(id) {{end}} {{ end }}
}

// MapColumns implements interface ColumnMapper only partially.
func (e *{{.Entity}}) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm{{range .Columns}}.{{GoFuncNull .}}(&e.{{ToGoCamelCase .Field}}){{end}}.Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c { {{range .Columns}}
		case "{{.Field }}"{{ range ColumnAliases .Field}},"{{.}}"{{end}}:
			cm.{{GoFuncNull .}}(&e.{{ToGoCamelCase .Field}}){{end}}
		default:
			return errors.NewNotFoundf("[{{.Package}}] {{.Entity}} Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}

// Reset resets the struct to its empty fields. TODO implement.
func (e *{{.Entity}}) Reset() *{{.Entity}} {
	return e
}

// {{.Collection}} represents a collection type for DB table {{ .TableName }}
// Generated via dmlgen.
type {{.Collection}} struct {
	Data           []*{{.Entity}}
	BeforeMapColumns	func(uint64, *{{.Entity}}) error
	AfterMapColumns 	func(uint64, *{{.Entity}}) error
}

func (cc *{{.Collection}}) scanColumns(cm *dml.ColumnMap,e *{{.Entity}}, idx uint64) error {
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

// MapColumns implements dml.ColumnMapper interface
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
		e := New{{.Entity}}()
		if err := cc.scanColumns(cm, e, cm.Count); err != nil {
			return errors.WithStack(err)
		}
		cc.Data = append(cc.Data, e)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next() {
			switch c := cm.Column(); c {
			{{- range .ExtractColumns -}}
			case "{{.Field }}"{{ range ColumnAliases .Field}},"{{.}}"{{end}}:
				cm.Args = cm.Args.{{GoFuncNull .}}s(cc.{{ToGoCamelCase .Field}}s()...)
			{{- end}}
			{{- range .ExtractUniquifiedColumns}}
			case "{{.Field }}"{{ range ColumnAliases .Field}},"{{.}}"{{end}}:
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
{{ range .ExtractColumns }}
// {{ToGoCamelCase .Field}}s returns a slice or appends to a slice all values.
func (cc *{{$.Collection}}) {{ToGoCamelCase .Field}}s(ret ...{{GoTypeNull .}}) []{{GoTypeNull .}} {
	if ret == nil {
		ret = make([]{{GoTypeNull .}}, 0, len(cc.Data))
	}
	for _, e := range cc.Data {
		ret = append(ret, e.{{ToGoCamelCase .Field}})
	}
	return ret
} {{end}}

{{- range .ExtractUniquifiedColumns }}
// {{ToGoCamelCase .Field}}s returns a slice or appends to a slice only unique
// values.
func (cc *{{$.Collection}}) {{ToGoCamelCase .Field}}s(ret ...{{GoType .}}) []{{GoType .}} {
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
`
