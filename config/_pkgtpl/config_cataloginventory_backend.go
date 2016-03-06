// +build ignore

package cataloginventory

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// CataloginventoryOptionsCanSubtract => Decrease Stock When Order is Placed.
	// Path: cataloginventory/options/can_subtract
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsCanSubtract model.Bool

	// CataloginventoryOptionsCanBackInStock => Set Items' Status to be In Stock When Order is Cancelled.
	// Path: cataloginventory/options/can_back_in_stock
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsCanBackInStock model.Bool

	// CataloginventoryOptionsShowOutOfStock => Display Out of Stock Products.
	// Products will still be shown by direct product URLs.
	// Path: cataloginventory/options/show_out_of_stock
	// BackendModel: Magento\CatalogInventory\Model\Config\Backend\ShowOutOfStock
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsShowOutOfStock model.Bool

	// CataloginventoryOptionsStockThresholdQty => Only X left Threshold.
	// Path: cataloginventory/options/stock_threshold_qty
	CataloginventoryOptionsStockThresholdQty model.Str

	// CataloginventoryOptionsDisplayProductStockStatus => Display Products Availability in Stock on Storefront.
	// Path: cataloginventory/options/display_product_stock_status
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsDisplayProductStockStatus model.Bool

	// CataloginventoryItemOptionsManageStock => Manage Stock.
	// Changing can take some time due to processing whole catalog.
	// Path: cataloginventory/item_options/manage_stock
	// BackendModel: Magento\CatalogInventory\Model\Config\Backend\Managestock
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryItemOptionsManageStock model.Bool

	// CataloginventoryItemOptionsBackorders => Backorders.
	// Changing can take some time due to processing whole catalog.
	// Path: cataloginventory/item_options/backorders
	// BackendModel: Magento\CatalogInventory\Model\Config\Backend\Backorders
	// SourceModel: Magento\CatalogInventory\Model\Source\Backorders
	CataloginventoryItemOptionsBackorders model.Str

	// CataloginventoryItemOptionsMaxSaleQty => Maximum Qty Allowed in Shopping Cart.
	// Path: cataloginventory/item_options/max_sale_qty
	CataloginventoryItemOptionsMaxSaleQty model.Str

	// CataloginventoryItemOptionsMinQty => Out-of-Stock Threshold.
	// Path: cataloginventory/item_options/min_qty
	// BackendModel: Magento\CatalogInventory\Model\System\Config\Backend\Minqty
	CataloginventoryItemOptionsMinQty model.Str

	// CataloginventoryItemOptionsMinSaleQty => Minimum Qty Allowed in Shopping Cart.
	// Path: cataloginventory/item_options/min_sale_qty
	// BackendModel: Magento\CatalogInventory\Model\System\Config\Backend\Minsaleqty
	CataloginventoryItemOptionsMinSaleQty model.Str

	// CataloginventoryItemOptionsNotifyStockQty => Notify for Quantity Below.
	// Path: cataloginventory/item_options/notify_stock_qty
	CataloginventoryItemOptionsNotifyStockQty model.Str

	// CataloginventoryItemOptionsAutoReturn => Automatically Return Credit Memo Item to Stock.
	// Path: cataloginventory/item_options/auto_return
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryItemOptionsAutoReturn model.Bool

	// CataloginventoryItemOptionsEnableQtyIncrements => Enable Qty Increments.
	// Path: cataloginventory/item_options/enable_qty_increments
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryItemOptionsEnableQtyIncrements model.Bool

	// CataloginventoryItemOptionsQtyIncrements => Qty Increments.
	// Path: cataloginventory/item_options/qty_increments
	// BackendModel: Magento\CatalogInventory\Model\System\Config\Backend\Qtyincrements
	CataloginventoryItemOptionsQtyIncrements model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CataloginventoryOptionsCanSubtract = model.NewBool(`cataloginventory/options/can_subtract`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryOptionsCanBackInStock = model.NewBool(`cataloginventory/options/can_back_in_stock`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryOptionsShowOutOfStock = model.NewBool(`cataloginventory/options/show_out_of_stock`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryOptionsStockThresholdQty = model.NewStr(`cataloginventory/options/stock_threshold_qty`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryOptionsDisplayProductStockStatus = model.NewBool(`cataloginventory/options/display_product_stock_status`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsManageStock = model.NewBool(`cataloginventory/item_options/manage_stock`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsBackorders = model.NewStr(`cataloginventory/item_options/backorders`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsMaxSaleQty = model.NewStr(`cataloginventory/item_options/max_sale_qty`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsMinQty = model.NewStr(`cataloginventory/item_options/min_qty`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsMinSaleQty = model.NewStr(`cataloginventory/item_options/min_sale_qty`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsNotifyStockQty = model.NewStr(`cataloginventory/item_options/notify_stock_qty`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsAutoReturn = model.NewBool(`cataloginventory/item_options/auto_return`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsEnableQtyIncrements = model.NewBool(`cataloginventory/item_options/enable_qty_increments`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsQtyIncrements = model.NewStr(`cataloginventory/item_options/qty_increments`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
