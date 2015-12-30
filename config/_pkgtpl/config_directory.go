// +build ignore

package directory

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "currency",
		Label:     `Currency Setup`,
		SortOrder: 60,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Backend::currency
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "options",
				Label:     `Currency Options`,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: currency/options/base
						ID:        "base",
						Label:     `Base Currency`,
						Comment:   element.LongText(`Base currency is used for all online payment transactions. If you have more than one store view, the base currency scope is defined by the catalog price scope ("Catalog" > "Price" > "Catalog Price Scope").`),
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `USD`,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Base
						// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
					},

					&config.Field{
						// Path: currency/options/default
						ID:        "default",
						Label:     `Default Display Currency`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `USD`,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\DefaultCurrency
						// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
					},

					&config.Field{
						// Path: currency/options/allow
						ID:         "allow",
						Label:      `Allowed Currencies`,
						Type:       config.TypeMultiselect,
						SortOrder:  3,
						Visible:    config.VisibleYes,
						Scope:      scope.PermAll,
						CanBeEmpty: true,
						Default:    `USD,EUR`,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Allow
						// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency
					},
				),
			},

			&config.Group{
				ID:        "webservicex",
				Label:     `Webservicex`,
				SortOrder: 40,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: currency/webservicex/timeout
						ID:      "timeout",
						Label:   `Connection Timeout in Seconds`,
						Type:    config.TypeText,
						Visible: config.VisibleYes,
						Scope:   scope.NewPerm(scope.DefaultID),
						Default: 100,
					},
				),
			},

			&config.Group{
				ID:        "import",
				Label:     `Scheduled Import Settings`,
				SortOrder: 50,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: currency/import/enabled
						ID:        "enabled",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: currency/import/error_email
						ID:        "error_email",
						Label:     `Error Email Recipient`,
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: currency/import/error_email_identity
						ID:        "error_email_identity",
						Label:     `Error Email Sender`,
						Type:      config.TypeSelect,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `general`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: currency/import/error_email_template
						ID:        "error_email_template",
						Label:     `Error Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `currency_import_error_email_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: currency/import/frequency
						ID:        "frequency",
						Label:     `Frequency`,
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: currency/import/service
						ID:        "service",
						Label:     `Service`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Currency\Cron
						// SourceModel: Otnegam\Directory\Model\Currency\Import\Source\Service
					},

					&config.Field{
						// Path: currency/import/time
						ID:        "time",
						Label:     `Start Time`,
						Type:      config.TypeTime,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "currency",
				Label:     `Currency`,
				SortOrder: 50,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/currency/installed
						ID:         "installed",
						Label:      `Installed Currencies`,
						Type:       config.TypeMultiselect,
						SortOrder:  1,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID),
						CanBeEmpty: true,
						Default:    `AZN,AZM,AFN,ALL,DZD,AOA,ARS,AMD,AWG,AUD,BSD,BHD,BDT,BBD,BYR,BZD,BMD,BTN,BOB,BAM,BWP,BRL,GBP,BND,BGN,BUK,BIF,KHR,CAD,CVE,CZK,KYD,CLP,CNY,COP,KMF,CDF,CRC,HRK,CUP,DKK,DJF,DOP,XCD,EGP,SVC,GQE,ERN,EEK,ETB,EUR,FKP,FJD,GMD,GEK,GEL,GHS,GIP,GTQ,GNF,GYD,HTG,HNL,HKD,HUF,ISK,INR,IDR,IRR,IQD,ILS,JMD,JPY,JOD,KZT,KES,KWD,KGS,LAK,LVL,LBP,LSL,LRD,LYD,LTL,MOP,MKD,MGA,MWK,MYR,MVR,LSM,MRO,MUR,MXN,MDL,MNT,MAD,MZN,MMK,NAD,NPR,ANG,TRL,TRY,NZD,NIC,NGN,KPW,NOK,OMR,PKR,PAB,PGK,PYG,PEN,PHP,PLN,QAR,RHD,RON,ROL,RUB,RWF,SHP,STD,SAR,RSD,SCR,SLL,SGD,SKK,SBD,SOS,ZAR,KRW,LKR,SDG,SRD,SZL,SEK,CHF,SYP,TWD,TJS,TZS,THB,TOP,TTD,TND,TMM,USD,UGX,UAH,AED,UYU,UZS,VUV,VEB,VEF,VND,CHE,CHW,XOF,XPF,WST,YER,ZMK,ZWD`,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Locale
						// SourceModel: Otnegam\Config\Model\Config\Source\Locale\Currency\All
					},
				),
			},
		),
	},
	&config.Section{
		ID: "general",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "country",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/country/optional_zip_countries
						ID:         "optional_zip_countries",
						Label:      `Zip/Postal Code is Optional for`,
						Type:       config.TypeMultiselect,
						SortOrder:  3,
						Visible:    config.VisibleYes,
						Scope:      scope.NewPerm(scope.DefaultID),
						CanBeEmpty: true,
						Default:    `HK,IE,MO,PA,GB`,
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},
				),
			},

			&config.Group{
				ID:        "region",
				Label:     `State Options`,
				SortOrder: 4,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/region/state_required
						ID:        "state_required",
						Label:     `State is Required for`,
						Type:      config.TypeMultiselect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Directory\Model\Config\Source\Country
					},

					&config.Field{
						// Path: general/region/display_all
						ID:        "display_all",
						Label:     `Allow to Choose State if It is Optional for Country`,
						Type:      config.TypeSelect,
						SortOrder: 8,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID: "locale",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/locale/weight_unit
						ID:        "weight_unit",
						Label:     `Weight Unit`,
						Type:      config.TypeSelect,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `lbs`,
						// SourceModel: Otnegam\Directory\Model\Config\Source\WeightUnit
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "currency",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "import",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: currency/import/error_email
						ID:      `error_email`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "general",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "country",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/country/allow
						ID:      `allow`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `AF,AL,DZ,AS,AD,AO,AI,AQ,AG,AR,AM,AW,AU,AT,AX,AZ,BS,BH,BD,BB,BY,BE,BZ,BJ,BM,BL,BT,BO,BA,BW,BV,BR,IO,VG,BN,BG,BF,BI,KH,CM,CA,CD,CV,KY,CF,TD,CL,CN,CX,CC,CO,KM,CG,CK,CR,HR,CU,CY,CZ,DK,DJ,DM,DO,EC,EG,SV,GQ,ER,EE,ET,FK,FO,FJ,FI,FR,GF,PF,TF,GA,GM,GE,DE,GG,GH,GI,GR,GL,GD,GP,GU,GT,GN,GW,GY,HT,HM,HN,HK,HU,IS,IM,IN,ID,IR,IQ,IE,IL,IT,CI,JE,JM,JP,JO,KZ,KE,KI,KW,KG,LA,LV,LB,LS,LR,LY,LI,LT,LU,ME,MF,MO,MK,MG,MW,MY,MV,ML,MT,MH,MQ,MR,MU,YT,FX,MX,FM,MD,MC,MN,MS,MA,MZ,MM,NA,NR,NP,NL,AN,NC,NZ,NI,NE,NG,NU,NF,KP,MP,NO,OM,PK,PW,PA,PG,PY,PE,PH,PN,PL,PS,PT,PR,QA,RE,RO,RS,RU,RW,SH,KN,LC,PM,VC,WS,SM,ST,SA,SN,SC,SL,SG,SK,SI,SB,SO,ZA,GS,KR,ES,LK,SD,SR,SJ,SZ,SE,CH,SY,TL,TW,TJ,TZ,TH,TG,TK,TO,TT,TN,TR,TM,TC,TV,VI,UG,UA,AE,GB,US,UM,UY,UZ,VU,VA,VE,VN,WF,EH,YE,ZM,ZW`,
					},

					&config.Field{
						// Path: general/country/default
						ID:      `default`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `US`,
					},
				),
			},

			&config.Group{
				ID: "locale",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: general/locale/datetime_format_long
						ID:      `datetime_format_long`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `%A, %B %e %Y [%I:%M %p]`,
					},

					&config.Field{
						// Path: general/locale/datetime_format_medium
						ID:      `datetime_format_medium`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `%a, %b %e %Y [%I:%M %p]`,
					},

					&config.Field{
						// Path: general/locale/datetime_format_short
						ID:      `datetime_format_short`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `%m/%d/%y [%I:%M %p]`,
					},

					&config.Field{
						// Path: general/locale/date_format_long
						ID:      `date_format_long`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `%A, %B %e %Y`,
					},

					&config.Field{
						// Path: general/locale/date_format_medium
						ID:      `date_format_medium`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `%a, %b %e %Y`,
					},

					&config.Field{
						// Path: general/locale/date_format_short
						ID:      `date_format_short`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `%m/%d/%y`,
					},

					&config.Field{
						// Path: general/locale/language
						ID:      `language`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `en`,
					},

					&config.Field{
						// Path: general/locale/code
						ID:      `code`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `en_US`,
					},

					&config.Field{
						// Path: general/locale/timezone
						ID:      `timezone`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `America/Los_Angeles`,
					},
				),
			},
		),
	},
)
