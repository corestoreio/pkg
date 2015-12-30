// +build ignore

package cataloginventory

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID:        "cataloginventory",
		Label:     `Inventory`,
		SortOrder: 50,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_CatalogInventory::cataloginventory
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:        "options",
				Label:     `Stock Options`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: cataloginventory/options/can_subtract
						ID:        "can_subtract",
						Label:     `Decrease Stock When Order is Placed`,
						Type:      element.TypeSelect,
						SortOrder: 2,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: cataloginventory/options/can_back_in_stock
						ID:        "can_back_in_stock",
						Label:     `Set Items' Status to be In Stock When Order is Cancelled`,
						Type:      element.TypeSelect,
						SortOrder: 2,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: cataloginventory/options/show_out_of_stock
						ID:        "show_out_of_stock",
						Label:     `Display Out of Stock Products`,
						Comment:   element.LongText(`Products will still be shown by direct product URLs.`),
						Type:      element.TypeSelect,
						SortOrder: 3,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   false,
						// BackendModel: Otnegam\CatalogInventory\Model\Config\Backend\ShowOutOfStock
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: cataloginventory/options/stock_threshold_qty
						ID:        "stock_threshold_qty",
						Label:     `Only X left Threshold`,
						Type:      element.TypeText,
						SortOrder: 4,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
					},

					&element.Field{
						// Path: cataloginventory/options/display_product_stock_status
						ID:        "display_product_stock_status",
						Label:     `Display Products Availability in Stock on Storefront`,
						Type:      element.TypeSelect,
						SortOrder: 50,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&element.Group{
				ID:        "item_options",
				Label:     `Product Stock Options`,
				Comment:   element.LongText(`Please note that these settings apply to individual items in the cart, not to the entire cart.`),
				SortOrder: 10,
				Scope:     scope.PermAll,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: cataloginventory/item_options/manage_stock
						ID:        "manage_stock",
						Label:     `Manage Stock`,
						Comment:   element.LongText(`Changing can take some time due to processing whole catalog.`),
						Type:      element.TypeSelect,
						SortOrder: 1,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   true,
						// BackendModel: Otnegam\CatalogInventory\Model\Config\Backend\Managestock
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: cataloginventory/item_options/backorders
						ID:        "backorders",
						Label:     `Backorders`,
						Comment:   element.LongText(`Changing can take some time due to processing whole catalog.`),
						Type:      element.TypeSelect,
						SortOrder: 3,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   false,
						// BackendModel: Otnegam\CatalogInventory\Model\Config\Backend\Backorders
						// SourceModel: Otnegam\CatalogInventory\Model\Source\Backorders
					},

					&element.Field{
						// Path: cataloginventory/item_options/max_sale_qty
						ID:        "max_sale_qty",
						Label:     `Maximum Qty Allowed in Shopping Cart`,
						Type:      element.TypeText,
						SortOrder: 4,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   10000,
					},

					&element.Field{
						// Path: cataloginventory/item_options/min_qty
						ID:        "min_qty",
						Label:     `Out-of-Stock Threshold`,
						Type:      element.TypeText,
						SortOrder: 5,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\CatalogInventory\Model\System\Config\Backend\Minqty
					},

					&element.Field{
						// Path: cataloginventory/item_options/min_sale_qty
						ID:        "min_sale_qty",
						Label:     `Minimum Qty Allowed in Shopping Cart`,
						Type:      element.Type,
						SortOrder: 6,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   1,
						// BackendModel: Otnegam\CatalogInventory\Model\System\Config\Backend\Minsaleqty
					},

					&element.Field{
						// Path: cataloginventory/item_options/notify_stock_qty
						ID:        "notify_stock_qty",
						Label:     `Notify for Quantity Below`,
						Type:      element.TypeText,
						SortOrder: 7,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   1,
					},

					&element.Field{
						// Path: cataloginventory/item_options/auto_return
						ID:        "auto_return",
						Label:     `Automatically Return Credit Memo Item to Stock`,
						Type:      element.TypeSelect,
						SortOrder: 10,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: cataloginventory/item_options/enable_qty_increments
						ID:        "enable_qty_increments",
						Label:     `Enable Qty Increments`,
						Type:      element.TypeSelect,
						SortOrder: 8,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: cataloginventory/item_options/qty_increments
						ID:        "qty_increments",
						Label:     `Qty Increments`,
						Type:      element.TypeText,
						SortOrder: 9,
						Visible:   element.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   1,
						// BackendModel: Otnegam\CatalogInventory\Model\System\Config\Backend\Qtyincrements
					},
				),
			},
		),
	},
)
