// +build ignore

package theme

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID: "design",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "head",
					Label:     `HTML Head`,
					SortOrder: 20,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/head/shortcut_icon
							ID:        "shortcut_icon",
							Label:     `Favicon Icon`,
							Comment:   text.Long(`Allowed file types: ICO, PNG, GIF, JPG, JPEG, APNG, SVG. Not all browsers support all these formats!`),
							Type:      element.TypeImage,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Magento\Config\Model\Config\Backend\Image\Favicon
						},

						&element.Field{
							// Path: design/head/default_title
							ID:        "default_title",
							Label:     `Default Title`,
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/head/title_prefix
							ID:        "title_prefix",
							Label:     `Title Prefix`,
							Type:      element.TypeText,
							SortOrder: 12,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/head/title_suffix
							ID:        "title_suffix",
							Label:     `Title Suffix`,
							Type:      element.TypeText,
							SortOrder: 14,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/head/default_description
							ID:        "default_description",
							Label:     `Default Description`,
							Type:      element.TypeTextarea,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/head/default_keywords
							ID:        "default_keywords",
							Label:     `Default Keywords`,
							Type:      element.TypeTextarea,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/head/includes
							ID:        "includes",
							Label:     `Miscellaneous Scripts`,
							Comment:   text.Long(`This will be included before head closing tag in page HTML.`),
							Type:      element.TypeTextarea,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/head/demonotice
							ID:        "demonotice",
							Label:     `Display Demo Store Notice`,
							Type:      element.TypeSelect,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "search_engine_robots",
					Label:     `Search Engine Robots`,
					SortOrder: 25,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/search_engine_robots/default_robots
							ID:        "default_robots",
							Label:     `Default Robots`,
							Comment:   text.Long(`This will be included before head closing tag in page HTML.`),
							Type:      element.TypeSelect,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `INDEX,FOLLOW`,
							// SourceModel: Magento\Config\Model\Config\Source\Design\Robots
						},

						&element.Field{
							// Path: design/search_engine_robots/custom_instructions
							ID:        "custom_instructions",
							Label:     `Edit custom instruction of robots.txt File`,
							Type:      element.TypeTextarea,
							SortOrder: 70,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Magento\Config\Model\Config\Backend\Admin\Robots
						},

						&element.Field{
							// Path: design/search_engine_robots/reset_to_defaults
							ID:        "reset_to_defaults",
							Label:     `Reset to Defaults`,
							Comment:   text.Long(`This action will delete your custom instructions and reset robots.txt file to system's default settings.`),
							Type:      element.TypeButton,
							SortOrder: 80,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},

				&element.Group{
					ID:        "header",
					Label:     `Header`,
					SortOrder: 30,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/header/logo_src
							ID:        "logo_src",
							Label:     `Logo Image`,
							Comment:   text.Long(`Allowed file types:PNG, GIF, JPG, JPEG, SVG.`),
							Type:      element.TypeImage,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Magento\Config\Model\Config\Backend\Image\Logo
						},

						&element.Field{
							// Path: design/header/logo_width
							ID:        "logo_width",
							Label:     `Logo Image Width`,
							Type:      element.TypeText,
							SortOrder: 11,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/header/logo_height
							ID:        "logo_height",
							Label:     `Logo Image Height`,
							Type:      element.TypeText,
							SortOrder: 12,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/header/logo_alt
							ID:        "logo_alt",
							Label:     `Logo Image Alt`,
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/header/welcome
							ID:        "welcome",
							Label:     `Welcome Text`,
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},

				&element.Group{
					ID:        "footer",
					Label:     `Footer`,
					SortOrder: 40,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/footer/copyright
							ID:        "copyright",
							Label:     `Copyright`,
							Type:      element.TypeTextarea,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/footer/absolute_footer
							ID:        "absolute_footer",
							Label:     `Miscellaneous HTML`,
							Comment:   text.Long(`This will be displayed just before body closing tag.`),
							Type:      element.TypeTextarea,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "design",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "invalid_caches",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/invalid_caches/block_html
							ID:      `block_html`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						&element.Field{
							// Path: design/invalid_caches/layout
							ID:      `layout`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},

						&element.Field{
							// Path: design/invalid_caches/translate
							ID:      `translate`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},
					),
				},

				&element.Group{
					ID: "head",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/head/_value
							ID:      `_value`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"default_title":"Magento Commerce","default_description":"Default Description","default_keywords":"Magento, Varien, E-commerce","default_media_type":"text\/html","default_charset":"utf-8"}`,
						},

						&element.Field{
							// Path: design/head/_attribute
							ID:      `_attribute`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"translate":"default_description"}`,
						},
					),
				},

				&element.Group{
					ID: "search_engine_robots",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/search_engine_robots/default_custom_instructions
							ID:      `default_custom_instructions`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
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
					),
				},

				&element.Group{
					ID: "header",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/header/_value
							ID:      `_value`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"logo_alt":"Magento Commerce","welcome":"Default welcome msg!"}`,
						},

						&element.Field{
							// Path: design/header/_attribute
							ID:      `_attribute`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"translate":"welcome"}`,
						},
					),
				},

				&element.Group{
					ID: "footer",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/footer/_value
							ID:      `_value`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"copyright":"Copyright \u00a9 2015 Magento. All rights reserved."}`,
						},

						&element.Field{
							// Path: design/footer/_attribute
							ID:      `_attribute`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"translate":"copyright"}`,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "theme",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "customization",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: theme/customization/custom_css
							ID:      `custom_css`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `Magento\Theme\Model\Theme\Customization\File\CustomCss`,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "system",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "media_storage_configuration",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      `allowed_resources`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"site_favicons":"favicon"}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
