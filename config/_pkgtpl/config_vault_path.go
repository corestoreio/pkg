// +build ignore

package vault

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathPaymentVaultVaultPayment => Vault Provider.
// Specified provider should be enabled.
// SourceModel: Otnegam\Vault\Model\Adminhtml\Source\VaultProvidersMap
var PathPaymentVaultVaultPayment = model.NewStr(`payment/vault/vault_payment`, model.WithPkgCfg(PackageConfiguration))
