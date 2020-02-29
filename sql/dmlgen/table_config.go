package dmlgen

import (
	"fmt"
	"strings"

	"github.com/corestoreio/errors"
)

// TableConfig used in conjunction with WithTableConfig and
// WithTableConfigDefault to apply different configurations for a generated
// struct and its struct collection.
type TableConfig struct {
	// Encoders add method receivers for, each struct, compatible with the
	// interface declarations in the various encoding packages. Supported
	// encoder names are: json resp. easyjson.
	Encoders []string
	// StructTags enables struct tags proactively for the whole struct. Allowed
	// values are: bson, db, env, json, protobuf, toml, yaml and xml. For bson,
	// json, yaml and xml the omitempty attribute has been set. If you need a
	// different struct tag for a specifiv column you must set the option
	// CustomStructTags.
	StructTags []string
	// CustomStructTags allows to specify custom struct tags for a specific
	// column. The slice must be balanced, means index i sets to the column name
	// and index i+1 to the desired struct tag.
	// 		[]string{"column_a",`json: ",omitempty"`,"column_b","`xml:,omitempty`"}
	// In case when foreign keys should be referenced:
	// 		[]string{"FieldNameX",`faker: "-"`,"FieldNameY","`xml:field_name_y,omitempty`"}
	// TODO CustomStructTags should be appended to StructTags
	CustomStructTags []string // balanced slice
	// AppendCustomStructTags []string // balanced slice TODO maybe this additionally
	// Comment adds custom comments to each struct type. Useful when relying on
	// 3rd party JSON marshaler code generators like easyjson or ffjson. If
	// comment spans over multiple lines each line will be checked if it starts
	// with the comment identifier (//). If not, the identifier will be
	// prepended.
	Comment string
	// ColumnAliases specifies different names used for a column. For example
	// customer_entity.entity_id can also be sales_order.customer_id, hence a
	// Foreign Key. The alias would be just: entity_id:[]string{"customer_id"}.
	ColumnAliases map[string][]string // key=column name value a list of aliases
	// UniquifiedColumns specifies columns which are non primary/unique key one
	// but should have a dedicated function to extract their unique primitive
	// values as a slice. Not allowed are text, blob and binary.
	UniquifiedColumns []string
	// PrivateFields list struct field names which should be private to avoid
	// accidentally leaking through encoders. Appropriate getter/setter methods
	// get generated.
	PrivateFields []string
	// FeaturesInclude if set includes only those features, otherwise
	// everything. Some features can only be included on Default level and not
	// on a per table level.
	FeaturesInclude FeatureToggle
	// FeaturesExclude excludes features while field FeaturesInclude is empty.
	// Some features can only be excluded on Default level and not on a per
	// table level.
	FeaturesExclude FeatureToggle
	// FieldMapFn can map a dbIdentifier (database identifier) of the current
	// table to a new name. dbIdentifier is in most cases the column name and in
	// cases of foreign keys, it is the table name.
	FieldMapFn func(dbIdentifier string) (newName string)
	lastErr    error
}

func (to *TableConfig) applyEncoders(t *Table, g *Generator) {
	encoders := make([]string, 0, len(to.Encoders)+len(g.defaultTableConfig.Encoders))
	encoders = append(encoders, to.Encoders...)
	encoders = append(encoders, g.defaultTableConfig.Encoders...)

	for i := 0; i < len(encoders) && to.lastErr == nil; i++ {
		switch enc := encoders[i]; enc {
		case "json", "easyjson":
			t.HasEasyJSONMarshaler = true
		case "protobuf", "fbs":
			t.HasSerializer = true // for now leave it in. maybe later PB gets added to the struct tags.
		default:
			to.lastErr = errors.NotSupported.Newf("[dmlgen] WithTableConfig: Table %q Encoder %q not supported", t.Table.Name, enc)
		}
	}
}

