func TestNewTables(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	{{with .TestSQLDumpGlobPath}}defer dmltest.SQLDumpLoad(t, "{{.}}", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	})(){{end}}

	ctx := context.TODO()
	tbls, err := NewTables(ctx, db.DB)
	assert.NoError(t, err)

	tblNames := tbls.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{ {{- range $table := .Tables }}"{{ .TableName}}",{{- end}}}, tblNames)

	err = tbls.Validate(ctx)
	assert.NoError(t, err)
	ps := pseudo.MustNewService(0, &pseudo.Options{Lang: "de"},
		pseudo.WithTagFakeFunc("website_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		pseudo.WithTagFakeFunc("store_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
	)

	// TODO run those tests in parallel
	{{- range $table := .Tables }}
	t.Run("{{ToGoCamelCase .TableName}}_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableName{{ToGoCamelCase .TableName}})

		inStmt, err := ccd.Insert().BuildValues().Prepare(ctx) // Do not use Ignore() to suppress DB errors.
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 9; i++ {
			entityIn := new({{ToGoCamelCase .TableName}})
			assert.NoError(t, ps.FakeData(entityIn))

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.{{ToGoCamelCase .TableName}}_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new({{ToGoCamelCase .TableName}})
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "RowCount did not match")

			{{- range $col := $table.Columns }}
				{{if $col.IsString -}}
					assert.ExactlyLength(t, {{$col.CharMaxLength.Int64}}, &entityIn.{{ToGoCamelCase $col.Field}}, &entityOut.{{ToGoCamelCase $col.Field}}, "IDX%d: {{ToGoCamelCase $col.Field}} should match", lID)
				{{- else if not $col.IsSystemVersioned -}}
					assert.Exactly(t, entityIn.{{ToGoCamelCase $col.Field}}, entityOut.{{ToGoCamelCase $col.Field}}, "IDX%d: {{ToGoCamelCase $col.Field}} should match", lID)
				{{- end}}
			{{- end}}
		}
	})
	{{- end}}
}
