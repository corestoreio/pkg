package dmlgen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alecthomas/repr"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/util/codegen"
	"github.com/corestoreio/pkg/util/strs"
	"golang.org/x/tools/go/packages"
)

// ProtocOptions allows to modify the protoc CLI command.
type ProtocOptions struct {
	WorkingDirectory string
	ProtoGen         string // default go, options: gofast, gogo, gogofast, gogofaster and other installed proto generators

	ProtoPath []string // where to find other *.proto files; if empty the usual defaults apply
	GoOutPath string   // where to write the generated proto files; default "."
	GoOpt     []string // each entry generates a "--go_opt" argument
	// GoSourceRelative if the paths=source_relative flag is specified, the output
	// file is placed in the same relative directory as the input file. For
	// example, an input file protos/buzz.proto results in an output file at
	// protos/buzz.pb.go.
	GoSourceRelative bool

	GRPC               bool
	GRPCOpt            []string // each entry generates a "--go-grpc_opt" argument
	GRPCGatewayOutMap  []string // GRPC must be enabled in the above field
	GRPCGatewayOutPath string   // GRPC must be enabled in the above field

	SwaggerOutPath string
	CustomArgs     []string
	// TODO add validation plugin, either
	//  https://github.com/mwitkow/go-proto-validators as used in github.com/go_go/grpc-example/proto/example.proto
	//  This github.com/mwitkow/go-proto-validators seems dead.
	//  or https://github.com/envoyproxy/protoc-gen-validate
	//  Requirement: error messages must be translatable and maybe use an errors.Kind type
}

var defaultProtoPaths = make([]string, 0, 8)

