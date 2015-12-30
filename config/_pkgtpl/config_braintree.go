// +build ignore

package braintree

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
				ID:        "braintree_section",
				Label:     `Braintree`,
				Comment:   element.LongText(`Accept credit/debit cards and PayPal in your Otnegam store. No setup or monthly fees and your customers never leave your store to complete the purchase.`),
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields:    config.NewFieldSlice(),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "payment",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "braintree",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/braintree/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Braintree\Model\PaymentMethod`,
					},

					&config.Field{
						// Path: payment/braintree/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Credit Card (Braintree)`,
					},

					&config.Field{
						// Path: payment/braintree/payment_action
						ID:      `payment_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `authorize`,
					},

					&config.Field{
						// Path: payment/braintree/active
						ID:      `active`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintree/cctypes
						ID:      `cctypes`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `AE,VI,MC,DI,JCB`,
					},

					&config.Field{
						// Path: payment/braintree/useccv
						ID:      `useccv`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/braintree/verify_3dsecure
						ID:      `verify_3dsecure`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintree/order_status
						ID:      `order_status`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `processing`,
					},

					&config.Field{
						// Path: payment/braintree/environment
						ID:      `environment`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `sandbox`,
					},

					&config.Field{
						// Path: payment/braintree/allowspecific
						ID:      `allowspecific`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintree/fraudprotection
						ID:      `fraudprotection`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintree/capture_action
						ID:      `capture_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `invoice`,
					},

					&config.Field{
						// Path: payment/braintree/data_js
						ID:      `data_js`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://js.braintreegateway.com/v1/braintree-data.js`,
					},

					&config.Field{
						// Path: payment/braintree/public_key
						ID:      `public_key`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: payment/braintree/private_key
						ID:      `private_key`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: payment/braintree/duplicate_card
						ID:      `duplicate_card`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintree/masked_fields
						ID:      `masked_fields`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `cvv,number`,
					},

					&config.Field{
						// Path: payment/braintree/usecache
						ID:      `usecache`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintree/enable_cc_detection
						ID:      `enable_cc_detection`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},
				),
			},

			&config.Group{
				ID: "braintree_paypal",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/braintree_paypal/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Otnegam\Braintree\Model\PaymentMethod\PayPal`,
					},

					&config.Field{
						// Path: payment/braintree_paypal/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `PayPal (Braintree)`,
					},

					&config.Field{
						// Path: payment/braintree_paypal/payment_action
						ID:      `payment_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `authorize`,
					},

					&config.Field{
						// Path: payment/braintree_paypal/active
						ID:      `active`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintree_paypal/allowspecific
						ID:      `allowspecific`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintree_paypal/order_status
						ID:      `order_status`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `processing`,
					},

					&config.Field{
						// Path: payment/braintree_paypal/dispaly_on_shopping_cart
						ID:      `dispaly_on_shopping_cart`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/braintree_paypal/require_billing_address
						ID:      `require_billing_address`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},
				),
			},
		),
	},
)
