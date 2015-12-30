// +build ignore

package adminnotification

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSystemAdminnotificationUseHttps => Use HTTPS to Get Feed.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSystemAdminnotificationUseHttps = model.NewBool(`system/adminnotification/use_https`)

// PathSystemAdminnotificationFrequency => Update Frequency.
// SourceModel: Otnegam\AdminNotification\Model\Config\Source\Frequency
var PathSystemAdminnotificationFrequency = model.NewStr(`system/adminnotification/frequency`)

// PathSystemAdminnotificationLastUpdate => Last Update.
var PathSystemAdminnotificationLastUpdate = model.NewStr(`system/adminnotification/last_update`)
