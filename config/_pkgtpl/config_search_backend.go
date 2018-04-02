// +build ignore

package search

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
	// CatalogSearchEngine => Search Engine.
	// Path: catalog/search/engine
	// SourceModel: Magento\Search\Model\Adminhtml\System\Config\Source\Engine
	CatalogSearchEngine cfgmodel.Str

	// CatalogSearchSearchType => .
	// Path: catalog/search/search_type
	CatalogSearchSearchType cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogSearchEngine = cfgmodel.NewStr(`catalog/search/engine`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSearchSearchType = cfgmodel.NewStr(`catalog/search/search_type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
