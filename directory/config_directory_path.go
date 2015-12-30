// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/csfw/config/model"
)

// PathGeneralCountryAllow per store view allowed list of countries
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathGeneralCountryAllow = model.NewStringCSV(`general/country/allow`)

// PathGeneralCountryDefault returns the store view default configured country code.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathGeneralCountryDefault = model.NewStr(`general/country/default`)

// PathGeneralCountryEuCountries => European Union Countries.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathGeneralCountryEuCountries = model.NewStringCSV(`general/country/eu_countries`)

// PathGeneralCountryDestinations contains codes of the most used countries.
// Such countries can be shown on the top of the country list.
// SourceModel: Otnegam\Directory\Model\Config\Source\Country
var PathGeneralCountryDestinations = model.NewStringCSV(`general/country/destinations`)

// PathGeneralLocaleTimezone => Timezone.
// BackendModel: Otnegam\Config\Model\Config\Backend\Locale\Timezone
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Timezone
var PathGeneralLocaleTimezone = model.NewStr(`general/locale/timezone`)

// PathGeneralLocaleCode => Locale.
// SourceModel: Otnegam\Config\Model\Config\Source\Locale
var PathGeneralLocaleCode = model.NewStr(`general/locale/code`)

// PathGeneralLocaleFirstday => First Day of Week.
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Weekdays
var PathGeneralLocaleFirstday = model.NewStr(`general/locale/firstday`)

// PathGeneralLocaleWeekend => Weekend Days.
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Weekdays
var PathGeneralLocaleWeekend = model.NewStringCSV(`general/locale/weekend`)

// PathCurrencyOptionsBase => Base Currency.
// Base currency is used for all online payment transactions. If you have more
// than one store view, the base currency scope is defined by the catalog
// price scope ("Catalog" > "Price" > "Catalog Price Scope").
// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Base
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
var PathCurrencyOptionsBase = NewConfigCurrency(`currency/options/base`)

// PathCurrencyOptionsDefault => Default Display Currency.
// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\DefaultCurrency
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
var PathCurrencyOptionsDefault = NewConfigCurrency(`currency/options/default`)

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

// PathSystemCurrencyInstalled defines all installed and available currencies.
// Global Scope
// BackendModel: Otnegam\Config\Model\Config\Backend\Locale
// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency\All
var PathSystemCurrencyInstalled = model.NewStringCSV(`system/currency/installed`)

//	// SourceCurrencyAll used in Path: `system/currency/installed`,
//func (sca *SourceCurrencyAll) Options() valuelabel.Slice {
//	// Magento\Framework\Locale\Resolver
//	// 1. get all allowed currencies from the config
//	// 2. get slice of currency code and currency name and filter out all not-allowed currencies
//	// grep locale from general/locale/code scope::store for the current store ID
//	// the store locale greps the currencies from http://php.net/manual/en/class.resourcebundle.php
//	// in the correct language
//	storeLocale, err := sca.mc.Config.String(config.Path(PathDefaultLocale), config.ScopeStore(sca.mc.ScopeStore.StoreID()))
//
//	fmt.Printf("\nstoreLocale: %s\n Err %s\n", storeLocale, err)
//
//	return nil
//}

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
