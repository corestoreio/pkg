// +build ignore

package backup

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
	// SystemBackupEnabled => Enable Scheduled Backup.
	// Path: system/backup/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SystemBackupEnabled cfgmodel.Bool

	// SystemBackupType => Backup Type.
	// Path: system/backup/type
	// SourceModel: Magento\Backup\Model\Config\Source\Type
	SystemBackupType cfgmodel.Str

	// SystemBackupTime => Start Time.
	// Path: system/backup/time
	SystemBackupTime cfgmodel.Str

	// SystemBackupFrequency => Frequency.
	// Path: system/backup/frequency
	// BackendModel: Magento\Backup\Model\Config\Backend\Cron
	// SourceModel: Magento\Cron\Model\Config\Source\Frequency
	SystemBackupFrequency cfgmodel.Str

	// SystemBackupMaintenance => Maintenance Mode.
	// Please put your store into maintenance mode during backup.
	// Path: system/backup/maintenance
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SystemBackupMaintenance cfgmodel.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemBackupEnabled = cfgmodel.NewBool(`system/backup/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemBackupType = cfgmodel.NewStr(`system/backup/type`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemBackupTime = cfgmodel.NewStr(`system/backup/time`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemBackupFrequency = cfgmodel.NewStr(`system/backup/frequency`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemBackupMaintenance = cfgmodel.NewBool(`system/backup/maintenance`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
