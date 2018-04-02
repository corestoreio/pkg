// +build ignore

package review

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
	// CatalogReviewAllowGuest => Allow Guests to Write Reviews.
	// Path: catalog/review/allow_guest
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogReviewAllowGuest cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogReviewAllowGuest = cfgmodel.NewBool(`catalog/review/allow_guest`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
