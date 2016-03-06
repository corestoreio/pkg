// +build ignore

package catalogsearch

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// CatalogSeoSearchTerms => Popular Search Terms.
	// Path: catalog/seo/search_terms
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	CatalogSeoSearchTerms model.Bool

	// CatalogSearchEngine => .
	// Path: catalog/search/engine
	// BackendModel: Magento\CatalogSearch\Model\Adminhtml\System\Config\Backend\Engine
	CatalogSearchEngine model.Str

	// CatalogSearchMinQueryLength => Minimal Query Length.
	// Path: catalog/search/min_query_length
	CatalogSearchMinQueryLength model.Str

	// CatalogSearchMaxQueryLength => Maximum Query Length.
	// Path: catalog/search/max_query_length
	CatalogSearchMaxQueryLength model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogSeoSearchTerms = model.NewBool(`catalog/seo/search_terms`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSearchEngine = model.NewStr(`catalog/search/engine`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSearchMinQueryLength = model.NewStr(`catalog/search/min_query_length`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSearchMaxQueryLength = model.NewStr(`catalog/search/max_query_length`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