func init() {
	preDefinedPaths := [...]string{
		build.Default.GOPATH + "/src/",
		"vendor/github.com/grpc-ecosystem/grpc-gateway/",
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
	if po.GRPC {
		if po.GRPCGatewayOutMap == nil {
			po.GRPCGatewayOutMap = []string{
				"allow_patch_feature=false",
			}
		}
		if po.GRPCGatewayOutPath == "" {
			po.GRPCGatewayOutPath = "."
		}
	}
	if po.GoOutPath == "" {
		po.GoOutPath = "."
	} else {
		if err := os.MkdirAll(filepath.Clean(po.GoOutPath), 0751); err != nil {
			panic(err)
		}
	}
	if po.ProtoPath == nil {
		po.ProtoPath = append(po.ProtoPath, defaultProtoPaths...)
	}
	if po.ProtoGen == "" {
		po.ProtoGen = "go"
	}

	args := []string{
		"--" + po.ProtoGen + "_out=" + po.GoOutPath,
		"--proto_path", strings.Join(po.ProtoPath, ":"),
	}
	if po.GoSourceRelative {
		args = append(args, "--go_opt=paths=source_relative")
	}
	for _, o := range po.GoOpt {
		args = append(args, "--go_opt="+o)
	}
	if po.GRPC {
		args = append(args, "--go-grpc_out="+po.GoOutPath)
		if po.GoSourceRelative {
			args = append(args, "--go-grpc_opt=paths=source_relative")
		}
		for _, o := range po.GRPCOpt {
			args = append(args, "--go-grpc_opt="+o)
		}
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

// RunProtoc searches all *.proto files in the given path and calls protoc
// to generate the Go source code.
func RunProtoc(protoFilesPath string, po *ProtocOptions) error {
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

		if err := ioutil.WriteFile(file, fContent, 0o644); err != nil {
			return errors.WithStack(err)
		}
	}

	// build a mapper between DB and proto :-(

	return nil
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
		panic("not yet supported")
	case "", "default", "none":
		return nil // do nothing
	default:
		return errors.NotAcceptable.Newf("[dmlgen] Serializer %q not supported.", g.Serializer)
	}

	return nil
}

func (g *Generator) generateProto(w io.Writer) error {
	proto := codegen.NewProto(g.Package)

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
		proto.Pln(`option go_package = `, fmt.Sprintf("%q;", g.PackageSerializer))
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

					// extend here with a custom code option, if someone needs
					proto.Pln(serType, strs.ToGoCamelCase(c.Field), `=`, c.Pos, `; //`, c.Comment)
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

func buildDMLProtoMapper(dmlGoFilesFullImportPath, protoGoFilesFullImportPath string) error {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, dmlGoFilesFullImportPath, protoGoFilesFullImportPath)
	if err != nil {
		return fmt.Errorf("failed to packages.Load: %w", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		panic("errors occurred")
	}

	dmlgenTypes := map[string][]fieldTypeInfo{}
	protoTypes := map[string][]fieldTypeInfo{}
	var protoPackageName string
	var dmlgenPackageName string
	// Print the names of the source files
	// for each package listed on the command line.
	for _, pkg := range pkgs {

		repr.Println(pkg.PkgPath, pkg.ID, pkg.GoFiles)
		for _, sntx := range pkg.Syntax {

			structNameAndType := buildDMLProtoASTDecl(sntx.Decls)
			switch {
			case dmlGoFilesFullImportPath == pkg.ID:
				dmlgenPackageName = sntx.Name.Name
				for n, fti := range structNameAndType {
					dmlgenTypes[n] = append(dmlgenTypes[n], fti...)
				}
			case protoGoFilesFullImportPath == pkg.ID:
				protoPackageName = sntx.Name.Name
				for n, fti := range structNameAndType {
					protoTypes[n] = append(protoTypes[n], fti...)
				}
			}
		}
	}

	if protoPackageName == "" {
		return fmt.Errorf("protoPackageName cannot be empty or no files found to parse")
	}
	repr.Println(dmlgenTypes)
	repr.Println(protoTypes)

	// <generate converter from proto to dml>
	cg := codegen.NewGo(protoPackageName)
	cg.AddImports(dmlGoFilesFullImportPath)
	for protoStructName, protoFieldTypeInfos := range protoTypes {

		cg.Pln(`func (x *`, protoStructName, `) ToDBType(optional *`, dmlgenPackageName, `.`, protoStructName, `) *`, dmlgenPackageName, `.`, protoStructName, `{`)
		cg.In()
		cg.Pln(`if optional == nil { optional = new(`, dmlgenPackageName, `.`, protoStructName, `) }`)
		cg.In()

		dmlgenFieldTypeInfos, ok := dmlgenTypes[protoStructName]
		if !ok {
			return fmt.Errorf("proto struct %q not found in generated dml package %q", protoStructName, dmlgenPackageName)
		}

		for idx, pft := range protoFieldTypeInfos {
			dmlft := dmlgenFieldTypeInfos[idx]

			if pft.isPointer && pft.externalPkgName == "null" {
				cg.Pln(`optional.`, pft.fname, ` .Reset()`)
			}

			getter := codegen.SkipWS(`Get`, pft.fname, `()`)
			switch {
			case pft.isSlice:
				cg.Pln(`optional.`, dmlft.fname, ` = optional.`, dmlft.fname, `[:0]`)
				cg.Pln(`for _, d := range x.`, pft.fname, ` {`)
				{
					cg.Pln(`optional.`, dmlft.fname, ` = append(optional.`, pft.fname, `, d.ToDBType(nil))`)
				}
				cg.Pln(`}`)

			case pft.isPointer && pft.isStruct: // proto type is a pointer

				switch {
				case strings.HasSuffix(dmlft.ftype, "time.Time") && strings.HasSuffix(pft.ftype, "timestamppb.Timestamp"):
					cg.Pln(`optional.`, dmlft.fname, ` = x.`, getter, `.AsTime() `)

				case strings.HasSuffix(dmlft.ftype, "null.Time") && strings.HasSuffix(pft.ftype, "timestamppb.Timestamp"):
					cg.Pln(`optional.`, dmlft.fname, `.SetProto( x.`, getter, `)`)

				case strings.HasSuffix(dmlft.ftype, "null.Decimal") && strings.HasSuffix(pft.ftype, "null.Decimal"):
					cg.Pln(`optional.`, dmlft.fname, `.SetPtr( x.`, getter, `)`)

				default:
					cg.Pln(`optional.`, dmlft.fname, ` = *x.`, pft.fname, `// TODO BUG fix`, pft.ftype, "isStruct", pft.isStruct)
				}

			case pft.isPointer:

				if dmlgenFieldTypeInfos[idx].externalPkgName == "null" {
					cg.Pln(`optional.`, dmlft.fname, `.SetPtr( x.`, pft.fname, `)`, `//1`, pft.ftype, "isStruct", pft.isStruct)
				} else {
					cg.Pln(`optional.`, dmlft.fname, ` =   x.`, pft.fname, `//2`, dmlft.ftype, "=>", pft.ftype, "isStruct", pft.isStruct, pft.isSlice, dmlft.isSlice)
				}
			default:
				cg.Pln(`optional.`, dmlft.fname, ` = x.`, pft.fname, `//3`, pft.ftype, "isStruct", pft.isStruct)
			}
		}
		cg.Out()
		cg.Pln(`return optional`) // end &type{
		cg.Out()
		cg.Pln(`}`) // end func
	}

	mapperFileName := filepath.Join(build.Default.GOPATH, "src", protoGoFilesFullImportPath, "/mapper.go")
	f, err := os.Create(mapperFileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	return cg.GenerateFile(f)
	// </generate converter from proto to dml>

	return nil
}

type fieldTypeInfo struct {
	fname           string
	ftype           string
	externalPkgName string
	isPointer       bool
	isStruct        bool // if false, then a primitive and true then a struct type like *timestamppb.Timestamp
	isSlice         bool
}

func buildDMLProtoASTDecl(decls []ast.Decl) map[string][]fieldTypeInfo { // map[structName] []Fields
	ret := map[string][]fieldTypeInfo{}
	for _, node := range decls {
		switch node.(type) {
		case *ast.GenDecl:
			genDecl := node.(*ast.GenDecl)
		genDeclSpecsLOOP:
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)
					if !typeSpec.Name.IsExported() {
						continue genDeclSpecsLOOP
					}

					switch tst := typeSpec.Type.(type) {
					case *ast.StructType:

						fieldTypeInfos := []fieldTypeInfo{}

						for _, field := range tst.Fields.List {
							fieldType := fieldToType(field)
							for _, name := range field.Names {
								if isExported(name.Name) {
									fieldType.fname = name.Name
									fieldTypeInfos = append(fieldTypeInfos, fieldType)
								}
							} // end for
						} // end for
						ret[typeSpec.Name.Name] = fieldTypeInfos
					}
				}
			}
		}
	}
	return ret
}

