// +build ignore

package braintreetwo

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
				ID:        "braintreetwo_section",
				Label:     `BraintreeTwo`,
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
				ID: "braintreetwo",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: payment/braintreetwo/model
						ID:      `model`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `BraintreeTwoFacade`,
					},

					&config.Field{
						// Path: payment/braintreetwo/title
						ID:      `title`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `Credit Card (BraintreeTwo)`,
					},

					&config.Field{
						// Path: payment/braintreetwo/payment_action
						ID:      `payment_action`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `authorize`,
					},

					&config.Field{
						// Path: payment/braintreetwo/active
						ID:      `active`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintreetwo/is_gateway
						ID:      `is_gateway`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/braintreetwo/can_use_checkout
						ID:      `can_use_checkout`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/braintreetwo/can_authorize
						ID:      `can_authorize`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/braintreetwo/can_capture
						ID:      `can_capture`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/braintreetwo/can_capture_partial
						ID:      `can_capture_partial`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},

					&config.Field{
						// Path: payment/braintreetwo/cctypes
						ID:      `cctypes`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `AE,VI,MC,DI,JCB,CUP,DN,MI`,
					},

					&config.Field{
						// Path: payment/braintreetwo/useccv
						ID:      `useccv`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `["1","1"]`,
					},

					&config.Field{
						// Path: payment/braintreetwo/cctypes_braintree_mapper
						ID:      `cctypes_braintree_mapper`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"american-express":"AE","discover":"DI","jcb":"JCB","mastercard":"MC","master-card":"MC","visa":"VI","maestro":"MI","diners-club":"DN","unionpay":"CUP"}`,
					},

					&config.Field{
						// Path: payment/braintreetwo/order_status
						ID:      `order_status`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `processing`,
					},

					&config.Field{
						// Path: payment/braintreetwo/environment
						ID:      `environment`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `sandbox`,
					},

					&config.Field{
						// Path: payment/braintreetwo/allowspecific
						ID:      `allowspecific`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: payment/braintreetwo/sdk_url
						ID:      `sdk_url`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://js.braintreegateway.com/js/braintree-2.17.4.min.js`,
					},

					&config.Field{
						// Path: payment/braintreetwo/public_key
						ID:      `public_key`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: payment/braintreetwo/private_key
						ID:      `private_key`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"_value":null,"_attribute":{"backend_model":"Otnegam\\Config\\Model\\Config\\Backend\\Encrypted"}}`,
					},

					&config.Field{
						// Path: payment/braintreetwo/masked_fields
						ID:      `masked_fields`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `cvv,number`,
					},

					&config.Field{
						// Path: payment/braintreetwo/privateInfoKeys
						ID:      `privateInfoKeys`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `avsPostalCodeResponseCode,avsStreetAddressResponseCode,cvvResponseCode,processorAuthorizationCode,processorResponseCode,processorResponseText,liabilityShifted,liabilityShiftPossible,riskDataId,riskDataDecision`,
					},

					&config.Field{
						// Path: payment/braintreetwo/paymentInfoKeys
						ID:      `paymentInfoKeys`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `cc_type,cc_number,avsPostalCodeResponseCode,avsStreetAddressResponseCode,cvvResponseCode,processorAuthorizationCode,processorResponseCode,processorResponseText,liabilityShifted,liabilityShiftPossible,riskDataId,riskDataDecision`,
					},

					&config.Field{
						// Path: payment/braintreetwo/can_use_internal
						ID:      `can_use_internal`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},
				),
			},
		),
	},
)
