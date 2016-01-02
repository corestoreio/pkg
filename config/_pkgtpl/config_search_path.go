// +build ignore

package search

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with PackageConfiguration.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// CatalogSearchEngine => Search Engine.
	// Path: catalog/search/engine
	// SourceModel: Otnegam\Search\Model\Adminhtml\System\Config\Source\Engine
	CatalogSearchEngine model.Str

	// CatalogSearchSearchType => .
	// Path: catalog/search/search_type
	CatalogSearchSearchType model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogSearchEngine = model.NewStr(`catalog/search/engine`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSearchSearchType = model.NewStr(`catalog/search/search_type`, model.WithPkgCfg(pkgCfg))

	return pp
}
