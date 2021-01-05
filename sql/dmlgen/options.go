package dmlgen

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/corestoreio/pkg/sql/dml"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
)

// Option represents a sortable option for the NewGenerator function. Each option
// function can be applied in a mixed order.
type Option struct {
	// sortOrder specifies the precedence of an option.
	sortOrder int
	fn        func(*Generator) error
}

// WithTableConfigDefault sets for all tables the same configuration but can be
// overwritten on a per table basis with function WithTableConfig. Current
// behaviour: Does not apply the default config to all tables but only to those
// which have set an empty
// `WithTableConfig("table_name",&dmlgen.TableConfig{})`.
func WithTableConfigDefault(opt TableConfig) (o Option) {
	o.sortOrder = 149 // must be applied before WithTableConfig
	o.fn = func(g *Generator) (err error) {
		g.defaultTableConfig = opt
		return opt.lastErr
	}
	return o
}

func defaultFieldMapFn(s string) string {
	return s
}

// WithTableConfig applies options to an existing table, identified by the table
// name used as map key. Options are custom struct or different encoders.
// Returns a not-found error if the table cannot be found in the `Generator` map.
func WithTableConfig(tableName string, opt *TableConfig) (o Option) {
	// Panic as early as possible.
	if len(opt.CustomStructTags)%2 == 1 {
		panic(errors.Fatal.Newf("[dmlgen] WithTableConfig: Table %q option CustomStructTags must be a balanced slice.", tableName))
	}
	o.sortOrder = 150
	o.fn = func(g *Generator) (err error) {
		t, ok := g.Tables[tableName]
		if t == nil || !ok {
			return errors.NotFound.Newf("[dmlgen] WithTableConfig: Table %q not found.", tableName)
		}
		opt.applyEncoders(t, g)
		opt.applyStructTags(t, g)
		opt.applyCustomStructTags(t)
		opt.applyPrivateFields(t)
		opt.applyComments(t)
		opt.applyColumnAliases(t)
		opt.applyUniquifiedColumns(t)
		t.featuresInclude = opt.FeaturesInclude | g.defaultTableConfig.FeaturesInclude
		t.featuresExclude = opt.FeaturesExclude | g.defaultTableConfig.FeaturesExclude
		t.fieldMapFn = opt.FieldMapFn
		if t.fieldMapFn == nil {
			t.fieldMapFn = defaultFieldMapFn
		}
		return opt.lastErr
	}
	return o
}

// ForeignKeyOptions applies to WithForeignKeyRelationships
type ForeignKeyOptions struct {
	IncludeRelationShips []string
	ExcludeRelationships []string
	// MToMRelations        bool
}

