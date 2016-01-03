// +build ignore

package productalert

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID: "catalog",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "productalert",
					Label:     `Product Alerts`,
					SortOrder: 250,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/productalert/allow_price
							ID:        "allow_price",
							Label:     `Allow Alert When Product Price Changes`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/productalert/allow_stock
							ID:        "allow_stock",
							Label:     `Allow Alert When Product Comes Back in Stock`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/productalert/email_price_template
							ID:        "email_price_template",
							Label:     `Price Alert Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `catalog_productalert_email_price_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: catalog/productalert/email_stock_template
							ID:        "email_stock_template",
							Label:     `Stock Alert Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `catalog_productalert_email_stock_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: catalog/productalert/email_identity
							ID:        "email_identity",
							Label:     `Alert Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `general`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},
					),
				},

				&element.Group{
					ID:        "productalert_cron",
					Label:     `Product Alerts Run Settings`,
					SortOrder: 260,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/productalert_cron/frequency
							ID:        "frequency",
							Label:     `Frequency`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Otnegam\Cron\Model\Config\Backend\Product\Alert
							// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
						},

						&element.Field{
							// Path: catalog/productalert_cron/time
							ID:        "time",
							Label:     `Start Time`,
							Type:      element.TypeTime,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
						},

						&element.Field{
							// Path: catalog/productalert_cron/error_email
							ID:        "error_email",
							Label:     `Error Email Recipient`,
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
						},

						&element.Field{
							// Path: catalog/productalert_cron/error_email_identity
							ID:        "error_email_identity",
							Label:     `Error Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `general`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: catalog/productalert_cron/error_email_template
							ID:        "error_email_template",
							Label:     `Error Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `catalog_productalert_cron_error_email_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "catalog",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "productalert_cron",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/productalert_cron/error_email
							ID:      `error_email`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
