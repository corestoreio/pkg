// +build ignore

package vault

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
	// PaymentVaultVaultPayment => Vault Provider.
	// Specified provider should be enabled.
	// Path: payment/vault/vault_payment
	// SourceModel: Otnegam\Vault\Model\Adminhtml\Source\VaultProvidersMap
	PaymentVaultVaultPayment model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.PaymentVaultVaultPayment = model.NewStr(`payment/vault/vault_payment`, model.WithConfigStructure(cfgStruct))

	return pp
}
