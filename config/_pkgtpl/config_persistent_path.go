// +build ignore

package persistent

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathPersistentOptionsEnabled => Enable Persistence.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPersistentOptionsEnabled = model.NewBool(`persistent/options/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathPersistentOptionsLifetime => Persistence Lifetime (seconds).
var PathPersistentOptionsLifetime = model.NewStr(`persistent/options/lifetime`, model.WithPkgCfg(PackageConfiguration))

// PathPersistentOptionsRememberEnabled => Enable "Remember Me".
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPersistentOptionsRememberEnabled = model.NewBool(`persistent/options/remember_enabled`, model.WithPkgCfg(PackageConfiguration))

// PathPersistentOptionsRememberDefault => "Remember Me" Default Value.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPersistentOptionsRememberDefault = model.NewBool(`persistent/options/remember_default`, model.WithPkgCfg(PackageConfiguration))

// PathPersistentOptionsLogoutClear => Clear Persistence on Sign Out.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPersistentOptionsLogoutClear = model.NewBool(`persistent/options/logout_clear`, model.WithPkgCfg(PackageConfiguration))

// PathPersistentOptionsShoppingCart => Persist Shopping Cart.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPersistentOptionsShoppingCart = model.NewBool(`persistent/options/shopping_cart`, model.WithPkgCfg(PackageConfiguration))
