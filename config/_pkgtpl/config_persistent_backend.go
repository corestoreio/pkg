// +build ignore

package persistent

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
	// PersistentOptionsEnabled => Enable Persistence.
	// Path: persistent/options/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PersistentOptionsEnabled model.Bool

	// PersistentOptionsLifetime => Persistence Lifetime (seconds).
	// Path: persistent/options/lifetime
	PersistentOptionsLifetime model.Str

	// PersistentOptionsRememberEnabled => Enable "Remember Me".
	// Path: persistent/options/remember_enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PersistentOptionsRememberEnabled model.Bool

	// PersistentOptionsRememberDefault => "Remember Me" Default Value.
	// Path: persistent/options/remember_default
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PersistentOptionsRememberDefault model.Bool

	// PersistentOptionsLogoutClear => Clear Persistence on Sign Out.
	// Path: persistent/options/logout_clear
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PersistentOptionsLogoutClear model.Bool

	// PersistentOptionsShoppingCart => Persist Shopping Cart.
	// Path: persistent/options/shopping_cart
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PersistentOptionsShoppingCart model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.PersistentOptionsEnabled = model.NewBool(`persistent/options/enabled`, model.WithConfigStructure(cfgStruct))
	pp.PersistentOptionsLifetime = model.NewStr(`persistent/options/lifetime`, model.WithConfigStructure(cfgStruct))
	pp.PersistentOptionsRememberEnabled = model.NewBool(`persistent/options/remember_enabled`, model.WithConfigStructure(cfgStruct))
	pp.PersistentOptionsRememberDefault = model.NewBool(`persistent/options/remember_default`, model.WithConfigStructure(cfgStruct))
	pp.PersistentOptionsLogoutClear = model.NewBool(`persistent/options/logout_clear`, model.WithConfigStructure(cfgStruct))
	pp.PersistentOptionsShoppingCart = model.NewBool(`persistent/options/shopping_cart`, model.WithConfigStructure(cfgStruct))

	return pp
}
