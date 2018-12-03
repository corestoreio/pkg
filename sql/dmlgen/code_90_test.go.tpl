func TestNewTables(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	defer dmltest.SQLDumpLoad(t, "test_*_tables.sql", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	})()

	ctx := context.TODO()
	tbls, err := NewTables(ctx, db.DB)
	assert.NoError(t, err)

	tblNames := tbls.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{ {{- range $table := .Tables }}"{{ .TableName}}",{{- end}}}, tblNames)

	err = tbls.Validate(ctx)
	assert.NoError(t, err)

	// TODO run those tests in parallel
	{{- range $table := .Tables }}
	t.Run("{{ToGoCamelCase .TableName}}_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableName{{ToGoCamelCase .TableName}})

		inStmt, err := ccd.Insert().Ignore().BuildValues().Prepare(ctx)
		assert.NoError(t, err, "%+v", err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 5; i++ {
			entityIn := new({{ToGoCamelCase .TableName}})
			assert.NoError(t, faker.FakeData(entityIn))

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.{{ToGoCamelCase .TableName}}_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new({{ToGoCamelCase .TableName}})
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err, "%+v", err)
			assert.Exactly(t, uint64(1), rowCount, "RowCount did not match")

			// assert.Exactly(t, entityIn.Scope, entityOut.Scope, "Scope did not match")
			// assert.Exactly(t, entityIn.ScopeID, entityOut.ScopeID, "ScopeID did not match")
			// assert.Exactly(t, entityIn.Expires, entityOut.Expires, "Expires did not match")
			// assert.Exactly(t, entityIn.Path, entityOut.Path, "Path did not match")
			// assert.Exactly(t, entityIn.Value, entityOut.Value, "Value did not match")
		}
	})
	{{- end}}
}
