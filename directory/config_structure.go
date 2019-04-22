// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
)

// NewConfigStructure global configuration structure for this package.
// Used in frontend (to display the user all the settings) and in
// backend (scope checks and default values). See the source code
// of this function for the overall available sections, groups and fields.
func NewConfigStructure() (config.Sections, error) {
	return config.MakeSectionsValidated(
		&config.Section{
			ID:        "currency",
			Label:     `Currency Setup`,
			SortOrder: 60,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::currency
			Groups: config.MakeGroups(
				&config.Group{
					ID:        "options",
					Label:     `Currency Options`,
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: config.MakeFields(
						&config.Field{
							// Path: currency/options/base
							ID:        "base",
							Label:     `Base Currency`,
							Comment:   `Base currency is used for all online payment transactions. If you have more than one store view, the base currency scope is defined by the catalog price scope ("Catalog" > "Price" > "Catalog Price Scope".`,
							Type:      config.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermWebsite,
							Default:   []byte(`USD`),
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\Base
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
						},

						&config.Field{
							// Path: currency/options/default
							ID:        "default",
							Label:     `Default Display Currency`,
							Type:      config.TypeSelect,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermStore,
							Default:   []byte(`USD`),
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\DefaultCurrency
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
						},

						&config.Field{
							// Path: currency/options/allow
							ID:         "allow",
							Label:      `Allowed Currencies`,
							Type:       config.TypeMultiselect,
							SortOrder:  3,
							Visible:    true,
							Scopes:     scope.PermStore,
							CanBeEmpty: true,
							Default:    []byte(`USD,EUR`),

							// BackendModel: Magento\Config\Model\Config\Backend\Currency\Allow
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
						},
					),
				},

				&config.Group{
					ID:        "webservicex",
					Label:     `Webservicex`,
					SortOrder: 40,
					Scopes:    scope.PermDefault,
					Fields: config.MakeFields(
						&config.Field{
							// Path: currency/webservicex/timeout
							ID:      "timeout",
							Label:   `Connection Timeout in Seconds`,
							Type:    config.TypeText,
							Visible: true,
							Scopes:  scope.PermDefault,
							Default: []byte(`100`),
						},
					),
				},

				&config.Group{
					ID:        "import",
					Label:     `Scheduled Import Settings`,
					SortOrder: 50,
					Scopes:    scope.PermDefault,
					Fields: config.MakeFields(
						&config.Field{
							// Path: currency/import/enabled
							ID:        "enabled",
							Label:     `Enabled`,
							Type:      config.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermStore,
							Default:   []byte(`false`),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&config.Field{
							// Path: currency/import/error_email
							ID:        "error_email",
							Label:     `Error Email Recipient`,
							Type:      config.TypeText,
							SortOrder: 5,
							Visible:   true,
							Scopes:    scope.PermStore,
						},

						&config.Field{
							// Path: currency/import/error_email_identity
							ID:        "error_email_identity",
							Label:     `Error Email Sender`,
							Type:      config.TypeSelect,
							SortOrder: 6,
							Visible:   true,
							Scopes:    scope.PermWebsite,
							Default:   []byte(`general`),
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						&config.Field{
							// Path: currency/import/error_email_template
							ID:        "error_email_template",
							Label:     `Error Email Template`,
							Comment:   `Email template chosen based on theme fallback when "Default" option is selected.`,
							Type:      config.TypeSelect,
							SortOrder: 7,
							Visible:   true,
							Scopes:    scope.PermWebsite,
							Default:   []byte(`currency_import_error_email_template`),
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						&config.Field{
							// Path: currency/import/frequency
							ID:        "frequency",
							Label:     `Frequency`,
							Type:      config.TypeSelect,
							SortOrder: 4,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Cron\Model\Config\Source\Frequency
						},

						&config.Field{
							// Path: currency/import/service
							ID:        "service",
							Label:     `Service`,
							Type:      config.TypeSelect,
							SortOrder: 2,
							Visible:   true,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\Cron
							// SourceModel: Magento\Directory\Model\Currency\Import\Source\Service
						},

						&config.Field{
							// Path: currency/import/time
							ID:        "time",
							Label:     `Start Time`,
							Type:      config.TypeTime,
							SortOrder: 3,
							Visible:   true,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
		&config.Section{
			ID: "system",
			Groups: config.MakeGroups(
				&config.Group{
					ID:        "currency",
					Label:     `Currency`,
					SortOrder: 50,
					Scopes:    scope.PermDefault,
					Fields: config.MakeFields(
						&config.Field{
							// Path: system/currency/installed
							ID:         "installed",
							Label:      `Installed Currencies`,
							Type:       config.TypeMultiselect,
							SortOrder:  1,
							Visible:    true,
							Scopes:     scope.PermDefault,
							CanBeEmpty: true,
							Default:    []byte(`AZN,AZM,AFN,ALL,DZD,AOA,ARS,AMD,AWG,AUD,BSD,BHD,BDT,BBD,BYR,BZD,BMD,BTN,BOB,BAM,BWP,BRL,GBP,BND,BGN,BUK,BIF,KHR,CAD,CVE,CZK,KYD,CLP,CNY,COP,KMF,CDF,CRC,HRK,CUP,DKK,DJF,DOP,XCD,EGP,SVC,GQE,ERN,EEK,ETB,EUR,FKP,FJD,GMD,GEK,GEL,GHS,GIP,GTQ,GNF,GYD,HTG,HNL,HKD,HUF,ISK,INR,IDR,IRR,IQD,ILS,JMD,JPY,JOD,KZT,KES,KWD,KGS,LAK,LVL,LBP,LSL,LRD,LYD,LTL,MOP,MKD,MGA,MWK,MYR,MVR,LSM,MRO,MUR,MXN,MDL,MNT,MAD,MZN,MMK,NAD,NPR,ANG,TRL,TRY,NZD,NIC,NGN,KPW,NOK,OMR,PKR,PAB,PGK,PYG,PEN,PHP,PLN,QAR,RHD,RON,ROL,RUB,RWF,SHP,STD,SAR,RSD,SCR,SLL,SGD,SKK,SBD,SOS,ZAR,KRW,LKR,SDG,SRD,SZL,SEK,CHF,SYP,TWD,TJS,TZS,THB,TOP,TTD,TND,TMM,USD,UGX,UAH,AED,UYU,UZS,VUV,VEB,VEF,VND,CHE,CHW,XOF,XPF,WST,YER,ZMK,ZWD`),

							// BackendModel: Magento\Config\Model\Config\Backend\Locale
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency\All
						},
					),
				},
			),
		},
		&config.Section{
			ID: "general",
			Groups: config.MakeGroups(
				&config.Group{
					ID:        "country",
					Label:     `Country Options`,
					SortOrder: 1,
					Scopes:    scope.PermStore,
					Fields: config.MakeFields(
						&config.Field{
							// Path: general/country/allow
							ID:         "allow",
							Label:      `Allow Countries`,
							Type:       config.TypeMultiselect,
							SortOrder:  2,
							Visible:    true,
							Scopes:     scope.PermStore,
							CanBeEmpty: true,
							Default:    []byte(`AF,AL,DZ,AS,AD,AO,AI,AQ,AG,AR,AM,AW,AU,AT,AX,AZ,BS,BH,BD,BB,BY,BE,BZ,BJ,BM,BL,BT,BO,BA,BW,BV,BR,IO,VG,BN,BG,BF,BI,KH,CM,CA,CD,CV,KY,CF,TD,CL,CN,CX,CC,CO,KM,CG,CK,CR,HR,CU,CY,CZ,DK,DJ,DM,DO,EC,EG,SV,GQ,ER,EE,ET,FK,FO,FJ,FI,FR,GF,PF,TF,GA,GM,GE,DE,GG,GH,GI,GR,GL,GD,GP,GU,GT,GN,GW,GY,HT,HM,HN,HK,HU,IS,IM,IN,ID,IR,IQ,IE,IL,IT,CI,JE,JM,JP,JO,KZ,KE,KI,KW,KG,LA,LV,LB,LS,LR,LY,LI,LT,LU,ME,MF,MO,MK,MG,MW,MY,MV,ML,MT,MH,MQ,MR,MU,YT,FX,MX,FM,MD,MC,MN,MS,MA,MZ,MM,NA,NR,NP,NL,AN,NC,NZ,NI,NE,NG,NU,NF,KP,MP,NO,OM,PK,PW,PA,PG,PY,PE,PH,PN,PL,PS,PT,PR,QA,RE,RO,RS,RU,RW,SH,KN,LC,PM,VC,WS,SM,ST,SA,SN,SC,SL,SG,SK,SI,SB,SO,ZA,GS,KR,ES,LK,SD,SR,SJ,SZ,SE,CH,SY,TL,TW,TJ,TZ,TH,TG,TK,TO,TT,TN,TR,TM,TC,TV,VI,UG,UA,AE,GB,US,UM,UY,UZ,VU,VA,VE,VN,WF,EH,YE,ZM,ZW`),

							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&config.Field{
							// Path: general/country/default
							ID:        "default",
							Label:     `Default Country`,
							Type:      config.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&config.Field{
							// Path: general/country/eu_countries
							ID:        "eu_countries",
							Label:     `European Union Countries`,
							Type:      config.TypeMultiselect,
							SortOrder: 30,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&config.Field{
							// Path: general/country/destinations
							ID:        "destinations",
							Label:     `Top destinations`,
							Comment:   `Contains codes of the most used countries. Such countries can be shown on the top of the country list.`,
							Type:      config.TypeMultiselect,
							SortOrder: 40,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},
						&config.Field{
							// Path: general/country/optional_zip_countries
							ID:         "optional_zip_countries",
							Label:      `Zip/Postal Code is Optional for`,
							Type:       config.TypeMultiselect,
							SortOrder:  3,
							Visible:    true,
							Scopes:     scope.PermDefault,
							CanBeEmpty: true,
							Default:    []byte(`HK,IE,MO,PA,GB`),

							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},
					),
				},

				&config.Group{
					ID:        "locale",
					Label:     `Locale Options`,
					SortOrder: 8,
					Scopes:    scope.PermStore,
					Fields: config.MakeFields(
						&config.Field{
							// Path: general/locale/timezone
							ID:        "timezone",
							Label:     `Timezone`,
							Type:      config.TypeSelect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermWebsite,
							Default:   []byte(`America/Los_Angeles`),
							// BackendModel: Magento\Config\Model\Config\Backend\Locale\Timezone
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Timezone
						},

						&config.Field{
							// Path: general/locale/code
							ID:        "code",
							Label:     `Locale`,
							Type:      config.TypeSelect,
							SortOrder: 5,
							Visible:   true,
							Scopes:    scope.PermStore,
							Default:   []byte(`en_US`),
							// SourceModel: Magento\Config\Model\Config\Source\Locale
						},

						&config.Field{
							// Path: general/locale/firstday
							ID:        "firstday",
							Label:     `First Day of Week`,
							Type:      config.TypeSelect,
							SortOrder: 10,
							Visible:   true,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Weekdays
						},

						&config.Field{
							// Path: general/locale/weekend
							ID:         "weekend",
							Label:      `Weekend Days`,
							Type:       config.TypeMultiselect,
							SortOrder:  15,
							Visible:    true,
							Scopes:     scope.PermStore,
							CanBeEmpty: true,
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Weekdays
						},
						&config.Field{
							// Path: general/locale/weight_unit
							ID:        "weight_unit",
							Label:     `Weight Unit`,
							Type:      config.TypeSelect,
							SortOrder: 7,
							Visible:   true,
							Scopes:    scope.PermStore,
							Default:   []byte(`lbs`),
							// SourceModel: Magento\Directory\Model\Config\Source\WeightUnit
						},
					),
				},
				&config.Group{
					ID:        "region",
					Label:     `State Options`,
					SortOrder: 4,
					Scopes:    scope.PermDefault,
					Fields: config.MakeFields(
						&config.Field{
							// Path: general/region/state_required
							ID:        "state_required",
							Label:     `State is Required for`,
							Type:      config.TypeMultiselect,
							SortOrder: 1,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&config.Field{
							// Path: general/region/display_all
							ID:        "display_all",
							Label:     `Allow to Choose State if It is Optional for Country`,
							Type:      config.TypeSelect,
							SortOrder: 8,
							Visible:   true,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&config.Section{
			ID: "general",
			Groups: config.MakeGroups(
				&config.Group{
					ID: "locale",
					Fields: config.MakeFields(
						&config.Field{
							// Path: general/locale/datetime_format_long
							ID:      "datetime_format_long",
							Type:    config.TypeHidden,
							Visible: false,
							Default: []byte(`%A, %B %e %Y [%I:%M %p]`),
						},

						&config.Field{
							// Path: general/locale/datetime_format_medium
							ID:      "datetime_format_medium",
							Type:    config.TypeHidden,
							Visible: false,
							Default: []byte(`%a, %b %e %Y [%I:%M %p]`),
						},

						&config.Field{
							// Path: general/locale/datetime_format_short
							ID:      "datetime_format_short",
							Type:    config.TypeHidden,
							Visible: false,
							Default: []byte(`%m/%d/%y [%I:%M %p]`),
						},

						&config.Field{
							// Path: general/locale/date_format_long
							ID:      "date_format_long",
							Type:    config.TypeHidden,
							Visible: false,
							Default: []byte(`%A, %B %e %Y`),
						},

						&config.Field{
							// Path: general/locale/date_format_medium
							ID:      "date_format_medium",
							Type:    config.TypeHidden,
							Visible: false,
							Default: []byte(`%a, %b %e %Y`),
						},

						&config.Field{
							// Path: general/locale/date_format_short
							ID:      "date_format_short",
							Type:    config.TypeHidden,
							Visible: false,
							Default: []byte(`%m/%d/%y`),
						},

						&config.Field{
							// Path: general/locale/language
							ID:      "language",
							Type:    config.TypeHidden,
							Visible: false,
							Default: []byte(`en`),
						},
					),
				},
			),
		},
	)
}
