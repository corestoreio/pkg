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
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID:        path.NewRoute("currency"),
			Label:     text.Chars(`Currency Setup`),
			SortOrder: 60,
			Scope:     scope.PermStore,
			Resource:  0, // Magento_Backend::currency
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("options"),
					Label:     text.Chars(`Currency Options`),
					SortOrder: 30,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: currency/options/base
							ID:        path.NewRoute("base"),
							Label:     text.Chars(`Base Currency`),
							Comment:   text.Chars(`Base currency is used for all online payment transactions. If you have more than one store view, the base currency scope is defined by the catalog price scope ("Catalog" > "Price" > "Catalog Price Scope").`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `USD`,
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\Base
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
						},

						&element.Field{
							// Path: currency/options/default
							ID:        path.NewRoute("default"),
							Label:     text.Chars(`Default Display Currency`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							Default:   `USD`,
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\DefaultCurrency
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
						},

						&element.Field{
							// Path: currency/options/allow
							ID:         path.NewRoute("allow"),
							Label:      text.Chars(`Allowed Currencies`),
							Type:       element.TypeMultiselect,
							SortOrder:  3,
							Visible:    element.VisibleYes,
							Scope:      scope.PermStore,
							CanBeEmpty: true,
							Default:    `USD,EUR`,
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\Allow
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("webservicex"),
					Label:     text.Chars(`Webservicex`),
					SortOrder: 40,
					Scope:     scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: currency/webservicex/timeout
							ID:      path.NewRoute("timeout"),
							Label:   text.Chars(`Connection Timeout in Seconds`),
							Type:    element.TypeText,
							Visible: element.VisibleYes,
							Scope:   scope.PermDefault,
							Default: 100,
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("import"),
					Label:     text.Chars(`Scheduled Import Settings`),
					SortOrder: 50,
					Scope:     scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: currency/import/enabled
							ID:        path.NewRoute("enabled"),
							Label:     text.Chars(`Enabled`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: currency/import/error_email
							ID:        path.NewRoute("error_email"),
							Label:     text.Chars(`Error Email Recipient`),
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},

						&element.Field{
							// Path: currency/import/error_email_identity
							ID:        path.NewRoute("error_email_identity"),
							Label:     text.Chars(`Error Email Sender`),
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: currency/import/error_email_template
							ID:        path.NewRoute("error_email_template"),
							Label:     text.Chars(`Error Email Template`),
							Comment:   text.Chars(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `currency_import_error_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: currency/import/frequency
							ID:        path.NewRoute("frequency"),
							Label:     text.Chars(`Frequency`),
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Cron\Model\Config\Source\Frequency
						},

						&element.Field{
							// Path: currency/import/service
							ID:        path.NewRoute("service"),
							Label:     text.Chars(`Service`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\Cron
							// SourceModel: Magento\Directory\Model\Currency\Import\Source\Service
						},

						&element.Field{
							// Path: currency/import/time
							ID:        path.NewRoute("time"),
							Label:     text.Chars(`Start Time`),
							Type:      element.TypeTime,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
						},
					),
				},
			),
		},
		&element.Section{
			ID: path.NewRoute("system"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("currency"),
					Label:     text.Chars(`Currency`),
					SortOrder: 50,
					Scope:     scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/currency/installed
							ID:         path.NewRoute("installed"),
							Label:      text.Chars(`Installed Currencies`),
							Type:       element.TypeMultiselect,
							SortOrder:  1,
							Visible:    element.VisibleYes,
							Scope:      scope.PermDefault,
							CanBeEmpty: true,
							Default:    `AZN,AZM,AFN,ALL,DZD,AOA,ARS,AMD,AWG,AUD,BSD,BHD,BDT,BBD,BYR,BZD,BMD,BTN,BOB,BAM,BWP,BRL,GBP,BND,BGN,BUK,BIF,KHR,CAD,CVE,CZK,KYD,CLP,CNY,COP,KMF,CDF,CRC,HRK,CUP,DKK,DJF,DOP,XCD,EGP,SVC,GQE,ERN,EEK,ETB,EUR,FKP,FJD,GMD,GEK,GEL,GHS,GIP,GTQ,GNF,GYD,HTG,HNL,HKD,HUF,ISK,INR,IDR,IRR,IQD,ILS,JMD,JPY,JOD,KZT,KES,KWD,KGS,LAK,LVL,LBP,LSL,LRD,LYD,LTL,MOP,MKD,MGA,MWK,MYR,MVR,LSM,MRO,MUR,MXN,MDL,MNT,MAD,MZN,MMK,NAD,NPR,ANG,TRL,TRY,NZD,NIC,NGN,KPW,NOK,OMR,PKR,PAB,PGK,PYG,PEN,PHP,PLN,QAR,RHD,RON,ROL,RUB,RWF,SHP,STD,SAR,RSD,SCR,SLL,SGD,SKK,SBD,SOS,ZAR,KRW,LKR,SDG,SRD,SZL,SEK,CHF,SYP,TWD,TJS,TZS,THB,TOP,TTD,TND,TMM,USD,UGX,UAH,AED,UYU,UZS,VUV,VEB,VEF,VND,CHE,CHW,XOF,XPF,WST,YER,ZMK,ZWD`,
							// BackendModel: Magento\Config\Model\Config\Backend\Locale
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency\All
						},
					),
				},
			),
		},
		&element.Section{
			ID: path.NewRoute("general"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("country"),
					Label:     text.Chars(`Country Options`),
					SortOrder: 1,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: general/country/allow
							ID:         path.NewRoute("allow"),
							Label:      text.Chars(`Allow Countries`),
							Type:       element.TypeMultiselect,
							SortOrder:  2,
							Visible:    element.VisibleYes,
							Scope:      scope.PermStore,
							CanBeEmpty: true,
							Default:    `AF,AL,DZ,AS,AD,AO,AI,AQ,AG,AR,AM,AW,AU,AT,AX,AZ,BS,BH,BD,BB,BY,BE,BZ,BJ,BM,BL,BT,BO,BA,BW,BV,BR,IO,VG,BN,BG,BF,BI,KH,CM,CA,CD,CV,KY,CF,TD,CL,CN,CX,CC,CO,KM,CG,CK,CR,HR,CU,CY,CZ,DK,DJ,DM,DO,EC,EG,SV,GQ,ER,EE,ET,FK,FO,FJ,FI,FR,GF,PF,TF,GA,GM,GE,DE,GG,GH,GI,GR,GL,GD,GP,GU,GT,GN,GW,GY,HT,HM,HN,HK,HU,IS,IM,IN,ID,IR,IQ,IE,IL,IT,CI,JE,JM,JP,JO,KZ,KE,KI,KW,KG,LA,LV,LB,LS,LR,LY,LI,LT,LU,ME,MF,MO,MK,MG,MW,MY,MV,ML,MT,MH,MQ,MR,MU,YT,FX,MX,FM,MD,MC,MN,MS,MA,MZ,MM,NA,NR,NP,NL,AN,NC,NZ,NI,NE,NG,NU,NF,KP,MP,NO,OM,PK,PW,PA,PG,PY,PE,PH,PN,PL,PS,PT,PR,QA,RE,RO,RS,RU,RW,SH,KN,LC,PM,VC,WS,SM,ST,SA,SN,SC,SL,SG,SK,SI,SB,SO,ZA,GS,KR,ES,LK,SD,SR,SJ,SZ,SE,CH,SY,TL,TW,TJ,TZ,TH,TG,TK,TO,TT,TN,TR,TM,TC,TV,VI,UG,UA,AE,GB,US,UM,UY,UZ,VU,VA,VE,VN,WF,EH,YE,ZM,ZW`,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: general/country/default
							ID:        path.NewRoute("default"),
							Label:     text.Chars(`Default Country`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: general/country/eu_countries
							ID:        path.NewRoute("eu_countries"),
							Label:     text.Chars(`European Union Countries`),
							Type:      element.TypeMultiselect,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: general/country/destinations
							ID:        path.NewRoute("destinations"),
							Label:     text.Chars(`Top destinations`),
							Comment:   text.Chars(`Contains codes of the most used countries. Such countries can be shown on the top of the country list.`),
							Type:      element.TypeMultiselect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},
						&element.Field{
							// Path: general/country/optional_zip_countries
							ID:         path.NewRoute("optional_zip_countries"),
							Label:      text.Chars(`Zip/Postal Code is Optional for`),
							Type:       element.TypeMultiselect,
							SortOrder:  3,
							Visible:    element.VisibleYes,
							Scope:      scope.PermDefault,
							CanBeEmpty: true,
							Default:    `HK,IE,MO,PA,GB`,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("locale"),
					Label:     text.Chars(`Locale Options`),
					SortOrder: 8,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: general/locale/timezone
							ID:        path.NewRoute("timezone"),
							Label:     text.Chars(`Timezone`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermWebsite,
							Default:   `America/Los_Angeles`,
							// BackendModel: Magento\Config\Model\Config\Backend\Locale\Timezone
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Timezone
						},

						&element.Field{
							// Path: general/locale/code
							ID:        path.NewRoute("code"),
							Label:     text.Chars(`Locale`),
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							Default:   `en_US`,
							// SourceModel: Magento\Config\Model\Config\Source\Locale
						},

						&element.Field{
							// Path: general/locale/firstday
							ID:        path.NewRoute("firstday"),
							Label:     text.Chars(`First Day of Week`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Weekdays
						},

						&element.Field{
							// Path: general/locale/weekend
							ID:         path.NewRoute("weekend"),
							Label:      text.Chars(`Weekend Days`),
							Type:       element.TypeMultiselect,
							SortOrder:  15,
							Visible:    element.VisibleYes,
							Scope:      scope.PermStore,
							CanBeEmpty: true,
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Weekdays
						},
						&element.Field{
							// Path: general/locale/weight_unit
							ID:        path.NewRoute("weight_unit"),
							Label:     text.Chars(`Weight Unit`),
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.PermStore,
							Default:   `lbs`,
							// SourceModel: Magento\Directory\Model\Config\Source\WeightUnit
						},
					),
				},
				&element.Group{
					ID:        path.NewRoute("region"),
					Label:     text.Chars(`State Options`),
					SortOrder: 4,
					Scope:     scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: general/region/state_required
							ID:        path.NewRoute("state_required"),
							Label:     text.Chars(`State is Required for`),
							Type:      element.TypeMultiselect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						&element.Field{
							// Path: general/region/display_all
							ID:        path.NewRoute("display_all"),
							Label:     text.Chars(`Allow to Choose State if It is Optional for Country`),
							Type:      element.TypeSelect,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scope:     scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: path.NewRoute("general"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: path.NewRoute("locale"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: general/locale/datetime_format_long
							ID:      path.NewRoute("datetime_format_long"),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%A, %B %e %Y [%I:%M %p]`,
						},

						&element.Field{
							// Path: general/locale/datetime_format_medium
							ID:      path.NewRoute("datetime_format_medium"),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%a, %b %e %Y [%I:%M %p]`,
						},

						&element.Field{
							// Path: general/locale/datetime_format_short
							ID:      path.NewRoute("datetime_format_short"),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%m/%d/%y [%I:%M %p]`,
						},

						&element.Field{
							// Path: general/locale/date_format_long
							ID:      path.NewRoute("date_format_long"),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%A, %B %e %Y`,
						},

						&element.Field{
							// Path: general/locale/date_format_medium
							ID:      path.NewRoute("date_format_medium"),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%a, %b %e %Y`,
						},

						&element.Field{
							// Path: general/locale/date_format_short
							ID:      path.NewRoute("date_format_short"),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%m/%d/%y`,
						},

						&element.Field{
							// Path: general/locale/language
							ID:      path.NewRoute("language"),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `en`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
