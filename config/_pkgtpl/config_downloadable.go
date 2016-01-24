// +build ignore

package downloadable

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
			ID: "catalog",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "downloadable",
					Label:     `Downloadable Product Options`,
					SortOrder: 600,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/downloadable/order_item_status
							ID:        "order_item_status",
							Label:     `Order Item Status to Enable Downloads`,
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   9,
							// SourceModel: Magento\Downloadable\Model\System\Config\Source\Orderitemstatus
						},

						&element.Field{
							// Path: catalog/downloadable/downloads_number
							ID:        "downloads_number",
							Label:     `Default Maximum Number of Downloads`,
							Type:      element.TypeText,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						},

						&element.Field{
							// Path: catalog/downloadable/shareable
							ID:        "shareable",
							Label:     `Shareable`,
							Type:      element.TypeSelect,
							SortOrder: 300,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/downloadable/samples_title
							ID:        "samples_title",
							Label:     `Default Sample Title`,
							Type:      element.TypeText,
							SortOrder: 400,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `Samples`,
						},

						&element.Field{
							// Path: catalog/downloadable/links_title
							ID:        "links_title",
							Label:     `Default Link Title`,
							Type:      element.TypeText,
							SortOrder: 500,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `Links`,
						},

						&element.Field{
							// Path: catalog/downloadable/links_target_new_window
							ID:        "links_target_new_window",
							Label:     `Open Links in New Window`,
							Type:      element.TypeSelect,
							SortOrder: 600,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/downloadable/content_disposition
							ID:        "content_disposition",
							Label:     `Use Content-Disposition`,
							Type:      element.TypeSelect,
							SortOrder: 700,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `inline`,
							// SourceModel: Magento\Downloadable\Model\System\Config\Source\Contentdisposition
						},

						&element.Field{
							// Path: catalog/downloadable/disable_guest_checkout
							ID:        "disable_guest_checkout",
							Label:     `Disable Guest Checkout if Cart Contains Downloadable Items`,
							Comment:   text.Long(`Guest checkout will only work with shareable.`),
							Type:      element.TypeSelect,
							SortOrder: 800,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
