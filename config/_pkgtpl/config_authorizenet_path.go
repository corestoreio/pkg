// +build ignore

package authorizenet

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
	// PaymentAuthorizenetDirectpostActive => Enabled.
	// Path: payment/authorizenet_directpost/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostActive model.Bool

	// PaymentAuthorizenetDirectpostPaymentAction => Payment Action.
	// Path: payment/authorizenet_directpost/payment_action
	// SourceModel: Otnegam\Authorizenet\Model\Source\PaymentAction
	PaymentAuthorizenetDirectpostPaymentAction model.Str

	// PaymentAuthorizenetDirectpostTitle => Title.
	// Path: payment/authorizenet_directpost/title
	PaymentAuthorizenetDirectpostTitle model.Str

	// PaymentAuthorizenetDirectpostLogin => API Login ID.
	// Path: payment/authorizenet_directpost/login
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	PaymentAuthorizenetDirectpostLogin model.Str

	// PaymentAuthorizenetDirectpostTransKey => Transaction Key.
	// Path: payment/authorizenet_directpost/trans_key
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	PaymentAuthorizenetDirectpostTransKey model.Str

	// PaymentAuthorizenetDirectpostTransMd5 => Merchant MD5.
	// Path: payment/authorizenet_directpost/trans_md5
	// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
	PaymentAuthorizenetDirectpostTransMd5 model.Str

	// PaymentAuthorizenetDirectpostOrderStatus => New Order Status.
	// Path: payment/authorizenet_directpost/order_status
	// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\Processing
	PaymentAuthorizenetDirectpostOrderStatus model.Str

	// PaymentAuthorizenetDirectpostTest => Test Mode.
	// Path: payment/authorizenet_directpost/test
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostTest model.Bool

	// PaymentAuthorizenetDirectpostCgiUrl => Gateway URL.
	// Path: payment/authorizenet_directpost/cgi_url
	PaymentAuthorizenetDirectpostCgiUrl model.Str

	// PaymentAuthorizenetDirectpostCgiUrlTd => Transaction Details URL.
	// Path: payment/authorizenet_directpost/cgi_url_td
	PaymentAuthorizenetDirectpostCgiUrlTd model.Str

	// PaymentAuthorizenetDirectpostCurrency => Accepted Currency.
	// Path: payment/authorizenet_directpost/currency
	// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
	PaymentAuthorizenetDirectpostCurrency model.Str

	// PaymentAuthorizenetDirectpostDebug => Debug.
	// Path: payment/authorizenet_directpost/debug
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostDebug model.Bool

	// PaymentAuthorizenetDirectpostEmailCustomer => Email Customer.
	// Path: payment/authorizenet_directpost/email_customer
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostEmailCustomer model.Bool

	// PaymentAuthorizenetDirectpostMerchantEmail => Merchant's Email.
	// Path: payment/authorizenet_directpost/merchant_email
	PaymentAuthorizenetDirectpostMerchantEmail model.Str

	// PaymentAuthorizenetDirectpostCctypes => Credit Card Types.
	// Path: payment/authorizenet_directpost/cctypes
	// SourceModel: Otnegam\Authorizenet\Model\Source\Cctype
	PaymentAuthorizenetDirectpostCctypes model.StringCSV

	// PaymentAuthorizenetDirectpostUseccv => Credit Card Verification.
	// Path: payment/authorizenet_directpost/useccv
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	PaymentAuthorizenetDirectpostUseccv model.Bool

	// PaymentAuthorizenetDirectpostAllowspecific => Payment from Applicable Countries.
	// Path: payment/authorizenet_directpost/allowspecific
	// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
	PaymentAuthorizenetDirectpostAllowspecific model.Str

	// PaymentAuthorizenetDirectpostSpecificcountry => Payment from Specific Countries.
	// Path: payment/authorizenet_directpost/specificcountry
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
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

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.PaymentAuthorizenetDirectpostActive = model.NewBool(`payment/authorizenet_directpost/active`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostPaymentAction = model.NewStr(`payment/authorizenet_directpost/payment_action`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostTitle = model.NewStr(`payment/authorizenet_directpost/title`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostLogin = model.NewStr(`payment/authorizenet_directpost/login`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostTransKey = model.NewStr(`payment/authorizenet_directpost/trans_key`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostTransMd5 = model.NewStr(`payment/authorizenet_directpost/trans_md5`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostOrderStatus = model.NewStr(`payment/authorizenet_directpost/order_status`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostTest = model.NewBool(`payment/authorizenet_directpost/test`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostCgiUrl = model.NewStr(`payment/authorizenet_directpost/cgi_url`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostCgiUrlTd = model.NewStr(`payment/authorizenet_directpost/cgi_url_td`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostCurrency = model.NewStr(`payment/authorizenet_directpost/currency`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostDebug = model.NewBool(`payment/authorizenet_directpost/debug`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostEmailCustomer = model.NewBool(`payment/authorizenet_directpost/email_customer`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostMerchantEmail = model.NewStr(`payment/authorizenet_directpost/merchant_email`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostCctypes = model.NewStringCSV(`payment/authorizenet_directpost/cctypes`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostUseccv = model.NewBool(`payment/authorizenet_directpost/useccv`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostAllowspecific = model.NewStr(`payment/authorizenet_directpost/allowspecific`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostSpecificcountry = model.NewStringCSV(`payment/authorizenet_directpost/specificcountry`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostMinOrderTotal = model.NewStr(`payment/authorizenet_directpost/min_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostMaxOrderTotal = model.NewStr(`payment/authorizenet_directpost/max_order_total`, model.WithPkgCfg(pkgCfg))
	pp.PaymentAuthorizenetDirectpostSortOrder = model.NewStr(`payment/authorizenet_directpost/sort_order`, model.WithPkgCfg(pkgCfg))

	return pp
}
