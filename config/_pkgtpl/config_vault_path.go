// +build ignore

package vault

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
	// PaymentVaultVaultPayment => Vault Provider.
	// Specified provider should be enabled.
	// Path: payment/vault/vault_payment
	// SourceModel: Otnegam\Vault\Model\Adminhtml\Source\VaultProvidersMap
	PaymentVaultVaultPayment model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.PaymentVaultVaultPayment = model.NewStr(`payment/vault/vault_payment`, model.WithPkgCfg(pkgCfg))

	return pp
}
