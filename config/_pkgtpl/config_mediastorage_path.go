// +build ignore

package mediastorage

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathSystemMediaStorageConfigurationMediaStorage => Media Storage.
// SourceModel: Otnegam\MediaStorage\Model\Config\Source\Storage\Media\Storage
var PathSystemMediaStorageConfigurationMediaStorage = model.NewStr(`system/media_storage_configuration/media_storage`)

// PathSystemMediaStorageConfigurationMediaDatabase => Select Media Database.
// BackendModel: Otnegam\MediaStorage\Model\Config\Backend\Storage\Media\Database
// SourceModel: Otnegam\MediaStorage\Model\Config\Source\Storage\Media\Database
var PathSystemMediaStorageConfigurationMediaDatabase = model.NewStr(`system/media_storage_configuration/media_database`)

// PathSystemMediaStorageConfigurationSynchronize => .
// After selecting a new media storage location, press the Synchronize button
// to transfer all media to that location. Media will not be available in the
// new location until the synchronization process is complete.
var PathSystemMediaStorageConfigurationSynchronize = model.NewStr(`system/media_storage_configuration/synchronize`)

// PathSystemMediaStorageConfigurationConfigurationUpdateTime => Environment Update Time.
var PathSystemMediaStorageConfigurationConfigurationUpdateTime = model.NewStr(`system/media_storage_configuration/configuration_update_time`)
