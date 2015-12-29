// +build ignore

package downloadable

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "downloadable",
				Label:     `Downloadable Product Options`,
				Comment:   ``,
				SortOrder: 600,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/downloadable/order_item_status`,
						ID:           "order_item_status",
						Label:        `Order Item Status to Enable Downloads`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      9,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Downloadable\Model\System\Config\Source\Orderitemstatus
					},

					&config.Field{
						// Path: `catalog/downloadable/downloads_number`,
						ID:           "downloads_number",
						Label:        `Default Maximum Number of Downloads`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    200,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      0,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/downloadable/shareable`,
						ID:           "shareable",
						Label:        `Shareable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    300,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/downloadable/samples_title`,
						ID:           "samples_title",
						Label:        `Default Sample Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    400,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Samples`,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/downloadable/links_title`,
						ID:           "links_title",
						Label:        `Default Link Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    500,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `Links`,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/downloadable/links_target_new_window`,
						ID:           "links_target_new_window",
						Label:        `Open Links in New Window`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    600,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/downloadable/content_disposition`,
						ID:           "content_disposition",
						Label:        `Use Content-Disposition`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    700,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `inline`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Downloadable\Model\System\Config\Source\Contentdisposition
					},

					&config.Field{
						// Path: `catalog/downloadable/disable_guest_checkout`,
						ID:           "disable_guest_checkout",
						Label:        `Disable Guest Checkout if Cart Contains Downloadable Items`,
						Comment:      `Guest checkout will only work with shareable.`,
						Type:         config.TypeSelect,
						SortOrder:    800,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:      true,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
)
