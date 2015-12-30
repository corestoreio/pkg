// +build ignore

package backup

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "backup",
				Label:     `Scheduled Backup Settings`,
				SortOrder: 500,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/backup/enabled
						ID:        "enabled",
						Label:     `Enable Scheduled Backup`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: system/backup/type
						ID:        "type",
						Label:     `Backup Type`,
						Type:      config.TypeSelect,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Backup\Model\Config\Source\Type
					},

					&config.Field{
						// Path: system/backup/time
						ID:        "time",
						Label:     `Start Time`,
						Type:      config.TypeTime,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},

					&config.Field{
						// Path: system/backup/frequency
						ID:        "frequency",
						Label:     `Frequency`,
						Type:      config.TypeSelect,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Backup\Model\Config\Backend\Cron
						// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: system/backup/maintenance
						ID:        "maintenance",
						Label:     `Maintenance Mode`,
						Comment:   element.LongText(`Please put your store into maintenance mode during backup.`),
						Type:      config.TypeSelect,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
)
