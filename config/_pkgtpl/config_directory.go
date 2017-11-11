// +build ignore

package directory

import (
	"github.com/corestoreio/cspkg/config/element"
	"github.com/corestoreio/cspkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		element.Section{
			ID:        "currency",
			Label:     `Currency Setup`,
			SortOrder: 60,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Backend::currency
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "options",
					Label:     `Currency Options`,
					SortOrder: 30,
					Scopes:    scope.PermStore,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: currency/options/base
							ID:        "base",
							Label:     `Base Currency`,
							Comment:   text.Long(`Base currency is used for all online payment transactions. If you have more than one store view, the base currency scope is defined by the catalog price scope ("Catalog" > "Price" > "Catalog Price Scope").`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `USD`,
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\Base
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
						},

						element.Field{
							// Path: currency/options/default
							ID:        "default",
							Label:     `Default Display Currency`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `USD`,
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\DefaultCurrency
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
						},

						element.Field{
							// Path: currency/options/allow
							ID:         "allow",
							Label:      `Allowed Currencies`,
							Type:       element.TypeMultiselect,
							SortOrder:  3,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermStore,
							CanBeEmpty: true,
							Default:    `USD,EUR`,
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\Allow
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency
						},
					),
				},

				element.Group{
					ID:        "webservicex",
					Label:     `Webservicex`,
					SortOrder: 40,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: currency/webservicex/timeout
							ID:      "timeout",
							Label:   `Connection Timeout in Seconds`,
							Type:    element.TypeText,
							Visible: element.VisibleYes,
							Scopes:  scope.PermDefault,
							Default: 100,
						},
					),
				},

				element.Group{
					ID:        "import",
					Label:     `Scheduled Import Settings`,
					SortOrder: 50,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: currency/import/enabled
							ID:        "enabled",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: currency/import/error_email
							ID:        "error_email",
							Label:     `Error Email Recipient`,
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: currency/import/error_email_identity
							ID:        "error_email_identity",
							Label:     `Error Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `general`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Identity
						},

						element.Field{
							// Path: currency/import/error_email_template
							ID:        "error_email_template",
							Label:     `Error Email Template`,
							Comment:   text.Long(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermWebsite,
							Default:   `currency_import_error_email_template`,
							// SourceModel: Magento\Config\Model\Config\Source\Email\Template
						},

						element.Field{
							// Path: currency/import/frequency
							ID:        "frequency",
							Label:     `Frequency`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Cron\Model\Config\Source\Frequency
						},

						element.Field{
							// Path: currency/import/service
							ID:        "service",
							Label:     `Service`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Currency\Cron
							// SourceModel: Magento\Directory\Model\Currency\Import\Source\Service
						},

						element.Field{
							// Path: currency/import/time
							ID:        "time",
							Label:     `Start Time`,
							Type:      element.TypeTime,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
		element.Section{
			ID: "system",
			Groups: element.NewGroupSlice(
				element.Group{
					ID:        "currency",
					Label:     `Currency`,
					SortOrder: 50,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: system/currency/installed
							ID:         "installed",
							Label:      `Installed Currencies`,
							Type:       element.TypeMultiselect,
							SortOrder:  1,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermDefault,
							CanBeEmpty: true,
							Default:    `AZN,AZM,AFN,ALL,DZD,AOA,ARS,AMD,AWG,AUD,BSD,BHD,BDT,BBD,BYR,BZD,BMD,BTN,BOB,BAM,BWP,BRL,GBP,BND,BGN,BUK,BIF,KHR,CAD,CVE,CZK,KYD,CLP,CNY,COP,KMF,CDF,CRC,HRK,CUP,DKK,DJF,DOP,XCD,EGP,SVC,GQE,ERN,EEK,ETB,EUR,FKP,FJD,GMD,GEK,GEL,GHS,GIP,GTQ,GNF,GYD,HTG,HNL,HKD,HUF,ISK,INR,IDR,IRR,IQD,ILS,JMD,JPY,JOD,KZT,KES,KWD,KGS,LAK,LVL,LBP,LSL,LRD,LYD,LTL,MOP,MKD,MGA,MWK,MYR,MVR,LSM,MRO,MUR,MXN,MDL,MNT,MAD,MZN,MMK,NAD,NPR,ANG,TRL,TRY,NZD,NIC,NGN,KPW,NOK,OMR,PKR,PAB,PGK,PYG,PEN,PHP,PLN,QAR,RHD,RON,ROL,RUB,RWF,SHP,STD,SAR,RSD,SCR,SLL,SGD,SKK,SBD,SOS,ZAR,KRW,LKR,SDG,SRD,SZL,SEK,CHF,SYP,TWD,TJS,TZS,THB,TOP,TTD,TND,TMM,USD,UGX,UAH,AED,UYU,UZS,VUV,VEB,VEF,VND,CHE,CHW,XOF,XPF,WST,YER,ZMK,ZWD`,
							// BackendModel: Magento\Config\Model\Config\Backend\Locale
							// SourceModel: Magento\Config\Model\Config\Source\Locale\Currency\All
						},
					),
				},
			),
		},
		element.Section{
			ID: "general",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "country",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: general/country/optional_zip_countries
							ID:         "optional_zip_countries",
							Label:      `Zip/Postal Code is Optional for`,
							Type:       element.TypeMultiselect,
							SortOrder:  3,
							Visible:    element.VisibleYes,
							Scopes:     scope.PermDefault,
							CanBeEmpty: true,
							Default:    `HK,IE,MO,PA,GB`,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},
					),
				},

				element.Group{
					ID:        "region",
					Label:     `State Options`,
					SortOrder: 4,
					Scopes:    scope.PermDefault,
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: general/region/state_required
							ID:        "state_required",
							Label:     `State is Required for`,
							Type:      element.TypeMultiselect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Directory\Model\Config\Source\Country
						},

						element.Field{
							// Path: general/region/display_all
							ID:        "display_all",
							Label:     `Allow to Choose State if It is Optional for Country`,
							Type:      element.TypeSelect,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				element.Group{
					ID: "locale",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: general/locale/weight_unit
							ID:        "weight_unit",
							Label:     `Weight Unit`,
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `lbs`,
							// SourceModel: Magento\Directory\Model\Config\Source\WeightUnit
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "currency",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "import",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: currency/import/error_email
							ID:      `error_email`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},
					),
				},
			),
		},
		element.Section{
			ID: "general",
			Groups: element.NewGroupSlice(
				element.Group{
					ID: "country",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: general/country/allow
							ID:      `allow`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `AF,AL,DZ,AS,AD,AO,AI,AQ,AG,AR,AM,AW,AU,AT,AX,AZ,BS,BH,BD,BB,BY,BE,BZ,BJ,BM,BL,BT,BO,BA,BW,BV,BR,IO,VG,BN,BG,BF,BI,KH,CM,CA,CD,CV,KY,CF,TD,CL,CN,CX,CC,CO,KM,CG,CK,CR,HR,CU,CY,CZ,DK,DJ,DM,DO,EC,EG,SV,GQ,ER,EE,ET,FK,FO,FJ,FI,FR,GF,PF,TF,GA,GM,GE,DE,GG,GH,GI,GR,GL,GD,GP,GU,GT,GN,GW,GY,HT,HM,HN,HK,HU,IS,IM,IN,ID,IR,IQ,IE,IL,IT,CI,JE,JM,JP,JO,KZ,KE,KI,KW,KG,LA,LV,LB,LS,LR,LY,LI,LT,LU,ME,MF,MO,MK,MG,MW,MY,MV,ML,MT,MH,MQ,MR,MU,YT,FX,MX,FM,MD,MC,MN,MS,MA,MZ,MM,NA,NR,NP,NL,AN,NC,NZ,NI,NE,NG,NU,NF,KP,MP,NO,OM,PK,PW,PA,PG,PY,PE,PH,PN,PL,PS,PT,PR,QA,RE,RO,RS,RU,RW,SH,KN,LC,PM,VC,WS,SM,ST,SA,SN,SC,SL,SG,SK,SI,SB,SO,ZA,GS,KR,ES,LK,SD,SR,SJ,SZ,SE,CH,SY,TL,TW,TJ,TZ,TH,TG,TK,TO,TT,TN,TR,TM,TC,TV,VI,UG,UA,AE,GB,US,UM,UY,UZ,VU,VA,VE,VN,WF,EH,YE,ZM,ZW`,
						},

						element.Field{
							// Path: general/country/default
							ID:      `default`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `US`,
						},
					),
				},

				element.Group{
					ID: "locale",
					Fields: element.NewFieldSlice(
						element.Field{
							// Path: general/locale/datetime_format_long
							ID:      `datetime_format_long`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%A, %B %e %Y [%I:%M %p]`,
						},

						element.Field{
							// Path: general/locale/datetime_format_medium
							ID:      `datetime_format_medium`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%a, %b %e %Y [%I:%M %p]`,
						},

						element.Field{
							// Path: general/locale/datetime_format_short
							ID:      `datetime_format_short`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%m/%d/%y [%I:%M %p]`,
						},

						element.Field{
							// Path: general/locale/date_format_long
							ID:      `date_format_long`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%A, %B %e %Y`,
						},

						element.Field{
							// Path: general/locale/date_format_medium
							ID:      `date_format_medium`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%a, %b %e %Y`,
						},

						element.Field{
							// Path: general/locale/date_format_short
							ID:      `date_format_short`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `%m/%d/%y`,
						},

						element.Field{
							// Path: general/locale/language
							ID:      `language`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `en`,
						},

						element.Field{
							// Path: general/locale/code
							ID:      `code`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `en_US`,
						},

						element.Field{
							// Path: general/locale/timezone
							ID:      `timezone`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `America/Los_Angeles`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
