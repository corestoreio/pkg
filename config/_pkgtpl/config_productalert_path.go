// +build ignore

package productalert

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogProductalertAllowPrice => Allow Alert When Product Price Changes.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogProductalertAllowPrice = model.NewBool(`catalog/productalert/allow_price`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductalertAllowStock => Allow Alert When Product Comes Back in Stock.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogProductalertAllowStock = model.NewBool(`catalog/productalert/allow_stock`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductalertEmailPriceTemplate => Price Alert Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCatalogProductalertEmailPriceTemplate = model.NewStr(`catalog/productalert/email_price_template`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductalertEmailStockTemplate => Stock Alert Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCatalogProductalertEmailStockTemplate = model.NewStr(`catalog/productalert/email_stock_template`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductalertEmailIdentity => Alert Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCatalogProductalertEmailIdentity = model.NewStr(`catalog/productalert/email_identity`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductalertCronFrequency => Frequency.
// BackendModel: Otnegam\Cron\Model\Config\Backend\Product\Alert
// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
var PathCatalogProductalertCronFrequency = model.NewStr(`catalog/productalert_cron/frequency`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductalertCronTime => Start Time.
var PathCatalogProductalertCronTime = model.NewStr(`catalog/productalert_cron/time`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductalertCronErrorEmail => Error Email Recipient.
var PathCatalogProductalertCronErrorEmail = model.NewStr(`catalog/productalert_cron/error_email`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductalertCronErrorEmailIdentity => Error Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCatalogProductalertCronErrorEmailIdentity = model.NewStr(`catalog/productalert_cron/error_email_identity`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogProductalertCronErrorEmailTemplate => Error Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCatalogProductalertCronErrorEmailTemplate = model.NewStr(`catalog/productalert_cron/error_email_template`, model.WithPkgCfg(PackageConfiguration))
