// +build ignore

package mediastorage

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "system",
		SortOrder: 900,
		Scope:     scope.PermAll,
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "media_storage_configuration",
				Label:     `Storage Configuration for Media`,
				SortOrder: 900,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/media_storage_configuration/media_storage
						ID:        "media_storage",
						Label:     `Media Storage`,
						Type:      config.TypeSelect,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\MediaStorage\Model\Config\Source\Storage\Media\Storage
					},

					&config.Field{
						// Path: system/media_storage_configuration/media_database
						ID:        "media_database",
						Label:     `Select Media Database`,
						Type:      config.TypeSelect,
						SortOrder: 200,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\MediaStorage\Model\Config\Backend\Storage\Media\Database
						// SourceModel: Otnegam\MediaStorage\Model\Config\Source\Storage\Media\Database
					},

					&config.Field{
						// Path: system/media_storage_configuration/synchronize
						ID:        "synchronize",
						Comment:   element.LongText(`After selecting a new media storage location, press the Synchronize button to transfer all media to that location. Media will not be available in the new location until the synchronization process is complete.`),
						Type:      config.TypeButton,
						SortOrder: 300,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},

					&config.Field{
						// Path: system/media_storage_configuration/configuration_update_time
						ID:        "configuration_update_time",
						Label:     `Environment Update Time`,
						Type:      config.TypeText,
						SortOrder: 400,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},
				),
			},
		),
	},
)
