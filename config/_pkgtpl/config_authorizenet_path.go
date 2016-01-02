// +build ignore

package authorizenet

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathPaymentAuthorizenetDirectpostActive => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentAuthorizenetDirectpostActive = model.NewBool(`payment/authorizenet_directpost/active`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostPaymentAction => Payment Action.
// SourceModel: Otnegam\Authorizenet\Model\Source\PaymentAction
var PathPaymentAuthorizenetDirectpostPaymentAction = model.NewStr(`payment/authorizenet_directpost/payment_action`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostTitle => Title.
var PathPaymentAuthorizenetDirectpostTitle = model.NewStr(`payment/authorizenet_directpost/title`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostLogin => API Login ID.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathPaymentAuthorizenetDirectpostLogin = model.NewStr(`payment/authorizenet_directpost/login`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostTransKey => Transaction Key.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathPaymentAuthorizenetDirectpostTransKey = model.NewStr(`payment/authorizenet_directpost/trans_key`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostTransMd5 => Merchant MD5.
// BackendModel: Otnegam\Config\Model\Config\Backend\Encrypted
var PathPaymentAuthorizenetDirectpostTransMd5 = model.NewStr(`payment/authorizenet_directpost/trans_md5`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostOrderStatus => New Order Status.
// SourceModel: Otnegam\Sales\Model\Config\Source\Order\Status\Processing
var PathPaymentAuthorizenetDirectpostOrderStatus = model.NewStr(`payment/authorizenet_directpost/order_status`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostTest => Test Mode.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentAuthorizenetDirectpostTest = model.NewBool(`payment/authorizenet_directpost/test`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostCgiUrl => Gateway URL.
var PathPaymentAuthorizenetDirectpostCgiUrl = model.NewStr(`payment/authorizenet_directpost/cgi_url`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostCgiUrlTd => Transaction Details URL.
var PathPaymentAuthorizenetDirectpostCgiUrlTd = model.NewStr(`payment/authorizenet_directpost/cgi_url_td`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostCurrency => Accepted Currency.
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
var PathPaymentAuthorizenetDirectpostCurrency = model.NewStr(`payment/authorizenet_directpost/currency`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostDebug => Debug.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentAuthorizenetDirectpostDebug = model.NewBool(`payment/authorizenet_directpost/debug`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostEmailCustomer => Email Customer.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentAuthorizenetDirectpostEmailCustomer = model.NewBool(`payment/authorizenet_directpost/email_customer`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostMerchantEmail => Merchant's Email.
var PathPaymentAuthorizenetDirectpostMerchantEmail = model.NewStr(`payment/authorizenet_directpost/merchant_email`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostCctypes => Credit Card Types.
// SourceModel: Otnegam\Authorizenet\Model\Source\Cctype
var PathPaymentAuthorizenetDirectpostCctypes = model.NewStringCSV(`payment/authorizenet_directpost/cctypes`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostUseccv => Credit Card Verification.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathPaymentAuthorizenetDirectpostUseccv = model.NewBool(`payment/authorizenet_directpost/useccv`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostAllowspecific => Payment from Applicable Countries.
// SourceModel: Otnegam\Payment\Model\Config\Source\Allspecificcountries
var PathPaymentAuthorizenetDirectpostAllowspecific = model.NewStr(`payment/authorizenet_directpost/allowspecific`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostSpecificcountry => Payment from Specific Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathPaymentAuthorizenetDirectpostSpecificcountry = model.NewStringCSV(`payment/authorizenet_directpost/specificcountry`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostMinOrderTotal => Minimum Order Total.
var PathPaymentAuthorizenetDirectpostMinOrderTotal = model.NewStr(`payment/authorizenet_directpost/min_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostMaxOrderTotal => Maximum Order Total.
var PathPaymentAuthorizenetDirectpostMaxOrderTotal = model.NewStr(`payment/authorizenet_directpost/max_order_total`, model.WithPkgCfg(PackageConfiguration))

// PathPaymentAuthorizenetDirectpostSortOrder => Sort Order.
var PathPaymentAuthorizenetDirectpostSortOrder = model.NewStr(`payment/authorizenet_directpost/sort_order`, model.WithPkgCfg(PackageConfiguration))
