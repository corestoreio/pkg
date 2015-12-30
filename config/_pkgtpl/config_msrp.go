// +build ignore

package msrp

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID: "sales",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:        "msrp",
				Label:     `Minimum Advertised Price`,
				SortOrder: 110,
				Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: sales/msrp/enabled
						ID:        "enabled",
						Label:     `Enable MAP`,
						Comment:   element.LongText(`<strong style="color:red">Warning!</strong> Enabling MAP by default will hide all product prices on Storefront.`),
						Type:      element.TypeSelect,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: sales/msrp/display_price_type
						ID:        "display_price_type",
						Label:     `Display Actual Price`,
						Type:      element.TypeSelect,
						SortOrder: 30,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Msrp\Model\Product\Attribute\Source\Type
					},

					&element.Field{
						// Path: sales/msrp/explanation_message
						ID:        "explanation_message",
						Label:     `Default Popup Text Message`,
						Type:      element.TypeTextarea,
						SortOrder: 40,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Our price is lower than the manufacturer's "minimum advertised price." As a result, we cannot show you the price in catalog or the product page. <br /><br /> You have no obligation to purchase the product once you know the price. You can simply remove the item from your cart.`,
					},

					&element.Field{
						// Path: sales/msrp/explanation_message_whats_this
						ID:        "explanation_message_whats_this",
						Label:     `Default "What's This" Text Message`,
						Type:      element.TypeTextarea,
						SortOrder: 50,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Our price is lower than the manufacturer's "minimum advertised price." As a result, we cannot show you the price in catalog or the product page. <br /><br /> You have no obligation to purchase the product once you know the price. You can simply remove the item from your cart.`,
					},
				),
			},
		),
	},
)