func (to *TableConfig) applyStructTags(t *Table, g *Generator) {
	for h := 0; h < len(t.Table.Columns) && to.lastErr == nil; h++ {
		c := t.Table.Columns[h]
		var buf strings.Builder

		structTags := make([]string, 0, len(to.StructTags)+len(g.defaultTableConfig.StructTags))
		structTags = append(structTags, to.StructTags...)
		structTags = append(structTags, g.defaultTableConfig.StructTags...)

		for lst, i := len(structTags), 0; i < lst && to.lastErr == nil; i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			// Maybe some types in the struct for a table don't need at
			// all an omitempty so build in some logic which creates the
			// tags more thoughtfully.
			switch tagName := structTags[i]; tagName {
			case "bson":
				fmt.Fprintf(&buf, `bson:"%s,omitempty"`, c.Field)
			case "db":
				fmt.Fprintf(&buf, `db:"%s"`, c.Field)
			case "env":
				fmt.Fprintf(&buf, `env:"%s"`, c.Field)
			case "json":
				fmt.Fprintf(&buf, `json:"%s,omitempty"`, c.Field)
			case "toml":
				fmt.Fprintf(&buf, `toml:"%s"`, c.Field)
			case "yaml":
				fmt.Fprintf(&buf, `yaml:"%s,omitempty"`, c.Field)
			case "xml":
				fmt.Fprintf(&buf, `xml:"%s,omitempty"`, c.Field)
			case "max_len":
				l := c.CharMaxLength.Int64
				if c.Precision.Valid {
					l = c.Precision.Int64
				}
				if l > 0 {
					fmt.Fprintf(&buf, `max_len:"%d"`, l)
				}
			case "protobuf":
				// github.com/gogo/protobuf/protoc-gen-gogo/generator/generator.go#L1629 Generator.goTag
				// The tag is a string like "varint,2,opt,name=fieldname,def=7" that
				// identifies details of the field for the protocol buffer marshaling and unmarshaling
				// code.  The fields are:
				// 	wire encoding
				// 	protocol tag number
				// 	opt,req,rep for optional, required, or repeated
				// 	packed whether the encoding is "packed" (optional; repeated primitives only)
				// 	name= the original declared name
				// 	enum= the name of the enum type if it is an enum-typed field.
				// 	proto3 if this field is in a proto3 message
				// 	def= string representation of the default value, if any.
				// The default value must be in a representation that can be used at run-time
				// to generate the default value. Thus bools become 0 and 1, for instance.

				// CYS: not quite sure if struct tags are really needed
				// pbType := "TODO"
				// customType := ",customtype=github.com/gogo/protobuf/test.TODO"
				// fmt.Fprintf(&buf, `protobuf:"%s,%d,opt,name=%s%s"`, pbType, c.Pos, c.Field, customType)
			default:
				to.lastErr = errors.NotSupported.Newf("[dmlgen] WithTableConfig: Table %q Tag %q not supported", t.Table.Name, tagName)
			}
		}
		c.StructTag = buf.String()
	} // end Columns loop
}

func (to *TableConfig) applyCustomStructTags(t *Table) {
	for i := 0; i < len(to.CustomStructTags) && to.lastErr == nil; i += 2 {
		for _, c := range t.Table.Columns {
			if c.Field == to.CustomStructTags[i] {
				c.StructTag = to.CustomStructTags[i+1]
			}
		}
		if t.customStructTagFields == nil {
			t.customStructTagFields = make(map[string]string)
		}

		// copy data to handle foreign keys if they should have struct tags.
		// as key use the kcuce.ReferencedTableName.String and value the struct tag itself.
		t.customStructTagFields[to.CustomStructTags[i]] = "`" + to.CustomStructTags[i+1] + "`"
	}
}

func (to *TableConfig) applyPrivateFields(t *Table) {
	if len(to.PrivateFields) > 0 && t.privateFields == nil {
		t.privateFields = make(map[string]bool)
	}
	for _, pf := range to.PrivateFields {
		t.privateFields[pf] = true
	}
}

func (to *TableConfig) applyComments(t *Table) {
	var buf strings.Builder
	lines := strings.Split(to.Comment, "\n")
	for i := 0; i < len(lines) && to.lastErr == nil && to.Comment != ""; i++ {
		line := lines[i]
		if !strings.HasPrefix(line, "//") {
			buf.WriteString("// ")
		}
		buf.WriteString(line)
		buf.WriteByte('\n')
	}
	t.Comment = strings.TrimSpace(buf.String())
}

func (to *TableConfig) applyColumnAliases(t *Table) {
	if to.lastErr != nil {
		return
	}
	// With iteration looping it this way, we can easily proof if the
	// developer has correctly written the column name. We might have
	// more work here but the developer has a better experience when a
	// column can't be found.
	for colName, aliases := range to.ColumnAliases {
		found := false
		for _, col := range t.Table.Columns {
			if col.Field == colName {
				found = true
				col.Aliases = aliases
			}
		}
		if !found {
			to.lastErr = errors.NotFound.Newf("[dmlgen] WithTableConfig:ColumnAliases: For table %q the Column %q cannot be found.",
				t.Table.Name, colName)
			return
		}
	}
}

// skips text and blob and varbinary and json and geo
func (to *TableConfig) applyUniquifiedColumns(t *Table) {
	for i := 0; i < len(to.UniquifiedColumns) && to.lastErr == nil; i++ {
		cn := to.UniquifiedColumns[i]
		found := false
		for _, c := range t.Table.Columns {
			if c.Field == cn && !c.IsBlobDataType() {
				c.Uniquified = true
				found = true
			}
		}
		if !found {
			to.lastErr = errors.NotFound.Newf("[dmlgen] WithTableConfig:UniquifiedColumns: For table %q the Column %q cannot be found in the list of available columns or its data type is not allowed.",
				t.Table.Name, cn)
		}
	}
}
