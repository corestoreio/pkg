// +build ignore

package braintree

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
			ID: "payment",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "braintree_section",
					Label:     `Braintree`,
					Comment:   text.Long(`Accept credit/debit cards and PayPal in your Magento store. No setup or monthly fees and your customers never leave your store to complete the purchase.`),
					SortOrder: 2,
					Scope:     scope.PermStore,
					Fields:    element.NewFieldSlice(),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "payment",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "braintree",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: payment/braintree/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Braintree\Model\PaymentMethod`,
						},

						&element.Field{
							// Path: payment/braintree/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Credit Card (Braintree)`,
						},

						&element.Field{
							// Path: payment/braintree/payment_action
							ID:      `payment_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `authorize`,
						},

						&element.Field{
							// Path: payment/braintree/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintree/cctypes
							ID:      `cctypes`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `AE,VI,MC,DI,JCB`,
						},

						&element.Field{
							// Path: payment/braintree/useccv
							ID:      `useccv`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/braintree/verify_3dsecure
							ID:      `verify_3dsecure`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintree/order_status
							ID:      `order_status`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `processing`,
						},

						&element.Field{
							// Path: payment/braintree/environment
							ID:      `environment`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `sandbox`,
						},

						&element.Field{
							// Path: payment/braintree/allowspecific
							ID:      `allowspecific`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintree/fraudprotection
							ID:      `fraudprotection`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintree/capture_action
							ID:      `capture_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `invoice`,
						},

						&element.Field{
							// Path: payment/braintree/data_js
							ID:      `data_js`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://js.braintreegateway.com/v1/braintree-data.js`,
						},

						&element.Field{
							// Path: payment/braintree/public_key
							ID:      `public_key`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						&element.Field{
							// Path: payment/braintree/private_key
							ID:      `private_key`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						&element.Field{
							// Path: payment/braintree/duplicate_card
							ID:      `duplicate_card`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintree/masked_fields
							ID:      `masked_fields`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `cvv,number`,
						},

						&element.Field{
							// Path: payment/braintree/usecache
							ID:      `usecache`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintree/enable_cc_detection
							ID:      `enable_cc_detection`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},

				&element.Group{
					ID: "braintree_paypal",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: payment/braintree_paypal/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Braintree\Model\PaymentMethod\PayPal`,
						},

						&element.Field{
							// Path: payment/braintree_paypal/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `PayPal (Braintree)`,
						},

						&element.Field{
							// Path: payment/braintree_paypal/payment_action
							ID:      `payment_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `authorize`,
						},

						&element.Field{
							// Path: payment/braintree_paypal/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintree_paypal/allowspecific
							ID:      `allowspecific`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintree_paypal/order_status
							ID:      `order_status`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `processing`,
						},

						&element.Field{
							// Path: payment/braintree_paypal/dispaly_on_shopping_cart
							ID:      `dispaly_on_shopping_cart`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/braintree_paypal/require_billing_address
							ID:      `require_billing_address`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
