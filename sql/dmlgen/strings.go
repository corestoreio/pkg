package dmlgen

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/strs"
)

func pluralize(name string) string {
	tg := strs.ToGoCamelCase(name)
	switch {
	case strings.HasSuffix(name, "y"):
		return tg[:len(tg)-1] + "ies"
	case strings.HasSuffix(name, "ch"):
		return tg + "es"
	case strings.HasSuffix(name, "x"):
		return tg + "es"
	case strings.HasSuffix(name, "us"):
		return tg + "i" // status -> stati
	case strings.HasSuffix(name, "um"):
		return tg + "en" // datum -> daten
	case strings.HasSuffix(name, "s"):
		return tg + "Collection" // stupid case, better ideas?
	default:
		return tg + "s"
	}
}

// findUsedPackages checks for needed packages which we must import.
func findUsedPackages(file []byte, predefinedImportPaths []string) ([]string, error) {
	af, err := goparser.ParseFile(token.NewFileSet(), "cs_virtual_file.go", append([]byte("package temporarily_main\n\n"), file...), 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	idents := map[string]struct{}{}
	ast.Inspect(af, func(n ast.Node) bool {
		if nt, ok := n.(*ast.Ident); ok {
			idents[nt.Name] = struct{}{} // will contain too much info
			// we only need to know: pkg.TYPE
		}
		return true
	})

	ret := make([]string, 0, len(predefinedImportPaths))
	for _, path := range predefinedImportPaths {
		_, pkg := filepath.Split(path)
		if _, ok := idents[pkg]; ok {
			ret = append(ret, path)
		}
	}
	return ret, nil
}
