// +build ignore

package productalert

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// CatalogProductalertAllowPrice => Allow Alert When Product Price Changes.
	// Path: catalog/productalert/allow_price
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogProductalertAllowPrice model.Bool

	// CatalogProductalertAllowStock => Allow Alert When Product Comes Back in Stock.
	// Path: catalog/productalert/allow_stock
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogProductalertAllowStock model.Bool

	// CatalogProductalertEmailPriceTemplate => Price Alert Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: catalog/productalert/email_price_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	CatalogProductalertEmailPriceTemplate model.Str

	// CatalogProductalertEmailStockTemplate => Stock Alert Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: catalog/productalert/email_stock_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	CatalogProductalertEmailStockTemplate model.Str

	// CatalogProductalertEmailIdentity => Alert Email Sender.
	// Path: catalog/productalert/email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	CatalogProductalertEmailIdentity model.Str

	// CatalogProductalertCronFrequency => Frequency.
	// Path: catalog/productalert_cron/frequency
	// BackendModel: Otnegam\Cron\Model\Config\Backend\Product\Alert
	// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
	CatalogProductalertCronFrequency model.Str

	// CatalogProductalertCronTime => Start Time.
	// Path: catalog/productalert_cron/time
	CatalogProductalertCronTime model.Str

	// CatalogProductalertCronErrorEmail => Error Email Recipient.
	// Path: catalog/productalert_cron/error_email
	CatalogProductalertCronErrorEmail model.Str

	// CatalogProductalertCronErrorEmailIdentity => Error Email Sender.
	// Path: catalog/productalert_cron/error_email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	CatalogProductalertCronErrorEmailIdentity model.Str

	// CatalogProductalertCronErrorEmailTemplate => Error Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: catalog/productalert_cron/error_email_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	CatalogProductalertCronErrorEmailTemplate model.Str
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogProductalertAllowPrice = model.NewBool(`catalog/productalert/allow_price`, model.WithConfigStructure(cfgStruct))
	pp.CatalogProductalertAllowStock = model.NewBool(`catalog/productalert/allow_stock`, model.WithConfigStructure(cfgStruct))
	pp.CatalogProductalertEmailPriceTemplate = model.NewStr(`catalog/productalert/email_price_template`, model.WithConfigStructure(cfgStruct))
	pp.CatalogProductalertEmailStockTemplate = model.NewStr(`catalog/productalert/email_stock_template`, model.WithConfigStructure(cfgStruct))
	pp.CatalogProductalertEmailIdentity = model.NewStr(`catalog/productalert/email_identity`, model.WithConfigStructure(cfgStruct))
	pp.CatalogProductalertCronFrequency = model.NewStr(`catalog/productalert_cron/frequency`, model.WithConfigStructure(cfgStruct))
	pp.CatalogProductalertCronTime = model.NewStr(`catalog/productalert_cron/time`, model.WithConfigStructure(cfgStruct))
	pp.CatalogProductalertCronErrorEmail = model.NewStr(`catalog/productalert_cron/error_email`, model.WithConfigStructure(cfgStruct))
	pp.CatalogProductalertCronErrorEmailIdentity = model.NewStr(`catalog/productalert_cron/error_email_identity`, model.WithConfigStructure(cfgStruct))
	pp.CatalogProductalertCronErrorEmailTemplate = model.NewStr(`catalog/productalert_cron/error_email_template`, model.WithConfigStructure(cfgStruct))

	return pp
}
