// +build ignore

package directory

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
	// CurrencyOptionsBase => Base Currency.
	// Base currency is used for all online payment transactions. If you have more
	// than one store view, the base currency scope is defined by the catalog
	// price scope ("Catalog" > "Price" > "Catalog Price Scope").
	// Path: currency/options/base
	// BackendModel: Magento\Config\Model\Config\Backend\Currency\Base
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
	CurrencyOptionsBase cfgmodel.Str

	// CurrencyOptionsDefault => Default Display Currency.
	// Path: currency/options/default
	// BackendModel: Magento\Config\Model\Config\Backend\Currency\DefaultCurrency
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
	CurrencyOptionsDefault cfgmodel.Str

	// CurrencyOptionsAllow => Allowed Currencies.
	// Path: currency/options/allow
	// BackendModel: Magento\Config\Model\Config\Backend\Currency\Allow
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
	CurrencyOptionsAllow cfgmodel.StringCSV

	// CurrencyWebservicexTimeout => Connection Timeout in Seconds.
	// Path: currency/webservicex/timeout
	CurrencyWebservicexTimeout cfgmodel.Str

	// CurrencyImportEnabled => Enabled.
	// Path: currency/import/enabled
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CurrencyImportEnabled cfgmodel.Bool

	// CurrencyImportErrorEmail => Error Email Recipient.
	// Path: currency/import/error_email
	CurrencyImportErrorEmail cfgmodel.Str

	// CurrencyImportErrorEmailIdentity => Error Email Sender.
	// Path: currency/import/error_email_identity
	// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
	CurrencyImportErrorEmailIdentity cfgmodel.Str

	// CurrencyImportErrorEmailTemplate => Error Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: currency/import/error_email_template
	// SourceModel: Magento\Config\Model\Config\Source\Email\Template
	CurrencyImportErrorEmailTemplate cfgmodel.Str

	// CurrencyImportFrequency => Frequency.
	// Path: currency/import/frequency
	// SourceModel: Magento\Cron\Model\Config\Source\Frequency
	CurrencyImportFrequency cfgmodel.Str

	// CurrencyImportService => Service.
	// Path: currency/import/service
	// BackendModel: Magento\Config\Model\Config\Backend\Currency\Cron
	// SourceModel: Magento\Directory\Model\Currency\Import\Source\Service
	CurrencyImportService cfgmodel.Str

	// CurrencyImportTime => Start Time.
	// Path: currency/import/time
	CurrencyImportTime cfgmodel.Str

	// SystemCurrencyInstalled => Installed Currencies.
	// Path: system/currency/installed
	// BackendModel: Magento\Config\Model\Config\Backend\Locale
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency\All
	SystemCurrencyInstalled cfgmodel.StringCSV

	// GeneralCountryOptionalZipCountries => Zip/Postal Code is Optional for.
	// Path: general/country/optional_zip_countries
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	GeneralCountryOptionalZipCountries cfgmodel.StringCSV

	// GeneralRegionStateRequired => State is Required for.
	// Path: general/region/state_required
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	GeneralRegionStateRequired cfgmodel.StringCSV

	// GeneralRegionDisplayAll => Allow to Choose State if It is Optional for Country.
	// Path: general/region/display_all
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	GeneralRegionDisplayAll cfgmodel.Bool

	// GeneralLocaleWeightUnit => Weight Unit.
	// Path: general/locale/weight_unit
	// SourceModel: Magento\Directory\Model\Config\Source\WeightUnit
	GeneralLocaleWeightUnit cfgmodel.Str
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CurrencyOptionsBase = cfgmodel.NewStr(`currency/options/base`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyOptionsDefault = cfgmodel.NewStr(`currency/options/default`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyOptionsAllow = cfgmodel.NewStringCSV(`currency/options/allow`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyWebservicexTimeout = cfgmodel.NewStr(`currency/webservicex/timeout`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyImportEnabled = cfgmodel.NewBool(`currency/import/enabled`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyImportErrorEmail = cfgmodel.NewStr(`currency/import/error_email`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyImportErrorEmailIdentity = cfgmodel.NewStr(`currency/import/error_email_identity`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyImportErrorEmailTemplate = cfgmodel.NewStr(`currency/import/error_email_template`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyImportFrequency = cfgmodel.NewStr(`currency/import/frequency`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyImportService = cfgmodel.NewStr(`currency/import/service`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.CurrencyImportTime = cfgmodel.NewStr(`currency/import/time`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.SystemCurrencyInstalled = cfgmodel.NewStringCSV(`system/currency/installed`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GeneralCountryOptionalZipCountries = cfgmodel.NewStringCSV(`general/country/optional_zip_countries`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GeneralRegionStateRequired = cfgmodel.NewStringCSV(`general/region/state_required`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GeneralRegionDisplayAll = cfgmodel.NewBool(`general/region/display_all`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	pp.GeneralLocaleWeightUnit = cfgmodel.NewStr(`general/locale/weight_unit`, cfgmodel.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
