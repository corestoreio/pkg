package dmlgen

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/mailru/easyjson/bootstrap"
	"github.com/mailru/easyjson/parser"
)

// GenerateJSON creates the easysjon code for a specific file or a whole
// directory. argument `g` can be nil.
// deprecated as there are stdlib compatible packages which have the same speed without code gen.
func GenerateJSON(fileNameOrDirectory, buildTags string, g *bootstrap.Generator) error {
	fInfo, err := os.Stat(fileNameOrDirectory)
	if err != nil {
		return errors.WithStack(err)
	}

	p := new(parser.Parser)
	if err := p.Parse(fileNameOrDirectory, fInfo.IsDir()); err != nil {
		return errors.CorruptData.Newf("[dmlgen] Error parsing failed %q: %v", fileNameOrDirectory, err)
	}

	var outName string
	if fInfo.IsDir() {
		outName = filepath.Join(fileNameOrDirectory, p.PkgName+"_easyjson.go")
	} else {
		if s := strings.TrimSuffix(fileNameOrDirectory, ".go"); s == fileNameOrDirectory {
			return errors.NotAcceptable.Newf("[dmlgen] GenerateJSON: Filename must end in '.go'")
		} else {
			outName = s + "_easyjson.go"
		}
	}

	if len(p.StructNames) == 0 {
		return errors.NotFound.Newf("[dmlgen] Can't find any StructNames in the Go files of %q", fileNameOrDirectory)
	}

	if g == nil {
		g = &bootstrap.Generator{
			BuildTags:             buildTags,
			PkgPath:               p.PkgPath,
			PkgName:               p.PkgName,
			Types:                 p.StructNames,
			SnakeCase:             true,
			LowerCamelCase:        true,
			NoStdMarshalers:       false,
			DisallowUnknownFields: false,
			OmitEmpty:             true,
			LeaveTemps:            false,
			OutName:               outName,
			StubsOnly:             false,
			NoFormat:              false,
		}
	} else {
		g.Types = p.StructNames
	}
	if err := g.Run(); err != nil {
		return errors.Fatal.Newf("[dmlgen] easyJSON: Bootstrap failed: %v", err)
	}
	return nil
}
