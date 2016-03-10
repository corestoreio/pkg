// +build ignore

package persistent

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
	// PersistentOptionsEnabled => Enable Persistence.
	// Path: persistent/options/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PersistentOptionsEnabled cfgmodel.Bool

	// PersistentOptionsLifetime => Persistence Lifetime (seconds).
	// Path: persistent/options/lifetime
	PersistentOptionsLifetime cfgmodel.Str

	// PersistentOptionsRememberEnabled => Enable "Remember Me".
	// Path: persistent/options/remember_enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PersistentOptionsRememberEnabled cfgmodel.Bool

	// PersistentOptionsRememberDefault => "Remember Me" Default Value.
	// Path: persistent/options/remember_default
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PersistentOptionsRememberDefault cfgmodel.Bool

	// PersistentOptionsLogoutClear => Clear Persistence on Sign Out.
	// Path: persistent/options/logout_clear
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PersistentOptionsLogoutClear cfgmodel.Bool

	// PersistentOptionsShoppingCart => Persist Shopping Cart.
	// Path: persistent/options/shopping_cart
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PersistentOptionsShoppingCart cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.PersistentOptionsEnabled = cfgmodel.NewBool(`persistent/options/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PersistentOptionsLifetime = cfgmodel.NewStr(`persistent/options/lifetime`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PersistentOptionsRememberEnabled = cfgmodel.NewBool(`persistent/options/remember_enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PersistentOptionsRememberDefault = cfgmodel.NewBool(`persistent/options/remember_default`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PersistentOptionsLogoutClear = cfgmodel.NewBool(`persistent/options/logout_clear`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PersistentOptionsShoppingCart = cfgmodel.NewBool(`persistent/options/shopping_cart`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
