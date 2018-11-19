// <code_entity.go.tpl>

// {{.Entity}} represents a single row for DB table `{{.TableName}}`.
// Auto generated.
{{.Comment -}}
{{- if .JsonMarshaler }}//easyjson:json{{end}}
type {{.Entity}} struct {
{{range .Columns}}{{ToGoCamelCase .Field}} {{GoTypeNull .}}
		{{- if ne .StructTag "" -}}`{{.StructTag}}`{{- end}} {{.GoComment}}
{{end}} }

// New{{.Entity}} creates a new pointer with pre-initialized fields. Auto
// generated.
func New{{.Entity}}() *{{.Entity}} {
	return &{{.Entity}}{}
}

// AssignLastInsertID updates the increment ID field with the last inserted ID
// from an INSERT operation. Implements dml.InsertIDAssigner. Auto generated.
func (e *{{.Entity}}) AssignLastInsertID(id int64) {
	{{range .Columns}}{{if and .IsPK .IsAutoIncrement}} e.{{ToGoCamelCase .Field}} = {{GoType .}}(id)
	{{end}}{{end -}}
}

// MapColumns implements interface ColumnMapper only partially. Auto generated.
func (e *{{.Entity}}) MapColumns(cm *dml.ColumnMap) error {
	if cm.Mode() == dml.ColumnMapEntityReadAll {
		return cm{{range .Columns}}.{{GoFuncNull .}}(&e.{{ToGoCamelCase .Field}}){{end}}.Err()
	}
	for cm.Next() {
		switch c := cm.Column(); c { {{range .Columns}}
			case "{{.Field}}"{{range .Aliases}},"{{.}}"{{end}}:
				cm.{{GoFuncNull .}}(&e.{{ToGoCamelCase .Field}}){{end}}
			default:
				return errors.NotFound.Newf("[{{.Package}}] {{.Entity}} Column %q not found", c)
		}
	}
	return errors.WithStack(cm.Err())
}
// </code_entity.go.tpl>

