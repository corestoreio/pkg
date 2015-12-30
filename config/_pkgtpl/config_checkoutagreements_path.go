// +build ignore

package checkoutagreements

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCheckoutOptionsEnableAgreements => Enable Terms and Conditions.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCheckoutOptionsEnableAgreements = model.NewBool(`checkout/options/enable_agreements`)
