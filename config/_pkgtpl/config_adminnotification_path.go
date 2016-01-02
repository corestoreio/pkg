// +build ignore

package adminnotification

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSystemAdminnotificationUseHttps => Use HTTPS to Get Feed.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSystemAdminnotificationUseHttps = model.NewBool(`system/adminnotification/use_https`, model.WithPkgCfg(PackageConfiguration))

// PathSystemAdminnotificationFrequency => Update Frequency.
// SourceModel: Otnegam\AdminNotification\Model\Config\Source\Frequency
var PathSystemAdminnotificationFrequency = model.NewStr(`system/adminnotification/frequency`, model.WithPkgCfg(PackageConfiguration))

// PathSystemAdminnotificationLastUpdate => Last Update.
var PathSystemAdminnotificationLastUpdate = model.NewStr(`system/adminnotification/last_update`, model.WithPkgCfg(PackageConfiguration))
