// +build ignore

package cataloginventory

import "github.com/corestoreio/csfw/config"

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "cataloginventory",
		Label:     "Inventory",
		SortOrder: 50,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "options",
				Label:     `Stock Options`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `cataloginventory/options/can_subtract`,
						ID:           "can_subtract",
						Label:        `Decrease Stock When Order is Placed`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `cataloginventory/options/can_back_in_stock`,
						ID:           "can_back_in_stock",
						Label:        `Set Items' Status to be In Stock When Order is Cancelled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `cataloginventory/options/show_out_of_stock`,
						ID:           "show_out_of_stock",
						Label:        `Display Out of Stock Products`,
						Comment:      `Products will still be shown by direct product URLs.`,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      false,
						BackendModel: nil, // Magento\CatalogInventory\Model\Config\Backend\ShowOutOfStock
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `cataloginventory/options/stock_threshold_qty`,
						ID:           "stock_threshold_qty",
						Label:        `Only X left Threshold`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
						Default:      0,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `cataloginventory/options/display_product_stock_status`,
						ID:           "display_product_stock_status",
						Label:        `Display products availability in stock on Storefront.`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    50,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "item_options",
				Label:     `Product Stock Options`,
				Comment:   `Please note that these settings apply to individual items in the cart, not to the entire cart.`,
				SortOrder: 10,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `cataloginventory/item_options/manage_stock`,
						ID:           "manage_stock",
						Label:        `Manage Stock`,
						Comment:      `Changing can take some time due to processing whole catalog.`,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      true,
						BackendModel: nil, // Magento\CatalogInventory\Model\Config\Backend\Managestock
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `cataloginventory/item_options/backorders`,
						ID:           "backorders",
						Label:        `Backorders`,
						Comment:      `Changing can take some time due to processing whole catalog.`,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      false,
						BackendModel: nil, // Magento\CatalogInventory\Model\Config\Backend\Backorders
						SourceModel:  nil, // Magento\CatalogInventory\Model\Source\Backorders
					},

					&config.Field{
						// Path: `cataloginventory/item_options/max_sale_qty`,
						ID:           "max_sale_qty",
						Label:        `Maximum Qty Allowed in Shopping Cart`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      10000,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `cataloginventory/item_options/min_qty`,
						ID:           "min_qty",
						Label:        `Out-of-Stock Threshold`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      0,
						BackendModel: nil, // Magento\CatalogInventory\Model\System\Config\Backend\Minqty
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `cataloginventory/item_options/min_sale_qty`,
						ID:           "min_sale_qty",
						Label:        `Minimum Qty Allowed in Shopping Cart`,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    6,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      1,
						BackendModel: nil, // Magento\CatalogInventory\Model\System\Config\Backend\Minsaleqty
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `cataloginventory/item_options/notify_stock_qty`,
						ID:           "notify_stock_qty",
						Label:        `Notify for Quantity Below`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    7,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      1,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `cataloginventory/item_options/auto_return`,
						ID:           "auto_return",
						Label:        `Automatically Return Credit Memo Item to Stock`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `cataloginventory/item_options/enable_qty_increments`,
						ID:           "enable_qty_increments",
						Label:        `Enable Qty Increments`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    8,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `cataloginventory/item_options/qty_increments`,
						ID:           "qty_increments",
						Label:        `Qty Increments`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    9,
						Visible:      config.VisibleYes,
						Scope:        config.NewScopePerm(config.ScopeDefaultID),
						Default:      1,
						BackendModel: nil, // Magento\CatalogInventory\Model\System\Config\Backend\Qtyincrements
						SourceModel:  nil,
					},
				},
			},
		},
	},
)
