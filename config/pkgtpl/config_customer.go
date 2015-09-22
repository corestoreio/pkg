// +build ignore

package customer

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "customer",
		Label:     "Customer Configuration",
		SortOrder: 130,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "account_share",
				Label:     `Account Sharing Options`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/account_share/scope`,
						ID:           "scope",
						Label:        `Share Customer Accounts`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      true,
						BackendModel: nil, // Magento\Customer\Model\Config\Share
						SourceModel:  nil, // Magento\Customer\Model\Config\Share
					},
				},
			},

			&config.Group{
				ID:        "create_account",
				Label:     `Create New Account Options`,
				Comment:   ``,
				SortOrder: 20,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/create_account/auto_group_assign`,
						ID:           "auto_group_assign",
						Label:        `Enable Automatic Assignment to Customer Group`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `customer/create_account/tax_calculation_address_type`,
						ID:           "tax_calculation_address_type",
						Label:        `Tax Calculation Based On`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `billing`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Customer\Model\Config\Source\Address\Type
					},

					&config.Field{
						// Path: `customer/create_account/default_group`,
						ID:           "default_group",
						Label:        `Default Group`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Customer\Model\Config\Source\Group
					},

					&config.Field{
						// Path: `customer/create_account/viv_domestic_group`,
						ID:           "viv_domestic_group",
						Label:        `Group for Valid VAT ID - Domestic`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Customer\Model\Config\Source\Group
					},

					&config.Field{
						// Path: `customer/create_account/viv_intra_union_group`,
						ID:           "viv_intra_union_group",
						Label:        `Group for Valid VAT ID - Intra-Union`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Customer\Model\Config\Source\Group
					},

					&config.Field{
						// Path: `customer/create_account/viv_invalid_group`,
						ID:           "viv_invalid_group",
						Label:        `Group for Invalid VAT ID`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Customer\Model\Config\Source\Group
					},

					&config.Field{
						// Path: `customer/create_account/viv_error_group`,
						ID:           "viv_error_group",
						Label:        `Validation Error Group`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    55,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Customer\Model\Config\Source\Group
					},

					&config.Field{
						// Path: `customer/create_account/viv_on_each_transaction`,
						ID:           "viv_on_each_transaction",
						Label:        `Validate on Each Transaction`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    56,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `customer/create_account/viv_disable_auto_group_assign_default`,
						ID:           "viv_disable_auto_group_assign_default",
						Label:        `Default Value for Disable Automatic Group Changes Based on VAT ID`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    57,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil, // Magento\Customer\Model\Config\Backend\CreateAccount\DisableAutoGroupAssignDefault
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `customer/create_account/vat_frontend_visibility`,
						ID:           "vat_frontend_visibility",
						Label:        `Show VAT Number on Storefront`,
						Comment:      `To show VAT number on Storefront, set Show VAT Number on Storefront option to Yes.`,
						Type:         config.TypeSelect,
						SortOrder:    58,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `customer/create_account/email_domain`,
						ID:           "email_domain",
						Label:        `Default Email Domain`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `example.com`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/create_account/email_template`,
						ID:           "email_template",
						Label:        `Default Welcome Email`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `customer_create_account_email_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `customer/create_account/email_no_password_template`,
						ID:    "email_no_password_template",
						Label: `Default Welcome Email Without Password`,
						Comment: `This email will be sent instead of the Default Welcome Email, if a customer was created without password. <br /><br />
                        Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    75,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `customer_create_account_email_no_password_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `customer/create_account/email_identity`,
						ID:           "email_identity",
						Label:        `Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `general`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `customer/create_account/confirm`,
						ID:           "confirm",
						Label:        `Require Emails Confirmation`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `customer/create_account/email_confirmation_template`,
						ID:           "email_confirmation_template",
						Label:        `Confirmation Link Email`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `customer_create_account_email_confirmation_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `customer/create_account/email_confirmed_template`,
						ID:    "email_confirmed_template",
						Label: `Welcome Email`,
						Comment: `This email will be sent instead of the Default Welcome Email, after account confirmation. <br /><br />
                        Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    110,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `customer_create_account_email_confirmed_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `customer/create_account/generate_human_friendly_id`,
						ID:           "generate_human_friendly_id",
						Label:        `Generate Human-Friendly Customer ID`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    120,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "password",
				Label:     `Password Options`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/password/forgot_email_template`,
						ID:           "forgot_email_template",
						Label:        `Forgot Email Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `customer_password_forgot_email_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `customer/password/remind_email_template`,
						ID:           "remind_email_template",
						Label:        `Remind Email Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `customer_password_remind_email_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `customer/password/reset_password_template`,
						ID:           "reset_password_template",
						Label:        `Reset Password Template`,
						Comment:      `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `customer_password_reset_password_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `customer/password/forgot_email_identity`,
						ID:           "forgot_email_identity",
						Label:        `Password Template Email Sender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `support`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: `customer/password/reset_link_expiration_period`,
						ID:           "reset_link_expiration_period",
						Label:        `Recovery Link Expiration Period (days)`,
						Comment:      `Please enter a number 1 or greater in this field.`,
						Type:         config.TypeText,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      1,
						BackendModel: nil, // Magento\Customer\Model\Config\Backend\Password\Link\Expirationperiod
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "address",
				Label:     `Name and Address Options`,
				Comment:   ``,
				SortOrder: 40,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/address/street_lines`,
						ID:           "street_lines",
						Label:        `Number of Lines in a Street Address`,
						Comment:      `Leave empty for default (2). Valid range: 1-4`,
						Type:         config.Type,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      2,
						BackendModel: nil, // Magento\Customer\Model\Config\Backend\Address\Street
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/address/prefix_show`,
						ID:           "prefix_show",
						Label:        `Show Prefix`,
						Comment:      `The title that goes before name (Mr., Mrs., etc.)`,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Customer\Model\Config\Backend\Show\Address
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Nooptreq
					},

					&config.Field{
						// Path: `customer/address/prefix_options`,
						ID:           "prefix_options",
						Label:        `Prefix Dropdown Options`,
						Comment:      `Semicolon (;) separated values.<br/>Put semicolon in the beginning for empty first option.<br/>Leave empty for open text field.`,
						Type:         config.Type,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/address/middlename_show`,
						ID:           "middlename_show",
						Label:        `Show Middle Name (initial)`,
						Comment:      `Always optional.`,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Customer\Model\Config\Backend\Show\Address
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `customer/address/suffix_show`,
						ID:           "suffix_show",
						Label:        `Show Suffix`,
						Comment:      `The suffix that goes after name (Jr., Sr., etc.)`,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Customer\Model\Config\Backend\Show\Address
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Nooptreq
					},

					&config.Field{
						// Path: `customer/address/suffix_options`,
						ID:           "suffix_options",
						Label:        `Suffix Dropdown Options`,
						Comment:      `Semicolon (;) separated values.<br/>Put semicolon in the beginning for empty first option.<br/>Leave empty for open text field.`,
						Type:         config.Type,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/address/dob_show`,
						ID:           "dob_show",
						Label:        `Show Date of Birth`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Customer\Model\Config\Backend\Show\Customer
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Nooptreq
					},

					&config.Field{
						// Path: `customer/address/taxvat_show`,
						ID:           "taxvat_show",
						Label:        `Show Tax/VAT Number`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Customer\Model\Config\Backend\Show\Customer
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Nooptreq
					},

					&config.Field{
						// Path: `customer/address/gender_show`,
						ID:           "gender_show",
						Label:        `Show Gender`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Customer\Model\Config\Backend\Show\Customer
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Nooptreq
					},
				},
			},

			&config.Group{
				ID:        "startup",
				Label:     `Login Options`,
				Comment:   ``,
				SortOrder: 90,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/startup/redirect_dashboard`,
						ID:           "redirect_dashboard",
						Label:        `Redirect Customer to Account Dashboard after Logging in`,
						Comment:      `Customer will stay on the current page if "No" is selected.`,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "address_templates",
				Label:     `Address Templates`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/address_templates/text`,
						ID:        "text",
						Label:     `Text`,
						Comment:   ``,
						Type:      config.TypeTextarea,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
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
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/address_templates/oneline`,
						ID:           "oneline",
						Label:        `Text One Line`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `{{depend prefix}}{{var prefix}} {{/depend}}{{var firstname}} {{depend middlename}}{{var middlename}} {{/depend}}{{var lastname}}{{depend suffix}} {{var suffix}}{{/depend}}, {{var street}}, {{var city}}, {{var region}} {{var postcode}}, {{var country}}`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/address_templates/html`,
						ID:        "html",
						Label:     `HTML`,
						Comment:   ``,
						Type:      config.TypeTextarea,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
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
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `customer/address_templates/pdf`,
						ID:        "pdf",
						Label:     `PDF`,
						Comment:   ``,
						Type:      config.TypeTextarea,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
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
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "online_customers",
				Label:     `Online Customers Options`,
				Comment:   ``,
				SortOrder: 10,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/online_customers/online_minutes_interval`,
						ID:           "online_minutes_interval",
						Label:        `Online Minutes Interval`,
						Comment:      `Leave empty for default (15 minutes).`,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "general",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "store_information",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     nil,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `general/store_information/validate_vat_number`,
						ID:           "validate_vat_number",
						Label:        ``,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    62,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "restriction",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     nil,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `general/restriction/autocomplete_on_storefront`,
						ID:           "autocomplete_on_storefront",
						Label:        `Enable Autocomplete on login/forgot password forms`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    65,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "customer",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "default",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/default/group`,
						ID:      "group",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: true,
					},
				},
			},

			&config.Group{
				ID: "address",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `customer/address/prefix_show`,
						ID:      "prefix_show",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `customer/address/prefix_options`,
						ID:      "prefix_options",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `customer/address/middlename_show`,
						ID:      "middlename_show",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `customer/address/suffix_show`,
						ID:      "suffix_show",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `customer/address/suffix_options`,
						ID:      "suffix_options",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `customer/address/dob_show`,
						ID:      "dob_show",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `customer/address/gender_show`,
						ID:      "gender_show",
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
