// +build ignore

package theme

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "design",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "head",
				Label:     `HTML Head`,
				Comment:   ``,
				SortOrder: 20,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/head/shortcut_icon`,
						ID:           "shortcut_icon",
						Label:        `Favicon Icon`,
						Comment:      `Allowed file types: ICO, PNG, GIF, JPG, JPEG, APNG, SVG. Not all browsers support all these formats!`,
						Type:         config.TypeImage,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Image\Favicon
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/head/default_title`,
						ID:           "default_title",
						Label:        `Default Title`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/head/title_prefix`,
						ID:           "title_prefix",
						Label:        `Title Prefix`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    12,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/head/title_suffix`,
						ID:           "title_suffix",
						Label:        `Title Suffix`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    14,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/head/default_description`,
						ID:           "default_description",
						Label:        `Default Description`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/head/default_keywords`,
						ID:           "default_keywords",
						Label:        `Default Keywords`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/head/includes`,
						ID:           "includes",
						Label:        `Miscellaneous Scripts`,
						Comment:      `This will be included before head closing tag in page HTML.`,
						Type:         config.TypeTextarea,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/head/demonotice`,
						ID:           "demonotice",
						Label:        `Display Demo Store Notice`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "search_engine_robots",
				Label:     `Search Engine Robots`,
				Comment:   ``,
				SortOrder: 25,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/search_engine_robots/default_robots`,
						ID:           "default_robots",
						Label:        `Default Robots`,
						Comment:      `This will be included before head closing tag in page HTML.`,
						Type:         config.TypeSelect,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `INDEX,FOLLOW`,
						BackendModel: nil,
						// SourceModel:  nil, // Magento\Config\Model\Config\Source\Design\Robots
					},

					&config.Field{
						// Path: `design/search_engine_robots/custom_instructions`,
						ID:           "custom_instructions",
						Label:        `Edit custom instruction of robots.txt File`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    70,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Admin\Robots
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/search_engine_robots/reset_to_defaults`,
						ID:           "reset_to_defaults",
						Label:        `Reset to Defaults`,
						Comment:      `This action will delete your custom instructions and reset robots.txt file to system's default settings.`,
						Type:         config.TypeButton,
						SortOrder:    80,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "header",
				Label:     `Header`,
				Comment:   ``,
				SortOrder: 30,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/header/logo_src`,
						ID:           "logo_src",
						Label:        `Logo Image`,
						Comment:      `Allowed file types:PNG, GIF, JPG, JPEG, SVG.`,
						Type:         config.TypeImage,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Image\Logo
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/header/logo_width`,
						ID:           "logo_width",
						Label:        `Logo Image Width`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    11,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/header/logo_height`,
						ID:           "logo_height",
						Label:        `Logo Image Height`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    12,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/header/logo_alt`,
						ID:           "logo_alt",
						Label:        `Logo Image Alt`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/header/welcome`,
						ID:           "welcome",
						Label:        `Welcome Text`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "footer",
				Label:     `Footer`,
				Comment:   ``,
				SortOrder: 40,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/footer/copyright`,
						ID:           "copyright",
						Label:        `Copyright`,
						Comment:      ``,
						Type:         config.TypeTextarea,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/footer/absolute_footer`,
						ID:           "absolute_footer",
						Label:        `Miscellaneous HTML`,
						Comment:      `This will be displayed just before body closing tag.`,
						Type:         config.TypeTextarea,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						// SourceModel:  nil,
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "design",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "head",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/head/_value`,
						ID:      "_value",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"default_title":"Magento Commerce","default_description":"Default Description","default_keywords":"Magento, Varien, E-commerce","default_media_type":"text\/html","default_charset":"utf-8"}`,
					},

					&config.Field{
						// Path: `design/head/_attribute`,
						ID:      "_attribute",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"translate":"default_description"}`,
					},
				},
			},

			&config.Group{
				ID: "search_engine_robots",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/search_engine_robots/default_custom_instructions`,
						ID:      "default_custom_instructions",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `
User-agent: *
Disallow: /index.php/
Disallow: /*?
Disallow: /checkout/
Disallow: /app/
Disallow: /lib/
Disallow: /*.php$
Disallow: /pkginfo/
Disallow: /report/
Disallow: /var/
Disallow: /catalog/
Disallow: /customer/
Disallow: /sendfriend/
Disallow: /review/
Disallow: /*SID=
                    `,
					},
				},
			},

			&config.Group{
				ID: "header",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/header/_value`,
						ID:      "_value",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"logo_alt":"Magento Commerce","welcome":"Default welcome msg!"}`,
					},

					&config.Field{
						// Path: `design/header/_attribute`,
						ID:      "_attribute",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"translate":"welcome"}`,
					},
				},
			},

			&config.Group{
				ID: "footer",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/footer/_value`,
						ID:      "_value",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"copyright":"Copyright \u00a9 2015 Magento. All rights reserved."}`,
					},

					&config.Field{
						// Path: `design/footer/_attribute`,
						ID:      "_attribute",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"translate":"copyright"}`,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "theme",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "customization",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `theme/customization/custom_css`,
						ID:      "custom_css",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `Magento\Theme\Model\Theme\Customization\File\CustomCss`,
					},
				},
			},
		},
	},
	&config.Section{
		ID: "system",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "media_storage_configuration",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `system/media_storage_configuration/allowed_resources`,
						ID:      "allowed_resources",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"site_favicons":"favicon"}`,
					},
				},
			},
		},
	},
)
