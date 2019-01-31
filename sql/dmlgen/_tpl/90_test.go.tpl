func TestNewTables(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	{{with .TestSQLDumpGlobPath}}defer dmltest.SQLDumpLoad(t, "{{.}}", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	}).Deferred(){{end}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	tbls, err := NewTables(ctx, ddl.WithConnPool(db))
	assert.NoError(t, err)

	tblNames := tbls.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{ {{- range $table := .Tables }}"{{ .TableName}}",{{- end}}}, tblNames)

	err = tbls.Validate(ctx)
	assert.NoError(t, err)
	var ps *pseudo.Service
	ps = pseudo.MustNewService(0, &pseudo.Options{Lang: "de",FloatMaxDecimals:6},
		pseudo.WithTagFakeFunc("website_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		pseudo.WithTagFakeFunc("store_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		{{- CustomCode "pseudo.MustNewService.Option" -}}
	)

	// TODO run those tests in parallel
	{{- range $table := .Tables }}
	t.Run("{{GoCamel .TableName}}_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableName{{GoCamel .TableName}})

		entINSERT := tbl.Insert().BuildValues()
		entINSERTStmtA := entINSERT.PrepareWithArgs(ctx)

		entSELECT := tbl.SelectByPK("*")
		entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders() // WithArgs generates the cached SQL string with key ""

		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where(
		{{range .Columns}}{{if and .IsPK .IsAutoIncrement}} dml.Column("{{.Field}}").LessOrEqual().Int(10),
		{{end}}{{end -}}).ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)

		for i := 0; i < 9; i++ {
			entIn := new({{GoCamel .TableName}})
			if err := ps.FakeData(entIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.{{GoCamel .TableName}}_Entity")(entINSERTStmtA.Record("", entIn).ExecContext(ctx))
			entINSERTStmtA.Reset()

			entOut := new({{GoCamel .TableName}})
			rowCount, err := entSELECTStmtA.Int64s(lID).Load(ctx, entOut)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)

			{{- range $col := $table.Columns }}
				{{if $col.IsString -}}
					assert.ExactlyLength(t, {{$col.CharMaxLength.Int64}}, &entIn.{{$table.GoCamelMaybePrivate $col.Field}}, &entOut.{{$table.GoCamelMaybePrivate $col.Field}}, "IDX%d: {{$table.GoCamelMaybePrivate $col.Field}} should match", lID)
				{{- else if not $col.IsSystemVersioned -}}
					assert.Exactly(t, entIn.{{$table.GoCamelMaybePrivate $col.Field}}, entOut.{{$table.GoCamelMaybePrivate $col.Field}}, "IDX%d: {{$table.GoCamelMaybePrivate $col.Field}} should match", lID)
				{{- end}}
			{{- end}}
		}
		dmltest.Close(t, entINSERTStmtA)

		entCol := New{{$table.CollectionName}}()
		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)

		entINSERTStmtA = entINSERT.WithCacheKey("row_count_%d", len(entCol.Data)).Replace().SetRowCount(len(entCol.Data)).PrepareWithArgs(ctx)

		lID := dmltest.CheckLastInsertID(t, "Error: {{$table.CollectionName}}")(entINSERTStmtA.Record("", entCol).ExecContext(ctx))
		dmltest.Close(t, entINSERTStmtA)
		t.Logf("Last insert ID into: %d", lID)
		t.Logf("INSERT queries: %#v", entINSERT.CachedQueries())
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
	})
	{{- end}}
}