// WithForeignKeyRelationships analyzes the foreign keys which points to a table
// and adds them as a struct field name. For example:
// customer_address_entity.parent_id is a foreign key to
// customer_entity.entity_id hence the generated struct CustomerEntity has a new
// field which gets named CustomerAddressEntityCollection, pointing to type
// CustomerAddressEntityCollection. includeRelationShips and
// excludeRelationships must be a balanced slice in the notation of
// "table1.column1","table2.column2". For example:
// 		"customer_entity.store_id", "store.store_id" which means that the
// struct CustomerEntity won't have a field to the Store struct (1:1
// relationship) (in case of excluding). The reverse case can also be added
// 		"store.store_id", "customer_entity.store_id" which means that the
// Store struct won't or will have a field pointing to the
// CustomerEntityCollection (1:M relationship).
// Setting includeRelationShips to nil will include all relationships.
// Wildcards are supported.
func WithForeignKeyRelationships(ctx context.Context, db dml.Querier, o ForeignKeyOptions) (opt Option) {
	opt.sortOrder = 210 // must run at the end or where the end is near ;-)
	opt.fn = func(g *Generator) (err error) {
		if len(o.ExcludeRelationships)%2 == 1 {
			return errors.Fatal.Newf("[dmlgen] excludeRelationships must be balanced slice. Read the doc.")
		}
		if len(o.IncludeRelationShips)%2 == 1 {
			return errors.Fatal.Newf("[dmlgen] includeRelationShips must be balanced slice. Read the doc.")
		}
		if db != nil {
			g.kcu, err = ddl.LoadKeyColumnUsage(ctx, db, g.sortedTableNames()...)
			if err != nil {
				return errors.WithStack(err)
			}
			g.kcuRev = ddl.ReverseKeyColumnUsage(g.kcu)

			g.krs, err = ddl.GenerateKeyRelationships(ctx, db, g.kcu)
			if err != nil {
				return errors.WithStack(err)
			}
		} else if isDebug() {
			println("DEBUG[WithForeignKeyRelationships] db dml.Querier is nil. Did not load LoadKeyColumnUsage and GenerateKeyRelationships.")
		}

		// TODO(CSC) maybe excludeRelationships can contain a wild card to disable
		//  all embedded structs from/to a type. e.g. "customer_entity.website_id",
		//  "store_website.website_id", for CustomerEntity would become
		//  "*.website_id", "store_website.website_id", to disable all tables which
		//  have a foreign key to store_website.
		// TODO optimize code.

		g.krsExclude = make(map[string]bool, len(o.ExcludeRelationships)/2)
		for i := 0; i < len(o.ExcludeRelationships); i += 2 {
			mainTable := strings.Split(o.ExcludeRelationships[i], ".") // mainTable.mainColumn
			mainTab := mainTable[0]
			mainCol := mainTable[1]
			referencedTable := strings.Split(o.ExcludeRelationships[i+1], ".") // referencedTable.referencedColumn
			referencedTab := referencedTable[0]
			referencedCol := referencedTable[1]

			var buf strings.Builder
			buf.WriteString(mainTab)
			if mainCol == "*" {
				g.krsExclude[buf.String()] = true
			}
			buf.WriteByte('.')
			buf.WriteString(mainCol)
			if referencedTab == "*" && referencedCol == "*" {
				g.krsExclude[buf.String()] = true
			}

			buf.WriteByte(':')
			buf.WriteString(referencedTab)
			if referencedCol == "*" {
				g.krsExclude[buf.String()] = true
			}
			buf.WriteByte('.')
			buf.WriteString(referencedCol)
			g.krsExclude[buf.String()] = true
		}

		if len(o.IncludeRelationShips) > 0 {
			g.krsInclude = make(map[string]bool, len(o.IncludeRelationShips)/2)
			for i := 0; i < len(o.IncludeRelationShips); i += 2 {
				mainTable := strings.Split(o.IncludeRelationShips[i], ".") // mainTable.mainColumn
				mainTab := mainTable[0]
				mainCol := mainTable[1]
				referencedTable := strings.Split(o.IncludeRelationShips[i+1], ".") // referencedTable.referencedColumn
				referencedTab := referencedTable[0]
				referencedCol := referencedTable[1]

				var buf strings.Builder
				buf.WriteString(mainTab)
				if mainCol == "*" {
					g.krsInclude[buf.String()] = true
				}
				buf.WriteByte('.')
				buf.WriteString(mainCol)
				if referencedTab == "*" && referencedCol == "*" {
					g.krsInclude[buf.String()] = true
				}
				buf.WriteByte(':')
				buf.WriteString(referencedTab)
				if referencedCol == "*" {
					g.krsInclude[buf.String()] = true
				}
				buf.WriteByte('.')
				buf.WriteString(referencedCol)
				g.krsInclude[buf.String()] = true
			}
		}

		if isDebug() {
			g.krs.Debug(os.Stdout)
			debugMapSB("krsInclude", g.krsInclude)
			debugMapSB("krsExclude", g.krsExclude)
		}

		return nil
	}
	return opt
}

// WithTable sets a table and its columns. Allows to overwrite a table fetched
// with function WithTablesFromDB. Argument `options` can be set to "overwrite"
// and/or "view". Each option is its own slice entry.
func WithTable(tableName string, columns ddl.Columns, options ...string) (opt Option) {
	checkAutoIncrement := func(previousSetting uint8) uint8 {
		if previousSetting > 0 {
			return previousSetting
		}
		for _, o := range options {
			if strings.ToLower(o) == "view" {
				return 1 // nope
			}
		}
		for _, c := range columns {
			if c.IsAutoIncrement() {
				return 2 // yes
			}
		}
		return 1 // nope
	}

	opt.sortOrder = 10
	opt.fn = func(g *Generator) error {
		isOverwrite := len(options) > 0 && options[0] == "overwrite"
		t, ok := g.Tables[tableName]
		if ok && isOverwrite {
			for ci, ct := range t.Table.Columns {
				for _, cc := range columns {
					if ct.Field == cc.Field {
						t.Table.Columns[ci] = cc
					}
				}
			}
		} else {
			t = &Table{
				Table: ddl.NewTable(tableName, columns...),
				debug: isDebug(),
			}
		}
		t.HasAutoIncrement = checkAutoIncrement(t.HasAutoIncrement)
		g.Tables[tableName] = t
		return nil
	}
	return opt
}

