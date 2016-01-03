// +build ignore

package search

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
	// CatalogSearchEngine => Search Engine.
	// Path: catalog/search/engine
	// SourceModel: Otnegam\Search\Model\Adminhtml\System\Config\Source\Engine
	CatalogSearchEngine model.Str

	// CatalogSearchSearchType => .
	// Path: catalog/search/search_type
	CatalogSearchSearchType model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogSearchEngine = model.NewStr(`catalog/search/engine`, model.WithConfigStructure(cfgStruct))
	pp.CatalogSearchSearchType = model.NewStr(`catalog/search/search_type`, model.WithConfigStructure(cfgStruct))

	return pp
}
