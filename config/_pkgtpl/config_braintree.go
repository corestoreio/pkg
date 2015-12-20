// +build ignore

package braintree

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "payment",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:    "braintree",
				Label: `Braintree`,
				Comment: `Accept credit/debit cards and PayPal in your Magento store. No setup or monthly fees and your customers never leave your store to complete the purchase.
                    <a href="https://www.braintreegateway.com/login" target="_blank">Click here to login to your existing Braintree account</a>. Or to setup a new account and accept payments on your website, <a href="https://apply.braintreegateway.com/signup/us" target="_blank">click here to signup for a Braintree account</a>.`,
				SortOrder: 25,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/braintree/active`,
						ID:           "active",
						Label:        `Enabled Braintree`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree/active_braintree_pay_pal`,
						ID:           "active_braintree_pay_pal",
						Label:        `Enabled PayPal through Braintree`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    11,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Credit Card (Braintree)`,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree/environment`,
						ID:           "environment",
						Label:        `Environment`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `sandbox`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Braintree\Model\Source\Environment
					},

					&config.Field{
						// Path: `payment/braintree/payment_action`,
						ID:           "payment_action",
						Label:        `Payment Action`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `authorize`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Braintree\Model\Source\PaymentAction
					},

					&config.Field{
						// Path: `payment/braintree/merchant_account_id`,
						ID:           "merchant_account_id",
						Label:        `Merchant Account ID`,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree/merchant_id`,
						ID:           "merchant_id",
						Label:        `Merchant ID`,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree/public_key`,
						ID:           "public_key",
						Label:        `Public Key`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree/private_key`,
						ID:           "private_key",
						Label:        `Private Key`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    110,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree/debug`,
						ID:           "debug",
						Label:        `Debug`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree/capture_action`,
						ID:           "capture_action",
						Label:        `Capture action`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `invoice`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Braintree\Model\Source\CaptureAction
					},

					&config.Field{
						// Path: `payment/braintree/order_status`,
						ID:           "order_status",
						Label:        `New Order Status`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `processing`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Sales\Model\Config\Source\Order\Status\Processing
					},

					&config.Field{
						// Path: `payment/braintree/use_vault`,
						ID:           "use_vault",
						Label:        `Use Vault`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    130,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree/duplicate_card`,
						ID:           "duplicate_card",
						Label:        `Allow Duplicate Cards`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    140,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree/useccv`,
						ID:           "useccv",
						Label:        `CVV Verification`,
						Comment:      `Be sure to Enable AVS and/or CVV in Your Braintree Account in Settings/Processing Section`,
						Type:         config.TypeSelect,
						SortOrder:    150,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree/cctypes`,
						ID:           "cctypes",
						Label:        `Credit Card Types`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    160,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `AE,VI,MC,DI,JCB`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Braintree\Model\Source\CcType
					},

					&config.Field{
						// Path: `payment/braintree/enable_cc_detection`,
						ID:           "enable_cc_detection",
						Label:        `Enable Credit Card auto-detection on Storefront`,
						Comment:      `Typing in a credit card number will automatically select the credit card type`,
						Type:         config.TypeSelect,
						SortOrder:    170,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree/fraudprotection`,
						ID:    "fraudprotection",
						Label: `Advanced Fraud Protection`,
						Comment: `Be sure to Enable Advanced Fraud Protection in Your Braintree Account in
                        Settings/Processing Section`,
						Type:         config.TypeSelect,
						SortOrder:    180,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree/kount_id`,
						ID:           "kount_id",
						Label:        `Your Kount ID`,
						Comment:      `Used for direct fraud tool integration. Make sure you also contact <a href="mailto:accounts@braintreepayments.com">accounts@braintreepayments.com</a> to setup your Kount account.`,
						Type:         config.TypeTextarea,
						SortOrder:    185,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree/usecache`,
						ID:           "usecache",
						Label:        `Use Cache`,
						Comment:      `Some of results will be cached to improve performance. Magento cache have to be enabled`,
						Type:         config.TypeSelect,
						SortOrder:    190,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    230,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree/allowspecific`,
						ID:           "allowspecific",
						Label:        `Payment from Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeAllowspecific,
						SortOrder:    200,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      0,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `payment/braintree/specificcountry`,
						ID:           "specificcountry",
						Label:        `Payment from Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    210,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Braintree\Model\System\Config\Source\Country
					},

					&config.Field{
						// Path: `payment/braintree/countrycreditcard`,
						ID:           "countrycreditcard",
						Label:        `Country Specific Credit Card Types`,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    220,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Braintree\Model\System\Config\Backend\Countrycreditcard
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree/verify_3dsecure`,
						ID:           "verify_3dsecure",
						Label:        `3d Secure Verification`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    150,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "braintree_paypal",
				Label:     `PayPal through Braintree`,
				Comment:   ``,
				SortOrder: 40,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/braintree_paypal/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `PayPal (Braintree)`,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree_paypal/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree_paypal/merchant_name_override`,
						ID:           "merchant_name_override",
						Label:        `Override Merchant Name`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/braintree_paypal/payment_action`,
						ID:           "payment_action",
						Label:        `Payment Action`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `authorize`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Braintree\Model\Source\PaymentAction
					},

					&config.Field{
						// Path: `payment/braintree_paypal/order_status`,
						ID:           "order_status",
						Label:        `New Order Status`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `processing`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Sales\Model\Config\Source\Order\Status\Processing
					},

					&config.Field{
						// Path: `payment/braintree_paypal/allowspecific`,
						ID:           "allowspecific",
						Label:        `Payment from Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeAllowspecific,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      0,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `payment/braintree_paypal/specificcountry`,
						ID:           "specificcountry",
						Label:        `Payment from Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Braintree\Model\System\Config\Source\Country
					},

					&config.Field{
						// Path: `payment/braintree_paypal/display_on_shopping_cart`,
						ID:           "display_on_shopping_cart",
						Label:        `Display on Shopping Cart`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree_paypal/require_billing_address`,
						ID:           "require_billing_address",
						Label:        `Require Customer's Billing Address`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree_paypal/allow_shipping_address_override`,
						ID:           "allow_shipping_address_override",
						Label:        `Allow to Edit Shipping Address Entered During Checkout on PayPal Side`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    110,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/braintree_paypal/debug`,
						ID:           "debug",
						Label:        `Debug`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    120,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "payment",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "braintree",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/braintree/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `Magento\Braintree\Model\PaymentMethod`,
					},

					&config.Field{
						// Path: `payment/braintree/data_js`,
						ID:      "data_js",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `https://js.braintreegateway.com/v1/braintree-data.js`,
					},

					&config.Field{
						// Path: `payment/braintree/masked_fields`,
						ID:      "masked_fields",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `cvv,number`,
					},
				},
			},

			&config.Group{
				ID: "braintree_paypal",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/braintree_paypal/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `Magento\Braintree\Model\PaymentMethod\PayPal`,
					},

					&config.Field{
						// Path: `payment/braintree_paypal/active`,
						ID:      "active",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `payment/braintree_paypal/dispaly_on_shopping_cart`,
						ID:      "dispaly_on_shopping_cart",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: true,
					},
				},
			},
		},
	},
)
