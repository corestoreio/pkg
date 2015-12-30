// +build ignore

package directory

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCurrencyOptionsBase => Base Currency.
// Base currency is used for all online payment transactions. If you have more
// than one store view, the base currency scope is defined by the catalog
// price scope ("Catalog" > "Price" > "Catalog Price Scope").
// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Base
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
var PathCurrencyOptionsBase = model.NewStr(`currency/options/base`)

// PathCurrencyOptionsDefault => Default Display Currency.
// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\DefaultCurrency
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
var PathCurrencyOptionsDefault = model.NewStr(`currency/options/default`)

// PathCurrencyOptionsAllow => Allowed Currencies.
// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Allow
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
var PathCurrencyOptionsAllow = model.NewStringCSV(`currency/options/allow`)

// PathCurrencyWebservicexTimeout => Connection Timeout in Seconds.
var PathCurrencyWebservicexTimeout = model.NewStr(`currency/webservicex/timeout`)

// PathCurrencyImportEnabled => Enabled.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCurrencyImportEnabled = model.NewBool(`currency/import/enabled`)

// PathCurrencyImportErrorEmail => Error Email Recipient.
var PathCurrencyImportErrorEmail = model.NewStr(`currency/import/error_email`)

// PathCurrencyImportErrorEmailIdentity => Error Email Sender.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
var PathCurrencyImportErrorEmailIdentity = model.NewStr(`currency/import/error_email_identity`)

// PathCurrencyImportErrorEmailTemplate => Error Email Template.
// Email template chosen based on theme fallback when "Default" option is
// selected.
// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
var PathCurrencyImportErrorEmailTemplate = model.NewStr(`currency/import/error_email_template`)

// PathCurrencyImportFrequency => Frequency.
// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
var PathCurrencyImportFrequency = model.NewStr(`currency/import/frequency`)

// PathCurrencyImportService => Service.
// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Cron
// SourceModel: Otnegam\Directory\Model\Currency\Import\Source\Service
var PathCurrencyImportService = model.NewStr(`currency/import/service`)

// PathCurrencyImportTime => Start Time.
var PathCurrencyImportTime = model.NewStr(`currency/import/time`)

// PathSystemCurrencyInstalled => Installed Currencies.
// BackendModel: Otnegam\Config\Model\Config\Backend\Locale
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency\All
var PathSystemCurrencyInstalled = model.NewStringCSV(`system/currency/installed`)

// PathGeneralCountryOptionalZipCountries => Zip/Postal Code is Optional for.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathGeneralCountryOptionalZipCountries = model.NewStringCSV(`general/country/optional_zip_countries`)

// PathGeneralRegionStateRequired => State is Required for.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathGeneralRegionStateRequired = model.NewStringCSV(`general/region/state_required`)

// PathGeneralRegionDisplayAll => Allow to Choose State if It is Optional for Country.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathGeneralRegionDisplayAll = model.NewBool(`general/region/display_all`)

// PathGeneralLocaleWeightUnit => Weight Unit.
// SourceModel: Otnegam\Directory\Model\Config\Source\WeightUnit
var PathGeneralLocaleWeightUnit = model.NewStr(`general/locale/weight_unit`)
