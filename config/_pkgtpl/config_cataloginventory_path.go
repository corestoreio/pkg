// +build ignore

package cataloginventory

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCataloginventoryOptionsCanSubtract => Decrease Stock When Order is Placed.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCataloginventoryOptionsCanSubtract = model.NewBool(`cataloginventory/options/can_subtract`)

// PathCataloginventoryOptionsCanBackInStock => Set Items' Status to be In Stock When Order is Cancelled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCataloginventoryOptionsCanBackInStock = model.NewBool(`cataloginventory/options/can_back_in_stock`)

// PathCataloginventoryOptionsShowOutOfStock => Display Out of Stock Products.
// Products will still be shown by direct product URLs.
// BackendModel: Otnegam\CatalogInventory\Model\Config\Backend\ShowOutOfStock
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCataloginventoryOptionsShowOutOfStock = model.NewBool(`cataloginventory/options/show_out_of_stock`)

// PathCataloginventoryOptionsStockThresholdQty => Only X left Threshold.
var PathCataloginventoryOptionsStockThresholdQty = model.NewStr(`cataloginventory/options/stock_threshold_qty`)

// PathCataloginventoryOptionsDisplayProductStockStatus => Display Products Availability in Stock on Storefront.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCataloginventoryOptionsDisplayProductStockStatus = model.NewBool(`cataloginventory/options/display_product_stock_status`)

// PathCataloginventoryItemOptionsManageStock => Manage Stock.
// Changing can take some time due to processing whole catalog.
// BackendModel: Otnegam\CatalogInventory\Model\Config\Backend\Managestock
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCataloginventoryItemOptionsManageStock = model.NewBool(`cataloginventory/item_options/manage_stock`)

// PathCataloginventoryItemOptionsBackorders => Backorders.
// Changing can take some time due to processing whole catalog.
// BackendModel: Otnegam\CatalogInventory\Model\Config\Backend\Backorders
// SourceModel: Otnegam\CatalogInventory\Model\Source\Backorders
var PathCataloginventoryItemOptionsBackorders = model.NewStr(`cataloginventory/item_options/backorders`)

// PathCataloginventoryItemOptionsMaxSaleQty => Maximum Qty Allowed in Shopping Cart.
var PathCataloginventoryItemOptionsMaxSaleQty = model.NewStr(`cataloginventory/item_options/max_sale_qty`)

// PathCataloginventoryItemOptionsMinQty => Out-of-Stock Threshold.
// BackendModel: Otnegam\CatalogInventory\Model\System\Config\Backend\Minqty
var PathCataloginventoryItemOptionsMinQty = model.NewStr(`cataloginventory/item_options/min_qty`)

// PathCataloginventoryItemOptionsMinSaleQty => Minimum Qty Allowed in Shopping Cart.
// BackendModel: Otnegam\CatalogInventory\Model\System\Config\Backend\Minsaleqty
var PathCataloginventoryItemOptionsMinSaleQty = model.NewStr(`cataloginventory/item_options/min_sale_qty`)

// PathCataloginventoryItemOptionsNotifyStockQty => Notify for Quantity Below.
var PathCataloginventoryItemOptionsNotifyStockQty = model.NewStr(`cataloginventory/item_options/notify_stock_qty`)

// PathCataloginventoryItemOptionsAutoReturn => Automatically Return Credit Memo Item to Stock.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCataloginventoryItemOptionsAutoReturn = model.NewBool(`cataloginventory/item_options/auto_return`)

// PathCataloginventoryItemOptionsEnableQtyIncrements => Enable Qty Increments.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCataloginventoryItemOptionsEnableQtyIncrements = model.NewBool(`cataloginventory/item_options/enable_qty_increments`)

// PathCataloginventoryItemOptionsQtyIncrements => Qty Increments.
// BackendModel: Otnegam\CatalogInventory\Model\System\Config\Backend\Qtyincrements
var PathCataloginventoryItemOptionsQtyIncrements = model.NewStr(`cataloginventory/item_options/qty_increments`)
