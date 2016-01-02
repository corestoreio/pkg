// +build ignore

package backup

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSystemBackupEnabled => Enable Scheduled Backup.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSystemBackupEnabled = model.NewBool(`system/backup/enabled`, model.WithPkgCfg(PackageConfiguration))

// PathSystemBackupType => Backup Type.
// SourceModel: Otnegam\Backup\Model\Config\Source\Type
var PathSystemBackupType = model.NewStr(`system/backup/type`, model.WithPkgCfg(PackageConfiguration))

// PathSystemBackupTime => Start Time.
var PathSystemBackupTime = model.NewStr(`system/backup/time`, model.WithPkgCfg(PackageConfiguration))

// PathSystemBackupFrequency => Frequency.
// BackendModel: Otnegam\Backup\Model\Config\Backend\Cron
// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
var PathSystemBackupFrequency = model.NewStr(`system/backup/frequency`, model.WithPkgCfg(PackageConfiguration))

// PathSystemBackupMaintenance => Maintenance Mode.
// Please put your store into maintenance mode during backup.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathSystemBackupMaintenance = model.NewBool(`system/backup/maintenance`, model.WithPkgCfg(PackageConfiguration))
