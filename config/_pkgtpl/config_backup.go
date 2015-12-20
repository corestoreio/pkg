// +build ignore

package backup

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "system",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "backup",
				Label:     `Scheduled Backup Settings`,
				Comment:   ``,
				SortOrder: 500,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/backup/enabled`,
						ID:           "enabled",
						Label:        `Enable Scheduled Backup`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `system/backup/type`,
						ID:           "type",
						Label:        `Backup Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Backup\Model\Config\Source\Type
					},

					&config.Field{
						// Path: `system/backup/time`,
						ID:           "time",
						Label:        `Start Time`,
						Comment:      ``,
						Type:         config.TypeTime,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `system/backup/frequency`,
						ID:           "frequency",
						Label:        `Frequency`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil, // Magento\Backup\Model\Config\Backend\Cron
						// SourceModel:  nil, // Magento\Cron\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: `system/backup/maintenance`,
						ID:           "maintenance",
						Label:        `Maintenance Mode`,
						Comment:      `Please put your store into maintenance mode during backup.`,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
)
