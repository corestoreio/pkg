// +build ignore

package mediastorage

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "system",
		Label:     "",
		SortOrder: 900,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "media_storage_configuration",
				Label:     `Storage Configuration for Media`,
				Comment:   ``,
				SortOrder: 900,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/media_storage_configuration/media_storage`,
						ID:           "media_storage",
						Label:        `Media Storage`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\MediaStorage\Model\Config\Source\Storage\Media\Storage
					},

					&config.Field{
						// Path: `system/media_storage_configuration/media_database`,
						ID:           "media_database",
						Label:        `Select Media Database`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    200,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil, // Magento\MediaStorage\Model\Config\Backend\Storage\Media\Database
						// SourceModel:  nil, // Magento\MediaStorage\Model\Config\Source\Storage\Media\Database
					},

					&config.Field{
						// Path: `system/media_storage_configuration/synchronize`,
						ID:           "synchronize",
						Label:        ``,
						Comment:      `After selecting a new media storage location, press the Synchronize button to transfer all media to that location. Media will not be available in the new location until the synchronization process is complete.`,
						Type:         config.TypeButton,
						SortOrder:    300,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `system/media_storage_configuration/configuration_update_time`,
						ID:           "configuration_update_time",
						Label:        `Environment Update Time`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    400,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},
)
