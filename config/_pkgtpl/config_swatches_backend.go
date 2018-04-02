// +build ignore

package swatches

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
	// CatalogFrontendSwatchesPerProduct => Swatches per Product.
	// Path: catalog/frontend/swatches_per_product
	CatalogFrontendSwatchesPerProduct cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogFrontendSwatchesPerProduct = cfgmodel.NewStr(`catalog/frontend/swatches_per_product`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
