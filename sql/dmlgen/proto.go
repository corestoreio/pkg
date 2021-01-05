package dmlgen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/strs"
)

// ProtocOptions allows to modify the protoc CLI command.
type ProtocOptions struct {
	BuildTags          []string
	WorkingDirectory   string
	ProtoGen           string // default gofast, options: gogo, gogofast, gogofaster
	GRPC               bool
	GRPCGatewayOutMap  []string // GRPC must be enabled in the above field
	GRPCGatewayOutPath string   // GRPC must be enabled in the above field
	ProtoPath          []string
	GoGoOutPath        string
	GoGoOutMap         []string
	SwaggerOutPath     string
	CustomArgs         []string
	// TODO add validation plugin, either
	//  https://github.com/mwitkow/go-proto-validators as used in github.com/gogo/grpc-example/proto/example.proto
	//  This github.com/mwitkow/go-proto-validators seems dead.
	//  or https://github.com/envoyproxy/protoc-gen-validate
	//  Requirement: error messages must be translatable and maybe use an errors.Kind type
}

var defaultProtoPaths = make([]string, 0, 8)

func init() {
	preDefinedPaths := [...]string{
		build.Default.GOPATH + "/src/",
		build.Default.GOPATH + "/src/github.com/gogo/protobuf/protobuf/",
		build.Default.GOPATH + "/src/github.com/gogo/googleapis/",
		"vendor/github.com/grpc-ecosystem/grpc-gateway/",
		"vendor/github.com/gogo/googleapis/",
		"vendor/",
		".",
	}
	for _, pdp := range preDefinedPaths {
		if _, err := os.Stat(pdp); !os.IsNotExist(err) {
			defaultProtoPaths = append(defaultProtoPaths, pdp)
		}
	}
}

func (po *ProtocOptions) toArgs() []string {
	gogoOut := make([]string, 0, 4)
	if po.GRPC {
		gogoOut = append(gogoOut, "plugins=grpc")
		if po.GRPCGatewayOutMap == nil {
			po.GRPCGatewayOutMap = []string{
				"allow_patch_feature=false",
				"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types",
				"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types",
				"Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types",
				"Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api",
				"Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types",
			}
		}
		if po.GRPCGatewayOutPath == "" {
			po.GRPCGatewayOutPath = "."
		}
	}
	if po.GoGoOutPath == "" {
		po.GoGoOutPath = "."
	}
	if po.GoGoOutMap == nil {
		po.GoGoOutMap = []string{
			"Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api",
			"Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types",
			"Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types",
			"Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types",
			"Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types",
		}
	}
	gogoOut = append(gogoOut, po.GoGoOutMap...)

	if po.ProtoPath == nil {
		po.ProtoPath = append(po.ProtoPath, defaultProtoPaths...)
	}

	if po.ProtoGen == "" {
		po.ProtoGen = "gofast"
	} else {
		switch po.ProtoGen {
		case "gofast", "gogo", "gogofast", "gogofaster":
			// ok
		default:
			panic(fmt.Sprintf("[dmlgen] ProtoGen CLI command %q not supported, allowed: gofast, gogo, gogofast, gogofaster", po.ProtoGen))
		}
	}

	// To generate PHP Code replace `gogo_out` with `php_out`.
	// Java bit similar. Java has ~15k LOC, Go ~3.7k
	args := []string{
		"--" + po.ProtoGen + "_out", fmt.Sprintf("%s:%s", strings.Join(gogoOut, ","), po.GoGoOutPath),
		"--proto_path", strings.Join(po.ProtoPath, ":"),
	}
	if po.GRPC && len(po.GRPCGatewayOutMap) > 0 {
		args = append(args, "--grpc-gateway_out="+strings.Join(po.GRPCGatewayOutMap, ",")+":"+po.GRPCGatewayOutPath)
	}
	if po.SwaggerOutPath != "" {
		args = append(args, "--swagger_out="+po.SwaggerOutPath)
	}
	return append(args, po.CustomArgs...)
}

