// +build ignore

package mediastorage

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// SystemMediaStorageConfigurationMediaStorage => Media Storage.
	// Path: system/media_storage_configuration/media_storage
	// SourceModel: Magento\MediaStorage\Model\Config\Source\Storage\Media\Storage
	SystemMediaStorageConfigurationMediaStorage cfgmodel.Str

	// SystemMediaStorageConfigurationMediaDatabase => Select Media Database.
	// Path: system/media_storage_configuration/media_database
	// BackendModel: Magento\MediaStorage\Model\Config\Backend\Storage\Media\Database
	// SourceModel: Magento\MediaStorage\Model\Config\Source\Storage\Media\Database
	SystemMediaStorageConfigurationMediaDatabase cfgmodel.Str

	// SystemMediaStorageConfigurationSynchronize => .
	// After selecting a new media storage location, press the Synchronize button
	// to transfer all media to that location. Media will not be available in the
	// new location until the synchronization process is complete.
	// Path: system/media_storage_configuration/synchronize
	SystemMediaStorageConfigurationSynchronize cfgmodel.Str

	// SystemMediaStorageConfigurationConfigurationUpdateTime => Environment Update Time.
	// Path: system/media_storage_configuration/configuration_update_time
	SystemMediaStorageConfigurationConfigurationUpdateTime cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.SystemMediaStorageConfigurationMediaStorage = cfgmodel.NewStr(`system/media_storage_configuration/media_storage`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemMediaStorageConfigurationMediaDatabase = cfgmodel.NewStr(`system/media_storage_configuration/media_database`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemMediaStorageConfigurationSynchronize = cfgmodel.NewStr(`system/media_storage_configuration/synchronize`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemMediaStorageConfigurationConfigurationUpdateTime = cfgmodel.NewStr(`system/media_storage_configuration/configuration_update_time`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
