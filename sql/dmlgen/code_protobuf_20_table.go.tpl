// {{.Entity}} represents a single row for DB table `{{.TableName}}`. Auto generated.
message {{.Entity}} {
	{{- range .Columns}}
	{{SerializerType .}} {{.Field}} = {{.Pos}} [(gogoproto.customname)="{{ToGoCamelCase .Field}}" {{- SerializerCustomType .}}];
	{{- end}}
}

// {{.Collection}} represents multiple rows for DB table `{{.TableName}}`. Auto generated.
message {{.Collection}} {
	repeated {{.Entity}} Data = 1;
}
