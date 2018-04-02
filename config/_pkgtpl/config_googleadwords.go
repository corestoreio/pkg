// +build ignore

package googleadwords

import (
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.Sections

func init() {
	ConfigStructure = element.MustMakeSectionsValidate(
		element.Section{
			ID: "google",
			Groups: element.MakeGroups(
				element.Group{
					ID:        "adwords",
					Label:     `Google AdWords`,
					SortOrder: 15,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: google/adwords/active
							ID:        "active",
							Label:     `Enable`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: google/adwords/conversion_id
							ID:        "conversion_id",
							Label:     `Conversion ID`,
							Type:      element.TypeText,
							SortOrder: 11,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\GoogleAdwords\Model\Config\Backend\ConversionId
						},

						element.Field{
							// Path: google/adwords/conversion_language
							ID:        "conversion_language",
							Label:     `Conversion Language`,
							Type:      element.TypeSelect,
							SortOrder: 12,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `en`,
							// SourceModel: Magento\GoogleAdwords\Model\Config\Source\Language
						},

						element.Field{
							// Path: google/adwords/conversion_format
							ID:        "conversion_format",
							Label:     `Conversion Format`,
							Type:      element.TypeText,
							SortOrder: 13,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   2,
						},

						element.Field{
							// Path: google/adwords/conversion_color
							ID:        "conversion_color",
							Label:     `Conversion Color`,
							Type:      element.TypeText,
							SortOrder: 14,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `FFFFFF`,
							// BackendModel: Magento\GoogleAdwords\Model\Config\Backend\Color
						},

						element.Field{
							// Path: google/adwords/conversion_label
							ID:        "conversion_label",
							Label:     `Conversion Label`,
							Type:      element.TypeText,
							SortOrder: 15,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: google/adwords/conversion_value_type
							ID:        "conversion_value_type",
							Label:     `Conversion Value Type`,
							Type:      element.TypeSelect,
							SortOrder: 16,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\GoogleAdwords\Model\Config\Source\ValueType
						},

						element.Field{
							// Path: google/adwords/conversion_value
							ID:        "conversion_value",
							Label:     `Conversion Value`,
							Type:      element.TypeText,
							SortOrder: 17,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		element.Section{
			ID: "google",
			Groups: element.MakeGroups(
				element.Group{
					ID: "adwords",
					Fields: element.MakeFields(
						element.Field{
							// Path: google/adwords/languages
							ID:      `languages`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"ar":"ar","bg":"bg","ca":"ca","cs":"cs","da":"da","de":"de","el":"el","en":"en","es":"es","et":"et","fi":"fi","fr":"fr","hi":"hi","hr":"hr","hu":"hu","id":"id","is":"is","it":"it","iw":"iw","ja":"ja","ko":"ko","lt":"lt","lv":"lv","nl":"nl","no":"no","pl":"pl","pt":"pt","ro":"ro","ru":"ru","sk":"sk","sl":"sl","sr":"sr","sv":"sv","th":"th","tl":"tl","tr":"tr","uk":"uk","ur":"ur","vi":"vi","zh_TW":"zh_TW","zh_CN":"zh_CN"}`,
						},

						element.Field{
							// Path: google/adwords/language_convert
							ID:      `language_convert`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"zh_CN":"zh_Hans","zh_TW":"zh_Hant","iw":"he"}`,
						},

						element.Field{
							// Path: google/adwords/conversion_js_src
							ID:      `conversion_js_src`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://www.googleadservices.com/pagead/conversion.js`,
						},

						element.Field{
							// Path: google/adwords/conversion_img_src
							ID:      `conversion_img_src`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `https://www.googleadservices.com/pagead/conversion/%s/?label=%s&guid=ON&script=0`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
