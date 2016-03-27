// +build ignore

package email

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID: "design",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "email",
					Label:     `Emails`,
					SortOrder: 510,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: design/email/logo
							ID:        "logo",
							Label:     `Logo Image`,
							Comment:   text.Long(`Allowed file types: jpg, jpeg, gif, png. To optimize logo for high-resolution displays, upload an image that is 3x normal size and then specify 1x dimensions in width/height fields below.`),
							Type:      element.TypeImage,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Email\Logo
						},

						element.Field{
							// Path: design/email/logo_alt
							ID:        "logo_alt",
							Label:     `Logo Image Alt`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: design/email/logo_width
							ID:        "logo_width",
							Label:     `Logo Width`,
							Comment:   text.Long(`Only necessary if image has been uploaded above. Enter number of pixels, without appending "px".`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: design/email/logo_height
							ID:        "logo_height",
							Label:     `Logo Height`,
							Comment:   text.Long(`Only necessary if image has been uploaded above. Enter number of pixels, without appending "px".`),
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: design/email/header_template
							ID:        "header_template",
							Label:     `Header Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `design_email_header_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: design/email/footer_template
							ID:        "footer_template",
							Label:     `Footer Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `design_email_footer_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "system",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "media_storage_configuration",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      `allowed_resources`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"email_folder":"email"}`,
						},
					),
				},

				element.Group{
					ID: "emails",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: system/emails/forgot_email_template
							ID:      `forgot_email_template`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `system_emails_forgot_email_template`,
						},

						element.Field{
							// Path: system/emails/forgot_email_identity
							ID:      `forgot_email_identity`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `general`,
						},
					),
				},

				element.Group{
					ID: "smtp",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: system/smtp/disable
							ID:      `disable`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: system/smtp/host
							ID:      `host`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `localhost`,
						},

						element.Field{
							// Path: system/smtp/port
							ID:      `port`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 25,
						},
					),
				},
			),
		},
		element.Section{
			ID: "trans_email",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "ident_custom1",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_custom1/email
							ID:      `email`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `custom1@example.com`,
						},

						element.Field{
							// Path: trans_email/ident_custom1/name
							ID:      `name`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Custom 1`,
						},
					),
				},

				element.Group{
					ID: "ident_custom2",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_custom2/email
							ID:      `email`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `custom2@example.com`,
						},

						element.Field{
							// Path: trans_email/ident_custom2/name
							ID:      `name`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Custom 2`,
						},
					),
				},

				element.Group{
					ID: "ident_general",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_general/email
							ID:      `email`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `owner@example.com`,
						},

						element.Field{
							// Path: trans_email/ident_general/name
							ID:      `name`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Owner`,
						},
					),
				},

				element.Group{
					ID: "ident_sales",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_sales/email
							ID:      `email`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `sales@example.com`,
						},

						element.Field{
							// Path: trans_email/ident_sales/name
							ID:      `name`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Sales`,
						},
					),
				},

				element.Group{
					ID: "ident_support",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: trans_email/ident_support/email
							ID:      `email`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `support@example.com`,
						},

						element.Field{
							// Path: trans_email/ident_support/name
							ID:      `name`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `CustomerSupport`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
