// +build ignore

package authorizenet

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
	// PaymentAuthorizenetDirectpostActive => Enabled.
	// Path: payment/authorizenet_directpost/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostActive model.Bool

	// PaymentAuthorizenetDirectpostPaymentAction => Payment Action.
	// Path: payment/authorizenet_directpost/payment_action
	// SourceModel: Magento\Authorizenet\Model\Source\PaymentAction
	PaymentAuthorizenetDirectpostPaymentAction model.Str

	// PaymentAuthorizenetDirectpostTitle => Title.
	// Path: payment/authorizenet_directpost/title
	PaymentAuthorizenetDirectpostTitle model.Str

	// PaymentAuthorizenetDirectpostLogin => API Login ID.
	// Path: payment/authorizenet_directpost/login
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	PaymentAuthorizenetDirectpostLogin model.Str

	// PaymentAuthorizenetDirectpostTransKey => Transaction Key.
	// Path: payment/authorizenet_directpost/trans_key
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	PaymentAuthorizenetDirectpostTransKey model.Str

	// PaymentAuthorizenetDirectpostTransMd5 => Merchant MD5.
	// Path: payment/authorizenet_directpost/trans_md5
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	PaymentAuthorizenetDirectpostTransMd5 model.Str

	// PaymentAuthorizenetDirectpostOrderStatus => New Order Status.
	// Path: payment/authorizenet_directpost/order_status
	// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\Processing
	PaymentAuthorizenetDirectpostOrderStatus model.Str

	// PaymentAuthorizenetDirectpostTest => Test Mode.
	// Path: payment/authorizenet_directpost/test
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostTest model.Bool

	// PaymentAuthorizenetDirectpostCgiUrl => Gateway URL.
	// Path: payment/authorizenet_directpost/cgi_url
	PaymentAuthorizenetDirectpostCgiUrl model.Str

	// PaymentAuthorizenetDirectpostCgiUrlTd => Transaction Details URL.
	// Path: payment/authorizenet_directpost/cgi_url_td
	PaymentAuthorizenetDirectpostCgiUrlTd model.Str

	// PaymentAuthorizenetDirectpostCurrency => Accepted Currency.
	// Path: payment/authorizenet_directpost/currency
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
	PaymentAuthorizenetDirectpostCurrency model.Str

	// PaymentAuthorizenetDirectpostDebug => Debug.
	// Path: payment/authorizenet_directpost/debug
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostDebug model.Bool

	// PaymentAuthorizenetDirectpostEmailCustomer => Email Customer.
	// Path: payment/authorizenet_directpost/email_customer
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostEmailCustomer model.Bool

	// PaymentAuthorizenetDirectpostMerchantEmail => Merchant's Email.
	// Path: payment/authorizenet_directpost/merchant_email
	PaymentAuthorizenetDirectpostMerchantEmail model.Str

	// PaymentAuthorizenetDirectpostCctypes => Credit Card Types.
	// Path: payment/authorizenet_directpost/cctypes
	// SourceModel: Magento\Authorizenet\Model\Source\Cctype
	PaymentAuthorizenetDirectpostCctypes model.StringCSV

	// PaymentAuthorizenetDirectpostUseccv => Credit Card Verification.
	// Path: payment/authorizenet_directpost/useccv
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostUseccv model.Bool

	// PaymentAuthorizenetDirectpostAllowspecific => Payment from Applicable Countries.
	// Path: payment/authorizenet_directpost/allowspecific
	// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
	PaymentAuthorizenetDirectpostAllowspecific model.Str

	// PaymentAuthorizenetDirectpostSpecificcountry => Payment from Specific Countries.
	// Path: payment/authorizenet_directpost/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	PaymentAuthorizenetDirectpostSpecificcountry model.StringCSV

	// PaymentAuthorizenetDirectpostMinOrderTotal => Minimum Order Total.
	// Path: payment/authorizenet_directpost/min_order_total
	PaymentAuthorizenetDirectpostMinOrderTotal model.Str

	// PaymentAuthorizenetDirectpostMaxOrderTotal => Maximum Order Total.
	// Path: payment/authorizenet_directpost/max_order_total
	PaymentAuthorizenetDirectpostMaxOrderTotal model.Str

	// PaymentAuthorizenetDirectpostSortOrder => Sort Order.
	// Path: payment/authorizenet_directpost/sort_order
	PaymentAuthorizenetDirectpostSortOrder model.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.PaymentAuthorizenetDirectpostActive = model.NewBool(`payment/authorizenet_directpost/active`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostPaymentAction = model.NewStr(`payment/authorizenet_directpost/payment_action`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostTitle = model.NewStr(`payment/authorizenet_directpost/title`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostLogin = model.NewStr(`payment/authorizenet_directpost/login`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostTransKey = model.NewStr(`payment/authorizenet_directpost/trans_key`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostTransMd5 = model.NewStr(`payment/authorizenet_directpost/trans_md5`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostOrderStatus = model.NewStr(`payment/authorizenet_directpost/order_status`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostTest = model.NewBool(`payment/authorizenet_directpost/test`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostCgiUrl = model.NewStr(`payment/authorizenet_directpost/cgi_url`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostCgiUrlTd = model.NewStr(`payment/authorizenet_directpost/cgi_url_td`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostCurrency = model.NewStr(`payment/authorizenet_directpost/currency`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostDebug = model.NewBool(`payment/authorizenet_directpost/debug`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostEmailCustomer = model.NewBool(`payment/authorizenet_directpost/email_customer`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostMerchantEmail = model.NewStr(`payment/authorizenet_directpost/merchant_email`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostCctypes = model.NewStringCSV(`payment/authorizenet_directpost/cctypes`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostUseccv = model.NewBool(`payment/authorizenet_directpost/useccv`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostAllowspecific = model.NewStr(`payment/authorizenet_directpost/allowspecific`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostSpecificcountry = model.NewStringCSV(`payment/authorizenet_directpost/specificcountry`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostMinOrderTotal = model.NewStr(`payment/authorizenet_directpost/min_order_total`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostMaxOrderTotal = model.NewStr(`payment/authorizenet_directpost/max_order_total`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostSortOrder = model.NewStr(`payment/authorizenet_directpost/sort_order`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