// WithTablesFromDB queries the information_schema table and loads the column
// definition of the provided `tables` slice. It adds the tables to the `Generator`
// map. Once added a call to WithTableConfig can add additional configurations.
func WithTablesFromDB(ctx context.Context, db *dml.ConnPool, tables ...string) (opt Option) {
	checkAutoIncrement := func(oneTbl *ddl.Table) uint8 {
		if oneTbl.IsView() {
			return 1 // nope
		}
		for _, c := range oneTbl.Columns {
			if c.IsAutoIncrement() {
				return 2 // yes
			}
		}
		return 1 // nope
	}

	opt.sortOrder = 1
	opt.fn = func(g *Generator) error {
		nt, err := ddl.NewTables(ddl.WithLoadTables(ctx, db.DB, tables...))
		if err != nil {
			return errors.WithStack(err)
		}

		if len(tables) == 0 {
			tables = nt.Tables() // use all tables from the DB
		}

		for _, tblName := range tables {
			oneTbl := nt.MustTable(tblName)
			g.Tables[tblName] = &Table{
				Table:            oneTbl,
				HasAutoIncrement: checkAutoIncrement(oneTbl),
				debug:            isDebug(),
			}
		}
		return nil
	}
	return opt
}

// SerializerConfig applies various optional settings to WithProtobuf and/or
// WithFlatbuffers and/or WithTypeScript.
type SerializerConfig struct {
	PackageImportPath string
	AdditionalHeaders []string
}

// WithProtobuf enables protocol buffers as a serialization method. Argument
// headerOptions is optional. Heads up: This function also sets the internal
// Serializer field and all types will get adjusted to the minimum protobuf
// types. E.g. uint32 minimum instead of uint8/uint16. So if the Generator gets
// created multiple times to separate the creation of code, the WithProtobuf
// function must get set for Generator objects. See package store.
func WithProtobuf(sc *SerializerConfig) (opt Option) {
	_, pkg := filepath.Split(sc.PackageImportPath)

	opt.sortOrder = 110
	opt.fn = func(g *Generator) error {
		g.Serializer = "protobuf"
		g.PackageSerializer = pkg
		g.PackageSerializerImportPath = sc.PackageImportPath
		g.SerializerHeaderOptions = []string{
			"(gogoproto.typedecl_all) = false",
			"(gogoproto.goproto_getters_all) = false",
			"(gogoproto.unmarshaler_all) = true",
			"(gogoproto.marshaler_all) = true",
			"(gogoproto.sizer_all) = true",
			"(gogoproto.goproto_unrecognized_all) = false",
		}
		g.SerializerHeaderOptions = append(g.SerializerHeaderOptions, sc.AdditionalHeaders...)
		return nil
	}
	return opt
}

// WithBuildTags adds your build tags to the file header. Each argument
// represents a build tag line.
func WithBuildTags(lines ...string) (opt Option) {
	opt.sortOrder = 112
	opt.fn = func(g *Generator) error {
		g.BuildTags = append(g.BuildTags, lines...)
		return nil
	}
	return opt
}

// WithCustomCode inserts at the marker position your custom Go code. For
// available markers search these .go files for the map access of field
// *Generator.customCode. An example got written in
// TestGenerate_Tables_Protobuf_Json. If the marker does not exists or has a
// typo, no error gets reported and no code gets written.
func WithCustomCode(marker, code string) (opt Option) {
	opt.sortOrder = 112
	opt.fn = func(g *Generator) error {
		if g.customCode == nil {
			g.customCode = make(map[string]func(*Generator, *Table, io.Writer))
		}
		g.customCode[marker] = func(_ *Generator, _ *Table, w io.Writer) {
			w.Write([]byte(code))
		}
		return nil
	}
	return opt
}

// WithCustomCodeFunc same as WithCustomCode but allows access to meta data. The
// func fn takes as first argument the main Generator where access to package
// global configuration is possible. If the scope of the marker is within a
// table, then argument t gets set, otherwise it is nil. The output must be
// written to w.
func WithCustomCodeFunc(marker string, fn func(g *Generator, t *Table, w io.Writer)) (opt Option) {
	opt.sortOrder = 113
	opt.fn = func(g *Generator) error {
		if g.customCode == nil {
			g.customCode = make(map[string]func(*Generator, *Table, io.Writer))
		}
		g.customCode[marker] = fn
		return nil
	}
	return opt
}
