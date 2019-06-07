// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/config/element"
	"github.com/corestoreio/pkg/storage/text"
	"github.com/corestoreio/pkg/store/scope"
)

// MustNewConfigStructure same as NewConfigStructure() but panics on error.
func MustNewConfigStructure() element.Sections {
	ss, err := NewConfigStructure()
	if err != nil {
		panic(err)
	}
	return ss
}

// NewConfigStructure global configuration structure for this package.
// Used in frontend (to display the user all the settings) and in
// backend (scope checks and default values). See the source code
// of this function for the overall available sections, groups and fields.
func NewConfigStructure() (element.Sections, error) {
	return element.MakeSectionsValidated(
		element.Section{
			ID:        cfgpath.MakeRoute("catalog"),
			Label:     text.Chars(`Catalog`),
			SortOrder: 40,
			Scopes:    scope.PermStore,
			Resource:  0, // Magento_Catalog::config_catalog
			Groups: element.MakeGroups(
				element.Group{
					ID:        cfgpath.MakeRoute("fields_masks"),
					Label:     text.Chars(`Product Fields Auto-Generation`),
					SortOrder: 90,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/fields_masks/sku
							ID:        cfgpath.MakeRoute("sku"),
							Label:     text.Chars(`Mask for SKU`),
							Comment:   text.Chars(`Use {{name}} as Product Name placeholder`),
							Type:      element.TypeText,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `{{name}}`,
						},

						element.Field{
							// Path: catalog/fields_masks/meta_title
							ID:        cfgpath.MakeRoute("meta_title"),
							Label:     text.Chars(`Mask for Meta Title`),
							Comment:   text.Chars(`Use {{name}} as Product Name placeholder`),
							Type:      element.TypeText,
							SortOrder: 20,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `{{name}}`,
						},

						element.Field{
							// Path: catalog/fields_masks/meta_keyword
							ID:        cfgpath.MakeRoute("meta_keyword"),
							Label:     text.Chars(`Mask for Meta Keywords`),
							Comment:   text.Chars(`Use {{name}} as Product Name or {{sku}} as Product SKU placeholders`),
							Type:      element.TypeText,
							SortOrder: 30,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   `{{name}}`,
						},

						element.Field{
							// Path: catalog/fields_masks/meta_description
							ID:        cfgpath.MakeRoute("meta_description"),
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

				element.Group{
					ID:        cfgpath.MakeRoute("frontend"),
					Label:     text.Chars(`Storefront`),
					SortOrder: 100,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/frontend/list_mode
							ID:        cfgpath.MakeRoute("list_mode"),
							Label:     text.Chars(`List Mode`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `grid-list`,
							// SourceModel: Magento\Catalog\Model\Config\Source\ListMode
						},

						element.Field{
							// Path: catalog/frontend/grid_per_page_values
							ID:        cfgpath.MakeRoute("grid_per_page_values"),
							Label:     text.Chars(`Products per Page on Grid Allowed Values`),
							Comment:   text.Chars(`Comma-separated.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `9,15,30`,
						},

						element.Field{
							// Path: catalog/frontend/grid_per_page
							ID:        cfgpath.MakeRoute("grid_per_page"),
							Label:     text.Chars(`Products per Page on Grid Default Value`),
							Comment:   text.Chars(`Must be in the allowed values list`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   9,
						},

						element.Field{
							// Path: catalog/frontend/list_per_page_values
							ID:        cfgpath.MakeRoute("list_per_page_values"),
							Label:     text.Chars(`Products per Page on List Allowed Values`),
							Comment:   text.Chars(`Comma-separated.`),
							Type:      element.TypeText,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `5,10,15,20,25`,
						},

						element.Field{
							// Path: catalog/frontend/list_per_page
							ID:        cfgpath.MakeRoute("list_per_page"),
							Label:     text.Chars(`Products per Page on List Default Value`),
							Comment:   text.Chars(`Must be in the allowed values list`),
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   10,
						},

						element.Field{
							// Path: catalog/frontend/flat_catalog_category
							ID:        cfgpath.MakeRoute("flat_catalog_category"),
							Label:     text.Chars(`Use Flat Catalog Category`),
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							Default:   false,
							// BackendModel: Magento\Catalog\Model\Indexer\Category\Flat\System\Config\Mode
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/frontend/flat_catalog_product
							ID:        cfgpath.MakeRoute("flat_catalog_product"),
							Label:     text.Chars(`Use Flat Catalog Product`),
							Type:      element.TypeSelect,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
							// BackendModel: Magento\Catalog\Model\Indexer\Product\Flat\System\Config\Mode
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/frontend/default_sort_by
							ID:        cfgpath.MakeRoute("default_sort_by"),
							Label:     text.Chars(`Product Listing Sort by`),
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `position`,
							// SourceModel: Magento\Catalog\Model\Config\Source\ListSort
						},

						element.Field{
							// Path: catalog/frontend/list_allow_all
							ID:        cfgpath.MakeRoute("list_allow_all"),
							Label:     text.Chars(`Allow All Products per Page`),
							Comment:   text.Chars(`Whether to show "All" option in the "Show X Per Page" dropdown`),
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/frontend/parse_url_directives
							ID:        cfgpath.MakeRoute("parse_url_directives"),
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

				element.Group{
					ID:        cfgpath.MakeRoute("placeholder"),
					Label:     text.Chars(`Product Image Placeholders`),
					SortOrder: 300,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/placeholder/placeholder
							ID:        cfgpath.MakeRoute("placeholder"),
							Type:      element.TypeImage,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Image
						},
					),
				},

				element.Group{
					ID:        cfgpath.MakeRoute("seo"),
					Label:     text.Chars(`Search Engine Optimization`),
					SortOrder: 500,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/seo/title_separator
							ID:        cfgpath.MakeRoute("title_separator"),
							Label:     text.Chars(`Page Title Separator`),
							Type:      element.TypeText,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `-`,
						},

						element.Field{
							// Path: catalog/seo/category_canonical_tag
							ID:        cfgpath.MakeRoute("category_canonical_tag"),
							Label:     text.Chars(`Use Canonical Link Meta Tag For Categories`),
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   false,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/seo/product_canonical_tag
							ID:        cfgpath.MakeRoute("product_canonical_tag"),
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

				element.Group{
					ID:        cfgpath.MakeRoute("price"),
					Label:     text.Chars(`Price`),
					SortOrder: 400,
					Scopes:    scope.PermDefault,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/price/scope
							ID:        cfgpath.MakeRoute("scope"),
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

				element.Group{
					ID:        cfgpath.MakeRoute("navigation"),
					Label:     text.Chars(`Category Top Navigation`),
					SortOrder: 500,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/navigation/max_depth
							ID:        cfgpath.MakeRoute("max_depth"),
							Label:     text.Chars(`Maximal Depth`),
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermDefault,
						},
					),
				},

				element.Group{
					ID:        cfgpath.MakeRoute("custom_options"),
					Label:     text.Chars(`Date & Time Custom Options`),
					SortOrder: 700,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/custom_options/use_calendar
							ID:        cfgpath.MakeRoute("use_calendar"),
							Label:     text.Chars(`Use JavaScript Calendar`),
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						element.Field{
							// Path: catalog/custom_options/date_fields_order
							ID:        cfgpath.MakeRoute("date_fields_order"),
							Label:     text.Chars(`Date Fields Order`),
							Type:      element.TypeSelect,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `m,d,y`,
						},

						element.Field{
							// Path: catalog/custom_options/time_format
							ID:        cfgpath.MakeRoute("time_format"),
							Label:     text.Chars(`Time Format`),
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							Default:   `12h`,
							// SourceModel: Magento\Catalog\Model\Config\Source\TimeFormat
						},

						element.Field{
							// Path: catalog/custom_options/year_range
							ID:        cfgpath.MakeRoute("year_range"),
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
		element.Section{
			ID: cfgpath.MakeRoute("design"),
			Groups: element.MakeGroups(
				element.Group{
					ID:        cfgpath.MakeRoute("watermark"),
					Label:     text.Chars(`Product Image Watermarks`),
					SortOrder: 400,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: design/watermark/size
							ID:        cfgpath.MakeRoute("size"),
							Label:     text.Chars(`Watermark Default Size`),
							Comment:   text.Chars(`Example format: 200x300.`),
							Type:      element.TypeText,
							SortOrder: 100,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: design/watermark/imageOpacity
							ID:        cfgpath.MakeRoute("imageOpacity"),
							Label:     text.Chars(`Watermark Opacity, Percent`),
							Type:      element.TypeText,
							SortOrder: 150,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
						},

						element.Field{
							// Path: design/watermark/image
							ID:        cfgpath.MakeRoute("image"),
							Label:     text.Chars(`Watermark`),
							Comment:   text.Chars(`Allowed file types: jpeg, gif, png.`),
							Type:      element.TypeImage,
							SortOrder: 200,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Config\Model\Config\Backend\Image
						},

						element.Field{
							// Path: design/watermark/position
							ID:        cfgpath.MakeRoute("position"),
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
		element.Section{
			ID: cfgpath.MakeRoute("cms"),
			Groups: element.MakeGroups(
				element.Group{
					ID: cfgpath.MakeRoute("wysiwyg"),
					Fields: element.MakeFields(
						element.Field{
							// Path: cms/wysiwyg/use_static_urls_in_catalog
							ID:        cfgpath.MakeRoute("use_static_urls_in_catalog"),
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
		element.Section{
			ID: cfgpath.MakeRoute("rss"),
			Groups: element.MakeGroups(
				element.Group{
					ID:        cfgpath.MakeRoute("catalog"),
					Label:     text.Chars(`Catalog`),
					SortOrder: 3,
					Scopes:    scope.PermStore,
					Fields: element.MakeFields(
						element.Field{
							// Path: rss/catalog/new
							ID:        cfgpath.MakeRoute("new"),
							Label:     text.Chars(`New Products`),
							Type:      element.TypeSelect,
							SortOrder: 10,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},

						element.Field{
							// Path: rss/catalog/special
							ID:        cfgpath.MakeRoute("special"),
							Label:     text.Chars(`Special Products`),
							Type:      element.TypeSelect,
							SortOrder: 11,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
						},

						element.Field{
							// Path: rss/catalog/category
							ID:        cfgpath.MakeRoute("category"),
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
		element.Section{
			ID: cfgpath.MakeRoute("catalog"),
			Groups: element.MakeGroups(
				element.Group{
					ID: cfgpath.MakeRoute("product"),
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/product/flat
							ID:      cfgpath.MakeRoute(`flat`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"max_index_count":"64"}`,
						},

						element.Field{
							// Path: catalog/product/default_tax_group
							ID:      cfgpath.MakeRoute(`default_tax_group`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: 2,
						},
					),
				},

				element.Group{
					ID: cfgpath.MakeRoute("seo"),
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/seo/product_url_suffix
							ID:      cfgpath.MakeRoute(`product_url_suffix`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `.html`,
						},

						element.Field{
							// Path: catalog/seo/category_url_suffix
							ID:      cfgpath.MakeRoute(`category_url_suffix`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `.html`,
						},

						element.Field{
							// Path: catalog/seo/product_use_categories
							ID:      cfgpath.MakeRoute(`product_use_categories`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: false,
						},

						element.Field{
							// Path: catalog/seo/save_rewrites_history
							ID:      cfgpath.MakeRoute(`save_rewrites_history`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: true,
						},
					),
				},

				element.Group{
					ID: cfgpath.MakeRoute("custom_options"),
					Fields: element.MakeFields(
						element.Field{
							// Path: catalog/custom_options/forbidden_extensions
							ID:      cfgpath.MakeRoute(`forbidden_extensions`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `php,exe`,
						},
					),
				},
			),
		},
		element.Section{
			ID: cfgpath.MakeRoute("system"),
			Groups: element.MakeGroups(
				element.Group{
					ID: cfgpath.MakeRoute("media_storage_configuration"),
					Fields: element.MakeFields(
						element.Field{
							// Path: system/media_storage_configuration/allowed_resources
							ID:      cfgpath.MakeRoute(`allowed_resources`),
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"tmp_images_folder":"tmp","catalog_images_folder":"catalog","product_custom_options_fodler":"custom_options"}`,
						},
					),
				},
			),
		},
	)
}
