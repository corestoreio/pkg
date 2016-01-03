// +build ignore

package cataloginventory

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// CataloginventoryOptionsCanSubtract => Decrease Stock When Order is Placed.
	// Path: cataloginventory/options/can_subtract
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsCanSubtract model.Bool

	// CataloginventoryOptionsCanBackInStock => Set Items' Status to be In Stock When Order is Cancelled.
	// Path: cataloginventory/options/can_back_in_stock
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsCanBackInStock model.Bool

	// CataloginventoryOptionsShowOutOfStock => Display Out of Stock Products.
	// Products will still be shown by direct product URLs.
	// Path: cataloginventory/options/show_out_of_stock
	// BackendModel: Otnegam\CatalogInventory\Model\Config\Backend\ShowOutOfStock
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsShowOutOfStock model.Bool

	// CataloginventoryOptionsStockThresholdQty => Only X left Threshold.
	// Path: cataloginventory/options/stock_threshold_qty
	CataloginventoryOptionsStockThresholdQty model.Str

	// CataloginventoryOptionsDisplayProductStockStatus => Display Products Availability in Stock on Storefront.
	// Path: cataloginventory/options/display_product_stock_status
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsDisplayProductStockStatus model.Bool

	// CataloginventoryItemOptionsManageStock => Manage Stock.
	// Changing can take some time due to processing whole catalog.
	// Path: cataloginventory/item_options/manage_stock
	// BackendModel: Otnegam\CatalogInventory\Model\Config\Backend\Managestock
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CataloginventoryItemOptionsManageStock model.Bool

	// CataloginventoryItemOptionsBackorders => Backorders.
	// Changing can take some time due to processing whole catalog.
	// Path: cataloginventory/item_options/backorders
	// BackendModel: Otnegam\CatalogInventory\Model\Config\Backend\Backorders
	// SourceModel: Otnegam\CatalogInventory\Model\Source\Backorders
	CataloginventoryItemOptionsBackorders model.Str

	// CataloginventoryItemOptionsMaxSaleQty => Maximum Qty Allowed in Shopping Cart.
	// Path: cataloginventory/item_options/max_sale_qty
	CataloginventoryItemOptionsMaxSaleQty model.Str

	// CataloginventoryItemOptionsMinQty => Out-of-Stock Threshold.
	// Path: cataloginventory/item_options/min_qty
	// BackendModel: Otnegam\CatalogInventory\Model\System\Config\Backend\Minqty
	CataloginventoryItemOptionsMinQty model.Str

	// CataloginventoryItemOptionsMinSaleQty => Minimum Qty Allowed in Shopping Cart.
	// Path: cataloginventory/item_options/min_sale_qty
	// BackendModel: Otnegam\CatalogInventory\Model\System\Config\Backend\Minsaleqty
	CataloginventoryItemOptionsMinSaleQty model.Str

	// CataloginventoryItemOptionsNotifyStockQty => Notify for Quantity Below.
	// Path: cataloginventory/item_options/notify_stock_qty
	CataloginventoryItemOptionsNotifyStockQty model.Str

	// CataloginventoryItemOptionsAutoReturn => Automatically Return Credit Memo Item to Stock.
	// Path: cataloginventory/item_options/auto_return
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CataloginventoryItemOptionsAutoReturn model.Bool

	// CataloginventoryItemOptionsEnableQtyIncrements => Enable Qty Increments.
	// Path: cataloginventory/item_options/enable_qty_increments
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CataloginventoryItemOptionsEnableQtyIncrements model.Bool

	// CataloginventoryItemOptionsQtyIncrements => Qty Increments.
	// Path: cataloginventory/item_options/qty_increments
	// BackendModel: Otnegam\CatalogInventory\Model\System\Config\Backend\Qtyincrements
	CataloginventoryItemOptionsQtyIncrements model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CataloginventoryOptionsCanSubtract = model.NewBool(`cataloginventory/options/can_subtract`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryOptionsCanBackInStock = model.NewBool(`cataloginventory/options/can_back_in_stock`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryOptionsShowOutOfStock = model.NewBool(`cataloginventory/options/show_out_of_stock`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryOptionsStockThresholdQty = model.NewStr(`cataloginventory/options/stock_threshold_qty`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryOptionsDisplayProductStockStatus = model.NewBool(`cataloginventory/options/display_product_stock_status`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryItemOptionsManageStock = model.NewBool(`cataloginventory/item_options/manage_stock`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryItemOptionsBackorders = model.NewStr(`cataloginventory/item_options/backorders`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryItemOptionsMaxSaleQty = model.NewStr(`cataloginventory/item_options/max_sale_qty`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryItemOptionsMinQty = model.NewStr(`cataloginventory/item_options/min_qty`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryItemOptionsMinSaleQty = model.NewStr(`cataloginventory/item_options/min_sale_qty`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryItemOptionsNotifyStockQty = model.NewStr(`cataloginventory/item_options/notify_stock_qty`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryItemOptionsAutoReturn = model.NewBool(`cataloginventory/item_options/auto_return`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryItemOptionsEnableQtyIncrements = model.NewBool(`cataloginventory/item_options/enable_qty_increments`, model.WithConfigStructure(cfgStruct))
	pp.CataloginventoryItemOptionsQtyIncrements = model.NewStr(`cataloginventory/item_options/qty_increments`, model.WithConfigStructure(cfgStruct))

	return pp
}
