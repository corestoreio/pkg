// +build ignore

package googleadwords

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "google",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "adwords",
				Label:     `Google AdWords`,
				SortOrder: 15,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: google/adwords/active
						ID:        "active",
						Label:     `Enable`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: google/adwords/conversion_id
						ID:        "conversion_id",
						Label:     `Conversion ID`,
						Type:      config.TypeText,
						SortOrder: 11,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\GoogleAdwords\Model\Config\Backend\ConversionId
					},

					&config.Field{
						// Path: google/adwords/conversion_language
						ID:        "conversion_language",
						Label:     `Conversion Language`,
						Type:      config.TypeSelect,
						SortOrder: 12,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `en`,
						// SourceModel: Otnegam\GoogleAdwords\Model\Config\Source\Language
					},

					&config.Field{
						// Path: google/adwords/conversion_format
						ID:        "conversion_format",
						Label:     `Conversion Format`,
						Type:      config.TypeText,
						SortOrder: 13,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   2,
					},

					&config.Field{
						// Path: google/adwords/conversion_color
						ID:        "conversion_color",
						Label:     `Conversion Color`,
						Type:      config.TypeText,
						SortOrder: 14,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `FFFFFF`,
						// BackendModel: Otnegam\GoogleAdwords\Model\Config\Backend\Color
					},

					&config.Field{
						// Path: google/adwords/conversion_label
						ID:        "conversion_label",
						Label:     `Conversion Label`,
						Type:      config.TypeText,
						SortOrder: 15,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: google/adwords/conversion_value_type
						ID:        "conversion_value_type",
						Label:     `Conversion Value Type`,
						Type:      config.TypeSelect,
						SortOrder: 16,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\GoogleAdwords\Model\Config\Source\ValueType
					},

					&config.Field{
						// Path: google/adwords/conversion_value
						ID:        "conversion_value",
						Label:     `Conversion Value`,
						Type:      config.TypeText,
						SortOrder: 17,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "google",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "adwords",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: google/adwords/languages
						ID:      `languages`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"ar":"ar","bg":"bg","ca":"ca","cs":"cs","da":"da","de":"de","el":"el","en":"en","es":"es","et":"et","fi":"fi","fr":"fr","hi":"hi","hr":"hr","hu":"hu","id":"id","is":"is","it":"it","iw":"iw","ja":"ja","ko":"ko","lt":"lt","lv":"lv","nl":"nl","no":"no","pl":"pl","pt":"pt","ro":"ro","ru":"ru","sk":"sk","sl":"sl","sr":"sr","sv":"sv","th":"th","tl":"tl","tr":"tr","uk":"uk","ur":"ur","vi":"vi","zh_TW":"zh_TW","zh_CN":"zh_CN"}`,
					},

					&config.Field{
						// Path: google/adwords/language_convert
						ID:      `language_convert`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"zh_CN":"zh_Hans","zh_TW":"zh_Hant","iw":"he"}`,
					},

					&config.Field{
						// Path: google/adwords/conversion_js_src
						ID:      `conversion_js_src`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://www.googleadservices.com/pagead/conversion.js`,
					},

					&config.Field{
						// Path: google/adwords/conversion_img_src
						ID:      `conversion_img_src`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `https://www.googleadservices.com/pagead/conversion/%s/?label=%s&guid=ON&script=0`,
					},
				),
			},
		),
	},
)
