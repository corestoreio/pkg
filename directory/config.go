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
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/storage/csdb"
)

const (
	PathSystemCurrencyInstalled = "system/currency/installed"

	// PathCurrencyBase defines the app base currency code
	PathCurrencyBase    = "currency/options/base"
	PathCurrencyDefault = "currency/options/default"
	PathCurrencyAllow   = "currency/options/allow"

	// PathOptionalZipCountries lists ISO2 country codes which have optional Zip/Postal pre-configured
	PathOptionalZipCountries = "general/country/optional_zip_countries"
	// PathStatesRequired lists countries, for which state is required. No default values.
	PathStatesRequired = "general/region/state_required"
	// PathDisplayAllStates detects whether or not display the state for the country, if it is not required
	PathDisplayAllStates = "general/region/display_all"
	PathDefaultCountry   = "general/country/default"
	PathDefaultLocale    = "general/locale/code"
	PathDefaultTimezone  = "general/locale/timezone"
)

// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
var TableCollection csdb.TableStructurer

// PackageConfiguration contains the main configuration
var PackageConfiguration config.SectionSlice

func init() {
	PackageConfiguration = config.NewConfiguration(
		&config.Section{
			ID:        "currency",
			Label:     "Currency Setup",
			SortOrder: 60,
			Scope:     config.ScopePermAll,
			Groups: config.GroupSlice{
				&config.Group{
					ID:        "options",
					Label:     `Currency Options`,
					Comment:   ``,
					SortOrder: 30,
					Scope:     config.ScopePermAll,
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `currency/options/base`,
							ID:           "base",
							Label:        `Base Currency`,
							Comment:      `Base currency is used for all online payment transactions. If you have more than one store view, the base currency scope is defined by the catalog price scope ("Catalog" > "Price" > "Catalog Price Scope").`,
							Type:         config.TypeSelect,
							SortOrder:    1,
							Visible:      config.VisibleYes,
							Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
							Default:      `USD`,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Currency\Base
							SourceModel:  nil, // Magento\Config\Model\Config\Source\Locale\Currency
						},

						&config.Field{
							// Path: `currency/options/default`,
							ID:           "default",
							Label:        `Default Display Currency`,
							Comment:      ``,
							Type:         config.TypeSelect,
							SortOrder:    2,
							Visible:      config.VisibleYes,
							Scope:        config.ScopePermAll,
							Default:      `USD`,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Currency\DefaultCurrency
							SourceModel:  nil, // Magento\Config\Model\Config\Source\Locale\Currency
						},

						&config.Field{
							// Path: `currency/options/allow`,
							ID:           "allow",
							Label:        `Allowed Currencies`,
							Comment:      ``,
							Type:         config.TypeMultiselect,
							SortOrder:    3,
							Visible:      config.VisibleYes,
							Scope:        config.ScopePermAll,
							Default:      `USD,EUR`,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Currency\Allow
							SourceModel:  nil, // Magento\Config\Model\Config\Source\Locale\Currency
						},
					},
				},

				&config.Group{
					ID:        "webservicex",
					Label:     `Webservicex`,
					Comment:   ``,
					SortOrder: 40,
					Scope:     config.NewScopePerm(config.ScopeDefaultID),
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `currency/webservicex/timeout`,
							ID:           "timeout",
							Label:        `Connection Timeout in Seconds`,
							Comment:      ``,
							Type:         config.TypeText,
							SortOrder:    0,
							Visible:      config.VisibleYes,
							Scope:        config.NewScopePerm(config.ScopeDefaultID),
							Default:      100,
							BackendModel: nil,
							SourceModel:  nil,
						},
					},
				},

				&config.Group{
					ID:        "import",
					Label:     `Scheduled Import Settings`,
					Comment:   ``,
					SortOrder: 50,
					Scope:     config.NewScopePerm(config.ScopeDefaultID),
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `currency/import/enabled`,
							ID:           "enabled",
							Label:        `Enabled`,
							Comment:      ``,
							Type:         config.TypeSelect,
							SortOrder:    1,
							Visible:      config.VisibleYes,
							Scope:        config.ScopePermAll,
							Default:      false,
							BackendModel: nil,
							SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
						},

						&config.Field{
							// Path: `currency/import/error_email`,
							ID:           "error_email",
							Label:        `Error Email Recipient`,
							Comment:      ``,
							Type:         config.TypeText,
							SortOrder:    5,
							Visible:      config.VisibleYes,
							Scope:        config.ScopePermAll,
							Default:      nil,
							BackendModel: nil,
							SourceModel:  nil,
						},

						&config.Field{
							// Path: `currency/import/error_email_identity`,
							ID:           "error_email_identity",
							Label:        `Error Email Sender`,
							Comment:      ``,
							Type:         config.TypeSelect,
							SortOrder:    6,
							Visible:      config.VisibleYes,
							Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
							Default:      `general`,
							BackendModel: nil,
							SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Identity
						},

						&config.Field{
							// Path: `currency/import/error_email_template`,
							ID:           "error_email_template",
							Label:        `Error Email Template`,
							Comment:      ``,
							Type:         config.TypeSelect,
							SortOrder:    7,
							Visible:      config.VisibleYes,
							Scope:        config.NewScopePerm(config.ScopeDefaultID, config.ScopeWebsiteID),
							Default:      `currency_import_error_email_template`,
							BackendModel: nil,
							SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
						},

						&config.Field{
							// Path: `currency/import/frequency`,
							ID:           "frequency",
							Label:        `Frequency`,
							Comment:      ``,
							Type:         config.TypeSelect,
							SortOrder:    4,
							Visible:      config.VisibleYes,
							Scope:        config.ScopePermAll,
							Default:      nil,
							BackendModel: nil,
							SourceModel:  nil, // Magento\Cron\Model\Config\Source\Frequency
						},

						&config.Field{
							// Path: `currency/import/service`,
							ID:           "service",
							Label:        `Service`,
							Comment:      ``,
							Type:         config.TypeSelect,
							SortOrder:    2,
							Visible:      config.VisibleYes,
							Scope:        config.ScopePermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Currency\Cron
							SourceModel:  nil, // Magento\Directory\Model\Currency\Import\Source\Service
						},

						&config.Field{
							// Path: `currency/import/time`,
							ID:           "time",
							Label:        `Start Time`,
							Comment:      ``,
							Type:         config.TypeTime,
							SortOrder:    3,
							Visible:      config.VisibleYes,
							Scope:        config.ScopePermAll,
							Default:      nil,
							BackendModel: nil,
							SourceModel:  nil,
						},
					},
				},
			},
		},
		&config.Section{
			ID:    "system",
			Scope: config.NewScopePerm(),
			Groups: config.GroupSlice{
				&config.Group{
					ID:        "currency",
					Label:     `Currency`,
					Comment:   ``,
					SortOrder: 50,
					Scope:     config.NewScopePerm(config.ScopeDefaultID),
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `system/currency/installed`,
							ID:           "installed",
							Label:        `Installed Currencies`,
							Comment:      ``,
							Type:         config.TypeMultiselect,
							SortOrder:    1,
							Visible:      config.VisibleYes,
							Scope:        config.NewScopePerm(config.ScopeDefaultID),
							Default:      `AZN,AZM,AFN,ALL,DZD,AOA,ARS,AMD,AWG,AUD,BSD,BHD,BDT,BBD,BYR,BZD,BMD,BTN,BOB,BAM,BWP,BRL,GBP,BND,BGN,BUK,BIF,KHR,CAD,CVE,CZK,KYD,CLP,CNY,COP,KMF,CDF,CRC,HRK,CUP,DKK,DJF,DOP,XCD,EGP,SVC,GQE,ERN,EEK,ETB,EUR,FKP,FJD,GMD,GEK,GEL,GHS,GIP,GTQ,GNF,GYD,HTG,HNL,HKD,HUF,ISK,INR,IDR,IRR,IQD,ILS,JMD,JPY,JOD,KZT,KES,KWD,KGS,LAK,LVL,LBP,LSL,LRD,LYD,LTL,MOP,MKD,MGA,MWK,MYR,MVR,LSM,MRO,MUR,MXN,MDL,MNT,MAD,MZN,MMK,NAD,NPR,ANG,TRL,TRY,NZD,NIC,NGN,KPW,NOK,OMR,PKR,PAB,PGK,PYG,PEN,PHP,PLN,QAR,RHD,RON,ROL,RUB,RWF,SHP,STD,SAR,RSD,SCR,SLL,SGD,SKK,SBD,SOS,ZAR,KRW,LKR,SDG,SRD,SZL,SEK,CHF,SYP,TWD,TJS,TZS,THB,TOP,TTD,TND,TMM,USD,UGX,UAH,AED,UYU,UZS,VUV,VEB,VEF,VND,CHE,CHW,XOF,XPF,WST,YER,ZMK,ZWD`,
							BackendModel: nil,                    // Magento\Config\Model\Config\Backend\Locale
							SourceModel:  NewSourceCurrencyAll(), // Magento\Config\Model\Config\Source\Locale\Currency\All
						},
					},
				},
			},
		},
		&config.Section{
			ID:    "general",
			Scope: config.NewScopePerm(),
			Groups: config.GroupSlice{
				&config.Group{
					ID:    "country",
					Scope: config.NewScopePerm(),
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `general/country/optional_zip_countries`,
							ID:           "optional_zip_countries",
							Label:        `Zip/Postal Code is Optional for`,
							Comment:      ``,
							Type:         config.TypeMultiselect,
							SortOrder:    3,
							Visible:      config.VisibleYes,
							Scope:        config.NewScopePerm(config.ScopeDefaultID),
							Default:      `HK,IE,MO,PA,GB`,
							BackendModel: nil,
							SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
						},
					},
				},

				&config.Group{
					ID:        "region",
					Label:     `State Options`,
					SortOrder: 4,
					Scope:     config.NewScopePerm(config.ScopeDefaultID),
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `general/region/state_required`,
							ID:           "state_required",
							Label:        `State is Required for`,
							Comment:      ``,
							Type:         config.TypeMultiselect,
							SortOrder:    1,
							Visible:      config.VisibleYes,
							Scope:        config.NewScopePerm(config.ScopeDefaultID),
							Default:      nil,
							BackendModel: nil,
							SourceModel:  nil, // Magento\Directory\Model\Config\Source\Country
						},

						&config.Field{
							// Path: `general/region/display_all`,
							ID:           "display_all",
							Label:        `Allow to Choose State if It is Optional for Country`,
							Comment:      ``,
							Type:         config.TypeSelect,
							SortOrder:    8,
							Visible:      config.VisibleYes,
							Scope:        config.NewScopePerm(config.ScopeDefaultID),
							Default:      nil,
							BackendModel: nil,
							SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
						},
					},
				},
			},
		},

		// Hidden Configuration
		&config.Section{
			ID: "general",
			Groups: config.GroupSlice{
				&config.Group{
					ID: "country",
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `general/country/allow`,
							ID:      "allow",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `AF,AL,DZ,AS,AD,AO,AI,AQ,AG,AR,AM,AW,AU,AT,AX,AZ,BS,BH,BD,BB,BY,BE,BZ,BJ,BM,BL,BT,BO,BA,BW,BV,BR,IO,VG,BN,BG,BF,BI,KH,CM,CA,CD,CV,KY,CF,TD,CL,CN,CX,CC,CO,KM,CG,CK,CR,HR,CU,CY,CZ,DK,DJ,DM,DO,EC,EG,SV,GQ,ER,EE,ET,FK,FO,FJ,FI,FR,GF,PF,TF,GA,GM,GE,DE,GG,GH,GI,GR,GL,GD,GP,GU,GT,GN,GW,GY,HT,HM,HN,HK,HU,IS,IM,IN,ID,IR,IQ,IE,IL,IT,CI,JE,JM,JP,JO,KZ,KE,KI,KW,KG,LA,LV,LB,LS,LR,LY,LI,LT,LU,ME,MF,MO,MK,MG,MW,MY,MV,ML,MT,MH,MQ,MR,MU,YT,FX,MX,FM,MD,MC,MN,MS,MA,MZ,MM,NA,NR,NP,NL,AN,NC,NZ,NI,NE,NG,NU,NF,KP,MP,NO,OM,PK,PW,PA,PG,PY,PE,PH,PN,PL,PS,PT,PR,QA,RE,RO,RS,RU,RW,SH,KN,LC,PM,VC,WS,SM,ST,SA,SN,SC,SL,SG,SK,SI,SB,SO,ZA,GS,KR,ES,LK,SD,SR,SJ,SZ,SE,CH,SY,TL,TW,TJ,TZ,TH,TG,TK,TO,TT,TN,TR,TM,TC,TV,VI,UG,UA,AE,GB,US,UM,UY,UZ,VU,VA,VE,VN,WF,EH,YE,ZM,ZW`,
						},

						&config.Field{
							// Path: `general/country/default`,
							ID:      "default",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `US`,
						},
					},
				},

				&config.Group{
					ID: "locale",
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `general/locale/datetime_format_long`,
							ID:      "datetime_format_long",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `%A, %B %e %Y [%I:%M %p]`,
						},

						&config.Field{
							// Path: `general/locale/datetime_format_medium`,
							ID:      "datetime_format_medium",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `%a, %b %e %Y [%I:%M %p]`,
						},

						&config.Field{
							// Path: `general/locale/datetime_format_short`,
							ID:      "datetime_format_short",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `%m/%d/%y [%I:%M %p]`,
						},

						&config.Field{
							// Path: `general/locale/date_format_long`,
							ID:      "date_format_long",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `%A, %B %e %Y`,
						},

						&config.Field{
							// Path: `general/locale/date_format_medium`,
							ID:      "date_format_medium",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `%a, %b %e %Y`,
						},

						&config.Field{
							// Path: `general/locale/date_format_short`,
							ID:      "date_format_short",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `%m/%d/%y`,
						},

						&config.Field{
							// Path: `general/locale/language`,
							ID:      "language",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `en`,
						},

						&config.Field{
							// Path: `general/locale/code`,
							ID:      "code",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `en_US`,
						},

						&config.Field{
							// Path: `general/locale/timezone`,
							ID:      "timezone",
							Type:    config.TypeHidden,
							Visible: config.VisibleNo,
							Scope:   config.NewScopePerm(config.ScopeDefaultID), // @todo search for that
							Default: `America/Los_Angeles`,
						},
					},
				},
			},
		},
	)
}
