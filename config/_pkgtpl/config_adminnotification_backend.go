// +build ignore

package adminnotification

import (
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// SystemAdminnotificationUseHttps => Use HTTPS to Get Feed.
	// Path: system/adminnotification/use_https
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SystemAdminnotificationUseHttps cfgmodel.Bool

	// SystemAdminnotificationFrequency => Update Frequency.
	// Path: system/adminnotification/frequency
	// SourceModel: Magento\AdminNotification\Model\Config\Source\Frequency
	SystemAdminnotificationFrequency cfgmodel.Str

	// SystemAdminnotificationLastUpdate => Last Update.
	// Path: system/adminnotification/last_update
	SystemAdminnotificationLastUpdate cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemAdminnotificationUseHttps = cfgmodel.NewBool(`system/adminnotification/use_https`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemAdminnotificationFrequency = cfgmodel.NewStr(`system/adminnotification/frequency`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemAdminnotificationLastUpdate = cfgmodel.NewStr(`system/adminnotification/last_update`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
