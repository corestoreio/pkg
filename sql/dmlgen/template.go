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

const TplEntity = `// {{.Entity}} represents a type for DB table {{.Tick}}{{.TableName}}{{.Tick}}
// Generated via dmlgen.
type {{.Entity}} struct {
	{{ range .Columns }}{{ToGoCamelCase .Field}} {{MySQLToGoType .}} {{ $.Tick }}json:"{{.Field}},omitempty"{{ $.Tick }} {{.GoComment}}
{{ end }} }

// New{{.Entity}} creates a new pointer with pre-initialized fields.
func New{{.Entity}}() *{{.Entity}} {
	return &{{.Entity}}{}
}

// AssignLastInsertID updates the increment ID field with the last inserted ID
// from an INSERT operation. Implements dml.InsertIDAssigner
func (e *{{.Entity}}) AssignLastInsertID(id int64) {
	{{ range .Columns }}{{if .IsPK}} e.{{ToGoCamelCase .Field}} = {{MySQLToGoType .}}(id) {{end}} {{ end }}
}

// MapColumns implements interface ColumnMapper only partially.
func (e *{{.Entity}}) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm{{range .Columns}}.{{GoTypeFuncName .}}(&e.{{ToGoCamelCase .Field}}){{end}}.Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c { {{range .Columns}}
		case "{{.Field }}"{{ range ColumnAliases .Field}},"{{.}}"{{end}}:
			cm.{{GoTypeFuncName .}}(&e.{{ToGoCamelCase .Field}}){{end}}
		default:
			return errors.NewNotFoundf("[{{.Package}}] {{.Entity}} Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}
`

const TplCollection = `// {{.Collection}} represents a collection type for DB table {{ .TableName }}
// Generated via dmlgen.
type {{.Collection}} struct {
	Data           []*{{.Entity}}
	BeforeMapColumns	func(uint64, *{{.Entity}}) error
	AfterMapColumns 	func(uint64, *{{.Entity}}) error
}

// MapColumns implements dml.ColumnMapper interface
func (cc *{{.Collection}}) MapColumns(cm *dml.ColumnMap) error {
	switch m := cm.Mode(); m {
	case dml.ColumnMapEntityReadAll, dml.ColumnMapEntityReadSet:
		for i, p := range cc.Data {
			if err := cc.BeforeMapColumns(uint64(i), p); err != nil {
				return errors.WithStack(err)
			}
			if err := p.MapColumns(cm); err != nil {
				return errors.WithStack(err)
			}
			if err := cc.AfterMapColumns(uint64(i), p); err != nil {
				return errors.WithStack(err)
			}
		}
	case dml.ColumnMapScan:
		if cm.Count == 0 {
			cc.Data = cc.Data[:0]
		}
		p := New{{.Entity}}()

		if err := cc.BeforeMapColumns(cm.Count, p); err != nil {
			return errors.WithStack(err)
		}
		if err := p.MapColumns(cm); err != nil {
			return errors.WithStack(err)
		}
		if err := cc.AfterMapColumns(cm.Count, p); err != nil {
			return errors.WithStack(err)
		}

		cc.Data = append(cc.Data, p)
	case dml.ColumnMapCollectionReadSet:
		for cm.Next() {
			switch c := cm.Column(); c { {{range .Columns}}
			case "{{.Field }}"{{ range ColumnAliases .Field}},"{{.}}"{{end}}:
				cm.Args = cm.Args.{{GoTypeFuncName .}}(cc.{{ToGoCamelCase .Field}}s()...){{end}}
			default:
				return errors.NewNotFoundf("[{{.Package}}] {{.Collection}} Column %q not found", c)
			}
		}
	default:
		return errors.NewNotSupportedf("[dml] Unknown Mode: %q", string(m))
	}
	return cm.Err()
}

`
