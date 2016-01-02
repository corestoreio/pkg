// +build ignore

package persistent

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

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.PersistentOptionsEnabled = model.NewBool(`persistent/options/enabled`, model.WithPkgCfg(pkgCfg))
	pp.PersistentOptionsLifetime = model.NewStr(`persistent/options/lifetime`, model.WithPkgCfg(pkgCfg))
	pp.PersistentOptionsRememberEnabled = model.NewBool(`persistent/options/remember_enabled`, model.WithPkgCfg(pkgCfg))
	pp.PersistentOptionsRememberDefault = model.NewBool(`persistent/options/remember_default`, model.WithPkgCfg(pkgCfg))
	pp.PersistentOptionsLogoutClear = model.NewBool(`persistent/options/logout_clear`, model.WithPkgCfg(pkgCfg))
	pp.PersistentOptionsShoppingCart = model.NewBool(`persistent/options/shopping_cart`, model.WithPkgCfg(pkgCfg))

	return pp
}
