// +build ignore

package braintreetwo

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
					ID:        "braintreetwo_section",
					Label:     `BraintreeTwo`,
					Comment:   text.Long(`Accept credit/debit cards and PayPal in your Magento store. No setup or monthly fees and your customers never leave your store to complete the purchase.`),
					SortOrder: 2,
					Scopes:    scope.PermStore,
					Fields:    element.NewFieldSlice(),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "payment",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "braintreetwo",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: payment/braintreetwo/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `BraintreeTwoFacade`,
						},

						&element.Field{
							// Path: payment/braintreetwo/title
							ID:      `title`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Credit Card (BraintreeTwo)`,
						},

						&element.Field{
							// Path: payment/braintreetwo/payment_action
							ID:      `payment_action`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `authorize`,
						},

						&element.Field{
							// Path: payment/braintreetwo/active
							ID:      `active`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintreetwo/is_gateway
							ID:      `is_gateway`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/braintreetwo/can_use_checkout
							ID:      `can_use_checkout`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/braintreetwo/can_authorize
							ID:      `can_authorize`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/braintreetwo/can_capture
							ID:      `can_capture`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/braintreetwo/can_capture_partial
							ID:      `can_capture_partial`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/braintreetwo/cctypes
							ID:      `cctypes`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `AE,VI,MC,DI,JCB,CUP,DN,MI`,
						},

						&element.Field{
							// Path: payment/braintreetwo/useccv
							ID:      `useccv`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `["1","1"]`,
						},

						&element.Field{
							// Path: payment/braintreetwo/cctypes_braintree_mapper
							ID:      `cctypes_braintree_mapper`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"american-express":"AE","discover":"DI","jcb":"JCB","mastercard":"MC","master-card":"MC","visa":"VI","maestro":"MI","diners-club":"DN","unionpay":"CUP"}`,
						},

						&element.Field{
							// Path: payment/braintreetwo/order_status
							ID:      `order_status`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `processing`,
						},

						&element.Field{
							// Path: payment/braintreetwo/environment
							ID:      `environment`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `sandbox`,
						},

						&element.Field{
							// Path: payment/braintreetwo/allowspecific
							ID:      `allowspecific`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: payment/braintreetwo/sdk_url
							ID:      `sdk_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://js.braintreegateway.com/js/braintree-2.17.4.min.js`,
						},

						&element.Field{
							// Path: payment/braintreetwo/public_key
							ID:      `public_key`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						&element.Field{
							// Path: payment/braintreetwo/private_key
							ID:      `private_key`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"_value":null,"_attribute":{"backend_model":"Magento\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
						},

						&element.Field{
							// Path: payment/braintreetwo/masked_fields
							ID:      `masked_fields`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `cvv,number`,
						},

						&element.Field{
							// Path: payment/braintreetwo/privateInfoKeys
							ID:      `privateInfoKeys`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `avsPostalCodeResponseCode,avsStreetAddressResponseCode,cvvResponseCode,processorAuthorizationCode,processorResponseCode,processorResponseText,liabilityShifted,liabilityShiftPossible,riskDataId,riskDataDecision`,
						},

						&element.Field{
							// Path: payment/braintreetwo/paymentInfoKeys
							ID:      `paymentInfoKeys`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `cc_type,cc_number,avsPostalCodeResponseCode,avsStreetAddressResponseCode,cvvResponseCode,processorAuthorizationCode,processorResponseCode,processorResponseText,liabilityShifted,liabilityShiftPossible,riskDataId,riskDataDecision`,
						},

						&element.Field{
							// Path: payment/braintreetwo/can_use_internal
							ID:      `can_use_internal`,
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
