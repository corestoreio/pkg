// +build ignore

package giftmessage

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
			ID: "sales",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "gift_options",
					Label:     `Gift Options`,
					SortOrder: 100,
					Scope:     scope.PermWebsite,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/gift_options/allow_order
							ID:        "allow_order",
							Label:     `Allow Gift Messages on Order Level`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sales/gift_options/allow_items
							ID:        "allow_items",
							Label:     `Allow Gift Messages for Order Items`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "sales",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "gift_messages",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sales/gift_messages/allow_items
							ID:      `allow_items`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: sales/gift_messages/allow_order
							ID:      `allow_order`,
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
