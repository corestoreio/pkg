// +build ignore

package offlinepayments

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathPaymentCheckmoActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentCheckmoActive = model.NewBool(`payment/checkmo/active`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
var PathPaymentCheckmoOrderStatus = model.NewStr(`payment/checkmo/order_status`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoSortOrder => Sort Order.
var PathPaymentCheckmoSortOrder = model.NewStr(`payment/checkmo/sort_order`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoTitle => Title.
var PathPaymentCheckmoTitle = model.NewStr(`payment/checkmo/title`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentCheckmoAllowspecific = model.NewStr(`payment/checkmo/allowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoSpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentCheckmoSpecificcountry = model.NewStringCSV(`payment/checkmo/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoPayableTo => Make Check Payable to.
var PathPaymentCheckmoPayableTo = model.NewStr(`payment/checkmo/payable_to`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoMailingAddress => Send Check to.
var PathPaymentCheckmoMailingAddress = model.NewStr(`payment/checkmo/mailing_address`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoMinOrderTotal => Minimum Order Total.
var PathPaymentCheckmoMinOrderTotal = model.NewStr(`payment/checkmo/min_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoMaxOrderTotal => Maximum Order Total.
var PathPaymentCheckmoMaxOrderTotal = model.NewStr(`payment/checkmo/max_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCheckmoModel => .
var PathPaymentCheckmoModel = model.NewStr(`payment/checkmo/model`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentPurchaseorderActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentPurchaseorderActive = model.NewBool(`payment/purchaseorder/active`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentPurchaseorderOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
var PathPaymentPurchaseorderOrderStatus = model.NewStr(`payment/purchaseorder/order_status`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentPurchaseorderSortOrder => Sort Order.
var PathPaymentPurchaseorderSortOrder = model.NewStr(`payment/purchaseorder/sort_order`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentPurchaseorderTitle => Title.
var PathPaymentPurchaseorderTitle = model.NewStr(`payment/purchaseorder/title`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentPurchaseorderAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentPurchaseorderAllowspecific = model.NewStr(`payment/purchaseorder/allowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentPurchaseorderSpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentPurchaseorderSpecificcountry = model.NewStringCSV(`payment/purchaseorder/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentPurchaseorderMinOrderTotal => Minimum Order Total.
var PathPaymentPurchaseorderMinOrderTotal = model.NewStr(`payment/purchaseorder/min_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentPurchaseorderMaxOrderTotal => Maximum Order Total.
var PathPaymentPurchaseorderMaxOrderTotal = model.NewStr(`payment/purchaseorder/max_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentPurchaseorderModel => .
var PathPaymentPurchaseorderModel = model.NewStr(`payment/purchaseorder/model`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentBanktransferActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentBanktransferActive = model.NewBool(`payment/banktransfer/active`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentBanktransferTitle => Title.
var PathPaymentBanktransferTitle = model.NewStr(`payment/banktransfer/title`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentBanktransferOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
var PathPaymentBanktransferOrderStatus = model.NewStr(`payment/banktransfer/order_status`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentBanktransferAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentBanktransferAllowspecific = model.NewStr(`payment/banktransfer/allowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentBanktransferSpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentBanktransferSpecificcountry = model.NewStringCSV(`payment/banktransfer/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentBanktransferInstructions => Instructions.
var PathPaymentBanktransferInstructions = model.NewStr(`payment/banktransfer/instructions`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentBanktransferMinOrderTotal => Minimum Order Total.
var PathPaymentBanktransferMinOrderTotal = model.NewStr(`payment/banktransfer/min_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentBanktransferMaxOrderTotal => Maximum Order Total.
var PathPaymentBanktransferMaxOrderTotal = model.NewStr(`payment/banktransfer/max_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentBanktransferSortOrder => Sort Order.
var PathPaymentBanktransferSortOrder = model.NewStr(`payment/banktransfer/sort_order`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCashondeliveryActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentCashondeliveryActive = model.NewBool(`payment/cashondelivery/active`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCashondeliveryTitle => Title.
var PathPaymentCashondeliveryTitle = model.NewStr(`payment/cashondelivery/title`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCashondeliveryOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
var PathPaymentCashondeliveryOrderStatus = model.NewStr(`payment/cashondelivery/order_status`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCashondeliveryAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentCashondeliveryAllowspecific = model.NewStr(`payment/cashondelivery/allowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCashondeliverySpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentCashondeliverySpecificcountry = model.NewStringCSV(`payment/cashondelivery/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCashondeliveryInstructions => Instructions.
var PathPaymentCashondeliveryInstructions = model.NewStr(`payment/cashondelivery/instructions`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCashondeliveryMinOrderTotal => Minimum Order Total.
var PathPaymentCashondeliveryMinOrderTotal = model.NewStr(`payment/cashondelivery/min_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCashondeliveryMaxOrderTotal => Maximum Order Total.
var PathPaymentCashondeliveryMaxOrderTotal = model.NewStr(`payment/cashondelivery/max_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentCashondeliverySortOrder => Sort Order.
var PathPaymentCashondeliverySortOrder = model.NewStr(`payment/cashondelivery/sort_order`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentFreeActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentFreeActive = model.NewBool(`payment/free/active`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentFreeOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\Newprocessing
var PathPaymentFreeOrderStatus = model.NewStr(`payment/free/order_status`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentFreePaymentAction => Automatically Invoice All Items.
// SourceModel: Otnegam\Payment\Model\Source\Invoice
var PathPaymentFreePaymentAction = model.NewStr(`payment/free/payment_action`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentFreeSortOrder => Sort Order.
var PathPaymentFreeSortOrder = model.NewStr(`payment/free/sort_order`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentFreeTitle => Title.
var PathPaymentFreeTitle = model.NewStr(`payment/free/title`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentFreeAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentFreeAllowspecific = model.NewStr(`payment/free/allowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentFreeSpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentFreeSpecificcountry = model.NewStringCSV(`payment/free/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentFreeModel => .
var PathPaymentFreeModel = model.NewStr(`payment/free/model`, model.WithPkgCfg(PackageConfiguration))
