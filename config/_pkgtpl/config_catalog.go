// +build ignore

package catalog

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     "Catalog",
		SortOrder: 40,
		Scope:     scope.PermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "fields_masks",
				Label:     `Product Fields Auto-Generation`,
				Comment:   ``,
				SortOrder: 90,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/fields_masks/sku`,
						ID:           "sku",
						Label:        `Mask for SKU`,
						Comment:      `Use {{name}} as Product Name placeholder`,
						Type:         config.TypeText,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `{{name}}`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/fields_masks/meta_title`,
						ID:           "meta_title",
						Label:        `Mask for Meta Title`,
						Comment:      `Use {{name}} as Product Name placeholder`,
						Type:         config.TypeText,
						SortOrder:    20,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `{{name}}`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/fields_masks/meta_keyword`,
						ID:           "meta_keyword",
						Label:        `Mask for Meta Keywords`,
						Comment:      `Use {{name}} as Product Name or {{sku}} as Product SKU placeholders`,
						Type:         config.TypeText,
						SortOrder:    30,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `{{name}}`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/fields_masks/meta_description`,
						ID:           "meta_description",
						Label:        `Mask for Meta Description`,
						Comment:      `Use {{name}} and {{description}} as Product Name and Product Description placeholders`,
						Type:         config.TypeText,
						SortOrder:    40,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      `{{name}} {{description}}`,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "frontend",
				Label:     `Storefront`,
				Comment:   ``,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/frontend/list_mode`,
						ID:           "list_mode",
						Label:        `List Mode`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `grid-list`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Catalog\Model\Config\Source\ListMode
					},

					&config.Field{
						// Path: `catalog/frontend/grid_per_page_values`,
						ID:           "grid_per_page_values",
						Label:        `Products per Page on Grid Allowed Values`,
						Comment:      `Comma-separated.`,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `9,15,30`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/frontend/grid_per_page`,
						ID:           "grid_per_page",
						Label:        `Products per Page on Grid Default Value`,
						Comment:      `Must be in the allowed values list`,
						Type:         config.TypeText,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      9,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/frontend/list_per_page_values`,
						ID:           "list_per_page_values",
						Label:        `Products per Page on List Allowed Values`,
						Comment:      `Comma-separated.`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `5,10,15,20,25`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/frontend/list_per_page`,
						ID:           "list_per_page",
						Label:        `Products per Page on List Default Value`,
						Comment:      `Must be in the allowed values list`,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      10,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/frontend/flat_catalog_category`,
						ID:           "flat_catalog_category",
						Label:        `Use Flat Catalog Category`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      false,
						BackendModel: nil, // Magento\Catalog\Model\Indexer\Category\Flat\System\Config\Mode
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/frontend/flat_catalog_product`,
						ID:           "flat_catalog_product",
						Label:        `Use Flat Catalog Product`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil, // Magento\Catalog\Model\Indexer\Product\Flat\System\Config\Mode
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/frontend/default_sort_by`,
						ID:           "default_sort_by",
						Label:        `Product Listing Sort by`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    6,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `position`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Catalog\Model\Config\Source\ListSort
					},

					&config.Field{
						// Path: `catalog/frontend/list_allow_all`,
						ID:           "list_allow_all",
						Label:        `Allow All Products per Page`,
						Comment:      `Whether to show "All" option in the "Show X Per Page" dropdown`,
						Type:         config.TypeSelect,
						SortOrder:    6,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/frontend/parse_url_directives`,
						ID:           "parse_url_directives",
						Label:        `Allow Dynamic Media URLs in Products and Categories`,
						Comment:      `E.g. {{media url="path/to/image.jpg"}} {{skin url="path/to/picture.gif"}}. Dynamic directives parsing impacts catalog performance.`,
						Type:         config.TypeSelect,
						SortOrder:    200,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "placeholder",
				Label:     `Product Image Placeholders`,
				Comment:   ``,
				SortOrder: 300,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/placeholder/placeholder`,
						ID:           "placeholder",
						Label:        ``,
						Comment:      ``,
						Type:         config.TypeImage,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Image
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "seo",
				Label:     `Search Engine Optimization`,
				Comment:   ``,
				SortOrder: 500,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/seo/title_separator`,
						ID:           "title_separator",
						Label:        `Page Title Separator`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    6,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `-`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/seo/category_canonical_tag`,
						ID:           "category_canonical_tag",
						Label:        `Use Canonical Link Meta Tag For Categories`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    7,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/seo/product_canonical_tag`,
						ID:           "product_canonical_tag",
						Label:        `Use Canonical Link Meta Tag For Products`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    8,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},

			&config.Group{
				ID:        "price",
				Label:     `Price`,
				Comment:   ``,
				SortOrder: 400,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/price/scope`,
						ID:           "scope",
						Label:        `Catalog Price Scope`,
						Comment:      `This defines the base currency scope ("Currency Setup" > "Currency Options" > "Base Currency").`,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil, // Magento\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope
						SourceModel:  nil, // Magento\Catalog\Model\Config\Source\Price\Scope
					},
				},
			},

			&config.Group{
				ID:        "navigation",
				Label:     `Category Top Navigation`,
				Comment:   ``,
				SortOrder: 500,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/navigation/max_depth`,
						ID:           "max_depth",
						Label:        `Maximal Depth`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      0,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},

			&config.Group{
				ID:        "custom_options",
				Label:     `Date & Time Custom Options`,
				Comment:   ``,
				SortOrder: 700,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/custom_options/use_calendar`,
						ID:           "use_calendar",
						Label:        `Use JavaScript Calendar`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `catalog/custom_options/date_fields_order`,
						ID:           "date_fields_order",
						Label:        `Date Fields Order`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `m,d,y`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `catalog/custom_options/time_format`,
						ID:           "time_format",
						Label:        `Time Format`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `12h`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Catalog\Model\Config\Source\TimeFormat
					},

					&config.Field{
						// Path: `catalog/custom_options/year_range`,
						ID:           "year_range",
						Label:        `Year Range`,
						Comment:      `Please use a four-digit year format.`,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "design",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "watermark",
				Label:     `Product Image Watermarks`,
				Comment:   ``,
				SortOrder: 400,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `design/watermark/size`,
						ID:           "size",
						Label:        `Watermark Default Size`,
						Comment:      `Example format: 200x300.`,
						Type:         config.TypeText,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/watermark/imageOpacity`,
						ID:           "imageOpacity",
						Label:        `Watermark Opacity, Percent`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    150,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/watermark/image`,
						ID:           "image",
						Label:        `Watermark`,
						Comment:      `Allowed file types: jpeg, gif, png.`,
						Type:         config.TypeImage,
						SortOrder:    200,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Image
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `design/watermark/position`,
						ID:           "position",
						Label:        `Watermark Position`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    300,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Catalog\Model\Config\Source\Watermark\Position
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "cms",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "wysiwyg",
				Label:     ``,
				Comment:   ``,
				SortOrder: 0,
				Scope:     nil,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `cms/wysiwyg/use_static_urls_in_catalog`,
						ID:           "use_static_urls_in_catalog",
						Label:        `Use Static URLs for Media Content in WYSIWYG for Catalog`,
						Comment:      `This applies only to catalog products and categories. Media content will be inserted into the editor as a static URL. Media content is not updated if the system configuration base URL changes.`,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(scope.DefaultID),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "rss",
		Label:     "",
		SortOrder: 0,
		Scope:     nil,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "catalog",
				Label:     `Catalog`,
				Comment:   ``,
				SortOrder: 3,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `rss/catalog/new`,
						ID:           "new",
						Label:        `New Products`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    10,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: `rss/catalog/special`,
						ID:           "special",
						Label:        `Special Products`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    11,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: `rss/catalog/category`,
						ID:           "category",
						Label:        `Top Level Category`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    14,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Enabledisable
					},
				},
			},
		},
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "catalog",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "product",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/product/flat`,
						ID:      "flat",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `{"max_index_count":"64"}`,
					},

					&config.Field{
						// Path: `catalog/product/default_tax_group`,
						ID:      "default_tax_group",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: 2,
					},
				},
			},

			&config.Group{
				ID: "seo",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/seo/product_url_suffix`,
						ID:      "product_url_suffix",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `.html`,
					},

					&config.Field{
						// Path: `catalog/seo/category_url_suffix`,
						ID:      "category_url_suffix",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `.html`,
					},

					&config.Field{
						// Path: `catalog/seo/product_use_categories`,
						ID:      "product_use_categories",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: false,
					},

					&config.Field{
						// Path: `catalog/seo/save_rewrites_history`,
						ID:      "save_rewrites_history",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: true,
					},
				},
			},

			&config.Group{
				ID: "custom_options",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `catalog/custom_options/forbidden_extensions`,
						ID:      "forbidden_extensions",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(scope.DefaultID), // @todo search for that
						Default: `php,exe`,
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
						Default: `{"tmp_images_folder":"tmp","catalog_images_folder":"catalog","product_custom_options_fodler":"custom_options"}`,
					},
				},
			},
		},
	},
)