func (po *ProtocOptions) chdir() (deferred func(), _ error) {
	deferred = func() {}
	if po.WorkingDirectory != "" {
		oldWD, err := os.Getwd()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := os.Chdir(po.WorkingDirectory); err != nil {
			return nil, errors.Wrapf(err, "[dmlgen] Failed to chdir to %q", po.WorkingDirectory)
		}
		deferred = func() {
			_ = os.Chdir(oldWD)
		}
	}
	return deferred, nil
}

// GenerateProto searches all *.proto files in the given path and calls protoc
// to generate the Go source code.
func GenerateProto(protoFilesPath string, po *ProtocOptions) error {
	restoreFn, err := po.chdir()
	if err != nil {
		return errors.WithStack(err)
	}
	defer restoreFn()

	protoFilesPath = filepath.Clean(protoFilesPath)
	if ps := string(os.PathSeparator); !strings.HasSuffix(protoFilesPath, ps) {
		protoFilesPath += ps
	}

	protoFiles, err := filepath.Glob(protoFilesPath + "*.proto")
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] Can't access proto files in path %q", protoFilesPath)
	}

	cmd := exec.Command("protoc", append(po.toArgs(), protoFiles...)...)
	cmdStr := fmt.Sprintf("\ncd %s && %s\n\n", po.WorkingDirectory, cmd)
	if isDebug() {
		if po.WorkingDirectory == "" {
			po.WorkingDirectory = "."
		}
		print(cmdStr)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] %s%s", out, cmdStr)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		text := scanner.Text()
		if !strings.Contains(text, "WARNING") {
			return errors.WriteFailed.Newf("[dmlgen] protoc Error: %s", text)
		}
	}

	// what a hack: find all *.pb.go files and remove `import null
	// "github.com/corestoreio/pkg/storage/null"` because no other way to get
	// rid of the unused import or reference that import somehow in the
	// generated file :-( Once there's a better solution, remove this code.
	pbGoFiles, err := filepath.Glob(protoFilesPath + "*.pb.*go")
	if err != nil {
		return errors.Wrapf(err, "[dmlgen] Can't access pb.go files in path %q", protoFilesPath)
	}

	removeImports := [][]byte{
		[]byte("import null \"github.com/corestoreio/pkg/storage/null\"\n"),
		[]byte("null \"github.com/corestoreio/pkg/storage/null\"\n"),
	}
	for _, file := range pbGoFiles {
		fContent, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, ri := range removeImports {
			fContent = bytes.Replace(fContent, ri, nil, -1)
		}

		var buf bytes.Buffer
		for _, bt := range po.BuildTags {
			fmt.Fprintf(&buf, "// +build %s\n", bt)
		}
		if buf.Len() > 0 {
			buf.WriteByte('\n')
			buf.Write(fContent)
			fContent = buf.Bytes()
		}

		if err := ioutil.WriteFile(file, fContent, 0o644); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// serializerCustomType switches the default type from function serializerType
// to the new type. For now supports only protobuf.
func (g *Generator) serializerCustomType(c *ddl.Column) []string {
	pt := g.toSerializerType(c, true)
	var buf []string
	if pt == "google.protobuf.Timestamp" {
		buf = append(buf, "(gogoproto.stdtime)=true")
	}
	if pt == "bytes" {
		return nil // bytes can be null
	}
	if c.IsNull() || strings.IndexByte(pt, '.') > 0 /*whenever it is a custom type like null. or google.proto.timestamp*/ {
		// Indeed nullable Go Types must be not-nullable in HasSerializer because we
		// have a non-pointer struct type which contains the field Valid.
		// HasSerializer treats nullable fields as pointer fields, but that is
		// ridiculous.
		buf = append(buf, "(gogoproto.nullable)=false")
	}
	return buf
}

// GenerateSerializer writes the protocol buffer specifications into `w` and its test
// sources into wTest, if there are any tests.
func (g *Generator) GenerateSerializer(wMain, wTest io.Writer) error {
	switch g.Serializer {
	case "protobuf":
		if err := g.generateProto(wMain); err != nil {
			return errors.WithStack(err)
		}
	case "fbs":

	case "", "default", "none":
		return nil // do nothing
	default:
		return errors.NotAcceptable.Newf("[dmlgen] Serializer %q not supported.", g.Serializer)
	}

	return nil
}

func (g *Generator) generateProto(w io.Writer) error {
	pPkg := g.PackageSerializer
	if pPkg == "" {
		pPkg = g.Package
	}

	proto := codegen.NewProto(pPkg)
	proto.Pln(`import "github.com/gogo/protobuf/gogoproto/gogo.proto";`)

	const importTimeStamp = `import "google/protobuf/timestamp.proto";`
	proto.Pln(importTimeStamp)
	proto.Pln(`import "github.com/corestoreio/pkg/storage/null/null.proto";`)
	var hasGoPackageOption bool
	for _, o := range g.SerializerHeaderOptions {
		proto.Pln(`option ` + o + `;`)
		if !hasGoPackageOption {
			hasGoPackageOption = strings.Contains(o, "go_package")
		}
	}
	if !hasGoPackageOption {
		proto.Pln(`option go_package = `, fmt.Sprintf("%q;", pPkg))
	}

	var hasTimestampField bool
	for _, tblname := range g.sortedTableNames() {
		t := g.Tables[tblname] // must panic if table name not found

		fieldMapFn := g.defaultTableConfig.FieldMapFn
		if fieldMapFn == nil {
			fieldMapFn = t.fieldMapFn
		}
		if fieldMapFn == nil {
			fieldMapFn = defaultFieldMapFn
		}

		proto.C(t.EntityName(), `represents a single row for`, t.Table.Name, `DB table. Auto generated.`)
		if t.Table.TableComment != "" {
			proto.C("Table comment:", t.Table.TableComment)
		}
		proto.Pln(`message`, t.EntityName(), `{`)
		{
			proto.In()
			var lastColumnPos uint64
			t.Table.Columns.Each(func(c *ddl.Column) {
				if t.IsFieldPublic(c.Field) {
					serType := g.serializerType(c)
					if !hasTimestampField && strings.HasPrefix(serType, "google.protobuf.Timestamp") {
						hasTimestampField = true
					}
					var optionConcret string
					if options := g.serializerCustomType(c); len(options) > 0 {
						optionConcret = `[` + strings.Join(options, ",") + `]`
					}
					// extend here with a custom code option, if someone needs
					proto.Pln(serType, strs.ToGoCamelCase(c.Field), `=`, c.Pos, optionConcret+`;`)
					lastColumnPos = c.Pos
				}
			})
			lastColumnPos++

			if g.hasFeature(t.featuresInclude, t.featuresExclude, FeatureEntityRelationships) {

				// for debugging see Table.fnEntityStruct function. This code is only different in the Pln function.

				var hasAtLeastOneRelationShip int
				relationShipSeen := map[string]bool{}
				if kcuc, ok := g.kcu[t.Table.Name]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
					for _, kcuce := range kcuc.Data {
						if !kcuce.ReferencedTableName.Valid {
							continue
						}
						hasAtLeastOneRelationShip++
						// case ONE-TO-MANY
						isOneToMany := g.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						isRelationAllowed := g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						hasTable := g.Tables[kcuce.ReferencedTableName.Data] != nil
						if isOneToMany && hasTable && isRelationAllowed {
							proto.Pln(fieldMapFn(pluralize(kcuce.ReferencedTableName.Data)), fieldMapFn(pluralize(kcuce.ReferencedTableName.Data)),
								"=", lastColumnPos, ";",
								"// 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
							lastColumnPos++
						}

						// case ONE-TO-ONE
						isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						if isOneToOne && hasTable && isRelationAllowed {
							proto.Pln(fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)), fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)),
								"=", lastColumnPos, ";",
								"// 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
							lastColumnPos++
						}

						// case MANY-TO-MANY
						targetTbl, targetColumn := g.krs.ManyToManyTarget(kcuce.TableName, kcuce.ColumnName)
						// hasTable variable shall not be added because usually the link table does not get loaded.
						if isRelationAllowed && targetTbl != "" && targetColumn != "" {
							proto.Pln(fieldMapFn(pluralize(targetTbl)), " *", pluralize(targetTbl),
								t.customStructTagFields[targetTbl],
								"// M:N", kcuce.TableName+"."+kcuce.ColumnName, "via", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data,
								"=>", targetTbl+"."+targetColumn,
							)
						}
					}
				}

				if kcuc, ok := g.kcuRev[t.Table.Name]; ok { // kcu = keyColumnUsage && kcuc = keyColumnUsageCollection
					for _, kcuce := range kcuc.Data {
						if !kcuce.ReferencedTableName.Valid {
							continue
						}
						hasAtLeastOneRelationShip++
						// case ONE-TO-MANY
						isOneToMany := g.krs.IsOneToMany(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						isRelationAllowed := g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						hasTable := g.Tables[kcuce.ReferencedTableName.Data] != nil
						keySeen := fieldMapFn(pluralize(kcuce.ReferencedTableName.Data))
						relationShipSeenAlready := relationShipSeen[keySeen]
						// case ONE-TO-MANY
						if isRelationAllowed && isOneToMany && hasTable && !relationShipSeenAlready {
							proto.Pln(fieldMapFn(pluralize(kcuce.ReferencedTableName.Data)), fieldMapFn(pluralize(kcuce.ReferencedTableName.Data)),
								"=", lastColumnPos, ";",
								"// Reversed 1:M", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
							relationShipSeen[keySeen] = true
							lastColumnPos++
						}

						// case ONE-TO-ONE
						isOneToOne := g.krs.IsOneToOne(kcuce.TableName, kcuce.ColumnName, kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						if isRelationAllowed && isOneToOne && hasTable {
							proto.Pln(fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)), fieldMapFn(strs.ToGoCamelCase(kcuce.ReferencedTableName.Data)),
								"=", lastColumnPos, ";",
								"// Reversed 1:1", kcuce.TableName+"."+kcuce.ColumnName, "=>", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data)
							lastColumnPos++
						}

						// case MANY-TO-MANY
						targetTbl, targetColumn := g.krs.ManyToManyTarget(kcuce.ReferencedTableName.Data, kcuce.ReferencedColumnName.Data)
						if targetTbl != "" && targetColumn != "" {
							keySeen := fieldMapFn(pluralize(targetTbl))
							isRelationAllowed = g.isAllowedRelationship(kcuce.TableName, kcuce.ColumnName, targetTbl, targetColumn) &&
								!relationShipSeen[keySeen]
							relationShipSeen[keySeen] = true
						}

						// case MANY-TO-MANY
						// hasTable shall not be added because usually the link table does not get loaded.
						if isRelationAllowed && targetTbl != "" && targetColumn != "" {
							proto.Pln(fieldMapFn(pluralize(targetTbl)), fieldMapFn(pluralize(targetTbl)),
								"=", lastColumnPos, ";",
								"// Reversed M:N", kcuce.TableName+"."+kcuce.ColumnName, "via", kcuce.ReferencedTableName.Data+"."+kcuce.ReferencedColumnName.Data,
								"=>", targetTbl+"."+targetColumn,
							)
							lastColumnPos++
						}
					}
				}
			}
			proto.Out()
		}
		proto.Pln(`}`)

		proto.C(t.CollectionName(), `represents multiple rows for the`, t.Table.Name, `DB table. Auto generated.`)
		proto.Pln(`message`, t.CollectionName(), `{`)
		{
			proto.In()
			proto.Pln(`repeated`, t.EntityName(), `Data = 1;`)
			proto.Out()
		}
		proto.Pln(`}`)
	}

	if !hasTimestampField {
		// bit hacky to remove the import of timestamp proto but for now OK.
		removedImport := strings.ReplaceAll(proto.String(), importTimeStamp, "")
		proto.Reset()
		proto.WriteString(removedImport)
	}
	return proto.GenerateFile(w)
}
