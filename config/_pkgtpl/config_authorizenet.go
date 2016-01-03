// +build ignore

package authorizenet

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
					ID:        "authorizenet_directpost",
					Label:     `Authorize.net Direct Post`,
					SortOrder: 34,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: payment/authorizenet_directpost/active
							ID:        "active",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/payment_action
							ID:        "payment_action",
							Label:     `Payment Action`,
							Type:      element.TypeSelect,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `authorize`,
							// SourceModel: Otnegam\Authorizenet\Model\Source\PaymentAction
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/title
							ID:        "title",
							Label:     `Title`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `Credit Card Direct Post (Authorize.net)`,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/login
							ID:        "login",
							Label:     `API Login ID`,
							Type:      element.TypeObscure,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   nil,
							// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/trans_key
							ID:        "trans_key",
							Label:     `Transaction Key`,
							Type:      element.TypeObscure,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   nil,
							// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/trans_md5
							ID:        "trans_md5",
							Label:     `Merchant MD5`,
							Type:      element.TypeObscure,
							SortOrder: 60,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   nil,
							// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted @todo Otnegam\Config\Model\Config\Backend\Encrypted
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/order_status
							ID:        "order_status",
							Label:     `New Order Status`,
							Type:      element.TypeSelect,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `processing`,
							// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\Processing
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/test
							ID:        "test",
							Label:     `Test Mode`,
							Type:      element.TypeSelect,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   true,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/cgi_url
							ID:        "cgi_url",
							Label:     `Gateway URL`,
							Type:      element.TypeText,
							SortOrder: 90,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `https://secure.authorize.net/gateway/transact.dll`,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/cgi_url_td
							ID:        "cgi_url_td",
							Label:     `Transaction Details URL`,
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `https://api2.authorize.net/xml/v1/request.api`,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/currency
							ID:        "currency",
							Label:     `Accepted Currency`,
							Type:      element.TypeSelect,
							SortOrder: 110,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `USD`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/debug
							ID:        "debug",
							Label:     `Debug`,
							Type:      element.TypeSelect,
							SortOrder: 120,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/email_customer
							ID:        "email_customer",
							Label:     `Email Customer`,
							Type:      element.TypeSelect,
							SortOrder: 130,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/merchant_email
							ID:        "merchant_email",
							Label:     `Merchant's Email`,
							Type:      element.TypeText,
							SortOrder: 140,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/cctypes
							ID:        "cctypes",
							Label:     `Credit Card Types`,
							Type:      element.TypeMultiselect,
							SortOrder: 150,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `AE,VI,MC,DI`,
							// SourceModel: Otnegam\Authorizenet\Model\Source\Cctype
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/useccv
							ID:        "useccv",
							Label:     `Credit Card Verification`,
							Type:      element.TypeSelect,
							SortOrder: 160,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/allowspecific
							ID:        "allowspecific",
							Label:     `Payment from Applicable Countries`,
							Type:      element.TypeAllowspecific,
							SortOrder: 170,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/specificcountry
							ID:        "specificcountry",
							Label:     `Payment from Specific Countries`,
							Type:      element.TypeMultiselect,
							SortOrder: 180,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Otnegam\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/min_order_total
							ID:        "min_order_total",
							Label:     `Minimum Order Total`,
							Type:      element.TypeText,
							SortOrder: 190,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/max_order_total
							ID:        "max_order_total",
							Label:     `Maximum Order Total`,
							Type:      element.TypeText,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/sort_order
							ID:        "sort_order",
							Label:     `Sort Order`,
							Type:      element.TypeText,
							SortOrder: 210,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "payment",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "authorizenet_directpost",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: payment/authorizenet_directpost/merchant_email
							ID:      `merchant_email`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/model
							ID:      `model`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Otnegam\Authorizenet\Model\Directpost`,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/create_order_before
							ID:      `create_order_before`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/date_delim
							ID:      `date_delim`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `/`,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/ccfields
							ID:      `ccfields`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `x_card_code,x_exp_date,x_card_num`,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/place_order_url
							ID:      `place_order_url`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `authorizenet/directpost_payment/place`,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/cgi_url_test_mode
							ID:      `cgi_url_test_mode`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://test.authorize.net/gateway/transact.dll`,
						},

						&element.Field{
							// Path: payment/authorizenet_directpost/cgi_url_td_test_mode
							ID:      `cgi_url_td_test_mode`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://apitest.authorize.net/xml/v1/request.api`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
