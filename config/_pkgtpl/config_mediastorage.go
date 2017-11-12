// +build ignore

package mediastorage

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "system",
			SortOrder: 900,
			Scopes:    scope.PermStore,
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "media_storage_configuration",
					Label:     `Storage Configuration for Media`,
					SortOrder: 900,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: system/media_storage_configuration/media_storage
							ID:        "media_storage",
							Label:     `Media Storage`,
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\MediaStorage\Model\Config\Source\Storage\Media\Storage
						},

						element.Field{
							// Path: system/media_storage_configuration/media_database
							ID:        "media_database",
							Label:     `Select Media Database`,
							Type:      element.TypeSelect,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\MediaStorage\Model\Config\Backend\Storage\Media\Database
							// SourceModel: Magento\MediaStorage\Model\Config\Source\Storage\Media\Database
						},

						element.Field{
							// Path: system/media_storage_configuration/synchronize
							ID:        "synchronize",
							Comment:   text.Long(`After selecting a new media storage location, press the Synchronize button to transfer all media to that location. Media will not be available in the new location until the synchronization process is complete.`),
							Type:      element.TypeButton,
							SortOrder: 300,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},

						element.Field{
							// Path: system/media_storage_configuration/configuration_update_time
							ID:        "configuration_update_time",
							Label:     `Environment Update Time`,
							Type:      element.TypeText,
							SortOrder: 400,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
