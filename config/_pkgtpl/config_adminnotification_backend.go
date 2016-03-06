// +build ignore

package adminnotification

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
	// SystemAdminnotificationUseHttps => Use HTTPS to Get Feed.
	// Path: system/adminnotification/use_https
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SystemAdminnotificationUseHttps model.Bool

	// SystemAdminnotificationFrequency => Update Frequency.
	// Path: system/adminnotification/frequency
	// SourceModel: Magento\AdminNotification\Model\Config\Source\Frequency
	SystemAdminnotificationFrequency model.Str

	// SystemAdminnotificationLastUpdate => Last Update.
	// Path: system/adminnotification/last_update
	SystemAdminnotificationLastUpdate model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemAdminnotificationUseHttps = model.NewBool(`system/adminnotification/use_https`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemAdminnotificationFrequency = model.NewStr(`system/adminnotification/frequency`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemAdminnotificationLastUpdate = model.NewStr(`system/adminnotification/last_update`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
