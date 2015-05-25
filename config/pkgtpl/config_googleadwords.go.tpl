package googleadwords

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "google",
		Label:     "",
		SortOrder: 0,
		Scope:     config.NewScopePerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "adwords",
				Label:     `Google AdWords`,
				Comment:   ``,
				SortOrder: 15,
				Scope:     config.IDScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `google/adwords/active`,
						ID:           "active",
						Label:        `Enable`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `google/adwords/conversion_id`,
						ID:           "conversion_id",
						Label:        `Conversion ID`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    11,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      nil,
						BackendModel: nil, // Magento\GoogleAdwords\Model\Config\Backend\ConversionId
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `google/adwords/conversion_language`,
						ID:           "conversion_language",
						Label:        `Conversion Language`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    12,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      `en`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\GoogleAdwords\Model\Config\Source\Language
					},

					&config.Field{
						// Path: `google/adwords/conversion_format`,
						ID:           "conversion_format",
						Label:        `Conversion Format`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    13,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      2,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `google/adwords/conversion_color`,
						ID:           "conversion_color",
						Label:        `Conversion Color`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    14,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      `FFFFFF`,
						BackendModel: nil, // Magento\GoogleAdwords\Model\Config\Backend\Color
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `google/adwords/conversion_label`,
						ID:           "conversion_label",
						Label:        `Conversion Label`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    15,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `google/adwords/conversion_value_type`,
						ID:           "conversion_value_type",
						Label:        `Conversion Value Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    16,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\GoogleAdwords\Model\Config\Source\ValueType
					},

					&config.Field{
						// Path: `google/adwords/conversion_value`,
						ID:           "conversion_value",
						Label:        `Conversion Value`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    17,
						Visible:      true,
						Scope:        config.IDScopePermAll,
						Default:      0,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "google",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "adwords",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `google/adwords/languages`,
						ID:      "languages",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `{"ar":"ar","bg":"bg","ca":"ca","cs":"cs","da":"da","de":"de","el":"el","en":"en","es":"es","et":"et","fi":"fi","fr":"fr","hi":"hi","hr":"hr","hu":"hu","id":"id","is":"is","it":"it","iw":"iw","ja":"ja","ko":"ko","lt":"lt","lv":"lv","nl":"nl","no":"no","pl":"pl","pt":"pt","ro":"ro","ru":"ru","sk":"sk","sl":"sl","sr":"sr","sv":"sv","th":"th","tl":"tl","tr":"tr","uk":"uk","ur":"ur","vi":"vi","zh_TW":"zh_TW","zh_CN":"zh_CN"}`,
					},

					&config.Field{
						// Path: `google/adwords/language_convert`,
						ID:      "language_convert",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `{"zh_CN":"zh_Hans","zh_TW":"zh_Hant","iw":"he"}`,
					},

					&config.Field{
						// Path: `google/adwords/conversion_js_src`,
						ID:      "conversion_js_src",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `https://www.googleadservices.com/pagead/conversion.js`,
					},

					&config.Field{
						// Path: `google/adwords/conversion_img_src`,
						ID:      "conversion_img_src",
						Type:    config.TypeHidden,
						Visible: false,
						Scope:   config.NewScopePerm(config.IDScopeDefault), // @todo search for that
						Default: `https://www.googleadservices.com/pagead/conversion/%s/?label=%s&guid=ON&script=0`,
					},
				},
			},
		},
	},
)
