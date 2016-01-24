// +build ignore

package backup

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
	// SystemBackupEnabled => Enable Scheduled Backup.
	// Path: system/backup/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SystemBackupEnabled model.Bool

	// SystemBackupType => Backup Type.
	// Path: system/backup/type
	// SourceModel: Magento\Backup\Model\Config\Source\Type
	SystemBackupType model.Str

	// SystemBackupTime => Start Time.
	// Path: system/backup/time
	SystemBackupTime model.Str

	// SystemBackupFrequency => Frequency.
	// Path: system/backup/frequency
	// BackendModel: Magento\Backup\Model\Config\Backend\Cron
	// SourceModel: Magento\Cron\Model\Config\Source\Frequency
	SystemBackupFrequency model.Str

	// SystemBackupMaintenance => Maintenance Mode.
	// Please put your store into maintenance mode during backup.
	// Path: system/backup/maintenance
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	SystemBackupMaintenance model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemBackupEnabled = model.NewBool(`system/backup/enabled`, model.WithConfigStructure(cfgStruct))
	pp.SystemBackupType = model.NewStr(`system/backup/type`, model.WithConfigStructure(cfgStruct))
	pp.SystemBackupTime = model.NewStr(`system/backup/time`, model.WithConfigStructure(cfgStruct))
	pp.SystemBackupFrequency = model.NewStr(`system/backup/frequency`, model.WithConfigStructure(cfgStruct))
	pp.SystemBackupMaintenance = model.NewBool(`system/backup/maintenance`, model.WithConfigStructure(cfgStruct))

	return pp
}
