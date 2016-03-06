// +build ignore

package review

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
	// CatalogReviewAllowGuest => Allow Guests to Write Reviews.
	// Path: catalog/review/allow_guest
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogReviewAllowGuest model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogReviewAllowGuest = model.NewBool(`catalog/review/allow_guest`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
