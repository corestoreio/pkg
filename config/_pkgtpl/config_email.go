// +build ignore

package email

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "design",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "email",
				Label:     `Emails`,
				SortOrder: 510,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: design/email/logo
						ID:        "logo",
						Label:     `Logo Image`,
						Comment:   element.LongText(`Allowed file types: jpg, jpeg, gif, png. To optimize logo for high-resolution displays, upload an image that is 3x normal size and then specify 1x dimensions in width/height fields below.`),
						Type:      config.TypeImage,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Email\Logo
					},

					&config.Field{
						// Path: design/email/logo_alt
						ID:        "logo_alt",
						Label:     `Logo Image Alt`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: design/email/logo_width
						ID:        "logo_width",
						Label:     `Logo Width`,
						Comment:   element.LongText(`Only necessary if image has been uploaded above. Enter number of pixels, without appending "px".`),
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: design/email/logo_height
						ID:        "logo_height",
						Label:     `Logo Height`,
						Comment:   element.LongText(`Only necessary if image has been uploaded above. Enter number of pixels, without appending "px".`),
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: design/email/header_template
						ID:        "header_template",
						Label:     `Header Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `design_email_header_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: design/email/footer_template
						ID:        "footer_template",
						Label:     `Footer Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 60,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `design_email_footer_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "media_storage_configuration",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/media_storage_configuration/allowed_resources
						ID:      `allowed_resources`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"email_folder":"email"}`,
					},
				),
			},

			&config.Group{
				ID: "emails",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/emails/forgot_email_template
						ID:      `forgot_email_template`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `system_emails_forgot_email_template`,
					},

					&config.Field{
						// Path: system/emails/forgot_email_identity
						ID:      `forgot_email_identity`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `general`,
					},
				),
			},

			&config.Group{
				ID: "smtp",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/smtp/disable
						ID:      `disable`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: system/smtp/host
						ID:      `host`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `localhost`,
					},

					&config.Field{
						// Path: system/smtp/port
						ID:      `port`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 25,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "trans_email",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "ident_custom1",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_custom1/email
						ID:      `email`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `custom1@example.com`,
					},

					&config.Field{
						// Path: trans_email/ident_custom1/name
						ID:      `name`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Custom 1`,
					},
				),
			},

			&config.Group{
				ID: "ident_custom2",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_custom2/email
						ID:      `email`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `custom2@example.com`,
					},

					&config.Field{
						// Path: trans_email/ident_custom2/name
						ID:      `name`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Custom 2`,
					},
				),
			},

			&config.Group{
				ID: "ident_general",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_general/email
						ID:      `email`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `owner@example.com`,
					},

					&config.Field{
						// Path: trans_email/ident_general/name
						ID:      `name`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Owner`,
					},
				),
			},

			&config.Group{
				ID: "ident_sales",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_sales/email
						ID:      `email`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `sales@example.com`,
					},

					&config.Field{
						// Path: trans_email/ident_sales/name
						ID:      `name`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Sales`,
					},
				),
			},

			&config.Group{
				ID: "ident_support",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: trans_email/ident_support/email
						ID:      `email`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `support@example.com`,
					},

					&config.Field{
						// Path: trans_email/ident_support/name
						ID:      `name`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `CustomerSupport`,
					},
				),
			},
		),
	},
)
