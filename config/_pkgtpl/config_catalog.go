// +build ignore

package catalog

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "catalog",
		Label:     `Catalog`,
		SortOrder: 40,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Catalog::config_catalog
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "fields_masks",
				Label:     `Product Fields Auto-Generation`,
				SortOrder: 90,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/fields_masks/sku
						ID:        "sku",
						Label:     `Mask for SKU`,
						Comment:   element.LongText(`Use {{name}} as Product Name placeholder`),
						Type:      config.TypeText,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `{{name}}`,
					},

					&config.Field{
						// Path: catalog/fields_masks/meta_title
						ID:        "meta_title",
						Label:     `Mask for Meta Title`,
						Comment:   element.LongText(`Use {{name}} as Product Name placeholder`),
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `{{name}}`,
					},

					&config.Field{
						// Path: catalog/fields_masks/meta_keyword
						ID:        "meta_keyword",
						Label:     `Mask for Meta Keywords`,
						Comment:   element.LongText(`Use {{name}} as Product Name or {{sku}} as Product SKU placeholders`),
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `{{name}}`,
					},

					&config.Field{
						// Path: catalog/fields_masks/meta_description
						ID:        "meta_description",
						Label:     `Mask for Meta Description`,
						Comment:   element.LongText(`Use {{name}} and {{description}} as Product Name and Product Description placeholders`),
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   `{{name}} {{description}}`,
					},
				),
			},

			&config.Group{
				ID:        "frontend",
				Label:     `Storefront`,
				SortOrder: 100,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/frontend/list_mode
						ID:        "list_mode",
						Label:     `List Mode`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `grid-list`,
						// SourceModel: Otnegam\Catalog\Model\Config\Source\ListMode
					},

					&config.Field{
						// Path: catalog/frontend/grid_per_page_values
						ID:        "grid_per_page_values",
						Label:     `Products per Page on Grid Allowed Values`,
						Comment:   element.LongText(`Comma-separated.`),
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `9,15,30`,
					},

					&config.Field{
						// Path: catalog/frontend/grid_per_page
						ID:        "grid_per_page",
						Label:     `Products per Page on Grid Default Value`,
						Comment:   element.LongText(`Must be in the allowed values list`),
						Type:      config.TypeText,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   9,
					},

					&config.Field{
						// Path: catalog/frontend/list_per_page_values
						ID:        "list_per_page_values",
						Label:     `Products per Page on List Allowed Values`,
						Comment:   element.LongText(`Comma-separated.`),
						Type:      config.TypeText,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `5,10,15,20,25`,
					},

					&config.Field{
						// Path: catalog/frontend/list_per_page
						ID:        "list_per_page",
						Label:     `Products per Page on List Default Value`,
						Comment:   element.LongText(`Must be in the allowed values list`),
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   10,
					},

					&config.Field{
						// Path: catalog/frontend/flat_catalog_category
						ID:        "flat_catalog_category",
						Label:     `Use Flat Catalog Category`,
						Type:      config.TypeSelect,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						Default:   false,
						// BackendModel: Otnegam\Catalog\Model\Indexer\Category\Flat\System\Config\Mode
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/frontend/flat_catalog_product
						ID:        "flat_catalog_product",
						Label:     `Use Flat Catalog Product`,
						Type:      config.TypeSelect,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Catalog\Model\Indexer\Product\Flat\System\Config\Mode
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/frontend/default_sort_by
						ID:        "default_sort_by",
						Label:     `Product Listing Sort by`,
						Type:      config.TypeSelect,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `position`,
						// SourceModel: Otnegam\Catalog\Model\Config\Source\ListSort
					},

					&config.Field{
						// Path: catalog/frontend/list_allow_all
						ID:        "list_allow_all",
						Label:     `Allow All Products per Page`,
						Comment:   element.LongText(`Whether to show "All" option in the "Show X Per Page" dropdown`),
						Type:      config.TypeSelect,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/frontend/parse_url_directives
						ID:        "parse_url_directives",
						Label:     `Allow Dynamic Media URLs in Products and Categories`,
						Comment:   element.LongText(`E.g. {{media url="path/to/image.jpg"}} {{skin url="path/to/picture.gif"}}. Dynamic directives parsing impacts catalog performance.`),
						Type:      config.TypeSelect,
						SortOrder: 200,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   true,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "placeholder",
				Label:     `Product Image Placeholders`,
				SortOrder: 300,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/placeholder/placeholder
						ID:        "placeholder",
						Type:      config.TypeImage,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Image
					},
				),
			},

			&config.Group{
				ID:        "seo",
				Label:     `Search Engine Optimization`,
				SortOrder: 500,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/seo/title_separator
						ID:        "title_separator",
						Label:     `Page Title Separator`,
						Type:      config.TypeText,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `-`,
					},

					&config.Field{
						// Path: catalog/seo/category_canonical_tag
						ID:        "category_canonical_tag",
						Label:     `Use Canonical Link Meta Tag For Categories`,
						Type:      config.TypeSelect,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/seo/product_canonical_tag
						ID:        "product_canonical_tag",
						Label:     `Use Canonical Link Meta Tag For Products`,
						Type:      config.TypeSelect,
						SortOrder: 8,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},

			&config.Group{
				ID:        "price",
				Label:     `Price`,
				SortOrder: 400,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/price/scope
						ID:        "scope",
						Label:     `Catalog Price Scope`,
						Comment:   element.LongText(`This defines the base currency scope ("Currency Setup" > "Currency Options" > "Base Currency").`),
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// BackendModel: Otnegam\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope
						// SourceModel: Otnegam\Catalog\Model\Config\Source\Price\Scope
					},
				),
			},

			&config.Group{
				ID:        "navigation",
				Label:     `Category Top Navigation`,
				SortOrder: 500,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/navigation/max_depth
						ID:        "max_depth",
						Label:     `Maximal Depth`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
					},
				),
			},

			&config.Group{
				ID:        "custom_options",
				Label:     `Date & Time Custom Options`,
				SortOrder: 700,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/custom_options/use_calendar
						ID:        "use_calendar",
						Label:     `Use JavaScript Calendar`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/custom_options/date_fields_order
						ID:        "date_fields_order",
						Label:     `Date Fields Order`,
						Type:      config.TypeSelect,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `m,d,y`,
					},

					&config.Field{
						// Path: catalog/custom_options/time_format
						ID:        "time_format",
						Label:     `Time Format`,
						Type:      config.TypeSelect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `12h`,
						// SourceModel: Otnegam\Catalog\Model\Config\Source\TimeFormat
					},

					&config.Field{
						// Path: catalog/custom_options/year_range
						ID:        "year_range",
						Label:     `Year Range`,
						Comment:   element.LongText(`Please use a four-digit year format.`),
						Type:      config.TypeText,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "design",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "watermark",
				Label:     `Product Image Watermarks`,
				SortOrder: 400,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: design/watermark/size
						ID:        "size",
						Label:     `Watermark Default Size`,
						Comment:   element.LongText(`Example format: 200x300.`),
						Type:      config.TypeText,
						SortOrder: 100,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: design/watermark/imageOpacity
						ID:        "imageOpacity",
						Label:     `Watermark Opacity, Percent`,
						Type:      config.TypeText,
						SortOrder: 150,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: design/watermark/image
						ID:        "image",
						Label:     `Watermark`,
						Comment:   element.LongText(`Allowed file types: jpeg, gif, png.`),
						Type:      config.TypeImage,
						SortOrder: 200,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Config\Model\Config\Backend\Image
					},

					&config.Field{
						// Path: design/watermark/position
						ID:        "position",
						Label:     `Watermark Position`,
						Type:      config.TypeSelect,
						SortOrder: 300,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Catalog\Model\Config\Source\Watermark\Position
					},
				),
			},
		),
	},
	&config.Section{
		ID: "cms",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "wysiwyg",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: cms/wysiwyg/use_static_urls_in_catalog
						ID:        "use_static_urls_in_catalog",
						Label:     `Use Static URLs for Media Content in WYSIWYG for Catalog`,
						Comment:   element.LongText(`This applies only to catalog products and categories. Media content will be inserted into the editor as a static URL. Media content is not updated if the system configuration base URL changes.`),
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID),
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
	&config.Section{
		ID: "rss",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "catalog",
				Label:     `Catalog`,
				SortOrder: 3,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: rss/catalog/new
						ID:        "new",
						Label:     `New Products`,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: rss/catalog/special
						ID:        "special",
						Label:     `Special Products`,
						Type:      config.TypeSelect,
						SortOrder: 11,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},

					&config.Field{
						// Path: rss/catalog/category
						ID:        "category",
						Label:     `Top Level Category`,
						Type:      config.TypeSelect,
						SortOrder: 14,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "catalog",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "product",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/product/flat
						ID:      `flat`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"max_index_count":"64"}`,
					},

					&config.Field{
						// Path: catalog/product/default_tax_group
						ID:      `default_tax_group`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: 2,
					},
				),
			},

			&config.Group{
				ID: "seo",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/seo/product_url_suffix
						ID:      `product_url_suffix`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `.html`,
					},

					&config.Field{
						// Path: catalog/seo/category_url_suffix
						ID:      `category_url_suffix`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `.html`,
					},

					&config.Field{
						// Path: catalog/seo/product_use_categories
						ID:      `product_use_categories`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: false,
					},

					&config.Field{
						// Path: catalog/seo/save_rewrites_history
						ID:      `save_rewrites_history`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: true,
					},
				),
			},

			&config.Group{
				ID: "custom_options",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/custom_options/forbidden_extensions
						ID:      `forbidden_extensions`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `php,exe`,
					},
				),
			},
		),
	},
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "media_storage_configuration",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: system/media_storage_configuration/allowed_resources
						ID:      `allowed_resources`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"tmp_images_folder":"tmp","catalog_images_folder":"catalog","product_custom_options_fodler":"custom_options"}`,
					},
				),
			},
		),
	},
)
