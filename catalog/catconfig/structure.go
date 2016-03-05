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
	"github.com/corestoreio/csfw/config/path"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/corestoreio/csfw/store/scope"
)

// ConfigStructure global configuration structure for this package.
// Used in frontend and backend. See init() for details.
var ConfigStructure element.SectionSlice

func init() {
	ConfigStructure = element.MustNewConfiguration(
		&element.Section{
			ID:        path.NewRoute("catalog"),
			Label:     text.Chars(`Catalog`),
			SortOrder: 40,
			Scope:     scope.PermStore,
			Resource:  0, // Magento_Catalog::config_catalog
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("fields_masks"),
					Label:     text.Chars(`Product Fields Auto-Generation`),
					SortOrder: 90,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/fields_masks/sku
							ID:        path.NewRoute("sku"),
							Label:     text.Chars(`Mask for SKU`),
							Comment:   text.Chars(`Use {{name}} as Product Name placeholder`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `{{name}}`,
						},

						&element.Field{
							// Path: catalog/fields_masks/meta_title
							ID:        path.NewRoute("meta_title"),
							Label:     text.Chars(`Mask for Meta Title`),
							Comment:   text.Chars(`Use {{name}} as Product Name placeholder`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `{{name}}`,
						},

						&element.Field{
							// Path: catalog/fields_masks/meta_keyword
							ID:        path.NewRoute("meta_keyword"),
							Label:     text.Chars(`Mask for Meta Keywords`),
							Comment:   text.Chars(`Use {{name}} as Product Name or {{sku}} as Product SKU placeholders`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `{{name}}`,
						},

						&element.Field{
							// Path: catalog/fields_masks/meta_description
							ID:        path.NewRoute("meta_description"),
							Label:     text.Chars(`Mask for Meta Description`),
							Comment:   text.Chars(`Use {{name}} and {{description}} as Product Name and Product Description placeholders`),
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `{{name}} {{description}}`,
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("frontend"),
					Label:     text.Chars(`Storefront`),
					SortOrder: 100,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/frontend/list_mode
							ID:        path.NewRoute("list_mode"),
							Label:     text.Chars(`List Mode`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `grid-list`,
							// SourceModel: Magento\Catalog\Model\Config\Source\ListMode
						},

						&element.Field{
							// Path: catalog/frontend/grid_per_page_values
							ID:        path.NewRoute("grid_per_page_values"),
							Label:     text.Chars(`Products per Page on Grid Allowed Values`),
							Comment:   text.Chars(`Comma-separated.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `9,15,30`,
						},

						&element.Field{
							// Path: catalog/frontend/grid_per_page
							ID:        path.NewRoute("grid_per_page"),
							Label:     text.Chars(`Products per Page on Grid Default Value`),
							Comment:   text.Chars(`Must be in the allowed values list`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   9,
						},

						&element.Field{
							// Path: catalog/frontend/list_per_page_values
							ID:        path.NewRoute("list_per_page_values"),
							Label:     text.Chars(`Products per Page on List Allowed Values`),
							Comment:   text.Chars(`Comma-separated.`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `5,10,15,20,25`,
						},

						&element.Field{
							// Path: catalog/frontend/list_per_page
							ID:        path.NewRoute("list_per_page"),
							Label:     text.Chars(`Products per Page on List Default Value`),
							Comment:   text.Chars(`Must be in the allowed values list`),
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   10,
						},

						&element.Field{
							// Path: catalog/frontend/flat_catalog_category
							ID:        path.NewRoute("flat_catalog_category"),
							Label:     text.Chars(`Use Flat Catalog Category`),
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   false,
							// BackendModel: Magento\Catalog\Model\Indexer\Category\Flat\System\Config\Mode
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/frontend/flat_catalog_product
							ID:        path.NewRoute("flat_catalog_product"),
							Label:     text.Chars(`Use Flat Catalog Product`),
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Catalog\Model\Indexer\Product\Flat\System\Config\Mode
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/frontend/default_sort_by
							ID:        path.NewRoute("default_sort_by"),
							Label:     text.Chars(`Product Listing Sort by`),
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `position`,
							// SourceModel: Magento\Catalog\Model\Config\Source\ListSort
						},

						&element.Field{
							// Path: catalog/frontend/list_allow_all
							ID:        path.NewRoute("list_allow_all"),
							Label:     text.Chars(`Allow All Products per Page`),
							Comment:   text.Chars(`Whether to show "All" option in the "Show X Per Page" dropdown`),
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/frontend/parse_url_directives
							ID:        path.NewRoute("parse_url_directives"),
							Label:     text.Chars(`Allow Dynamic Media URLs in Products and Categories`),
							Comment:   text.Chars(`E.g. {{media url="path/to/image.jpg"}} {{skin url="path/to/picture.gif"}}. Dynamic directives parsing impacts catalog performance.`),
							Type:      element.TypeSelect,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   true,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("placeholder"),
					Label:     text.Chars(`Product Image Placeholders`),
					SortOrder: 300,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/placeholder/placeholder
							ID:        path.NewRoute("placeholder"),
							Type:      element.TypeImage,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Image
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("seo"),
					Label:     text.Chars(`Search Engine Optimization`),
					SortOrder: 500,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/seo/title_separator
							ID:        path.NewRoute("title_separator"),
							Label:     text.Chars(`Page Title Separator`),
							Type:      element.TypeText,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `-`,
						},

						&element.Field{
							// Path: catalog/seo/category_canonical_tag
							ID:        path.NewRoute("category_canonical_tag"),
							Label:     text.Chars(`Use Canonical Link Meta Tag For Categories`),
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/seo/product_canonical_tag
							ID:        path.NewRoute("product_canonical_tag"),
							Label:     text.Chars(`Use Canonical Link Meta Tag For Products`),
							Type:      element.TypeSelect,
							SortOrder: 8,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("price"),
					Label:     text.Chars(`Price`),
					SortOrder: 400,
					Scope:     scope.PermDefault,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/price/scope
							ID:        path.NewRoute("scope"),
							Label:     text.Chars(`Catalog Price Scope`),
							Comment:   text.Chars(`This defines the base currency scope ("Currency Setup" > "Currency Options" > "Base Currency").`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope
							// SourceModel: Magento\Catalog\Model\Config\Source\Price\Scope
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("navigation"),
					Label:     text.Chars(`Category Top Navigation`),
					SortOrder: 500,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/navigation/max_depth
							ID:        path.NewRoute("max_depth"),
							Label:     text.Chars(`Maximal Depth`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
					),
				},

				&element.Group{
					ID:        path.NewRoute("custom_options"),
					Label:     text.Chars(`Date & Time Custom Options`),
					SortOrder: 700,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/custom_options/use_calendar
							ID:        path.NewRoute("use_calendar"),
							Label:     text.Chars(`Use JavaScript Calendar`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/custom_options/date_fields_order
							ID:        path.NewRoute("date_fields_order"),
							Label:     text.Chars(`Date Fields Order`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `m,d,y`,
						},

						&element.Field{
							// Path: catalog/custom_options/time_format
							ID:        path.NewRoute("time_format"),
							Label:     text.Chars(`Time Format`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `12h`,
							// SourceModel: Magento\Catalog\Model\Config\Source\TimeFormat
						},

						&element.Field{
							// Path: catalog/custom_options/year_range
							ID:        path.NewRoute("year_range"),
							Label:     text.Chars(`Year Range`),
							Comment:   text.Chars(`Please use a four-digit year format.`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},
					),
				},
			),
		},
		&element.Section{
			ID: path.NewRoute("design"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("watermark"),
					Label:     text.Chars(`Product Image Watermarks`),
					SortOrder: 400,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: design/watermark/size
							ID:        path.NewRoute("size"),
							Label:     text.Chars(`Watermark Default Size`),
							Comment:   text.Chars(`Example format: 200x300.`),
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						&element.Field{
							// Path: design/watermark/imageOpacity
							ID:        path.NewRoute("imageOpacity"),
							Label:     text.Chars(`Watermark Opacity, Percent`),
							Type:      element.TypeText,
							SortOrder: 150,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						&element.Field{
							// Path: design/watermark/image
							ID:        path.NewRoute("image"),
							Label:     text.Chars(`Watermark`),
							Comment:   text.Chars(`Allowed file types: jpeg, gif, png.`),
							Type:      element.TypeImage,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Image
						},

						&element.Field{
							// Path: design/watermark/position
							ID:        path.NewRoute("position"),
							Label:     text.Chars(`Watermark Position`),
							Type:      element.TypeSelect,
							SortOrder: 300,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Catalog\Model\Config\Source\Watermark\Position
						},
					),
				},
			),
		},
		&element.Section{
			ID: path.NewRoute("cms"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: path.NewRoute("wysiwyg"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: cms/wysiwyg/use_static_urls_in_catalog
							ID:        path.NewRoute("use_static_urls_in_catalog"),
							Label:     text.Chars(`Use Static URLs for Media Content in WYSIWYG for Catalog`),
							Comment:   text.Chars(`This applies only to catalog products and categories. Media content will be inserted into the editor as a static URL. Media content is not updated if the system configuration base URL changes.`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
		&element.Section{
			ID: path.NewRoute("rss"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        path.NewRoute("catalog"),
					Label:     text.Chars(`Catalog`),
					SortOrder: 3,
					Scope:     scope.PermStore,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: rss/catalog/new
							ID:        path.NewRoute("new"),
							Label:     text.Chars(`New Products`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},

						&element.Field{
							// Path: rss/catalog/special
							ID:        path.NewRoute("special"),
							Label:     text.Chars(`Special Products`),
							Type:      element.TypeSelect,
							SortOrder: 11,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},

						&element.Field{
							// Path: rss/catalog/category
							ID:        path.NewRoute("category"),
							Label:     text.Chars(`Top Level Category`),
							Type:      element.TypeSelect,
							SortOrder: 14,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: path.NewRoute("catalog"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: path.NewRoute("product"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/product/flat
							ID:      path.NewRoute(`flat`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"max_index_count":"64"}`,
						},

						&element.Field{
							// Path: catalog/product/default_tax_group
							ID:      path.NewRoute(`default_tax_group`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 2,
						},
					),
				},

				&element.Group{
					ID: path.NewRoute("seo"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/seo/product_url_suffix
							ID:      path.NewRoute(`product_url_suffix`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `.html`,
						},

						&element.Field{
							// Path: catalog/seo/category_url_suffix
							ID:      path.NewRoute(`category_url_suffix`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `.html`,
						},

						&element.Field{
							// Path: catalog/seo/product_use_categories
							ID:      path.NewRoute(`product_use_categories`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						&element.Field{
							// Path: catalog/seo/save_rewrites_history
							ID:      path.NewRoute(`save_rewrites_history`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},

				&element.Group{
					ID: path.NewRoute("custom_options"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: catalog/custom_options/forbidden_extensions
							ID:      path.NewRoute(`forbidden_extensions`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `php,exe`,
						},
					),
				},
			),
		},
		&element.Section{
			ID: path.NewRoute("system"),
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: path.NewRoute("media_storage_configuration"),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      path.NewRoute(`allowed_resources`),
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
