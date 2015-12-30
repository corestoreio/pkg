// +build ignore

package downloadable

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "catalog",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "downloadable",
				Label:     `Downloadable Product Options`,
				SortOrder: 600,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/downloadable/order_item_status
						ID:        "order_item_status",
						Label:     `Order Item Status to Enable Downloads`,
						Type:      config.TypeSelect,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   9,
						// SourceModel: Otnegam\Downloadable\Model\System\Config\Source\Orderitemstatus
					},

					&config.Field{
						// Path: catalog/downloadable/downloads_number
						ID:        "downloads_number",
						Label:     `Default Maximum Number of Downloads`,
						Type:      config.TypeText,
						SortOrder: 200,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&config.Field{
						// Path: catalog/downloadable/shareable
						ID:        "shareable",
						Label:     `Shareable`,
						Type:      config.TypeSelect,
						SortOrder: 300,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/downloadable/samples_title
						ID:        "samples_title",
						Label:     `Default Sample Title`,
						Type:      config.TypeText,
						SortOrder: 400,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Samples`,
					},

					&config.Field{
						// Path: catalog/downloadable/links_title
						ID:        "links_title",
						Label:     `Default Link Title`,
						Type:      config.TypeText,
						SortOrder: 500,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `Links`,
					},

					&config.Field{
						// Path: catalog/downloadable/links_target_new_window
						ID:        "links_target_new_window",
						Label:     `Open Links in New Window`,
						Type:      config.TypeSelect,
						SortOrder: 600,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/downloadable/content_disposition
						ID:        "content_disposition",
						Label:     `Use Content-Disposition`,
						Type:      config.TypeSelect,
						SortOrder: 700,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `inline`,
						// SourceModel: Otnegam\Downloadable\Model\System\Config\Source\Contentdisposition
					},

					&config.Field{
						// Path: catalog/downloadable/disable_guest_checkout
						ID:        "disable_guest_checkout",
						Label:     `Disable Guest Checkout if Cart Contains Downloadable Items`,
						Comment:   element.LongText(`Guest checkout will only work with shareable.`),
						Type:      config.TypeSelect,
						SortOrder: 800,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
)
