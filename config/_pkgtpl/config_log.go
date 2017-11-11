// +build ignore

package log

import (
	"github.com/corestoreio/cspkg/config"
	"github.com/corestoreio/cspkg/store/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "customer",
		Label:     "",
		SortOrder: 130,
		Scope:     scope.PermStore,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "online_customers",
				Label:     `Online Customers Options`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     scope.PermDefault,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/online_customers/online_minutes_interval`,
						ID:           "online_minutes_interval",
						Label:        `Online Minutes Interval`,
						Comment:      `Leave empty for default (15 minutes).`,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermDefault,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "system",
		Label:     "",
		SortOrder: 0,
		Scope:     scope.PermStore,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "log",
				Label:     `Log Cleaning`,
				Comment:   ``,
				SortOrder: 200,
				Scope:     scope.PermDefault,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/log/clean_after_day`,
						ID:           "clean_after_day",
						Label:        `Save Log, Days`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermDefault,
						Default:      180,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `system/log/enabled`,
						ID:           "enabled",
						Label:        `Enable Log Cleaning`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermDefault,
						Default:      false,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `system/log/time`,
						ID:           "time",
						Label:        `Start Time`,
						Comment:      ``,
						Type:         config.TypeTime,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.PermDefault,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `system/log/frequency`,
						ID:           "frequency",
						Label:        `Frequency`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        scope.PermDefault,
						Default:      `D`,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Log\Cron
						// SourceModel:  nil, // Magento\Cron\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: `system/log/error_email`,
						ID:           "error_email",
						Label:        `Error Email Recipient`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.PermDefault,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `system/log/error_email_identity`,
						ID:           "error_email_identity",
						Label:        `Error Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    6,
						Visible:      config.VisibleYes,
						Scope:        scope.PermDefault,
						Default:      `general`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `system/log/error_email_template`,
						ID:           "error_email_template",
						Label:        `Error Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    7,
						Visible:      config.VisibleYes,
						Scope:        scope.PermDefault,
						Default:      `system_log_error_email_template`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "log",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "visitor",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `log/visitor/online_update_frequency`,
						ID:      "online_update_frequency",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.PermDefault, // @todo search for that
						Default: 60,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "system",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "log",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/log/time`,
						ID:      "time",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.PermDefault, // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `system/log/error_email`,
						ID:      "error_email",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.PermDefault, // @todo search for that
						Default: nil,
					},
				},
			},
		},
	},
)
