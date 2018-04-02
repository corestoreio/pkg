// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package directory

import (
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

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
	CurrencyOptionsBase ConfigCurrency

	// CurrencyOptionsDefault => Default Display Currency.
	// Path: currency/options/default
	// BackendModel: Magento\Config\Model\Config\Backend\Currency\DefaultCurrency
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
	CurrencyOptionsDefault ConfigCurrency

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
	// Defines all installed and available currencies.
	// Path: system/currency/installed
	// BackendModel: Magento\Config\Model\Config\Backend\Locale
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency\All
	//
	// TODO:
	//	// SourceCurrencyAll used in Path: `system/currency/installed`,
	//func (sca *SourceCurrencyAll) Options() source.Slice {
	//	// Magento\Framework\Locale\Resolver
	//	// 1. get all allowed currencies from the config
	//	// 2. get slice of currency code and currency name and filter out all not-allowed currencies
	//	// grep locale from general/locale/code scope::store for the current store ID
	//	// the store locale greps the currencies from http://php.net/manual/en/class.resourcebundle.php
	//	// in the correct language
	//	storeLocale, err := sca.mc.Config.String(config.Path(PathDefaultLocale), config.ScopeStore(sca.mc.ScopeStore.StoreID()))
	//	return nil
	//}
	SystemCurrencyInstalled cfgmodel.StringCSV

	// GeneralCountryOptionalZipCountries => Zip/Postal Code is Optional for.
	// Path: general/country/optional_zip_countries
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	GeneralCountryOptionalZipCountries cfgmodel.StringCSV

	// PathGeneralCountryAllow per store view allowed list of countries
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	GeneralCountryAllow cfgmodel.StringCSV

	// PathGeneralCountryDefault returns the store view default configured country code.
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	GeneralCountryDefault cfgmodel.Str

	// PathGeneralCountryEuCountries => European Union Countries.
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	GeneralCountryEuCountries cfgmodel.StringCSV

	// PathGeneralCountryDestinations contains codes of the most used countries.
	// Such countries can be shown on the top of the country list.
	// SourceModel: Magento\Directory\Model\Config\Source\Country
	GeneralCountryDestinations cfgmodel.StringCSV

	// PathGeneralLocaleTimezone => Timezone.
	// BackendModel: Magento\Config\Model\Config\Backend\Locale\Timezone
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Timezone
	GeneralLocaleTimezone cfgmodel.Str

	// PathGeneralLocaleCode => Locale.
	// SourceModel: Magento\Config\Model\Config\Source\Locale
	GeneralLocaleCode cfgmodel.Str

	// PathGeneralLocaleFirstday => First Day of Week.
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Weekdays
	GeneralLocaleFirstday cfgmodel.Str

	// PathGeneralLocaleWeekend => Weekend Days.
	// SourceModel: Magento\Config\Model\Config\Source\Locale\Weekdays
	GeneralLocaleWeekend cfgmodel.StringCSV

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

// NewBackend initializes the global configuration models containing the
// cfgpath.Route variable to the appropriate entry.
// The function Load() will be executed to apply the Sections
// to all models. See Load() for more details.
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).Load(cfgStruct)
}

// Load creates the configuration models for each PkgBackend field.
// Internal mutex will protect the fields during loading.
// The argument Sections will be applied to all models.
func (pp *PkgBackend) Load(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()

	opt := cfgmodel.WithFieldFromSectionSlice(cfgStruct)

	pp.CurrencyOptionsBase = NewConfigCurrency(`currency/options/base`, opt)
	pp.CurrencyOptionsDefault = NewConfigCurrency(`currency/options/default`, opt)
	pp.CurrencyOptionsAllow = cfgmodel.NewStringCSV(`currency/options/allow`, opt)
	pp.CurrencyWebservicexTimeout = cfgmodel.NewStr(`currency/webservicex/timeout`, opt)
	pp.CurrencyImportEnabled = cfgmodel.NewBool(`currency/import/enabled`, opt)
	pp.CurrencyImportErrorEmail = cfgmodel.NewStr(`currency/import/error_email`, opt)
	pp.CurrencyImportErrorEmailIdentity = cfgmodel.NewStr(`currency/import/error_email_identity`, opt)
	pp.CurrencyImportErrorEmailTemplate = cfgmodel.NewStr(`currency/import/error_email_template`, opt)
	pp.CurrencyImportFrequency = cfgmodel.NewStr(`currency/import/frequency`, opt)
	pp.CurrencyImportService = cfgmodel.NewStr(`currency/import/service`, opt)
	pp.CurrencyImportTime = cfgmodel.NewStr(`currency/import/time`, opt)
	pp.SystemCurrencyInstalled = cfgmodel.NewStringCSV(`system/currency/installed`, opt)
	pp.GeneralCountryOptionalZipCountries = cfgmodel.NewStringCSV(`general/country/optional_zip_countries`, opt)
	pp.GeneralCountryAllow = cfgmodel.NewStringCSV(`general/country/allow`, opt)
	pp.GeneralCountryDefault = cfgmodel.NewStr(`general/country/default`, opt)
	pp.GeneralCountryEuCountries = cfgmodel.NewStringCSV(`general/country/eu_countries`, opt)
	pp.GeneralCountryDestinations = cfgmodel.NewStringCSV(`general/country/destinations`, opt)
	pp.GeneralLocaleTimezone = cfgmodel.NewStr(`general/locale/timezone`, opt)
	pp.GeneralLocaleCode = cfgmodel.NewStr(`general/locale/code`, opt)
	pp.GeneralLocaleFirstday = cfgmodel.NewStr(`general/locale/firstday`, opt)
	pp.GeneralLocaleWeekend = cfgmodel.NewStringCSV(`general/locale/weekend`, opt)
	pp.GeneralRegionStateRequired = cfgmodel.NewStringCSV(`general/region/state_required`, opt)
	pp.GeneralRegionDisplayAll = cfgmodel.NewBool(`general/region/display_all`, opt)
	pp.GeneralLocaleWeightUnit = cfgmodel.NewStr(`general/locale/weight_unit`, opt)

	return pp
}
