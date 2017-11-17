// {{.Entity}} represents a single row for DB table `{{.TableName}}`. Auto generated.
message {{.Entity}} {
	{{- range .Columns}}
	{{ProtoType .}} {{.Field}} = {{.Pos}} [(gogoproto.customname)="{{ToGoCamelCase .Field}}" {{- ProtoCustomType .}}];
	{{- end}}
}

// {{.Collection}} represents multiple rows for DB table `{{.TableName}}`. Auto generated.
message {{.Collection}} {
	repeated {{.Entity}} Data = 1;
}
