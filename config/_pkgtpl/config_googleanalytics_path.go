// +build ignore

package googleanalytics

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathGoogleAnalyticsActive => Enable.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathGoogleAnalyticsActive = model.NewBool(`google/analytics/active`)

// PathGoogleAnalyticsAccount => Account Number.
var PathGoogleAnalyticsAccount = model.NewStr(`google/analytics/account`)
