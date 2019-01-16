// {{.Entity}} represents a single row for DB table `{{.TableName}}`. Auto generated.
message {{.Entity}} {
	{{- range .Columns}}{{ if IsFieldPublic .Field }}
	{{SerializerType .}} {{.Field}} = {{.Pos}} [(gogoproto.customname)="{{GoCamel .Field}}" {{- SerializerCustomType .}}];{{end}}
	{{- end}}
}

// {{.Collection}} represents multiple rows for DB table `{{.TableName}}`. Auto generated.
message {{.Collection}} {
	repeated {{.Entity}} Data = 1;
}
