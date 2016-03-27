// +build ignore

package paypal

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
			ID: "payment",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "paypal_notice",
					SortOrder: 3,
					Scopes:    scope.PermStore,
					HelpURL:   text.Long(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup`),
					Fields:    element.NewFieldSlice(),
				},

				element.Group{
					ID:        "account",
					Label:     `Merchant Location`,
					SortOrder: 1,
					Scopes:    scope.PermWebsite,
					Fields: element.NewFieldSlice(
						element.Field{
							ConfigPath: `paypal/general/merchant_country`, // Original: payment/account/merchant_country
							ID:         "merchant_country",
							Label:      `Merchant Country`,
							Comment:    text.Long(`If not specified, Default Country from General Config will be used`),
							Type:       element.TypeSelect,
							SortOrder:  5,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermWebsite,
							// BackendModel: Magento\Paypal\Model\System\Config\Backend\MerchantCountry
							// SourceModel: Magento\Paypal\Model\System\Config\Source\MerchantCountry
						},
					),
				},
			),
		},
		element.Section{
			ID: "payment_all_paypal",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:     "payments_pro_hosted_solution_without_bml",
					Fields: element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_us",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "paypal_group_all_in_one",
					Label:     `PayPal All-in-One Payment Solutions&nbsp;&nbsp;<i>Accept and process credit cards and PayPal payments.</i>`,
					Comment:   text.Long(`Choose a secure bundled payment solution for your business.`),
					SortOrder: 10,
					Scopes:    scope.PermStore,
					HelpURL:   text.Long(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup&subTab=newlineup`),
					Fields:    element.NewFieldSlice(),
				},

				element.Group{
					ID:        "paypal_payment_gateways",
					Label:     `PayPal Payment Gateways`,
					Comment:   text.Long(`Process payments using your own internet merchant account.`),
					SortOrder: 15,
					Scopes:    scope.PermStore,
					HelpURL:   text.Long(`https://merchant.paypal.com/cgi-bin/marketingweb?cmd=_render-content`),
					Fields:    element.NewFieldSlice(),
				},

				element.Group{
					ID:        "paypal_alternative_payment_methods",
					Label:     `PayPal Express Checkout`,
					Comment:   text.Long(`Add another payment method to your existing solution or as a stand-alone option.`),
					SortOrder: 20,
					Scopes:    scope.PermStore,
					HelpURL:   text.Long(`https://merchant.paypal.com/cgi-bin/marketingweb?cmd=_render-content`),
					Fields:    element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_gb",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "paypal_group_all_in_one",
					Label:     `PayPal All-in-One Payment Solutions&nbsp;&nbsp;<i>Accept and process credit cards and PayPal payments.</i>`,
					Comment:   text.Long(`Choose a secure bundled payment solution for your business.`),
					SortOrder: 10,
					Scopes:    scope.PermStore,
					HelpURL:   text.Long(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup&subTab=newlineup`),
					Fields:    element.NewFieldSlice(),
				},

				element.Group{
					ID:        "paypal_alternative_payment_methods",
					Label:     `PayPal Express Checkout`,
					Comment:   text.Long(`Add another payment method to your existing solution or as a stand-alone option.`),
					SortOrder: 20,
					Scopes:    scope.PermStore,
					HelpURL:   text.Long(`https://merchant.paypal.com/cgi-bin/marketingweb?cmd=_render-content`),
					Fields:    element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_de",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "paypal_payment_solutions",
					Label:     `PayPal Payment Solutions`,
					Comment:   text.Long(`Add another payment method to your existing solution or as a stand-alone option.`),
					SortOrder: 10,
					Scopes:    scope.PermStore,
					HelpURL:   text.Long(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup&subTab=newlineup`),
					Fields:    element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_other",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "paypal_payment_solutions",
					Label:     `PayPal Payment Solutions`,
					Comment:   text.Long(`Add another payment method to your existing solution or as a stand-alone option.`),
					SortOrder: 10,
					Scopes:    scope.PermStore,
					HelpURL:   text.Long(`https://www.paypal-marketing.com/emarketing/partner/na/merchantlineup/home.page#mainTab=checkoutlineup&subTab=newlineup`),
					Fields:    element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_ca",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:     "paypal_payment_solutions",
					Fields: element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_au",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:     "paypal_payment_solutions",
					Fields: element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_jp",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:     "paypal_payment_solutions",
					Fields: element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_fr",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:     "paypal_payment_solutions",
					Fields: element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_it",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:     "paypal_payment_solutions",
					Fields: element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_es",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:     "paypal_payment_solutions",
					Fields: element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_hk",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:     "paypal_payment_solutions",
					Fields: element.NewFieldSlice(),
				},
			),
		},
		element.Section{
			ID: "payment_nz",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:     "paypal_payment_solutions",
					Fields: element.NewFieldSlice(),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "paypal",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "style",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: paypal/style/logo
							ID:      `logo`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},
					),
				},

				element.Group{
					ID: "wpp",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: paypal/wpp/api_password
							ID:      `api_password`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: paypal/wpp/api_signature
							ID:      `api_signature`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: paypal/wpp/api_username
							ID:      `api_username`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: paypal/wpp/button_flavor
							ID:      `button_flavor`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `dynamic`,
						},
					),
				},

				element.Group{
					ID: "wpuk",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: paypal/wpuk/user
							ID:      `user`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: paypal/wpuk/pwd
							ID:      `pwd`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},
					),
				},

				element.Group{
					ID: "fetch_reports",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: paypal/fetch_reports/ftp_login
							ID:      `ftp_login`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: paypal/fetch_reports/ftp_password
							ID:      `ftp_password`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: paypal/fetch_reports/schedule
							ID:      `schedule`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: paypal/fetch_reports/time
							ID:      `time`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `00,00,00`,
						},
					),
				},
			),
		},
		element.Section{
			ID: "payment",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "paypal_express",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: payment/paypal_express/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Paypal\Model\Express`,
						},

						element.Field{
							// Path: payment/paypal_express/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `PayPal Express Checkout`,
						},

						element.Field{
							// Path: payment/paypal_express/payment_action
							ID:      `payment_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Authorization`,
						},

						element.Field{
							// Path: payment/paypal_express/solution_type
							ID:      `solution_type`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Mark`,
						},

						element.Field{
							// Path: payment/paypal_express/line_items_enabled
							ID:      `line_items_enabled`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/paypal_express/visible_on_cart
							ID:      `visible_on_cart`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/paypal_express/visible_on_product
							ID:      `visible_on_product`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/paypal_express/allow_ba_signup
							ID:      `allow_ba_signup`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `never`,
						},

						element.Field{
							// Path: payment/paypal_express/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal`,
						},

						element.Field{
							// Path: payment/paypal_express/authorization_honor_period
							ID:      `authorization_honor_period`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 3,
						},

						element.Field{
							// Path: payment/paypal_express/order_valid_period
							ID:      `order_valid_period`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 29,
						},

						element.Field{
							// Path: payment/paypal_express/child_authorization_number
							ID:      `child_authorization_number`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/paypal_express/verify_peer
							ID:      `verify_peer`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/paypal_express/skip_order_review_step
							ID:      `skip_order_review_step`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},

				element.Group{
					ID: "paypal_express_bml",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: payment/paypal_express_bml/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Paypal\Model\Bml`,
						},

						element.Field{
							// Path: payment/paypal_express_bml/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `PayPal Credit`,
						},

						element.Field{
							// Path: payment/paypal_express_bml/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal`,
						},
					),
				},

				element.Group{
					ID: "payflow_express",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: payment/payflow_express/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `PayPal Express Checkout Payflow Edition`,
						},

						element.Field{
							// Path: payment/payflow_express/payment_action
							ID:      `payment_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Authorization`,
						},

						element.Field{
							// Path: payment/payflow_express/line_items_enabled
							ID:      `line_items_enabled`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_express/visible_on_cart
							ID:      `visible_on_cart`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_express/visible_on_product
							ID:      `visible_on_product`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_express/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal`,
						},

						element.Field{
							// Path: payment/payflow_express/verify_peer
							ID:      `verify_peer`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_express/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Paypal\Model\PayflowExpress`,
						},
					),
				},

				element.Group{
					ID: "payflow_express_bml",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: payment/payflow_express_bml/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Paypal\Model\Payflow\Bml`,
						},

						element.Field{
							// Path: payment/payflow_express_bml/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `PayPal Credit`,
						},

						element.Field{
							// Path: payment/payflow_express_bml/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal`,
						},
					),
				},

				element.Group{
					ID: "payflowpro",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: payment/payflowpro/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Paypal\Model\Payflow\Transparent`,
						},

						element.Field{
							// Path: payment/payflowpro/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Credit Card`,
						},

						element.Field{
							// Path: payment/payflowpro/payment_action
							ID:      `payment_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Authorization`,
						},

						element.Field{
							// Path: payment/payflowpro/cctypes
							ID:      `cctypes`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `AE,VI`,
						},

						element.Field{
							// Path: payment/payflowpro/useccv
							ID:      `useccv`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflowpro/tender
							ID:      `tender`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `C`,
						},

						element.Field{
							// Path: payment/payflowpro/verbosity
							ID:      `verbosity`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `MEDIUM`,
						},

						element.Field{
							// Path: payment/payflowpro/user
							ID:      `user`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: payment/payflowpro/pwd
							ID:      `pwd`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: payment/payflowpro/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal`,
						},

						element.Field{
							// Path: payment/payflowpro/verify_peer
							ID:      `verify_peer`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflowpro/date_delim
							ID:      `date_delim`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						element.Field{
							// Path: payment/payflowpro/ccfields
							ID:      `ccfields`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `csc,expdate,acct`,
						},

						element.Field{
							// Path: payment/payflowpro/place_order_url
							ID:      `place_order_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal/transparent/requestSecureToken`,
						},

						element.Field{
							// Path: payment/payflowpro/cgi_url_test_mode
							ID:      `cgi_url_test_mode`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://pilot-payflowlink.paypal.com`,
						},

						element.Field{
							// Path: payment/payflowpro/cgi_url
							ID:      `cgi_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://payflowlink.paypal.com`,
						},

						element.Field{
							// Path: payment/payflowpro/transaction_url_test_mode
							ID:      `transaction_url_test_mode`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://pilot-payflowpro.paypal.com`,
						},

						element.Field{
							// Path: payment/payflowpro/transaction_url
							ID:      `transaction_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://payflowpro.paypal.com`,
						},

						element.Field{
							// Path: payment/payflowpro/avs_street
							ID:      `avs_street`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: payment/payflowpro/avs_zip
							ID:      `avs_zip`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: payment/payflowpro/avs_international
							ID:      `avs_international`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: payment/payflowpro/avs_security_code
							ID:      `avs_security_code`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflowpro/cc_year_length
							ID:      `cc_year_length`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 2,
						},
					),
				},

				element.Group{
					ID: "paypal_billing_agreement",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: payment/paypal_billing_agreement/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/paypal_billing_agreement/allow_billing_agreement_wizard
							ID:      `allow_billing_agreement_wizard`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/paypal_billing_agreement/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Paypal\Model\Method\Agreement`,
						},

						element.Field{
							// Path: payment/paypal_billing_agreement/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `PayPal Billing Agreement`,
						},

						element.Field{
							// Path: payment/paypal_billing_agreement/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal`,
						},

						element.Field{
							// Path: payment/paypal_billing_agreement/verify_peer
							ID:      `verify_peer`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},

				element.Group{
					ID: "payflow_link",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: payment/payflow_link/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Paypal\Model\Payflowlink`,
						},

						element.Field{
							// Path: payment/payflow_link/payment_action
							ID:      `payment_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Authorization`,
						},

						element.Field{
							// Path: payment/payflow_link/verbosity
							ID:      `verbosity`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `HIGH`,
						},

						element.Field{
							// Path: payment/payflow_link/user
							ID:      `user`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: payment/payflow_link/pwd
							ID:      `pwd`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: payment/payflow_link/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal`,
						},

						element.Field{
							// Path: payment/payflow_link/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Credit Card`,
						},

						element.Field{
							// Path: payment/payflow_link/partner
							ID:      `partner`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `PayPal`,
						},

						element.Field{
							// Path: payment/payflow_link/csc_required
							ID:      `csc_required`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_link/csc_editable
							ID:      `csc_editable`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_link/url_method
							ID:      `url_method`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `GET`,
						},

						element.Field{
							// Path: payment/payflow_link/email_confirmation
							ID:      `email_confirmation`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: payment/payflow_link/verify_peer
							ID:      `verify_peer`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_link/transaction_url_test_mode
							ID:      `transaction_url_test_mode`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://pilot-payflowpro.paypal.com`,
						},

						element.Field{
							// Path: payment/payflow_link/transaction_url
							ID:      `transaction_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://payflowpro.paypal.com`,
						},

						element.Field{
							// Path: payment/payflow_link/cgi_url_test_mode
							ID:      `cgi_url_test_mode`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://pilot-payflowlink.paypal.com`,
						},

						element.Field{
							// Path: payment/payflow_link/cgi_url
							ID:      `cgi_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://payflowlink.paypal.com`,
						},
					),
				},

				element.Group{
					ID: "payflow_advanced",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: payment/payflow_advanced/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Paypal\Model\Payflowadvanced`,
						},

						element.Field{
							// Path: payment/payflow_advanced/payment_action
							ID:      `payment_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Authorization`,
						},

						element.Field{
							// Path: payment/payflow_advanced/verbosity
							ID:      `verbosity`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `HIGH`,
						},

						element.Field{
							// Path: payment/payflow_advanced/user
							ID:      `user`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `[{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}},"PayPal"]`,
						},

						element.Field{
							// Path: payment/payflow_advanced/pwd
							ID:      `pwd`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						element.Field{
							// Path: payment/payflow_advanced/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal`,
						},

						element.Field{
							// Path: payment/payflow_advanced/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Credit Card`,
						},

						element.Field{
							// Path: payment/payflow_advanced/partner
							ID:      `partner`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `PayPal`,
						},

						element.Field{
							// Path: payment/payflow_advanced/vendor
							ID:      `vendor`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `PayPal`,
						},

						element.Field{
							// Path: payment/payflow_advanced/csc_required
							ID:      `csc_required`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_advanced/csc_editable
							ID:      `csc_editable`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_advanced/url_method
							ID:      `url_method`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `GET`,
						},

						element.Field{
							// Path: payment/payflow_advanced/email_confirmation
							ID:      `email_confirmation`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: payment/payflow_advanced/verify_peer
							ID:      `verify_peer`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						element.Field{
							// Path: payment/payflow_advanced/transaction_url_test_mode
							ID:      `transaction_url_test_mode`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://pilot-payflowpro.paypal.com`,
						},

						element.Field{
							// Path: payment/payflow_advanced/transaction_url
							ID:      `transaction_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://payflowpro.paypal.com`,
						},

						element.Field{
							// Path: payment/payflow_advanced/cgi_url_test_mode
							ID:      `cgi_url_test_mode`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://pilot-payflowlink.paypal.com`,
						},

						element.Field{
							// Path: payment/payflow_advanced/cgi_url
							ID:      `cgi_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://payflowlink.paypal.com`,
						},
					),
				},

				element.Group{
					ID: "hosted_pro",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: payment/hosted_pro/cfgmodel
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Paypal\Model\Hostedpro`,
						},

						element.Field{
							// Path: payment/hosted_pro/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Payment by cards or by PayPal account`,
						},

						element.Field{
							// Path: payment/hosted_pro/payment_action
							ID:      `payment_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Authorization`,
						},

						element.Field{
							// Path: payment/hosted_pro/group
							ID:      `group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `paypal`,
						},

						element.Field{
							// Path: payment/hosted_pro/display_ec
							ID:      `display_ec`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: payment/hosted_pro/verify_peer
							ID:      `verify_peer`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
