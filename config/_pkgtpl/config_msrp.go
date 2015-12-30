// +build ignore

package msrp

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "sales",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "msrp",
				Label:     `Minimum Advertised Price`,
				SortOrder: 110,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sales/msrp/enabled
						ID:        "enabled",
						Label:     `Enable MAP`,
						Comment:   element.LongText(`<strong style="color:red">Warning!</strong> Enabling MAP by default will hide all product prices on Storefront.`),
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: sales/msrp/display_price_type
						ID:        "display_price_type",
						Label:     `Display Actual Price`,
						Type:      config.TypeSelect,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Msrp\Model\Product\Attribute\Source\Type
					},

					&config.Field{
						// Path: sales/msrp/explanation_message
						ID:        "explanation_message",
						Label:     `Default Popup Text Message`,
						Type:      config.TypeTextarea,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Our price is lower than the manufacturer's "minimum advertised price." As a result, we cannot show you the price in catalog or the product page. <br /><br /> You have no obligation to purchase the product once you know the price. You can simply remove the item from your cart.`,
					},

					&config.Field{
						// Path: sales/msrp/explanation_message_whats_this
						ID:        "explanation_message_whats_this",
						Label:     `Default "What's This" Text Message`,
						Type:      config.TypeTextarea,
						SortOrder: 50,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Our price is lower than the manufacturer's "minimum advertised price." As a result, we cannot show you the price in catalog or the product page. <br /><br /> You have no obligation to purchase the product once you know the price. You can simply remove the item from your cart.`,
					},
				),
			},
		),
	},
)
