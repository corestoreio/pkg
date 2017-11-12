// +build ignore

package cataloginventory

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CataloginventoryOptionsCanSubtract => Decrease Stock When Order is Placed.
	// Path: cataloginventory/options/can_subtract
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsCanSubtract cfgmodel.Bool

	// CataloginventoryOptionsCanBackInStock => Set Items' Status to be In Stock When Order is Cancelled.
	// Path: cataloginventory/options/can_back_in_stock
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsCanBackInStock cfgmodel.Bool

	// CataloginventoryOptionsShowOutOfStock => Display Out of Stock Products.
	// Products will still be shown by direct product URLs.
	// Path: cataloginventory/options/show_out_of_stock
	// BackendModel: Magento\CatalogInventory\Model\Config\Backend\ShowOutOfStock
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsShowOutOfStock cfgmodel.Bool

	// CataloginventoryOptionsStockThresholdQty => Only X left Threshold.
	// Path: cataloginventory/options/stock_threshold_qty
	CataloginventoryOptionsStockThresholdQty cfgmodel.Str

	// CataloginventoryOptionsDisplayProductStockStatus => Display Products Availability in Stock on Storefront.
	// Path: cataloginventory/options/display_product_stock_status
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryOptionsDisplayProductStockStatus cfgmodel.Bool

	// CataloginventoryItemOptionsManageStock => Manage Stock.
	// Changing can take some time due to processing whole catalog.
	// Path: cataloginventory/item_options/manage_stock
	// BackendModel: Magento\CatalogInventory\Model\Config\Backend\Managestock
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryItemOptionsManageStock cfgmodel.Bool

	// CataloginventoryItemOptionsBackorders => Backorders.
	// Changing can take some time due to processing whole catalog.
	// Path: cataloginventory/item_options/backorders
	// BackendModel: Magento\CatalogInventory\Model\Config\Backend\Backorders
	// SourceModel: Magento\CatalogInventory\Model\Source\Backorders
	CataloginventoryItemOptionsBackorders cfgmodel.Str

	// CataloginventoryItemOptionsMaxSaleQty => Maximum Qty Allowed in Shopping Cart.
	// Path: cataloginventory/item_options/max_sale_qty
	CataloginventoryItemOptionsMaxSaleQty cfgmodel.Str

	// CataloginventoryItemOptionsMinQty => Out-of-Stock Threshold.
	// Path: cataloginventory/item_options/min_qty
	// BackendModel: Magento\CatalogInventory\Model\System\Config\Backend\Minqty
	CataloginventoryItemOptionsMinQty cfgmodel.Str

	// CataloginventoryItemOptionsMinSaleQty => Minimum Qty Allowed in Shopping Cart.
	// Path: cataloginventory/item_options/min_sale_qty
	// BackendModel: Magento\CatalogInventory\Model\System\Config\Backend\Minsaleqty
	CataloginventoryItemOptionsMinSaleQty cfgmodel.Str

	// CataloginventoryItemOptionsNotifyStockQty => Notify for Quantity Below.
	// Path: cataloginventory/item_options/notify_stock_qty
	CataloginventoryItemOptionsNotifyStockQty cfgmodel.Str

	// CataloginventoryItemOptionsAutoReturn => Automatically Return Credit Memo Item to Stock.
	// Path: cataloginventory/item_options/auto_return
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryItemOptionsAutoReturn cfgmodel.Bool

	// CataloginventoryItemOptionsEnableQtyIncrements => Enable Qty Increments.
	// Path: cataloginventory/item_options/enable_qty_increments
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CataloginventoryItemOptionsEnableQtyIncrements cfgmodel.Bool

	// CataloginventoryItemOptionsQtyIncrements => Qty Increments.
	// Path: cataloginventory/item_options/qty_increments
	// BackendModel: Magento\CatalogInventory\Model\System\Config\Backend\Qtyincrements
	CataloginventoryItemOptionsQtyIncrements cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CataloginventoryOptionsCanSubtract = cfgmodel.NewBool(`cataloginventory/options/can_subtract`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryOptionsCanBackInStock = cfgmodel.NewBool(`cataloginventory/options/can_back_in_stock`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryOptionsShowOutOfStock = cfgmodel.NewBool(`cataloginventory/options/show_out_of_stock`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryOptionsStockThresholdQty = cfgmodel.NewStr(`cataloginventory/options/stock_threshold_qty`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryOptionsDisplayProductStockStatus = cfgmodel.NewBool(`cataloginventory/options/display_product_stock_status`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsManageStock = cfgmodel.NewBool(`cataloginventory/item_options/manage_stock`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsBackorders = cfgmodel.NewStr(`cataloginventory/item_options/backorders`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsMaxSaleQty = cfgmodel.NewStr(`cataloginventory/item_options/max_sale_qty`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsMinQty = cfgmodel.NewStr(`cataloginventory/item_options/min_qty`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsMinSaleQty = cfgmodel.NewStr(`cataloginventory/item_options/min_sale_qty`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsNotifyStockQty = cfgmodel.NewStr(`cataloginventory/item_options/notify_stock_qty`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsAutoReturn = cfgmodel.NewBool(`cataloginventory/item_options/auto_return`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsEnableQtyIncrements = cfgmodel.NewBool(`cataloginventory/item_options/enable_qty_increments`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CataloginventoryItemOptionsQtyIncrements = cfgmodel.NewStr(`cataloginventory/item_options/qty_increments`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
