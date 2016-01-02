// +build ignore

package offlinepayments

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with PackageConfiguration.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// PaymentCheckmoActive => Enabled.
	// Path: payment/checkmo/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentCheckmoActive model.Bool

	// PaymentCheckmoOrderStatus => New Order Status.
	// Path: payment/checkmo/order_status
	// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
	PaymentCheckmoOrderStatus model.Str

	// PaymentCheckmoSortOrder => Sort Order.
	// Path: payment/checkmo/sort_order
	PaymentCheckmoSortOrder model.Str

	// PaymentCheckmoTitle => Title.
	// Path: payment/checkmo/title
	PaymentCheckmoTitle model.Str

	// PaymentCheckmoAllowspecific => Payment from Applicable Countries.
	// Path: payment/checkmo/allowspecific
	// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
	PaymentCheckmoAllowspecific model.Str

	// PaymentCheckmoSpecificcountry => Payment from Specific Countries.
	// Path: payment/checkmo/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	PaymentCheckmoSpecificcountry model.StringCSV

	// PaymentCheckmoPayableTo => Make Check Payable to.
	// Path: payment/checkmo/payable_to
	PaymentCheckmoPayableTo model.Str

	// PaymentCheckmoMailingAddress => Send Check to.
	// Path: payment/checkmo/mailing_address
	PaymentCheckmoMailingAddress model.Str

	// PaymentCheckmoMinOrderTotal => Minimum Order Total.
	// Path: payment/checkmo/min_order_total
	PaymentCheckmoMinOrderTotal model.Str

	// PaymentCheckmoMaxOrderTotal => Maximum Order Total.
	// Path: payment/checkmo/max_order_total
	PaymentCheckmoMaxOrderTotal model.Str

	// PaymentCheckmoModel => .
	// Path: payment/checkmo/model
	PaymentCheckmoModel model.Str

	// PaymentPurchaseorderActive => Enabled.
	// Path: payment/purchaseorder/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentPurchaseorderActive model.Bool

	// PaymentPurchaseorderOrderStatus => New Order Status.
	// Path: payment/purchaseorder/order_status
	// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
	PaymentPurchaseorderOrderStatus model.Str

	// PaymentPurchaseorderSortOrder => Sort Order.
	// Path: payment/purchaseorder/sort_order
	PaymentPurchaseorderSortOrder model.Str

	// PaymentPurchaseorderTitle => Title.
	// Path: payment/purchaseorder/title
	PaymentPurchaseorderTitle model.Str

	// PaymentPurchaseorderAllowspecific => Payment from Applicable Countries.
	// Path: payment/purchaseorder/allowspecific
	// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
	PaymentPurchaseorderAllowspecific model.Str

	// PaymentPurchaseorderSpecificcountry => Payment from Specific Countries.
	// Path: payment/purchaseorder/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	PaymentPurchaseorderSpecificcountry model.StringCSV

	// PaymentPurchaseorderMinOrderTotal => Minimum Order Total.
	// Path: payment/purchaseorder/min_order_total
	PaymentPurchaseorderMinOrderTotal model.Str

	// PaymentPurchaseorderMaxOrderTotal => Maximum Order Total.
	// Path: payment/purchaseorder/max_order_total
	PaymentPurchaseorderMaxOrderTotal model.Str

	// PaymentPurchaseorderModel => .
	// Path: payment/purchaseorder/model
	PaymentPurchaseorderModel model.Str

	// PaymentBanktransferActive => Enabled.
	// Path: payment/banktransfer/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentBanktransferActive model.Bool

	// PaymentBanktransferTitle => Title.
	// Path: payment/banktransfer/title
	PaymentBanktransferTitle model.Str

	// PaymentBanktransferOrderStatus => New Order Status.
	// Path: payment/banktransfer/order_status
	// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
	PaymentBanktransferOrderStatus model.Str

	// PaymentBanktransferAllowspecific => Payment from Applicable Countries.
	// Path: payment/banktransfer/allowspecific
	// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
	PaymentBanktransferAllowspecific model.Str

	// PaymentBanktransferSpecificcountry => Payment from Specific Countries.
	// Path: payment/banktransfer/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	PaymentBanktransferSpecificcountry model.StringCSV

	// PaymentBanktransferInstructions => Instructions.
	// Path: payment/banktransfer/instructions
	PaymentBanktransferInstructions model.Str

	// PaymentBanktransferMinOrderTotal => Minimum Order Total.
	// Path: payment/banktransfer/min_order_total
	PaymentBanktransferMinOrderTotal model.Str

	// PaymentBanktransferMaxOrderTotal => Maximum Order Total.
	// Path: payment/banktransfer/max_order_total
	PaymentBanktransferMaxOrderTotal model.Str

	// PaymentBanktransferSortOrder => Sort Order.
	// Path: payment/banktransfer/sort_order
	PaymentBanktransferSortOrder model.Str

	// PaymentCashondeliveryActive => Enabled.
	// Path: payment/cashondelivery/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentCashondeliveryActive model.Bool

	// PaymentCashondeliveryTitle => Title.
	// Path: payment/cashondelivery/title
	PaymentCashondeliveryTitle model.Str

	// PaymentCashondeliveryOrderStatus => New Order Status.
	// Path: payment/cashondelivery/order_status
	// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
	PaymentCashondeliveryOrderStatus model.Str

	// PaymentCashondeliveryAllowspecific => Payment from Applicable Countries.
	// Path: payment/cashondelivery/allowspecific
	// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
	PaymentCashondeliveryAllowspecific model.Str

	// PaymentCashondeliverySpecificcountry => Payment from Specific Countries.
	// Path: payment/cashondelivery/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	PaymentCashondeliverySpecificcountry model.StringCSV

	// PaymentCashondeliveryInstructions => Instructions.
	// Path: payment/cashondelivery/instructions
	PaymentCashondeliveryInstructions model.Str

	// PaymentCashondeliveryMinOrderTotal => Minimum Order Total.
	// Path: payment/cashondelivery/min_order_total
	PaymentCashondeliveryMinOrderTotal model.Str

	// PaymentCashondeliveryMaxOrderTotal => Maximum Order Total.
	// Path: payment/cashondelivery/max_order_total
	PaymentCashondeliveryMaxOrderTotal model.Str

	// PaymentCashondeliverySortOrder => Sort Order.
	// Path: payment/cashondelivery/sort_order
	PaymentCashondeliverySortOrder model.Str

	// PaymentFreeActive => Enabled.
	// Path: payment/free/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentFreeActive model.Bool

	// PaymentFreeOrderStatus => New Order Status.
	// Path: payment/free/order_status
	// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\Newprocessing
	PaymentFreeOrderStatus model.Str

	// PaymentFreePaymentAction => Automatically Invoice All Items.
	// Path: payment/free/payment_action
	// SourceModel: Otnegam\Payment\Model\Source\Invoice
	PaymentFreePaymentAction model.Str

	// PaymentFreeSortOrder => Sort Order.
	// Path: payment/free/sort_order
	PaymentFreeSortOrder model.Str

	// PaymentFreeTitle => Title.
	// Path: payment/free/title
	PaymentFreeTitle model.Str

	// PaymentFreeAllowspecific => Payment from Applicable Countries.
	// Path: payment/free/allowspecific
	// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
	PaymentFreeAllowspecific model.Str

	// PaymentFreeSpecificcountry => Payment from Specific Countries.
	// Path: payment/free/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	PaymentFreeSpecificcountry model.StringCSV

	// PaymentFreeModel => .
	// Path: payment/free/model
	PaymentFreeModel model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.PaymentCheckmoActive = model.NewBool(`payment/checkmo/active`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoOrderStatus = model.NewStr(`payment/checkmo/order_status`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoSortOrder = model.NewStr(`payment/checkmo/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoTitle = model.NewStr(`payment/checkmo/title`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoAllowspecific = model.NewStr(`payment/checkmo/allowspecific`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoSpecificcountry = model.NewStringCSV(`payment/checkmo/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoPayableTo = model.NewStr(`payment/checkmo/payable_to`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoMailingAddress = model.NewStr(`payment/checkmo/mailing_address`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoMinOrderTotal = model.NewStr(`payment/checkmo/min_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoMaxOrderTotal = model.NewStr(`payment/checkmo/max_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCheckmoModel = model.NewStr(`payment/checkmo/model`, model.WithPkgCfg(pkgCfg))
	pp.PaymentPurchaseorderActive = model.NewBool(`payment/purchaseorder/active`, model.WithPkgCfg(pkgCfg))
	pp.PaymentPurchaseorderOrderStatus = model.NewStr(`payment/purchaseorder/order_status`, model.WithPkgCfg(pkgCfg))
	pp.PaymentPurchaseorderSortOrder = model.NewStr(`payment/purchaseorder/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.PaymentPurchaseorderTitle = model.NewStr(`payment/purchaseorder/title`, model.WithPkgCfg(pkgCfg))
	pp.PaymentPurchaseorderAllowspecific = model.NewStr(`payment/purchaseorder/allowspecific`, model.WithPkgCfg(pkgCfg))
	pp.PaymentPurchaseorderSpecificcountry = model.NewStringCSV(`payment/purchaseorder/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.PaymentPurchaseorderMinOrderTotal = model.NewStr(`payment/purchaseorder/min_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentPurchaseorderMaxOrderTotal = model.NewStr(`payment/purchaseorder/max_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentPurchaseorderModel = model.NewStr(`payment/purchaseorder/model`, model.WithPkgCfg(pkgCfg))
	pp.PaymentBanktransferActive = model.NewBool(`payment/banktransfer/active`, model.WithPkgCfg(pkgCfg))
	pp.PaymentBanktransferTitle = model.NewStr(`payment/banktransfer/title`, model.WithPkgCfg(pkgCfg))
	pp.PaymentBanktransferOrderStatus = model.NewStr(`payment/banktransfer/order_status`, model.WithPkgCfg(pkgCfg))
	pp.PaymentBanktransferAllowspecific = model.NewStr(`payment/banktransfer/allowspecific`, model.WithPkgCfg(pkgCfg))
	pp.PaymentBanktransferSpecificcountry = model.NewStringCSV(`payment/banktransfer/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.PaymentBanktransferInstructions = model.NewStr(`payment/banktransfer/instructions`, model.WithPkgCfg(pkgCfg))
	pp.PaymentBanktransferMinOrderTotal = model.NewStr(`payment/banktransfer/min_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentBanktransferMaxOrderTotal = model.NewStr(`payment/banktransfer/max_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentBanktransferSortOrder = model.NewStr(`payment/banktransfer/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCashondeliveryActive = model.NewBool(`payment/cashondelivery/active`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCashondeliveryTitle = model.NewStr(`payment/cashondelivery/title`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCashondeliveryOrderStatus = model.NewStr(`payment/cashondelivery/order_status`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCashondeliveryAllowspecific = model.NewStr(`payment/cashondelivery/allowspecific`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCashondeliverySpecificcountry = model.NewStringCSV(`payment/cashondelivery/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCashondeliveryInstructions = model.NewStr(`payment/cashondelivery/instructions`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCashondeliveryMinOrderTotal = model.NewStr(`payment/cashondelivery/min_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCashondeliveryMaxOrderTotal = model.NewStr(`payment/cashondelivery/max_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentCashondeliverySortOrder = model.NewStr(`payment/cashondelivery/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.PaymentFreeActive = model.NewBool(`payment/free/active`, model.WithPkgCfg(pkgCfg))
	pp.PaymentFreeOrderStatus = model.NewStr(`payment/free/order_status`, model.WithPkgCfg(pkgCfg))
	pp.PaymentFreePaymentAction = model.NewStr(`payment/free/payment_action`, model.WithPkgCfg(pkgCfg))
	pp.PaymentFreeSortOrder = model.NewStr(`payment/free/sort_order`, model.WithPkgCfg(pkgCfg))
	pp.PaymentFreeTitle = model.NewStr(`payment/free/title`, model.WithPkgCfg(pkgCfg))
	pp.PaymentFreeAllowspecific = model.NewStr(`payment/free/allowspecific`, model.WithPkgCfg(pkgCfg))
	pp.PaymentFreeSpecificcountry = model.NewStringCSV(`payment/free/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.PaymentFreeModel = model.NewStr(`payment/free/model`, model.WithPkgCfg(pkgCfg))

	return pp
}
