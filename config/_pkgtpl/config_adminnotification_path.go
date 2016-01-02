// +build ignore

package adminnotification

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
	// SystemAdminnotificationUseHttps => Use HTTPS to Get Feed.
	// Path: system/adminnotification/use_https
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	SystemAdminnotificationUseHttps model.Bool

	// SystemAdminnotificationFrequency => Update Frequency.
	// Path: system/adminnotification/frequency
	// SourceModel: Otnegam\AdminNotification\Model\Config\Source\Frequency
	SystemAdminnotificationFrequency model.Str

	// SystemAdminnotificationLastUpdate => Last Update.
	// Path: system/adminnotification/last_update
	SystemAdminnotificationLastUpdate model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemAdminnotificationUseHttps = model.NewBool(`system/adminnotification/use_https`, model.WithPkgCfg(pkgCfg))
	pp.SystemAdminnotificationFrequency = model.NewStr(`system/adminnotification/frequency`, model.WithPkgCfg(pkgCfg))
	pp.SystemAdminnotificationLastUpdate = model.NewStr(`system/adminnotification/last_update`, model.WithPkgCfg(pkgCfg))

	return pp
}
