// +build ignore

package productalert

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "catalog",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "productalert",
				Label:     `Product Alerts`,
				SortOrder: 250,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/productalert/allow_price
						ID:        "allow_price",
						Label:     `Allow Alert When Product Price Changes`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/productalert/allow_stock
						ID:        "allow_stock",
						Label:     `Allow Alert When Product Comes Back in Stock`,
						Type:      config.TypeSelect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/productalert/email_price_template
						ID:        "email_price_template",
						Label:     `Price Alert Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `catalog_productalert_email_price_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: catalog/productalert/email_stock_template
						ID:        "email_stock_template",
						Label:     `Stock Alert Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `catalog_productalert_email_stock_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: catalog/productalert/email_identity
						ID:        "email_identity",
						Label:     `Alert Email Sender`,
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `general`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},
				),
			},

			&config.Group{
				ID:        "productalert_cron",
				Label:     `Product Alerts Run Settings`,
				SortOrder: 260,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/productalert_cron/frequency
						ID:        "frequency",
						Label:     `Frequency`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Cron\Model\Config\Backend\Product\Alert
						// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: catalog/productalert_cron/time
						ID:        "time",
						Label:     `Start Time`,
						Type:      config.TypeTime,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},

					&config.Field{
						// Path: catalog/productalert_cron/error_email
						ID:        "error_email",
						Label:     `Error Email Recipient`,
						Type:      config.TypeText,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},

					&config.Field{
						// Path: catalog/productalert_cron/error_email_identity
						ID:        "error_email_identity",
						Label:     `Error Email Sender`,
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `general`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: catalog/productalert_cron/error_email_template
						ID:        "error_email_template",
						Label:     `Error Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `catalog_productalert_cron_error_email_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "catalog",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "productalert_cron",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/productalert_cron/error_email
						ID:      `error_email`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},
				),
			},
		),
	},
)
