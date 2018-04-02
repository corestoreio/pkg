// +build ignore

package authorizenet

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
	// PaymentAuthorizenetDirectpostActive => Enabled.
	// Path: payment/authorizenet_directpost/active
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostActive cfgmodel.Bool

	// PaymentAuthorizenetDirectpostPaymentAction => Payment Action.
	// Path: payment/authorizenet_directpost/payment_action
	// SourceModel: Magento\Authorizenet\Model\Source\PaymentAction
	PaymentAuthorizenetDirectpostPaymentAction cfgmodel.Str

	// PaymentAuthorizenetDirectpostTitle => Title.
	// Path: payment/authorizenet_directpost/title
	PaymentAuthorizenetDirectpostTitle cfgmodel.Str

	// PaymentAuthorizenetDirectpostLogin => API Login ID.
	// Path: payment/authorizenet_directpost/login
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	PaymentAuthorizenetDirectpostLogin cfgmodel.Str

	// PaymentAuthorizenetDirectpostTransKey => Transaction Key.
	// Path: payment/authorizenet_directpost/trans_key
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	PaymentAuthorizenetDirectpostTransKey cfgmodel.Str

	// PaymentAuthorizenetDirectpostTransMd5 => Merchant MD5.
	// Path: payment/authorizenet_directpost/trans_md5
	// BackendModel: Magento\Config\Model\Config\Backend\Encrypted
	PaymentAuthorizenetDirectpostTransMd5 cfgmodel.Str

	// PaymentAuthorizenetDirectpostOrderStatus => New Order Status.
	// Path: payment/authorizenet_directpost/order_status
	// SourceModel: Magento\Sales\Model\Config\Source\Order\Status\Processing
	PaymentAuthorizenetDirectpostOrderStatus cfgmodel.Str

	// PaymentAuthorizenetDirectpostTest => Test Mode.
	// Path: payment/authorizenet_directpost/test
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostTest cfgmodel.Bool

	// PaymentAuthorizenetDirectpostCgiUrl => Gateway URL.
	// Path: payment/authorizenet_directpost/cgi_url
	PaymentAuthorizenetDirectpostCgiUrl cfgmodel.Str

	// PaymentAuthorizenetDirectpostCgiUrlTd => Transaction Details URL.
	// Path: payment/authorizenet_directpost/cgi_url_td
	PaymentAuthorizenetDirectpostCgiUrlTd cfgmodel.Str

	// PaymentAuthorizenetDirectpostCurrency => Accepted Currency.
	// Path: payment/authorizenet_directpost/currency
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
	PaymentAuthorizenetDirectpostCurrency cfgmodel.Str

	// PaymentAuthorizenetDirectpostDebug => Debug.
	// Path: payment/authorizenet_directpost/debug
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostDebug cfgmodel.Bool

	// PaymentAuthorizenetDirectpostEmailCustomer => Email Customer.
	// Path: payment/authorizenet_directpost/email_customer
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostEmailCustomer cfgmodel.Bool

	// PaymentAuthorizenetDirectpostMerchantEmail => Merchant's Email.
	// Path: payment/authorizenet_directpost/merchant_email
	PaymentAuthorizenetDirectpostMerchantEmail cfgmodel.Str

	// PaymentAuthorizenetDirectpostCctypes => Credit Card Types.
	// Path: payment/authorizenet_directpost/cctypes
	// SourceModel: Magento\Authorizenet\Model\Source\Cctype
	PaymentAuthorizenetDirectpostCctypes cfgmodel.StringCSV

	// PaymentAuthorizenetDirectpostUseccv => Credit Card Verification.
	// Path: payment/authorizenet_directpost/useccv
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostUseccv cfgmodel.Bool

	// PaymentAuthorizenetDirectpostAllowspecific => Payment from Applicable Countries.
	// Path: payment/authorizenet_directpost/allowspecific
	// SourceModel: Magento\Payment\Model\Config\Source\Allspecificcountries
	PaymentAuthorizenetDirectpostAllowspecific cfgmodel.Str

	// PaymentAuthorizenetDirectpostSpecificcountry => Payment from Specific Countries.
	// Path: payment/authorizenet_directpost/specificcountry
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	PaymentAuthorizenetDirectpostSpecificcountry cfgmodel.StringCSV

	// PaymentAuthorizenetDirectpostMinOrderTotal => Minimum Order Total.
	// Path: payment/authorizenet_directpost/min_order_total
	PaymentAuthorizenetDirectpostMinOrderTotal cfgmodel.Str

	// PaymentAuthorizenetDirectpostMaxOrderTotal => Maximum Order Total.
	// Path: payment/authorizenet_directpost/max_order_total
	PaymentAuthorizenetDirectpostMaxOrderTotal cfgmodel.Str

	// PaymentAuthorizenetDirectpostSortOrder => Sort Order.
	// Path: payment/authorizenet_directpost/sort_order
	PaymentAuthorizenetDirectpostSortOrder cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.PaymentAuthorizenetDirectpostActive = cfgmodel.NewBool(`payment/authorizenet_directpost/active`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostPaymentAction = cfgmodel.NewStr(`payment/authorizenet_directpost/payment_action`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostTitle = cfgmodel.NewStr(`payment/authorizenet_directpost/title`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostLogin = cfgmodel.NewStr(`payment/authorizenet_directpost/login`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostTransKey = cfgmodel.NewStr(`payment/authorizenet_directpost/trans_key`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostTransMd5 = cfgmodel.NewStr(`payment/authorizenet_directpost/trans_md5`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostOrderStatus = cfgmodel.NewStr(`payment/authorizenet_directpost/order_status`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostTest = cfgmodel.NewBool(`payment/authorizenet_directpost/test`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostCgiUrl = cfgmodel.NewStr(`payment/authorizenet_directpost/cgi_url`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostCgiUrlTd = cfgmodel.NewStr(`payment/authorizenet_directpost/cgi_url_td`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostCurrency = cfgmodel.NewStr(`payment/authorizenet_directpost/currency`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostDebug = cfgmodel.NewBool(`payment/authorizenet_directpost/debug`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostEmailCustomer = cfgmodel.NewBool(`payment/authorizenet_directpost/email_customer`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostMerchantEmail = cfgmodel.NewStr(`payment/authorizenet_directpost/merchant_email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostCctypes = cfgmodel.NewStringCSV(`payment/authorizenet_directpost/cctypes`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostUseccv = cfgmodel.NewBool(`payment/authorizenet_directpost/useccv`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostAllowspecific = cfgmodel.NewStr(`payment/authorizenet_directpost/allowspecific`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostSpecificcountry = cfgmodel.NewStringCSV(`payment/authorizenet_directpost/specificcountry`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostMinOrderTotal = cfgmodel.NewStr(`payment/authorizenet_directpost/min_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostMaxOrderTotal = cfgmodel.NewStr(`payment/authorizenet_directpost/max_order_total`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.PaymentAuthorizenetDirectpostSortOrder = cfgmodel.NewStr(`payment/authorizenet_directpost/sort_order`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
