// +build ignore

package productalert

import (
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CatalogProductalertAllowPrice => Allow Alert When Product Price Changes.
	// Path: catalog/productalert/allow_price
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogProductalertAllowPrice cfgmodel.Bool

	// CatalogProductalertAllowStock => Allow Alert When Product Comes Back in Stock.
	// Path: catalog/productalert/allow_stock
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogProductalertAllowStock cfgmodel.Bool

	// CatalogProductalertEmailPriceTemplate => Price Alert Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: catalog/productalert/email_price_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CatalogProductalertEmailPriceTemplate cfgmodel.Str

	// CatalogProductalertEmailStockTemplate => Stock Alert Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: catalog/productalert/email_stock_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CatalogProductalertEmailStockTemplate cfgmodel.Str

	// CatalogProductalertEmailIdentity => Alert Email Sender.
	// Path: catalog/productalert/email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	CatalogProductalertEmailIdentity cfgmodel.Str

	// CatalogProductalertCronFrequency => Frequency.
	// Path: catalog/productalert_cron/frequency
	// BackendModel: Magento\Cron\Model\Config\Backend\Product\Alert
	// SourceModel: Magento\Cron\Model\Config\Source\Frequency
	CatalogProductalertCronFrequency cfgmodel.Str

	// CatalogProductalertCronTime => Start Time.
	// Path: catalog/productalert_cron/time
	CatalogProductalertCronTime cfgmodel.Str

	// CatalogProductalertCronErrorEmail => Error Email Recipient.
	// Path: catalog/productalert_cron/error_email
	CatalogProductalertCronErrorEmail cfgmodel.Str

	// CatalogProductalertCronErrorEmailIdentity => Error Email Sender.
	// Path: catalog/productalert_cron/error_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	CatalogProductalertCronErrorEmailIdentity cfgmodel.Str

	// CatalogProductalertCronErrorEmailTemplate => Error Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: catalog/productalert_cron/error_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CatalogProductalertCronErrorEmailTemplate cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogProductalertAllowPrice = cfgmodel.NewBool(`catalog/productalert/allow_price`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductalertAllowStock = cfgmodel.NewBool(`catalog/productalert/allow_stock`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductalertEmailPriceTemplate = cfgmodel.NewStr(`catalog/productalert/email_price_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductalertEmailStockTemplate = cfgmodel.NewStr(`catalog/productalert/email_stock_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductalertEmailIdentity = cfgmodel.NewStr(`catalog/productalert/email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductalertCronFrequency = cfgmodel.NewStr(`catalog/productalert_cron/frequency`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductalertCronTime = cfgmodel.NewStr(`catalog/productalert_cron/time`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductalertCronErrorEmail = cfgmodel.NewStr(`catalog/productalert_cron/error_email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductalertCronErrorEmailIdentity = cfgmodel.NewStr(`catalog/productalert_cron/error_email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogProductalertCronErrorEmailTemplate = cfgmodel.NewStr(`catalog/productalert_cron/error_email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
