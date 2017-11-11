// +build ignore

package vault

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// PaymentVaultVaultPayment => Vault Provider.
	// Specified provider should be enabled.
	// Path: payment/vault/vault_payment
	// SourceModel: Magento\Vault\Model\Adminhtml\Source\VaultProvidersMap
	PaymentVaultVaultPayment cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.PaymentVaultVaultPayment = cfgmodel.NewStr(`payment/vault/vault_payment`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
