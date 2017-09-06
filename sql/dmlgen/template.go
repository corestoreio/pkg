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
	rc dml.RowConvert
	{{ range .Columns }}{{.Field | ToGoCamelCase}} {{.GoPrimitive}} {{ $.Tick }}json:",omitempty"{{ $.Tick }} {{.GoComment}}
{{ end }} }

// AssignLastInsertID updates the increment ID field with the last inserted ID
// from an INSERT operation. Implements dml.InsertIDAssigner
func (e *{{.Entity}}) AssignLastInsertID(id int64) {
	{{ range .Columns }}{{if .IsPK}} e.{{. | FieldName}} = uint64(id) {{end}}
{{ end }}
}

// RowScan loads a single row from a SELECT statement returning only one row.
func (e *{{.Entity}}) RowScan(r *sql.Rows) error {
	if err := e.rc.Scan(r); err != nil {
		return errors.WithStack(err)
	}
	return e.assign(&e.rc)
}

func (e *{{.Entity}}) assign(rc *dml.RowConvert) (err error) {
	for i, c := range rc.Columns {
		b := rc.Index(i)
		switch c { {{ range .Columns }}
			case "{{.Field }}":
				e.{{. | FieldName}}, err = b.{{.RowConvertName}}(){{end}}
		default:
			return errors.NewNotFoundf("[{{.Package}}] Column %q not found", c)
		}
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (e *{{.Entity}}) AppendArgs(args dml.Arguments, columns []string) (_ dml.Arguments, err error) {
	l := len(columns)
	if l == 1 {
		return e.appendArgs(args, columns[0])
	}
	if l == 0 {
		return args.Uint64(e.ID).Str(e.Name).NullString(e.Email), nil // except auto inc column ;-)
	}
	for _, col := range columns {
		if args, err = e.appendArgs(args, col); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, err
}

func (e *{{.Entity}}) appendArgs(args dml.Arguments, column string) (_ dml.Arguments, err error) {
	switch column { {{ range .Columns }}
	case "{{.Field }}":
		args = args.{{.RowConvertName}}(e.{{. | FieldName}}){{end}}
	default:
		return nil, errors.NewNotFoundf("[dml_test] dmlPerson Column %q not found", column)
	}
	return args, err
}

`

const TplCollection = `// {{.Collection}} represents a collection type for DB table {{ .TableName }}
// Generated via dmlgen.
type {{.Collection}} struct {
	rc        	   dml.RowConvert
	Data           []*{{.Entity}}
	EventAfterScan []func(*{{.Entity}})
}

func (c *{{.Collection}}) AppendArgs(args Arguments, columns []string) (_ Arguments, err error) {
	if len(columns) != 1 {
		// INSERT STATEMENT requesting all columns or specific columns
		for _, p := range c.Data {
			if args, err = p.AppendArgs(args, columns); err != nil {
				return nil, errors.WithStack(err)
			}
		}
		return args, err
	}

	// SELECT, DELETE or UPDATE or INSERT with one column
	column := columns[0]
	var ids []uint64
	var names []string
	var emails []NullString
	for _, p := range c.Data {
		switch column {
		case "id":
			ids = append(ids, p.ID)
		case "name":
			names = append(names, p.Name)
		case "email":
			emails = append(emails, p.Email)
			// case "key": don't add key, it triggers a test failure condition
		default:
			return nil, errors.NewNotFoundf("[dml_test] dmlPerson Column %q not found", column)
		}
	}

	switch column {
	case "id":
		args = args.Uint64s(ids...)
	case "name":
		args = args.Strs(names...)
	case "email":
		args = args.NullString(emails...)
	}

	return args, nil
}

func (c *{{.Collection}}) RowScan(r *sql.Rows) error {
	if err := c.rc.Scan(r); err != nil {
		return errors.WithStack(err)
	}

	p := new(dmlPerson)
	if err := p.assign(&c.rc); err != nil {
		return errors.WithStack(err)
	}
	c.Data = append(c.Data, p)
	return nil
}

`
