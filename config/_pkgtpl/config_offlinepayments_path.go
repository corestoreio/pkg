// +build ignore

package offlinepayments

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathPaymentCheckmoActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentCheckmoActive = model.NewBool(`payment/checkmo/active`)

// PathPaymentCheckmoOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
var PathPaymentCheckmoOrderStatus = model.NewStr(`payment/checkmo/order_status`)

// PathPaymentCheckmoSortOrder => Sort Order.
var PathPaymentCheckmoSortOrder = model.NewStr(`payment/checkmo/sort_order`)

// PathPaymentCheckmoTitle => Title.
var PathPaymentCheckmoTitle = model.NewStr(`payment/checkmo/title`)

// PathPaymentCheckmoAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentCheckmoAllowspecific = model.NewStr(`payment/checkmo/allowspecific`)

// PathPaymentCheckmoSpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentCheckmoSpecificcountry = model.NewStringCSV(`payment/checkmo/specificcountry`)

// PathPaymentCheckmoPayableTo => Make Check Payable to.
var PathPaymentCheckmoPayableTo = model.NewStr(`payment/checkmo/payable_to`)

// PathPaymentCheckmoMailingAddress => Send Check to.
var PathPaymentCheckmoMailingAddress = model.NewStr(`payment/checkmo/mailing_address`)

// PathPaymentCheckmoMinOrderTotal => Minimum Order Total.
var PathPaymentCheckmoMinOrderTotal = model.NewStr(`payment/checkmo/min_order_total`)

// PathPaymentCheckmoMaxOrderTotal => Maximum Order Total.
var PathPaymentCheckmoMaxOrderTotal = model.NewStr(`payment/checkmo/max_order_total`)

// PathPaymentCheckmoModel => .
var PathPaymentCheckmoModel = model.NewStr(`payment/checkmo/model`)

// PathPaymentPurchaseorderActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentPurchaseorderActive = model.NewBool(`payment/purchaseorder/active`)

// PathPaymentPurchaseorderOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
var PathPaymentPurchaseorderOrderStatus = model.NewStr(`payment/purchaseorder/order_status`)

// PathPaymentPurchaseorderSortOrder => Sort Order.
var PathPaymentPurchaseorderSortOrder = model.NewStr(`payment/purchaseorder/sort_order`)

// PathPaymentPurchaseorderTitle => Title.
var PathPaymentPurchaseorderTitle = model.NewStr(`payment/purchaseorder/title`)

// PathPaymentPurchaseorderAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentPurchaseorderAllowspecific = model.NewStr(`payment/purchaseorder/allowspecific`)

// PathPaymentPurchaseorderSpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentPurchaseorderSpecificcountry = model.NewStringCSV(`payment/purchaseorder/specificcountry`)

// PathPaymentPurchaseorderMinOrderTotal => Minimum Order Total.
var PathPaymentPurchaseorderMinOrderTotal = model.NewStr(`payment/purchaseorder/min_order_total`)

// PathPaymentPurchaseorderMaxOrderTotal => Maximum Order Total.
var PathPaymentPurchaseorderMaxOrderTotal = model.NewStr(`payment/purchaseorder/max_order_total`)

// PathPaymentPurchaseorderModel => .
var PathPaymentPurchaseorderModel = model.NewStr(`payment/purchaseorder/model`)

// PathPaymentBanktransferActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentBanktransferActive = model.NewBool(`payment/banktransfer/active`)

// PathPaymentBanktransferTitle => Title.
var PathPaymentBanktransferTitle = model.NewStr(`payment/banktransfer/title`)

// PathPaymentBanktransferOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
var PathPaymentBanktransferOrderStatus = model.NewStr(`payment/banktransfer/order_status`)

// PathPaymentBanktransferAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentBanktransferAllowspecific = model.NewStr(`payment/banktransfer/allowspecific`)

// PathPaymentBanktransferSpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentBanktransferSpecificcountry = model.NewStringCSV(`payment/banktransfer/specificcountry`)

// PathPaymentBanktransferInstructions => Instructions.
var PathPaymentBanktransferInstructions = model.NewStr(`payment/banktransfer/instructions`)

// PathPaymentBanktransferMinOrderTotal => Minimum Order Total.
var PathPaymentBanktransferMinOrderTotal = model.NewStr(`payment/banktransfer/min_order_total`)

// PathPaymentBanktransferMaxOrderTotal => Maximum Order Total.
var PathPaymentBanktransferMaxOrderTotal = model.NewStr(`payment/banktransfer/max_order_total`)

// PathPaymentBanktransferSortOrder => Sort Order.
var PathPaymentBanktransferSortOrder = model.NewStr(`payment/banktransfer/sort_order`)

// PathPaymentCashondeliveryActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentCashondeliveryActive = model.NewBool(`payment/cashondelivery/active`)

// PathPaymentCashondeliveryTitle => Title.
var PathPaymentCashondeliveryTitle = model.NewStr(`payment/cashondelivery/title`)

// PathPaymentCashondeliveryOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\NewStatus
var PathPaymentCashondeliveryOrderStatus = model.NewStr(`payment/cashondelivery/order_status`)

// PathPaymentCashondeliveryAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentCashondeliveryAllowspecific = model.NewStr(`payment/cashondelivery/allowspecific`)

// PathPaymentCashondeliverySpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentCashondeliverySpecificcountry = model.NewStringCSV(`payment/cashondelivery/specificcountry`)

// PathPaymentCashondeliveryInstructions => Instructions.
var PathPaymentCashondeliveryInstructions = model.NewStr(`payment/cashondelivery/instructions`)

// PathPaymentCashondeliveryMinOrderTotal => Minimum Order Total.
var PathPaymentCashondeliveryMinOrderTotal = model.NewStr(`payment/cashondelivery/min_order_total`)

// PathPaymentCashondeliveryMaxOrderTotal => Maximum Order Total.
var PathPaymentCashondeliveryMaxOrderTotal = model.NewStr(`payment/cashondelivery/max_order_total`)

// PathPaymentCashondeliverySortOrder => Sort Order.
var PathPaymentCashondeliverySortOrder = model.NewStr(`payment/cashondelivery/sort_order`)

// PathPaymentFreeActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentFreeActive = model.NewBool(`payment/free/active`)

// PathPaymentFreeOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\Newprocessing
var PathPaymentFreeOrderStatus = model.NewStr(`payment/free/order_status`)

// PathPaymentFreePaymentAction => Automatically Invoice All Items.
// SourceModel: Otnegam\Payment\Model\Source\Invoice
var PathPaymentFreePaymentAction = model.NewStr(`payment/free/payment_action`)

// PathPaymentFreeSortOrder => Sort Order.
var PathPaymentFreeSortOrder = model.NewStr(`payment/free/sort_order`)

// PathPaymentFreeTitle => Title.
var PathPaymentFreeTitle = model.NewStr(`payment/free/title`)

// PathPaymentFreeAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentFreeAllowspecific = model.NewStr(`payment/free/allowspecific`)

// PathPaymentFreeSpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentFreeSpecificcountry = model.NewStringCSV(`payment/free/specificcountry`)

// PathPaymentFreeModel => .
var PathPaymentFreeModel = model.NewStr(`payment/free/model`)
