// +build ignore

package swatches

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// CatalogFrontendSwatchesPerProduct => Swatches per Product.
	// Path: catalog/frontend/swatches_per_product
	CatalogFrontendSwatchesPerProduct model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogFrontendSwatchesPerProduct = model.NewStr(`catalog/frontend/swatches_per_product`, model.WithConfigStructure(cfgStruct))

	return pp
}
