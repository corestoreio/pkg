// +build ignore

package productalert

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogProductalertAllowPrice => Allow Alert When Product Price Changes.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogProductalertAllowPrice = model.NewBool(`catalog/productalert/allow_price`)

// PathCatalogProductalertAllowStock => Allow Alert When Product Comes Back in Stock.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogProductalertAllowStock = model.NewBool(`catalog/productalert/allow_stock`)

// PathCatalogProductalertEmailPriceTemplate => Price Alert Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCatalogProductalertEmailPriceTemplate = model.NewStr(`catalog/productalert/email_price_template`)

// PathCatalogProductalertEmailStockTemplate => Stock Alert Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCatalogProductalertEmailStockTemplate = model.NewStr(`catalog/productalert/email_stock_template`)

// PathCatalogProductalertEmailIdentity => Alert Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCatalogProductalertEmailIdentity = model.NewStr(`catalog/productalert/email_identity`)

// PathCatalogProductalertCronFrequency => Frequency.
// BackendModel: Otnegam\Cron\Model\Config\Backend\Product\Alert
// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
var PathCatalogProductalertCronFrequency = model.NewStr(`catalog/productalert_cron/frequency`)

// PathCatalogProductalertCronTime => Start Time.
var PathCatalogProductalertCronTime = model.NewStr(`catalog/productalert_cron/time`)

// PathCatalogProductalertCronErrorEmail => Error Email Recipient.
var PathCatalogProductalertCronErrorEmail = model.NewStr(`catalog/productalert_cron/error_email`)

// PathCatalogProductalertCronErrorEmailIdentity => Error Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCatalogProductalertCronErrorEmailIdentity = model.NewStr(`catalog/productalert_cron/error_email_identity`)

// PathCatalogProductalertCronErrorEmailTemplate => Error Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCatalogProductalertCronErrorEmailTemplate = model.NewStr(`catalog/productalert_cron/error_email_template`)
