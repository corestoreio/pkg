// +build ignore

package centinel

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "payment_services",
		Label:     "Payment Services",
		SortOrder: 450,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "centinel",
				Label:     `3D Secure Credit Card Validation`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     scope.NewPerm(config.IDScopeDefault, config.IDScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `payment_services/centinel/processor_id`,
						ID:           "processor_id",
						Label:        `Processor ID`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment_services/centinel/merchant_id`,
						ID:           "merchant_id",
						Label:        `Merchant ID`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment_services/centinel/password`,
						ID:           "password",
						Label:        `Password`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `payment_services/centinel/test_mode`,
						ID:           "test_mode",
						Label:        `Test Mode`,
						Comment:      `This overrides any API URL that may be specified by a payment method.`,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `payment_services/centinel/debug`,
						ID:           "debug",
						Label:        `Debug Mode`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
)
