// +build ignore

package authorizenet

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
				ID:        "authorizenet_directpost",
				Label:     `Authorize.net Direct Post`,
				Comment:   ``,
				SortOrder: 34,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/authorizenet_directpost/active`,
						ID:           "active",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/payment_action`,
						ID:           "payment_action",
						Label:        `Payment Action`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `authorize`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Authorizenet\Model\Source\PaymentAction
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/title`,
						ID:           "title",
						Label:        `Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Credit Card Direct Post (Authorize.net)`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/login`,
						ID:           "login",
						Label:        `API Login ID`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/trans_key`,
						ID:           "trans_key",
						Label:        `Transaction Key`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/trans_md5`,
						ID:           "trans_md5",
						Label:        `Merchant MD5`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    60,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/order_status`,
						ID:           "order_status",
						Label:        `New Order Status`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `processing`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sales\Model\Config\Source\Order\Status\Processing
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/test`,
						ID:           "test",
						Label:        `Test Mode`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/cgi_url`,
						ID:           "cgi_url",
						Label:        `Gateway URL`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `https://secure.authorize.net/gateway/transact.dll`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/cgi_url_td`,
						ID:           "cgi_url_td",
						Label:        `Transaction Details Url`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `https://api2.authorize.net/xml/v1/request.api`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/currency`,
						ID:           "currency",
						Label:        `Accepted Currency`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    110,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `USD`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Locale\Currency
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/debug`,
						ID:           "debug",
						Label:        `Debug`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    120,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/email_customer`,
						ID:           "email_customer",
						Label:        `Email Customer`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    130,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/merchant_email`,
						ID:           "merchant_email",
						Label:        `Merchant's Email`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    140,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/cctypes`,
						ID:           "cctypes",
						Label:        `Credit Card Types`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    150,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      `AE,VI,MC,DI`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Authorizenet\Model\Source\Cctype
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/useccv`,
						ID:           "useccv",
						Label:        `Credit Card Verification`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    160,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/allowspecific`,
						ID:           "allowspecific",
						Label:        `Payment from Applicable Countries`,
						Comment:      ``,
						Type:         config.TypeAllowspecific,
						SortOrder:    170,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      0,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Payment\Model\Config\Source\Allspecificcountries
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/specificcountry`,
						ID:           "specificcountry",
						Label:        `Payment from Specific Countries`,
						Comment:      ``,
						Type:         config.TypeMultiselect,
						SortOrder:    180,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/min_order_total`,
						ID:           "min_order_total",
						Label:        `Minimum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    190,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/max_order_total`,
						ID:           "max_order_total",
						Label:        `Maximum Order Total`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    200,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/sort_order`,
						ID:           "sort_order",
						Label:        `Sort Order`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    210,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
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
				ID: "authorizenet_directpost",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment/authorizenet_directpost/merchant_email`,
						ID:      "merchant_email",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: nil,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/model`,
						ID:      "model",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `Magento\Authorizenet\Model\Directpost`,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/create_order_before`,
						ID:      "create_order_before",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: true,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/date_delim`,
						ID:      "date_delim",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `/`,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/ccfields`,
						ID:      "ccfields",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `x_card_code,x_exp_date,x_card_num`,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/place_order_url`,
						ID:      "place_order_url",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `authorizenet/directpost_payment/place`,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/cgi_url_test_mode`,
						ID:      "cgi_url_test_mode",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `https://test.authorize.net/gateway/transact.dll`,
					},

					&config.Field{
						// Path: `payment/authorizenet_directpost/cgi_url_td_test_mode`,
						ID:      "cgi_url_td_test_mode",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `https://apitest.authorize.net/xml/v1/request.api`,
					},
				},
			},
		},
	},
)
