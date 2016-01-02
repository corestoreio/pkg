// +build ignore

package catalogsearch

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
	// CatalogSeoSearchTerms => Popular Search Terms.
	// Path: catalog/seo/search_terms
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	CatalogSeoSearchTerms model.Bool

	// CatalogSearchEngine => .
	// Path: catalog/search/engine
	// BackendModel: Otnegam\CatalogSearch\Model\Adminhtml\System\Config\Backend\Engine
	CatalogSearchEngine model.Str

	// CatalogSearchMinQueryLength => Minimal Query Length.
	// Path: catalog/search/min_query_length
	CatalogSearchMinQueryLength model.Str

	// CatalogSearchMaxQueryLength => Maximum Query Length.
	// Path: catalog/search/max_query_length
	CatalogSearchMaxQueryLength model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogSeoSearchTerms = model.NewBool(`catalog/seo/search_terms`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSearchEngine = model.NewStr(`catalog/search/engine`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSearchMinQueryLength = model.NewStr(`catalog/search/min_query_length`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSearchMaxQueryLength = model.NewStr(`catalog/search/max_query_length`, model.WithPkgCfg(pkgCfg))

	return pp
}
