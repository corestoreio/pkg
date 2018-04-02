// +build ignore

package catalogsearch

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CatalogSeoSearchTerms => Popular Search Terms.
	// Path: catalog/seo/search_terms
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	CatalogSeoSearchTerms cfgmodel.Bool

	// CatalogSearchEngine => .
	// Path: catalog/search/engine
	// BackendModel: Magento\CatalogSearch\Model\Adminhtml\System\Config\Backend\Engine
	CatalogSearchEngine cfgmodel.Str

	// CatalogSearchMinQueryLength => Minimal Query Length.
	// Path: catalog/search/min_query_length
	CatalogSearchMinQueryLength cfgmodel.Str

	// CatalogSearchMaxQueryLength => Maximum Query Length.
	// Path: catalog/search/max_query_length
	CatalogSearchMaxQueryLength cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogSeoSearchTerms = cfgmodel.NewBool(`catalog/seo/search_terms`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSearchEngine = cfgmodel.NewStr(`catalog/search/engine`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSearchMinQueryLength = cfgmodel.NewStr(`catalog/search/min_query_length`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSearchMaxQueryLength = cfgmodel.NewStr(`catalog/search/max_query_length`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
