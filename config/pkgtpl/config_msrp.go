// +build ignore

package msrp

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "sales",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "msrp",
				Label:     `Minimum Advertised Price`,
				Comment:   ``,
				SortOrder: 110,
				Scope:     config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sales/msrp/enabled`,
						ID:           "enabled",
						Label:        `Enable MAP`,
						Comment:      `<strong style="color:red">Warning!</strong> Enabling MAP by default will hide all product prices on the front end.`,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sales/msrp/display_price_type`,
						ID:           "display_price_type",
						Label:        `Display Actual Price`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Msrp\Model\Product\Attribute\Source\Type
					},

					&config.Field{
						// Path: `sales/msrp/explanation_message`,
						ID:           "explanation_message",
						Label:        `Default Popup Text Message`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `Our price is lower than the manufacturer's "minimum advertised price." As a result, we cannot show you the price in catalog or the product page. <br /><br /> You have no obligation to purchase the product once you know the price. You can simply remove the item from your cart.`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sales/msrp/explanation_message_whats_this`,
						ID:           "explanation_message_whats_this",
						Label:        `Default "What's This" Text Message`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `Our price is lower than the manufacturer's "minimum advertised price." As a result, we cannot show you the price in catalog or the product page. <br /><br /> You have no obligation to purchase the product once you know the price. You can simply remove the item from your cart.`,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
)
