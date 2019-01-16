// {{.Entity}} represents a single row for DB table `{{.TableName}}`. Auto generated.
table {{.Entity}} {
	{{- range .Columns}}
		{{GoCamel .Field}}:{{SerializerType .}}; // {{.Field}}
	{{- end}}
}

// {{.Collection}} represents multiple rows for DB table `{{.TableName}}`. Auto generated.
table {{.Collection}} {
	Data:[{{.Entity}}];
}
