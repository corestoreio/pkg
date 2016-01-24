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

package catconfig

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
			ID:        "catalog",
			Label:     `Catalog`,
			SortOrder: 40,
			Scope:     scope.PermAll,
			Resource:  0, // Magento_Catalog::config_catalog
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "fields_masks",
					Label:     `Product Fields Auto-Generation`,
					SortOrder: 90,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/fields_masks/sku
							ID:        "sku",
							Label:     `Mask for SKU`,
							Comment:   element.LongText(`Use {{name}} as Product Name placeholder`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `{{name}}`,
						},

						&element.Field{
							// Path: catalog/fields_masks/meta_title
							ID:        "meta_title",
							Label:     `Mask for Meta Title`,
							Comment:   element.LongText(`Use {{name}} as Product Name placeholder`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `{{name}}`,
						},

						&element.Field{
							// Path: catalog/fields_masks/meta_keyword
							ID:        "meta_keyword",
							Label:     `Mask for Meta Keywords`,
							Comment:   element.LongText(`Use {{name}} as Product Name or {{sku}} as Product SKU placeholders`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `{{name}}`,
						},

						&element.Field{
							// Path: catalog/fields_masks/meta_description
							ID:        "meta_description",
							Label:     `Mask for Meta Description`,
							Comment:   element.LongText(`Use {{name}} and {{description}} as Product Name and Product Description placeholders`),
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   `{{name}} {{description}}`,
						},
					),
				},

				&element.Group{
					ID:        "frontend",
					Label:     `Storefront`,
					SortOrder: 100,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/frontend/list_mode
							ID:        "list_mode",
							Label:     `List Mode`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `grid-list`,
							// SourceModel: Magento\Catalog\Model\Config\Source\ListMode
						},

						&element.Field{
							// Path: catalog/frontend/grid_per_page_values
							ID:        "grid_per_page_values",
							Label:     `Products per Page on Grid Allowed Values`,
							Comment:   element.LongText(`Comma-separated.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `9,15,30`,
						},

						&element.Field{
							// Path: catalog/frontend/grid_per_page
							ID:        "grid_per_page",
							Label:     `Products per Page on Grid Default Value`,
							Comment:   element.LongText(`Must be in the allowed values list`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   9,
						},

						&element.Field{
							// Path: catalog/frontend/list_per_page_values
							ID:        "list_per_page_values",
							Label:     `Products per Page on List Allowed Values`,
							Comment:   element.LongText(`Comma-separated.`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `5,10,15,20,25`,
						},

						&element.Field{
							// Path: catalog/frontend/list_per_page
							ID:        "list_per_page",
							Label:     `Products per Page on List Default Value`,
							Comment:   element.LongText(`Must be in the allowed values list`),
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   10,
						},

						&element.Field{
							// Path: catalog/frontend/flat_catalog_category
							ID:        "flat_catalog_category",
							Label:     `Use Flat Catalog Category`,
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							Default:   false,
							// BackendModel: Magento\Catalog\Model\Indexer\Category\Flat\System\Config\Mode
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/frontend/flat_catalog_product
							ID:        "flat_catalog_product",
							Label:     `Use Flat Catalog Product`,
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Catalog\Model\Indexer\Product\Flat\System\Config\Mode
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/frontend/default_sort_by
							ID:        "default_sort_by",
							Label:     `Product Listing Sort by`,
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `position`,
							// SourceModel: Magento\Catalog\Model\Config\Source\ListSort
						},

						&element.Field{
							// Path: catalog/frontend/list_allow_all
							ID:        "list_allow_all",
							Label:     `Allow All Products per Page`,
							Comment:   element.LongText(`Whether to show "All" option in the "Show X Per Page" dropdown`),
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/frontend/parse_url_directives
							ID:        "parse_url_directives",
							Label:     `Allow Dynamic Media URLs in Products and Categories`,
							Comment:   element.LongText(`E.g. {{media url="path/to/image.jpg"}} {{skin url="path/to/picture.gif"}}. Dynamic directives parsing impacts catalog performance.`),
							Type:      element.TypeSelect,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "placeholder",
					Label:     `Product Image Placeholders`,
					SortOrder: 300,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/placeholder/placeholder
							ID:        "placeholder",
							Type:      element.TypeImage,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Magento\Config\Model\Config\Backend\Image
						},
					),
				},

				&element.Group{
					ID:        "seo",
					Label:     `Search Engine Optimization`,
					SortOrder: 500,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/seo/title_separator
							ID:        "title_separator",
							Label:     `Page Title Separator`,
							Type:      element.TypeText,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `-`,
						},

						&element.Field{
							// Path: catalog/seo/category_canonical_tag
							ID:        "category_canonical_tag",
							Label:     `Use Canonical Link Meta Tag For Categories`,
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/seo/product_canonical_tag
							ID:        "product_canonical_tag",
							Label:     `Use Canonical Link Meta Tag For Products`,
							Type:      element.TypeSelect,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        "price",
					Label:     `Price`,
					SortOrder: 400,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/price/scope
							ID:        "scope",
							Label:     `Catalog Price Scope`,
							Comment:   element.LongText(`This defines the base currency scope ("Currency Setup" > "Currency Options" > "Base Currency").`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// BackendModel: Magento\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope
							// SourceModel: Magento\Catalog\Model\Config\Source\Price\Scope
						},
					),
				},

				&element.Group{
					ID:        "navigation",
					Label:     `Category Top Navigation`,
					SortOrder: 500,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/navigation/max_depth
							ID:        "max_depth",
							Label:     `Maximal Depth`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
						},
					),
				},

				&element.Group{
					ID:        "custom_options",
					Label:     `Date & Time Custom Options`,
					SortOrder: 700,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/custom_options/use_calendar
							ID:        "use_calendar",
							Label:     `Use JavaScript Calendar`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/custom_options/date_fields_order
							ID:        "date_fields_order",
							Label:     `Date Fields Order`,
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `m,d,y`,
						},

						&element.Field{
							// Path: catalog/custom_options/time_format
							ID:        "time_format",
							Label:     `Time Format`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `12h`,
							// SourceModel: Magento\Catalog\Model\Config\Source\TimeFormat
						},

						&element.Field{
							// Path: catalog/custom_options/year_range
							ID:        "year_range",
							Label:     `Year Range`,
							Comment:   element.LongText(`Please use a four-digit year format.`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},
			),
		},
		&element.Section{
			ID: "design",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "watermark",
					Label:     `Product Image Watermarks`,
					SortOrder: 400,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/watermark/size
							ID:        "size",
							Label:     `Watermark Default Size`,
							Comment:   element.LongText(`Example format: 200x300.`),
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/watermark/imageOpacity
							ID:        "imageOpacity",
							Label:     `Watermark Opacity, Percent`,
							Type:      element.TypeText,
							SortOrder: 150,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: design/watermark/image
							ID:        "image",
							Label:     `Watermark`,
							Comment:   element.LongText(`Allowed file types: jpeg, gif, png.`),
							Type:      element.TypeImage,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Magento\Config\Model\Config\Backend\Image
						},

						&element.Field{
							// Path: design/watermark/position
							ID:        "position",
							Label:     `Watermark Position`,
							Type:      element.TypeSelect,
							SortOrder: 300,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Magento\Catalog\Model\Config\Source\Watermark\Position
						},
					),
				},
			),
		},
		&element.Section{
			ID: "cms",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "wysiwyg",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: cms/wysiwyg/use_static_urls_in_catalog
							ID:        "use_static_urls_in_catalog",
							Label:     `Use Static URLs for Media Content in WYSIWYG for Catalog`,
							Comment:   element.LongText(`This applies only to catalog products and categories. Media content will be inserted into the editor as a static URL. Media content is not updated if the system configuration base URL changes.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID),
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID: "rss",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "catalog",
					Label:     `Catalog`,
					SortOrder: 3,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: rss/catalog/new
							ID:        "new",
							Label:     `New Products`,
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},

						&element.Field{
							// Path: rss/catalog/special
							ID:        "special",
							Label:     `Special Products`,
							Type:      element.TypeSelect,
							SortOrder: 11,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},

						&element.Field{
							// Path: rss/catalog/category
							ID:        "category",
							Label:     `Top Level Category`,
							Type:      element.TypeSelect,
							SortOrder: 14,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "catalog",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "product",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/product/flat
							ID:      `flat`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"max_index_count":"64"}`,
						},

						&element.Field{
							// Path: catalog/product/default_tax_group
							ID:      `default_tax_group`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 2,
						},
					),
				},

				&element.Group{
					ID: "seo",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/seo/product_url_suffix
							ID:      `product_url_suffix`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `.html`,
						},

						&element.Field{
							// Path: catalog/seo/category_url_suffix
							ID:      `category_url_suffix`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `.html`,
						},

						&element.Field{
							// Path: catalog/seo/product_use_categories
							ID:      `product_use_categories`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: catalog/seo/save_rewrites_history
							ID:      `save_rewrites_history`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},

				&element.Group{
					ID: "custom_options",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/custom_options/forbidden_extensions
							ID:      `forbidden_extensions`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `php,exe`,
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
							Default: `{"tmp_images_folder":"tmp","catalog_images_folder":"catalog","product_custom_options_fodler":"custom_options"}`,
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
