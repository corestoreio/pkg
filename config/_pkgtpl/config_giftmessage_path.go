// +build ignore

package giftmessage

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
	// SalesGiftOptionsAllowOrder => Allow Gift Messages on Order Level.
	// Path: sales/gift_options/allow_order
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesGiftOptionsAllowOrder model.Bool

	// SalesGiftOptionsAllowItems => Allow Gift Messages for Order Items.
	// Path: sales/gift_options/allow_items
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SalesGiftOptionsAllowItems model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.SalesGiftOptionsAllowOrder = model.NewBool(`sales/gift_options/allow_order`, model.WithPkgCfg(pkgCfg))
	pp.SalesGiftOptionsAllowItems = model.NewBool(`sales/gift_options/allow_items`, model.WithPkgCfg(pkgCfg))

	return pp
}