// fieldToType returns the type name and whether if it's exported.
func fieldToType(f *ast.Field) fieldTypeInfo {
	switch arg := f.Type.(type) {
	case *ast.ArrayType:
		n := astNodeName(arg.Elt)

		_, isSlice := arg.Elt.(*ast.StarExpr) // special custom slice, not []byte or anything else.

		return fieldTypeInfo{ftype: "[]" + n, isSlice: isSlice}
	case *ast.Ellipsis:
		n := astNodeName(arg.Elt)
		return fieldTypeInfo{ftype: n}
	case *ast.FuncType:
		// Do not print the function signature to not overload the trace.
		return fieldTypeInfo{ftype: "func"}
	case *ast.Ident:
		return fieldTypeInfo{ftype: arg.Name}
	case *ast.InterfaceType:
		return fieldTypeInfo{ftype: "any"}
	case *ast.SelectorExpr:

		if ident, ok := arg.X.(*ast.Ident); ok && ident.Name != "" {
			return fieldTypeInfo{ftype: ident.Name + "." + arg.Sel.Name, externalPkgName: ident.Name, isStruct: true}
		}
		return fieldTypeInfo{ftype: arg.Sel.Name}
	case *ast.StarExpr:
		if sel, ok := arg.X.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name != "" {
				n := astNodeName(arg.X)
				return fieldTypeInfo{ftype: "*" + ident.Name + "." + n, externalPkgName: ident.Name, isPointer: true, isStruct: true}
			}
		}
		n := astNodeName(arg.X)
		return fieldTypeInfo{ftype: "*" + n, isPointer: true}
	case *ast.MapType:
		return fieldTypeInfo{ftype: fmt.Sprintf("map[%s]%s", astNodeName(arg.Key), astNodeName(arg.Value))}
	case *ast.ChanType:
		return fieldTypeInfo{ftype: fmt.Sprintf("chan %s", astNodeName(arg.Value))}
	default:
		return fieldTypeInfo{ftype: "<unknown>"}
	}
}

func isExported(name string) bool {
	switch name {
	case "bool", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64",
		"float32", "float64", "complex64", "complex128",
		"string", "int", "uint", "uintptr", "byte", "rune":
		return true
	}
	return token.IsExported(name)
}

func astNodeName(n ast.Node) string {
	switch t := n.(type) {
	case *ast.InterfaceType:
		return "any"
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	case *ast.StarExpr:
		return "*" + astNodeName(t.X)
	default:
		return "<unknown>"
	}
}
