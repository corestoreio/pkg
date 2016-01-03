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
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend

	// CurrencyOptionsBase => Base Currency.
	// Base currency is used for all online payment transactions. If you have more
	// than one store view, the base currency scope is defined by the catalog
	// price scope ("Catalog" > "Price" > "Catalog Price Scope").
	// Path: currency/options/base
	// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Base
	// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
	CurrencyOptionsBase ConfigCurrency

	// CurrencyOptionsDefault => Default Display Currency.
	// Path: currency/options/default
	// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\DefaultCurrency
	// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
	CurrencyOptionsDefault ConfigCurrency

	// CurrencyOptionsAllow => Allowed Currencies.
	// Path: currency/options/allow
	// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Allow
	// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
	CurrencyOptionsAllow model.StringCSV

	// CurrencyWebservicexTimeout => Connection Timeout in Seconds.
	// Path: currency/webservicex/timeout
	CurrencyWebservicexTimeout model.Str

	// CurrencyImportEnabled => Enabled.
	// Path: currency/import/enabled
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CurrencyImportEnabled model.Bool

	// CurrencyImportErrorEmail => Error Email Recipient.
	// Path: currency/import/error_email
	CurrencyImportErrorEmail model.Str

	// CurrencyImportErrorEmailIdentity => Error Email Sender.
	// Path: currency/import/error_email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	CurrencyImportErrorEmailIdentity model.Str

	// CurrencyImportErrorEmailTemplate => Error Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: currency/import/error_email_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	CurrencyImportErrorEmailTemplate model.Str

	// CurrencyImportFrequency => Frequency.
	// Path: currency/import/frequency
	// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
	CurrencyImportFrequency model.Str

	// CurrencyImportService => Service.
	// Path: currency/import/service
	// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Cron
	// SourceModel: Otnegam\Directory\Model\Currency\Import\Source\Service
	CurrencyImportService model.Str

	// CurrencyImportTime => Start Time.
	// Path: currency/import/time
	CurrencyImportTime model.Str

	// SystemCurrencyInstalled => Installed Currencies.
	// Defines all installed and available currencies.
	// Path: system/currency/installed
	// BackendModel: Otnegam\Config\Model\Config\Backend\Locale
	// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency\All
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
	SystemCurrencyInstalled model.StringCSV

	// GeneralCountryOptionalZipCountries => Zip/Postal Code is Optional for.
	// Path: general/country/optional_zip_countries
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	GeneralCountryOptionalZipCountries model.StringCSV

	// PathGeneralCountryAllow per store view allowed list of countries
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	GeneralCountryAllow model.StringCSV

	// PathGeneralCountryDefault returns the store view default configured country code.
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	GeneralCountryDefault model.Str

	// PathGeneralCountryEuCountries => European Union Countries.
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	GeneralCountryEuCountries model.StringCSV

	// PathGeneralCountryDestinations contains codes of the most used countries.
	// Such countries can be shown on the top of the country list.
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	GeneralCountryDestinations model.StringCSV

	// PathGeneralLocaleTimezone => Timezone.
	// BackendModel: Otnegam\Config\Model\Config\Backend\Locale\Timezone
	// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Timezone
	GeneralLocaleTimezone model.Str

	// PathGeneralLocaleCode => Locale.
	// SourceModel: Otnegam\Config\Model\Config\Source\Locale
	GeneralLocaleCode model.Str

	// PathGeneralLocaleFirstday => First Day of Week.
	// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Weekdays
	GeneralLocaleFirstday model.Str

	// PathGeneralLocaleWeekend => Weekend Days.
	// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Weekdays
	GeneralLocaleWeekend model.StringCSV

	// GeneralRegionStateRequired => State is Required for.
	// Path: general/region/state_required
	// SourceModel: Otnegam\Directory\Model\Config\Source\Country
	GeneralRegionStateRequired model.StringCSV

	// GeneralRegionDisplayAll => Allow to Choose State if It is Optional for Country.
	// Path: general/region/display_all
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	GeneralRegionDisplayAll model.Bool

	// GeneralLocaleWeightUnit => Weight Unit.
	// Path: general/locale/weight_unit
	// SourceModel: Otnegam\Directory\Model\Config\Source\WeightUnit
	GeneralLocaleWeightUnit model.Str
}

// NewBackend initializes the global Path variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CurrencyOptionsBase = NewConfigCurrency(`currency/options/base`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyOptionsDefault = NewConfigCurrency(`currency/options/default`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyOptionsAllow = model.NewStringCSV(`currency/options/allow`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyWebservicexTimeout = model.NewStr(`currency/webservicex/timeout`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyImportEnabled = model.NewBool(`currency/import/enabled`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyImportErrorEmail = model.NewStr(`currency/import/error_email`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyImportErrorEmailIdentity = model.NewStr(`currency/import/error_email_identity`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyImportErrorEmailTemplate = model.NewStr(`currency/import/error_email_template`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyImportFrequency = model.NewStr(`currency/import/frequency`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyImportService = model.NewStr(`currency/import/service`, model.WithConfigStructure(cfgStruct))
	pp.CurrencyImportTime = model.NewStr(`currency/import/time`, model.WithConfigStructure(cfgStruct))
	pp.SystemCurrencyInstalled = model.NewStringCSV(`system/currency/installed`, model.WithConfigStructure(cfgStruct))
	pp.GeneralCountryOptionalZipCountries = model.NewStringCSV(`general/country/optional_zip_countries`, model.WithConfigStructure(cfgStruct))
	pp.GeneralCountryAllow = model.NewStringCSV(`general/country/allow`, model.WithConfigStructure(cfgStruct))
	pp.GeneralCountryDefault = model.NewStr(`general/country/default`, model.WithConfigStructure(cfgStruct))
	pp.GeneralCountryEuCountries = model.NewStringCSV(`general/country/eu_countries`, model.WithConfigStructure(cfgStruct))
	pp.GeneralCountryDestinations = model.NewStringCSV(`general/country/destinations`, model.WithConfigStructure(cfgStruct))
	pp.GeneralLocaleTimezone = model.NewStr(`general/locale/timezone`, model.WithConfigStructure(cfgStruct))
	pp.GeneralLocaleCode = model.NewStr(`general/locale/code`, model.WithConfigStructure(cfgStruct))
	pp.GeneralLocaleFirstday = model.NewStr(`general/locale/firstday`, model.WithConfigStructure(cfgStruct))
	pp.GeneralLocaleWeekend = model.NewStringCSV(`general/locale/weekend`, model.WithConfigStructure(cfgStruct))
	pp.GeneralRegionStateRequired = model.NewStringCSV(`general/region/state_required`, model.WithConfigStructure(cfgStruct))
	pp.GeneralRegionDisplayAll = model.NewBool(`general/region/display_all`, model.WithConfigStructure(cfgStruct))
	pp.GeneralLocaleWeightUnit = model.NewStr(`general/locale/weight_unit`, model.WithConfigStructure(cfgStruct))

	return pp
}
