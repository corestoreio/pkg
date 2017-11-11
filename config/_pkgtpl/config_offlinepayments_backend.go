// +build ignore

package offlinepayments

import (
	"github.com/corestoreio/cspkg/config/cfgmodel"
	"github.com/corestoreio/cspkg/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// PaymentCheckmoActive => Enabled.
	// Path: payment/checkmo/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentCheckmoActive cfgmodel.Bool

	// PaymentCheckmoOrderStatus => New Order Status.
	// Path: payment/checkmo/order_status
	// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\NewStatus
	PaymentCheckmoOrderStatus cfgmodel.Str

	// PaymentCheckmoSortOrder => Sort Order.
	// Path: payment/checkmo/sort_order
	PaymentCheckmoSortOrder cfgmodel.Str

	// PaymentCheckmoTitle => Title.
	// Path: payment/checkmo/title
	PaymentCheckmoTitle cfgmodel.Str

	// PaymentCheckmoAllowspecific => Payment from Applicable Countries.
	// Path: payment/checkmo/allowspecific
	// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
	PaymentCheckmoAllowspecific cfgmodel.Str

	// PaymentCheckmoSpecificcountry => Payment from Specific Countries.
	// Path: payment/checkmo/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	PaymentCheckmoSpecificcountry cfgmodel.StringCSV

	// PaymentCheckmoPayableTo => Make Check Payable to.
	// Path: payment/checkmo/payable_to
	PaymentCheckmoPayableTo cfgmodel.Str

	// PaymentCheckmoMailingAddress => Send Check to.
	// Path: payment/checkmo/mailing_address
	PaymentCheckmoMailingAddress cfgmodel.Str

	// PaymentCheckmoMinOrderTotal => Minimum Order Total.
	// Path: payment/checkmo/min_order_total
	PaymentCheckmoMinOrderTotal cfgmodel.Str

	// PaymentCheckmoMaxOrderTotal => Maximum Order Total.
	// Path: payment/checkmo/max_order_total
	PaymentCheckmoMaxOrderTotal cfgmodel.Str

	// PaymentCheckmoModel => .
	// Path: payment/checkmo/cfgmodel
	PaymentCheckmoModel cfgmodel.Str

	// PaymentPurchaseorderActive => Enabled.
	// Path: payment/purchaseorder/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentPurchaseorderActive cfgmodel.Bool

	// PaymentPurchaseorderOrderStatus => New Order Status.
	// Path: payment/purchaseorder/order_status
	// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\NewStatus
	PaymentPurchaseorderOrderStatus cfgmodel.Str

	// PaymentPurchaseorderSortOrder => Sort Order.
	// Path: payment/purchaseorder/sort_order
	PaymentPurchaseorderSortOrder cfgmodel.Str

	// PaymentPurchaseorderTitle => Title.
	// Path: payment/purchaseorder/title
	PaymentPurchaseorderTitle cfgmodel.Str

	// PaymentPurchaseorderAllowspecific => Payment from Applicable Countries.
	// Path: payment/purchaseorder/allowspecific
	// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
	PaymentPurchaseorderAllowspecific cfgmodel.Str

	// PaymentPurchaseorderSpecificcountry => Payment from Specific Countries.
	// Path: payment/purchaseorder/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	PaymentPurchaseorderSpecificcountry cfgmodel.StringCSV

	// PaymentPurchaseorderMinOrderTotal => Minimum Order Total.
	// Path: payment/purchaseorder/min_order_total
	PaymentPurchaseorderMinOrderTotal cfgmodel.Str

	// PaymentPurchaseorderMaxOrderTotal => Maximum Order Total.
	// Path: payment/purchaseorder/max_order_total
	PaymentPurchaseorderMaxOrderTotal cfgmodel.Str

	// PaymentPurchaseorderModel => .
	// Path: payment/purchaseorder/cfgmodel
	PaymentPurchaseorderModel cfgmodel.Str

	// PaymentBanktransferActive => Enabled.
	// Path: payment/banktransfer/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentBanktransferActive cfgmodel.Bool

	// PaymentBanktransferTitle => Title.
	// Path: payment/banktransfer/title
	PaymentBanktransferTitle cfgmodel.Str

	// PaymentBanktransferOrderStatus => New Order Status.
	// Path: payment/banktransfer/order_status
	// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\NewStatus
	PaymentBanktransferOrderStatus cfgmodel.Str

	// PaymentBanktransferAllowspecific => Payment from Applicable Countries.
	// Path: payment/banktransfer/allowspecific
	// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
	PaymentBanktransferAllowspecific cfgmodel.Str

	// PaymentBanktransferSpecificcountry => Payment from Specific Countries.
	// Path: payment/banktransfer/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	PaymentBanktransferSpecificcountry cfgmodel.StringCSV

	// PaymentBanktransferInstructions => Instructions.
	// Path: payment/banktransfer/instructions
	PaymentBanktransferInstructions cfgmodel.Str

	// PaymentBanktransferMinOrderTotal => Minimum Order Total.
	// Path: payment/banktransfer/min_order_total
	PaymentBanktransferMinOrderTotal cfgmodel.Str

	// PaymentBanktransferMaxOrderTotal => Maximum Order Total.
	// Path: payment/banktransfer/max_order_total
	PaymentBanktransferMaxOrderTotal cfgmodel.Str

	// PaymentBanktransferSortOrder => Sort Order.
	// Path: payment/banktransfer/sort_order
	PaymentBanktransferSortOrder cfgmodel.Str

	// PaymentCashondeliveryActive => Enabled.
	// Path: payment/cashondelivery/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentCashondeliveryActive cfgmodel.Bool

	// PaymentCashondeliveryTitle => Title.
	// Path: payment/cashondelivery/title
	PaymentCashondeliveryTitle cfgmodel.Str

	// PaymentCashondeliveryOrderStatus => New Order Status.
	// Path: payment/cashondelivery/order_status
	// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\NewStatus
	PaymentCashondeliveryOrderStatus cfgmodel.Str

	// PaymentCashondeliveryAllowspecific => Payment from Applicable Countries.
	// Path: payment/cashondelivery/allowspecific
	// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
	PaymentCashondeliveryAllowspecific cfgmodel.Str

	// PaymentCashondeliverySpecificcountry => Payment from Specific Countries.
	// Path: payment/cashondelivery/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	PaymentCashondeliverySpecificcountry cfgmodel.StringCSV

	// PaymentCashondeliveryInstructions => Instructions.
	// Path: payment/cashondelivery/instructions
	PaymentCashondeliveryInstructions cfgmodel.Str

	// PaymentCashondeliveryMinOrderTotal => Minimum Order Total.
	// Path: payment/cashondelivery/min_order_total
	PaymentCashondeliveryMinOrderTotal cfgmodel.Str

	// PaymentCashondeliveryMaxOrderTotal => Maximum Order Total.
	// Path: payment/cashondelivery/max_order_total
	PaymentCashondeliveryMaxOrderTotal cfgmodel.Str

	// PaymentCashondeliverySortOrder => Sort Order.
	// Path: payment/cashondelivery/sort_order
	PaymentCashondeliverySortOrder cfgmodel.Str

	// PaymentFreeActive => Enabled.
	// Path: payment/free/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentFreeActive cfgmodel.Bool

	// PaymentFreeOrderStatus => New Order Status.
	// Path: payment/free/order_status
	// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\Newprocessing
	PaymentFreeOrderStatus cfgmodel.Str

	// PaymentFreePaymentAction => Automatically Invoice All Items.
	// Path: payment/free/payment_action
	// SourceModel: Magento\Payment\Model\Source\Invoice
	PaymentFreePaymentAction cfgmodel.Str

	// PaymentFreeSortOrder => Sort Order.
	// Path: payment/free/sort_order
	PaymentFreeSortOrder cfgmodel.Str

	// PaymentFreeTitle => Title.
	// Path: payment/free/title
	PaymentFreeTitle cfgmodel.Str

	// PaymentFreeAllowspecific => Payment from Applicable Countries.
	// Path: payment/free/allowspecific
	// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
	PaymentFreeAllowspecific cfgmodel.Str

	// PaymentFreeSpecificcountry => Payment from Specific Countries.
	// Path: payment/free/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	PaymentFreeSpecificcountry cfgmodel.StringCSV

	// PaymentFreeModel => .
	// Path: payment/free/cfgmodel
	PaymentFreeModel cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.PaymentCheckmoActive = cfgmodel.NewBool(`payment/checkmo/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoOrderStatus = cfgmodel.NewStr(`payment/checkmo/order_status`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoSortOrder = cfgmodel.NewStr(`payment/checkmo/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoTitle = cfgmodel.NewStr(`payment/checkmo/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoAllowspecific = cfgmodel.NewStr(`payment/checkmo/allowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoSpecificcountry = cfgmodel.NewStringCSV(`payment/checkmo/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoPayableTo = cfgmodel.NewStr(`payment/checkmo/payable_to`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoMailingAddress = cfgmodel.NewStr(`payment/checkmo/mailing_address`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoMinOrderTotal = cfgmodel.NewStr(`payment/checkmo/min_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoMaxOrderTotal = cfgmodel.NewStr(`payment/checkmo/max_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCheckmoModel = cfgmodel.NewStr(`payment/checkmo/cfgmodel`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentPurchaseorderActive = cfgmodel.NewBool(`payment/purchaseorder/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentPurchaseorderOrderStatus = cfgmodel.NewStr(`payment/purchaseorder/order_status`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentPurchaseorderSortOrder = cfgmodel.NewStr(`payment/purchaseorder/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentPurchaseorderTitle = cfgmodel.NewStr(`payment/purchaseorder/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentPurchaseorderAllowspecific = cfgmodel.NewStr(`payment/purchaseorder/allowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentPurchaseorderSpecificcountry = cfgmodel.NewStringCSV(`payment/purchaseorder/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentPurchaseorderMinOrderTotal = cfgmodel.NewStr(`payment/purchaseorder/min_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentPurchaseorderMaxOrderTotal = cfgmodel.NewStr(`payment/purchaseorder/max_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentPurchaseorderModel = cfgmodel.NewStr(`payment/purchaseorder/cfgmodel`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentBanktransferActive = cfgmodel.NewBool(`payment/banktransfer/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentBanktransferTitle = cfgmodel.NewStr(`payment/banktransfer/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentBanktransferOrderStatus = cfgmodel.NewStr(`payment/banktransfer/order_status`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentBanktransferAllowspecific = cfgmodel.NewStr(`payment/banktransfer/allowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentBanktransferSpecificcountry = cfgmodel.NewStringCSV(`payment/banktransfer/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentBanktransferInstructions = cfgmodel.NewStr(`payment/banktransfer/instructions`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentBanktransferMinOrderTotal = cfgmodel.NewStr(`payment/banktransfer/min_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentBanktransferMaxOrderTotal = cfgmodel.NewStr(`payment/banktransfer/max_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentBanktransferSortOrder = cfgmodel.NewStr(`payment/banktransfer/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCashondeliveryActive = cfgmodel.NewBool(`payment/cashondelivery/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCashondeliveryTitle = cfgmodel.NewStr(`payment/cashondelivery/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCashondeliveryOrderStatus = cfgmodel.NewStr(`payment/cashondelivery/order_status`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCashondeliveryAllowspecific = cfgmodel.NewStr(`payment/cashondelivery/allowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCashondeliverySpecificcountry = cfgmodel.NewStringCSV(`payment/cashondelivery/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCashondeliveryInstructions = cfgmodel.NewStr(`payment/cashondelivery/instructions`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCashondeliveryMinOrderTotal = cfgmodel.NewStr(`payment/cashondelivery/min_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCashondeliveryMaxOrderTotal = cfgmodel.NewStr(`payment/cashondelivery/max_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentCashondeliverySortOrder = cfgmodel.NewStr(`payment/cashondelivery/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentFreeActive = cfgmodel.NewBool(`payment/free/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentFreeOrderStatus = cfgmodel.NewStr(`payment/free/order_status`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentFreePaymentAction = cfgmodel.NewStr(`payment/free/payment_action`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentFreeSortOrder = cfgmodel.NewStr(`payment/free/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentFreeTitle = cfgmodel.NewStr(`payment/free/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentFreeAllowspecific = cfgmodel.NewStr(`payment/free/allowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentFreeSpecificcountry = cfgmodel.NewStringCSV(`payment/free/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentFreeModel = cfgmodel.NewStr(`payment/free/cfgmodel`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
