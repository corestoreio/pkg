// +build ignore

package paypal

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "payment",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "paypal_notice",
				SortOrder: 3,
				Scope:     scope.PermAll,
				HelpURL:   element.LongText(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup`),
				Fields:    config.NewFieldSlice(),
			},

			&config.Group{
				ID:        "account",
				Label:     `Merchant Location`,
				SortOrder: 1,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						ConfigPath: `paypal/general/merchant_country`, // Original: payment/account/merchant_country
						ID:         "merchant_country",
						Label:      `Merchant Country`,
						Comment:    element.LongText(`If not specified, Default Country from General Config will be used`),
						Type:       config.TypeSelect,
						SortOrder:  5,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// BackendModel: Otnegam\Paypal\Model\System\Config\Backend\MerchantCountry
						// SourceModel: Otnegam\Paypal\Model\System\Config\Source\MerchantCountry
					},
				),
			},
		),
	},
	&config.Section{
		ID: "payment_all_paypal",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:     "payments_pro_hosted_solution_without_bml",
				Fields: config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_us",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "paypal_group_all_in_one",
				Label:     `PayPal All-in-One Payment Solutions&nbsp;&nbsp;<i>Accept and process credit cards and PayPal payments.</i>`,
				Comment:   element.LongText(`Choose a secure bundled payment solution for your business.`),
				SortOrder: 10,
				Scope:     scope.PermAll,
				HelpURL:   element.LongText(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup&subTab=newlineup`),
				Fields:    config.NewFieldSlice(),
			},

			&config.Group{
				ID:        "paypal_payment_gateways",
				Label:     `PayPal Payment Gateways`,
				Comment:   element.LongText(`Process payments using your own internet merchant account.`),
				SortOrder: 15,
				Scope:     scope.PermAll,
				HelpURL:   element.LongText(`https://merchant.paypal.com/cgi-bin/marketingweb?cmd=_render-content`),
				Fields:    config.NewFieldSlice(),
			},

			&config.Group{
				ID:        "paypal_alternative_payment_methods",
				Label:     `PayPal Express Checkout`,
				Comment:   element.LongText(`Add another payment method to your existing solution or as a stand-alone option.`),
				SortOrder: 20,
				Scope:     scope.PermAll,
				HelpURL:   element.LongText(`https://merchant.paypal.com/cgi-bin/marketingweb?cmd=_render-content`),
				Fields:    config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_gb",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "paypal_group_all_in_one",
				Label:     `PayPal All-in-One Payment Solutions&nbsp;&nbsp;<i>Accept and process credit cards and PayPal payments.</i>`,
				Comment:   element.LongText(`Choose a secure bundled payment solution for your business.`),
				SortOrder: 10,
				Scope:     scope.PermAll,
				HelpURL:   element.LongText(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup&subTab=newlineup`),
				Fields:    config.NewFieldSlice(),
			},

			&config.Group{
				ID:        "paypal_alternative_payment_methods",
				Label:     `PayPal Express Checkout`,
				Comment:   element.LongText(`Add another payment method to your existing solution or as a stand-alone option.`),
				SortOrder: 20,
				Scope:     scope.PermAll,
				HelpURL:   element.LongText(`https://merchant.paypal.com/cgi-bin/marketingweb?cmd=_render-content`),
				Fields:    config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_de",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "paypal_payment_solutions",
				Label:     `PayPal Payment Solutions`,
				Comment:   element.LongText(`Add another payment method to your existing solution or as a stand-alone option.`),
				SortOrder: 10,
				Scope:     scope.PermAll,
				HelpURL:   element.LongText(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup&subTab=newlineup`),
				Fields:    config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_other",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "paypal_payment_solutions",
				Label:     `PayPal Payment Solutions`,
				Comment:   element.LongText(`Add another payment method to your existing solution or as a stand-alone option.`),
				SortOrder: 10,
				Scope:     scope.PermAll,
				HelpURL:   element.LongText(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup&subTab=newlineup`),
				Fields:    config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_ca",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:     "paypal_payment_solutions",
				Fields: config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_au",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:     "paypal_payment_solutions",
				Fields: config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_jp",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:     "paypal_payment_solutions",
				Fields: config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_fr",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:     "paypal_payment_solutions",
				Fields: config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_it",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:     "paypal_payment_solutions",
				Fields: config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_es",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:     "paypal_payment_solutions",
				Fields: config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_hk",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:     "paypal_payment_solutions",
				Fields: config.NewFieldSlice(),
			},
		),
	},
	&config.Section{
		ID: "payment_nz",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:     "paypal_payment_solutions",
				Fields: config.NewFieldSlice(),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "paypal",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "style",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: paypal/style/logo
						ID:      `logo`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},
				),
			},

			&config.Group{
				ID: "wpp",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: paypal/wpp/api_password
						ID:      `api_password`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: paypal/wpp/api_signature
						ID:      `api_signature`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: paypal/wpp/api_username
						ID:      `api_username`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: paypal/wpp/button_flavor
						ID:      `button_flavor`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `dynamic`,
					},
				),
			},

			&config.Group{
				ID: "wpuk",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: paypal/wpuk/user
						ID:      `user`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: paypal/wpuk/pwd
						ID:      `pwd`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},
				),
			},

			&config.Group{
				ID: "fetch_reports",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: paypal/fetch_reports/ftp_login
						ID:      `ftp_login`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: paypal/fetch_reports/ftp_password
						ID:      `ftp_password`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: paypal/fetch_reports/schedule
						ID:      `schedule`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: paypal/fetch_reports/time
						ID:      `time`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `00,00,00`,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "payment",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "paypal_express",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/paypal_express/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Paypal\Model\Express`,
					},

					&config.Field{
						// Path: payment/paypal_express/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `PayPal Express Checkout`,
					},

					&config.Field{
						// Path: payment/paypal_express/payment_action
						ID:      `payment_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Authorization`,
					},

					&config.Field{
						// Path: payment/paypal_express/solution_type
						ID:      `solution_type`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Mark`,
					},

					&config.Field{
						// Path: payment/paypal_express/line_items_enabled
						ID:      `line_items_enabled`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/paypal_express/visible_on_cart
						ID:      `visible_on_cart`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/paypal_express/visible_on_product
						ID:      `visible_on_product`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/paypal_express/allow_ba_signup
						ID:      `allow_ba_signup`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `never`,
					},

					&config.Field{
						// Path: payment/paypal_express/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal`,
					},

					&config.Field{
						// Path: payment/paypal_express/authorization_honor_period
						ID:      `authorization_honor_period`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 3,
					},

					&config.Field{
						// Path: payment/paypal_express/order_valid_period
						ID:      `order_valid_period`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 29,
					},

					&config.Field{
						// Path: payment/paypal_express/child_authorization_number
						ID:      `child_authorization_number`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/paypal_express/verify_peer
						ID:      `verify_peer`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/paypal_express/skip_order_review_step
						ID:      `skip_order_review_step`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},
				),
			},

			&config.Group{
				ID: "paypal_express_bml",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/paypal_express_bml/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Paypal\Model\Bml`,
					},

					&config.Field{
						// Path: payment/paypal_express_bml/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `PayPal Credit`,
					},

					&config.Field{
						// Path: payment/paypal_express_bml/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal`,
					},
				),
			},

			&config.Group{
				ID: "payflow_express",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/payflow_express/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `PayPal Express Checkout Payflow Edition`,
					},

					&config.Field{
						// Path: payment/payflow_express/payment_action
						ID:      `payment_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Authorization`,
					},

					&config.Field{
						// Path: payment/payflow_express/line_items_enabled
						ID:      `line_items_enabled`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_express/visible_on_cart
						ID:      `visible_on_cart`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_express/visible_on_product
						ID:      `visible_on_product`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_express/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal`,
					},

					&config.Field{
						// Path: payment/payflow_express/verify_peer
						ID:      `verify_peer`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_express/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Paypal\Model\PayflowExpress`,
					},
				),
			},

			&config.Group{
				ID: "payflow_express_bml",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/payflow_express_bml/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Paypal\Model\Payflow\Bml`,
					},

					&config.Field{
						// Path: payment/payflow_express_bml/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `PayPal Credit`,
					},

					&config.Field{
						// Path: payment/payflow_express_bml/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal`,
					},
				),
			},

			&config.Group{
				ID: "payflowpro",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/payflowpro/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Paypal\Model\Payflow\Transparent`,
					},

					&config.Field{
						// Path: payment/payflowpro/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Credit Card`,
					},

					&config.Field{
						// Path: payment/payflowpro/payment_action
						ID:      `payment_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Authorization`,
					},

					&config.Field{
						// Path: payment/payflowpro/cctypes
						ID:      `cctypes`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `AE,VI`,
					},

					&config.Field{
						// Path: payment/payflowpro/useccv
						ID:      `useccv`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflowpro/tender
						ID:      `tender`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `C`,
					},

					&config.Field{
						// Path: payment/payflowpro/verbosity
						ID:      `verbosity`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `MEDIUM`,
					},

					&config.Field{
						// Path: payment/payflowpro/user
						ID:      `user`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: payment/payflowpro/pwd
						ID:      `pwd`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: payment/payflowpro/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal`,
					},

					&config.Field{
						// Path: payment/payflowpro/verify_peer
						ID:      `verify_peer`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflowpro/date_delim
						ID:      `date_delim`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},

					&config.Field{
						// Path: payment/payflowpro/ccfields
						ID:      `ccfields`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `csc,expdate,acct`,
					},

					&config.Field{
						// Path: payment/payflowpro/place_order_url
						ID:      `place_order_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal/transparent/requestSecureToken`,
					},

					&config.Field{
						// Path: payment/payflowpro/cgi_url_test_mode
						ID:      `cgi_url_test_mode`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://pilot-payflowlink.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflowpro/cgi_url
						ID:      `cgi_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://payflowlink.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflowpro/transaction_url_test_mode
						ID:      `transaction_url_test_mode`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://pilot-payflowpro.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflowpro/transaction_url
						ID:      `transaction_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://payflowpro.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflowpro/avs_street
						ID:      `avs_street`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/payflowpro/avs_zip
						ID:      `avs_zip`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/payflowpro/avs_international
						ID:      `avs_international`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/payflowpro/avs_security_code
						ID:      `avs_security_code`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflowpro/cc_year_length
						ID:      `cc_year_length`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 2,
					},
				),
			},

			&config.Group{
				ID: "paypal_billing_agreement",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/paypal_billing_agreement/active
						ID:      `active`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/paypal_billing_agreement/allow_billing_agreement_wizard
						ID:      `allow_billing_agreement_wizard`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/paypal_billing_agreement/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Paypal\Model\Method\Agreement`,
					},

					&config.Field{
						// Path: payment/paypal_billing_agreement/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `PayPal Billing Agreement`,
					},

					&config.Field{
						// Path: payment/paypal_billing_agreement/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal`,
					},

					&config.Field{
						// Path: payment/paypal_billing_agreement/verify_peer
						ID:      `verify_peer`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},
				),
			},

			&config.Group{
				ID: "payflow_link",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/payflow_link/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Paypal\Model\Payflowlink`,
					},

					&config.Field{
						// Path: payment/payflow_link/payment_action
						ID:      `payment_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Authorization`,
					},

					&config.Field{
						// Path: payment/payflow_link/verbosity
						ID:      `verbosity`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `HIGH`,
					},

					&config.Field{
						// Path: payment/payflow_link/user
						ID:      `user`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: payment/payflow_link/pwd
						ID:      `pwd`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: payment/payflow_link/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal`,
					},

					&config.Field{
						// Path: payment/payflow_link/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Credit Card`,
					},

					&config.Field{
						// Path: payment/payflow_link/partner
						ID:      `partner`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `PayPal`,
					},

					&config.Field{
						// Path: payment/payflow_link/csc_required
						ID:      `csc_required`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_link/csc_editable
						ID:      `csc_editable`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_link/url_method
						ID:      `url_method`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `GET`,
					},

					&config.Field{
						// Path: payment/payflow_link/email_confirmation
						ID:      `email_confirmation`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/payflow_link/verify_peer
						ID:      `verify_peer`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_link/transaction_url_test_mode
						ID:      `transaction_url_test_mode`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://pilot-payflowpro.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflow_link/transaction_url
						ID:      `transaction_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://payflowpro.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflow_link/cgi_url_test_mode
						ID:      `cgi_url_test_mode`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://pilot-payflowlink.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflow_link/cgi_url
						ID:      `cgi_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://payflowlink.paypal.com`,
					},
				),
			},

			&config.Group{
				ID: "payflow_advanced",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/payflow_advanced/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Paypal\Model\Payflowadvanced`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/payment_action
						ID:      `payment_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Authorization`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/verbosity
						ID:      `verbosity`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `HIGH`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/user
						ID:      `user`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `[{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}},"PayPal"]`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/pwd
						ID:      `pwd`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Credit Card`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/partner
						ID:      `partner`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `PayPal`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/vendor
						ID:      `vendor`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `PayPal`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/csc_required
						ID:      `csc_required`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_advanced/csc_editable
						ID:      `csc_editable`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_advanced/url_method
						ID:      `url_method`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `GET`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/email_confirmation
						ID:      `email_confirmation`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/payflow_advanced/verify_peer
						ID:      `verify_peer`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/payflow_advanced/transaction_url_test_mode
						ID:      `transaction_url_test_mode`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://pilot-payflowpro.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/transaction_url
						ID:      `transaction_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://payflowpro.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/cgi_url_test_mode
						ID:      `cgi_url_test_mode`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://pilot-payflowlink.paypal.com`,
					},

					&config.Field{
						// Path: payment/payflow_advanced/cgi_url
						ID:      `cgi_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://payflowlink.paypal.com`,
					},
				),
			},

			&config.Group{
				ID: "hosted_pro",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/hosted_pro/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Paypal\Model\Hostedpro`,
					},

					&config.Field{
						// Path: payment/hosted_pro/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Payment by cards or by PayPal account`,
					},

					&config.Field{
						// Path: payment/hosted_pro/payment_action
						ID:      `payment_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Authorization`,
					},

					&config.Field{
						// Path: payment/hosted_pro/group
						ID:      `group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `paypal`,
					},

					&config.Field{
						// Path: payment/hosted_pro/display_ec
						ID:      `display_ec`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/hosted_pro/verify_peer
						ID:      `verify_peer`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},
				),
			},
		),
	},
)
