// +build ignore

package giftmessage

import (
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// SalesGiftOptionsAllowOrder => Allow Gift Messages on Order Level.
	// Path: sales/gift_options/allow_order
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesGiftOptionsAllowOrder cfgmodel.Bool

	// SalesGiftOptionsAllowItems => Allow Gift Messages for Order Items.
	// Path: sales/gift_options/allow_items
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SalesGiftOptionsAllowItems cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SalesGiftOptionsAllowOrder = cfgmodel.NewBool(`sales/gift_options/allow_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SalesGiftOptionsAllowItems = cfgmodel.NewBool(`sales/gift_options/allow_items`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
