// +build ignore

package vault

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
	// PaymentVaultVaultPayment => Vault Provider.
	// Specified provider should be enabled.
	// Path: payment/vault/vault_payment
	// SourceModel: Otnegam\Vault\Model\Adminhtml\Source\VaultProvidersMap
	PaymentVaultVaultPayment model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.PaymentVaultVaultPayment = model.NewStr(`payment/vault/vault_payment`, model.WithConfigStructure(cfgStruct))

	return pp
}
