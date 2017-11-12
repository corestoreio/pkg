// +build ignore

package customer

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
			ID:        "customer",
			Label:     `Customer Configuration`,
			SortOrder: 130,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Customer::config_customer
			Groups: element.NewGroupSlice(
				element.Group{
					ID:                    "account_share",
					Label:                 `Account Sharing Options`,
					SortOrder:             10,
					Scopes:                scope.PermDefault,
					HideInSingleStoreMode: true,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: customer/account_share/scope
							ID:        "scope",
							Label:     `Share Customer Accounts`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   true,
							// BackendModel: Magento\Customer\Model\Config\Share
							// SourceModel: Magento\Customer\Model\Config\Share
						},
					),
				},

				element.Group{
					ID:        "create_account",
					Label:     `Create New Account Options`,
					SortOrder: 20,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: customer/create_account/auto_group_assign
							ID:        "auto_group_assign",
							Label:     `Enable Automatic Assignment to Customer Group`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: customer/create_account/tax_calculation_address_type
							ID:        "tax_calculation_address_type",
							Label:     `Tax Calculation Based On`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `billing`,
							// SourceModel: Magento\Customer\Model\Config\Source\Address\Type
						},

						element.Field{
							// Path: customer/create_account/default_group
							ID:        "default_group",
							Label:     `Default Group`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Customer\Model\Config\Source\Group
						},

						element.Field{
							// Path: customer/create_account/viv_domestic_group
							ID:        "viv_domestic_group",
							Label:     `Group for Valid VAT ID - Domestic`,
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Customer\Model\Config\Source\Group
						},

						element.Field{
							// Path: customer/create_account/viv_intra_union_group
							ID:        "viv_intra_union_group",
							Label:     `Group for Valid VAT ID - Intra-Union`,
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Customer\Model\Config\Source\Group
						},

						element.Field{
							// Path: customer/create_account/viv_invalid_group
							ID:        "viv_invalid_group",
							Label:     `Group for Invalid VAT ID`,
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Customer\Model\Config\Source\Group
						},

						element.Field{
							// Path: customer/create_account/viv_error_group
							ID:        "viv_error_group",
							Label:     `Validation Error Group`,
							Type:      element.TypeSelect,
							SortOrder: 55,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Customer\Model\Config\Source\Group
						},

						element.Field{
							// Path: customer/create_account/viv_on_each_transaction
							ID:        "viv_on_each_transaction",
							Label:     `Validate on Each Transaction`,
							Type:      element.TypeSelect,
							SortOrder: 56,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: customer/create_account/viv_disable_auto_group_assign_default
							ID:        "viv_disable_auto_group_assign_default",
							Label:     `Default Value for Disable Automatic Group Changes Based on VAT ID`,
							Type:      element.TypeSelect,
							SortOrder: 57,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Customer\Model\Config\Backend\CreateAccount\DisableAutoGroupAssignDefault
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: customer/create_account/vat_frontend_visibility
							ID:        "vat_frontend_visibility",
							Label:     `Show VAT Number on Storefront`,
							Comment:   text.Long(`To show VAT number on Storefront, set Show VAT Number on Storefront option to Yes.`),
							Type:      element.TypeSelect,
							SortOrder: 58,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: customer/create_account/email_domain
							ID:        "email_domain",
							Label:     `Default Email Domain`,
							Type:      element.TypeText,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `example.com`,
						},

						element.Field{
							// Path: customer/create_account/email_template
							ID:        "email_template",
							Label:     `Default Welcome Email`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `customer_create_account_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: customer/create_account/email_no_password_template
							ID:        "email_no_password_template",
							Label:     `Default Welcome Email Without Password`,
							Comment:   text.Long(`This email will be sent instead of the Default Welcome Email, if a customer was created without password. <br /><br /> Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 75,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `customer_create_account_email_no_password_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: customer/create_account/email_identity
							ID:        "email_identity",
							Label:     `Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: customer/create_account/confirm
							ID:        "confirm",
							Label:     `Require Emails Confirmation`,
							Type:      element.TypeSelect,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: customer/create_account/email_confirmation_template
							ID:        "email_confirmation_template",
							Label:     `Confirmation Link Email`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `customer_create_account_email_confirmation_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: customer/create_account/email_confirmed_template
							ID:        "email_confirmed_template",
							Label:     `Welcome Email`,
							Comment:   text.Long(`This email will be sent instead of the Default Welcome Email, after account confirmation. <br /><br /> Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 110,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `customer_create_account_email_confirmed_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: customer/create_account/generate_human_friendly_id
							ID:        "generate_human_friendly_id",
							Label:     `Generate Human-Friendly Customer ID`,
							Type:      element.TypeSelect,
							SortOrder: 120,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        "password",
					Label:     `Password Options`,
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: customer/password/forgot_email_template
							ID:        "forgot_email_template",
							Label:     `Forgot Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `customer_password_forgot_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: customer/password/remind_email_template
							ID:        "remind_email_template",
							Label:     `Remind Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `customer_password_remind_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: customer/password/reset_password_template
							ID:        "reset_password_template",
							Label:     `Reset Password Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `customer_password_reset_password_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: customer/password/forgot_email_identity
							ID:        "forgot_email_identity",
							Label:     `Password Template Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `support`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: customer/password/reset_link_expiration_period
							ID:        "reset_link_expiration_period",
							Label:     `Recovery Link Expiration Period (days)`,
							Comment:   text.Long(`Please enter a number 1 or greater in this field.`),
							Type:      element.TypeText,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   1,
							// BackendModel: Magento\Customer\Model\Config\Backend\Password\Link\Expirationperiod
						},
					),
				},

				element.Group{
					ID:        "address",
					Label:     `Name and Address Options`,
					SortOrder: 40,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: customer/address/street_lines
							ID:        "street_lines",
							Label:     `Number of Lines in a Street Address`,
							Comment:   text.Long(`Leave empty for default (2). Valid range: 1-4`),
							Type:      element.Type,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   2,
							// BackendModel: Magento\Customer\Model\Config\Backend\Address\Street
						},

						element.Field{
							// Path: customer/address/prefix_show
							ID:        "prefix_show",
							Label:     `Show Prefix`,
							Comment:   text.Long(`The title that goes before name (Mr., Mrs., etc.)`),
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// BackendModel: Magento\Customer\Model\Config\Backend\Show\Address
							// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
						},

						element.Field{
							// Path: customer/address/prefix_options
							ID:        "prefix_options",
							Label:     `Prefix Dropdown Options`,
							Comment:   text.Long(`Semicolon (;) separated values.<br/>Put semicolon in the beginning for empty first option.<br/>Leave empty for open text field.`),
							Type:      element.Type,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: customer/address/middlename_show
							ID:        "middlename_show",
							Label:     `Show Middle Name (initial)`,
							Comment:   text.Long(`Always optional.`),
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// BackendModel: Magento\Customer\Model\Config\Backend\Show\Address
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: customer/address/suffix_show
							ID:        "suffix_show",
							Label:     `Show Suffix`,
							Comment:   text.Long(`The suffix that goes after name (Jr., Sr., etc.)`),
							Type:      element.TypeSelect,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// BackendModel: Magento\Customer\Model\Config\Backend\Show\Address
							// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
						},

						element.Field{
							// Path: customer/address/suffix_options
							ID:        "suffix_options",
							Label:     `Suffix Dropdown Options`,
							Comment:   text.Long(`Semicolon (;) separated values.<br/>Put semicolon in the beginning for empty first option.<br/>Leave empty for open text field.`),
							Type:      element.Type,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},

						element.Field{
							// Path: customer/address/dob_show
							ID:        "dob_show",
							Label:     `Show Date of Birth`,
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// BackendModel: Magento\Customer\Model\Config\Backend\Show\Customer
							// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
						},

						element.Field{
							// Path: customer/address/taxvat_show
							ID:        "taxvat_show",
							Label:     `Show Tax/VAT Number`,
							Type:      element.TypeSelect,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// BackendModel: Magento\Customer\Model\Config\Backend\Show\Customer
							// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
						},

						element.Field{
							// Path: customer/address/gender_show
							ID:        "gender_show",
							Label:     `Show Gender`,
							Type:      element.TypeSelect,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// BackendModel: Magento\Customer\Model\Config\Backend\Show\Customer
							// SourceModel: Magento\Config\Model\Config\Source\Nooptreq
						},
					),
				},

				element.Group{
					ID:        "startup",
					Label:     `Login Options`,
					SortOrder: 90,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: customer/startup/redirect_dashboard
							ID:        "redirect_dashboard",
							Label:     `Redirect Customer to Account Dashboard after Logging in`,
							Comment:   text.Long(`Customer will stay on the current page if "No" is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID:        "address_templates",
					Label:     `Address Templates`,
					SortOrder: 100,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: customer/address_templates/text
							ID:        "text",
							Label:     `Text`,
							Type:      element.TypeTextarea,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default: `{{depend prefix}}{{var prefix}} {{/depend}}{{var firstname}} {{depend middlename}}{{var middlename}} {{/depend}}{{var lastname}}{{depend suffix}} {{var suffix}}{{/depend}}
{{depend company}}{{var company}}{{/depend}}
{{if street1}}{{var street1}}
{{/if}}
{{depend street2}}{{var street2}}{{/depend}}
{{depend street3}}{{var street3}}{{/depend}}
{{depend street4}}{{var street4}}{{/depend}}
{{if city}}{{var city}},  {{/if}}{{if region}}{{var region}}, {{/if}}{{if postcode}}{{var postcode}}{{/if}}
{{var country}}
T: {{var telephone}}
{{depend fax}}F: {{var fax}}{{/depend}}
{{depend vat_id}}VAT: {{var vat_id}}{{/depend}}`,
						},

						element.Field{
							// Path: customer/address_templates/oneline
							ID:        "oneline",
							Label:     `Text One Line`,
							Type:      element.TypeTextarea,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `{{depend prefix}}{{var prefix}} {{/depend}}{{var firstname}} {{depend middlename}}{{var middlename}} {{/depend}}{{var lastname}}{{depend suffix}} {{var suffix}}{{/depend}}, {{var street}}, {{var city}}, {{var region}} {{var postcode}}, {{var country}}`,
						},

						element.Field{
							// Path: customer/address_templates/html
							ID:        "html",
							Label:     `HTML`,
							Type:      element.TypeTextarea,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default: `{{depend prefix}}{{var prefix}} {{/depend}}{{var firstname}} {{depend middlename}}{{var middlename}} {{/depend}}{{var lastname}}{{depend suffix}} {{var suffix}}{{/depend}}{{depend firstname}}<br/>{{/depend}}
{{depend company}}{{var company}}<br />{{/depend}}
{{if street1}}{{var street1}}<br />{{/if}}
{{depend street2}}{{var street2}}<br />{{/depend}}
{{depend street3}}{{var street3}}<br />{{/depend}}
{{depend street4}}{{var street4}}<br />{{/depend}}
{{if city}}{{var city}},  {{/if}}{{if region}}{{var region}}, {{/if}}{{if postcode}}{{var postcode}}{{/if}}<br/>
{{var country}}<br/>
{{depend telephone}}T: {{var telephone}}{{/depend}}
{{depend fax}}<br/>F: {{var fax}}{{/depend}}
{{depend vat_id}}<br/>VAT: {{var vat_id}}{{/depend}}`,
						},

						element.Field{
							// Path: customer/address_templates/pdf
							ID:        "pdf",
							Label:     `PDF`,
							Type:      element.TypeTextarea,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default: `{{depend prefix}}{{var prefix}} {{/depend}}{{var firstname}} {{depend middlename}}{{var middlename}} {{/depend}}{{var lastname}}{{depend suffix}} {{var suffix}}{{/depend}}|
{{depend company}}{{var company}}|{{/depend}}
{{if street1}}{{var street1}}
{{/if}}
{{depend street2}}{{var street2}}|{{/depend}}
{{depend street3}}{{var street3}}|{{/depend}}
{{depend street4}}{{var street4}}|{{/depend}}
{{if city}}{{var city}},|{{/if}}
{{if region}}{{var region}}, {{/if}}{{if postcode}}{{var postcode}}{{/if}}|
{{var country}}|
{{depend telephone}}T: {{var telephone}}{{/depend}}|
{{depend fax}}<br/>F: {{var fax}}{{/depend}}|
{{depend vat_id}}<br/>VAT: {{var vat_id}}{{/depend}}|`,
						},
					),
				},

				element.Group{
					ID:        "online_customers",
					Label:     `Online Customers Options`,
					SortOrder: 10,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: customer/online_customers/online_minutes_interval
							ID:        "online_minutes_interval",
							Label:     `Online Minutes Interval`,
							Comment:   text.Long(`Leave empty for default (15 minutes).`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
					),
				},
			),
		},
		element.Section{
			ID: "general",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "store_information",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: general/store_information/validate_vat_number
							ID:        "validate_vat_number",
							Type:      element.Type,
							SortOrder: 62,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
						},
					),
				},

				element.Group{
					ID: "restriction",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: general/restriction/autocomplete_on_storefront
							ID:        "autocomplete_on_storefront",
							Label:     `Enable Autocomplete on login/forgot password forms`,
							Type:      element.TypeSelect,
							SortOrder: 65,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "customer",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "default",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: customer/default/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},

				element.Group{
					ID: "address",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: customer/address/prefix_show
							ID:      `prefix_show`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						element.Field{
							// Path: customer/address/prefix_options
							ID:      `prefix_options`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						element.Field{
							// Path: customer/address/middlename_show
							ID:      `middlename_show`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						element.Field{
							// Path: customer/address/suffix_show
							ID:      `suffix_show`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						element.Field{
							// Path: customer/address/suffix_options
							ID:      `suffix_options`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						element.Field{
							// Path: customer/address/dob_show
							ID:      `dob_show`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						element.Field{
							// Path: customer/address/gender_show
							ID:      `gender_show`,
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
