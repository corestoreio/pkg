// +build ignore

package productalert

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "productalert",
				Label:     `Product Alerts`,
				Comment:   ``,
				SortOrder: 250,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/productalert/allow_price`,
						ID:           "allow_price",
						Label:        `Allow Alert When Product Price Changes`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/productalert/allow_stock`,
						ID:           "allow_stock",
						Label:        `Allow Alert When Product Comes Back in Stock`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/productalert/email_price_template`,
						ID:           "email_price_template",
						Label:        `Price Alert Email Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `catalog_productalert_email_price_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `catalog/productalert/email_stock_template`,
						ID:           "email_stock_template",
						Label:        `Stock Alert Email Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `catalog_productalert_email_stock_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `catalog/productalert/email_identity`,
						ID:           "email_identity",
						Label:        `Alert Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `general`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},
				},
			},

			&config.Group{
				ID:        "productalert_cron",
				Label:     `Product Alerts Run Settings`,
				Comment:   ``,
				SortOrder: 260,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/productalert_cron/frequency`,
						ID:           "frequency",
						Label:        `Frequency`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil, // Magento\Cron\Model\Config\Backend\Product\Alert
						SourceModel:  nil, // Magento\Cron\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: `catalog/productalert_cron/time`,
						ID:           "time",
						Label:        `Start Time`,
						Comment:      ``,
						Type:         config.TypeTime,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/productalert_cron/error_email`,
						ID:           "error_email",
						Label:        `Error Email Recipient`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/productalert_cron/error_email_identity`,
						ID:           "error_email_identity",
						Label:        `Error Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `general`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `catalog/productalert_cron/error_email_template`,
						ID:           "error_email_template",
						Label:        `Error Email Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `catalog_productalert_cron_error_email_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "catalog",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "productalert_cron",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/productalert_cron/error_email`,
						ID:      "error_email",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: nil,
					},
				},
			},
		},
	},
)
