// +build ignore

package msrp

import (
	"github.com/corestoreio/cspkg/config/element"
	"github.com/corestoreio/cspkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID: "sales",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "msrp",
					Label:     `Minimum Advertised Price`,
					SortOrder: 110,
					Scopes:    scope.PermWebsite,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: sales/msrp/enabled
							ID:        "enabled",
							Label:     `Enable MAP`,
							Comment:   text.Long(`<strong style="color:red">Warning!</strong> Enabling MAP by default will hide all product prices on Storefront.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: sales/msrp/display_price_type
							ID:        "display_price_type",
							Label:     `Display Actual Price`,
							Type:      element.TypeSelect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   true,
							// SourceModel: Magento\Msrp\Model\Product\Attribute\Source\Type
						},

						element.Field{
							// Path: sales/msrp/explanation_message
							ID:        "explanation_message",
							Label:     `Default Popup Text Message`,
							Type:      element.TypeTextarea,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Our price is lower than the manufacturer's "minimum advertised price." As a result, we cannot show you the price in catalog or the product page. <br /><br /> You have no obligation to purchase the product once you know the price. You can simply remove the item from your cart.`,
						},

						element.Field{
							// Path: sales/msrp/explanation_message_whats_this
							ID:        "explanation_message_whats_this",
							Label:     `Default "What's This" Text Message`,
							Type:      element.TypeTextarea,
							SortOrder: 50,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `Our price is lower than the manufacturer's "minimum advertised price." As a result, we cannot show you the price in catalog or the product page. <br /><br /> You have no obligation to purchase the product once you know the price. You can simply remove the item from your cart.`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
