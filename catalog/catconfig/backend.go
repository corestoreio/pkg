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
	"github.com/corestoreio/csfw/config/model"
)

// Backend will be initialized in the init() function together with ConfigStructure.
var Backend *PkgBackend

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	model.PkgBackend
	// CatalogFieldsMasksSku => Mask for SKU.
	// Use {{name}} as Product Name placeholder
	// Path: catalog/fields_masks/sku
	CatalogFieldsMasksSku model.Str

	// CatalogFieldsMasksMetaTitle => Mask for Meta Title.
	// Use {{name}} as Product Name placeholder
	// Path: catalog/fields_masks/meta_title
	CatalogFieldsMasksMetaTitle model.Str

	// CatalogFieldsMasksMetaKeyword => Mask for Meta Keywords.
	// Use {{name}} as Product Name or {{sku}} as Product SKU placeholders
	// Path: catalog/fields_masks/meta_keyword
	CatalogFieldsMasksMetaKeyword model.Str

	// CatalogFieldsMasksMetaDescription => Mask for Meta Description.
	// Use {{name}} and {{description}} as Product Name and Product Description
	// placeholders
	// Path: catalog/fields_masks/meta_description
	CatalogFieldsMasksMetaDescription model.Str

	// CatalogFrontendListMode => List Mode.
	// Path: catalog/frontend/list_mode
	// SourceModel: Magento\Catalog\Model\Config\Source\ListMode
	CatalogFrontendListMode model.Str

	// CatalogFrontendGridPerPageValues => Products per Page on Grid Allowed Values.
	// Comma-separated.
	// Path: catalog/frontend/grid_per_page_values
	CatalogFrontendGridPerPageValues model.Str

	// CatalogFrontendGridPerPage => Products per Page on Grid Default Value.
	// Must be in the allowed values list
	// Path: catalog/frontend/grid_per_page
	CatalogFrontendGridPerPage model.Str

	// CatalogFrontendListPerPageValues => Products per Page on List Allowed Values.
	// Comma-separated.
	// Path: catalog/frontend/list_per_page_values
	CatalogFrontendListPerPageValues model.Str

	// CatalogFrontendListPerPage => Products per Page on List Default Value.
	// Must be in the allowed values list
	// Path: catalog/frontend/list_per_page
	CatalogFrontendListPerPage model.Str

	// CatalogFrontendFlatCatalogCategory => Use Flat Catalog Category.
	// Path: catalog/frontend/flat_catalog_category
	// BackendModel: Magento\Catalog\Model\Indexer\Category\Flat\System\Config\Mode
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogFrontendFlatCatalogCategory model.Bool

	// CatalogFrontendFlatCatalogProduct => Use Flat Catalog Product.
	// Path: catalog/frontend/flat_catalog_product
	// BackendModel: Magento\Catalog\Model\Indexer\Product\Flat\System\Config\Mode
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogFrontendFlatCatalogProduct model.Bool

	// CatalogFrontendDefaultSortBy => Product Listing Sort by.
	// Path: catalog/frontend/default_sort_by
	// SourceModel: Magento\Catalog\Model\Config\Source\ListSort
	CatalogFrontendDefaultSortBy model.Str

	// CatalogFrontendListAllowAll => Allow All Products per Page.
	// Whether to show "All" option in the "Show X Per Page" dropdown
	// Path: catalog/frontend/list_allow_all
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogFrontendListAllowAll model.Bool

	// CatalogFrontendParseUrlDirectives => Allow Dynamic Media URLs in Products and Categories.
	// E.g. {{media url="path/to/image.jpg"}} {{skin url="path/to/picture.gif"}}.
	// Dynamic directives parsing impacts catalog performance.
	// Path: catalog/frontend/parse_url_directives
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogFrontendParseUrlDirectives model.Bool

	// CatalogPlaceholderPlaceholder => .
	// Path: catalog/placeholder/placeholder
	// BackendModel: Magento\Config\Model\Config\Backend\Image
	CatalogPlaceholderPlaceholder model.Str

	// CatalogSeoTitleSeparator => Page Title Separator.
	// Path: catalog/seo/title_separator
	CatalogSeoTitleSeparator model.Str

	// CatalogSeoCategoryCanonicalTag => Use Canonical Link Meta Tag For Categories.
	// Path: catalog/seo/category_canonical_tag
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogSeoCategoryCanonicalTag model.Bool

	// CatalogSeoProductCanonicalTag => Use Canonical Link Meta Tag For Products.
	// Path: catalog/seo/product_canonical_tag
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogSeoProductCanonicalTag model.Bool

	// CatalogPriceScope => Catalog Price Scope.
	// This defines the base currency scope ("Currency Setup" > "Currency Options"
	// > "Base Currency").
	// Path: catalog/price/scope
	// BackendModel: Magento\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope
	// SourceModel: Magento\Catalog\Model\Config\Source\Price\Scope
	CatalogPriceScope configPriceScope

	// CatalogNavigationMaxDepth => Maximal Depth.
	// Path: catalog/navigation/max_depth
	CatalogNavigationMaxDepth model.Str

	// CatalogCustomOptionsUseCalendar => Use JavaScript Calendar.
	// Path: catalog/custom_options/use_calendar
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogCustomOptionsUseCalendar model.Bool

	// CatalogCustomOptionsDateFieldsOrder => Date Fields Order.
	// Path: catalog/custom_options/date_fields_order
	CatalogCustomOptionsDateFieldsOrder model.Str

	// CatalogCustomOptionsTimeFormat => Time Format.
	// Path: catalog/custom_options/time_format
	// SourceModel: Magento\Catalog\Model\Config\Source\TimeFormat
	CatalogCustomOptionsTimeFormat model.Str

	// CatalogCustomOptionsYearRange => Year Range.
	// Please use a four-digit year format.
	// Path: catalog/custom_options/year_range
	CatalogCustomOptionsYearRange model.Str

	// DesignWatermarkSize => Watermark Default Size.
	// Example format: 200x300.
	// Path: design/watermark/size
	DesignWatermarkSize model.Str

	// DesignWatermarkImageOpacity => Watermark Opacity, Percent.
	// Path: design/watermark/imageOpacity
	DesignWatermarkImageOpacity model.Str

	// DesignWatermarkImage => Watermark.
	// Allowed file types: jpeg, gif, png.
	// Path: design/watermark/image
	// BackendModel: Magento\Config\Model\Config\Backend\Image
	DesignWatermarkImage model.Str

	// DesignWatermarkPosition => Watermark Position.
	// Path: design/watermark/position
	// SourceModel: Magento\Catalog\Model\Config\Source\Watermark\Position
	DesignWatermarkPosition model.Str

	// CmsWysiwygUseStaticUrlsInCatalog => Use Static URLs for Media Content in WYSIWYG for Catalog.
	// This applies only to catalog products and categories. Media content will be
	// inserted into the editor as a static URL. Media content is not updated if
	// the system configuration base URL changes.
	// Path: cms/wysiwyg/use_static_urls_in_catalog
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CmsWysiwygUseStaticUrlsInCatalog model.Bool

	// RssCatalogNew => New Products.
	// Path: rss/catalog/new
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssCatalogNew model.Bool

	// RssCatalogSpecial => Special Products.
	// Path: rss/catalog/special
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssCatalogSpecial model.Bool

	// RssCatalogCategory => Top Level Category.
	// Path: rss/catalog/category
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssCatalogCategory model.Bool
}

// NewBackend initializes the global Backend variable. See init()
func NewBackend(cfgStruct element.SectionSlice) *PkgBackend {
	return (&PkgBackend{}).init(cfgStruct)
}

func (pp *PkgBackend) init(cfgStruct element.SectionSlice) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogFieldsMasksSku = model.NewStr(`catalog/fields_masks/sku`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFieldsMasksMetaTitle = model.NewStr(`catalog/fields_masks/meta_title`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFieldsMasksMetaKeyword = model.NewStr(`catalog/fields_masks/meta_keyword`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFieldsMasksMetaDescription = model.NewStr(`catalog/fields_masks/meta_description`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendListMode = model.NewStr(`catalog/frontend/list_mode`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendGridPerPageValues = model.NewStr(`catalog/frontend/grid_per_page_values`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendGridPerPage = model.NewStr(`catalog/frontend/grid_per_page`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendListPerPageValues = model.NewStr(`catalog/frontend/list_per_page_values`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendListPerPage = model.NewStr(`catalog/frontend/list_per_page`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendFlatCatalogCategory = model.NewBool(`catalog/frontend/flat_catalog_category`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendFlatCatalogProduct = model.NewBool(`catalog/frontend/flat_catalog_product`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendDefaultSortBy = model.NewStr(`catalog/frontend/default_sort_by`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendListAllowAll = model.NewBool(`catalog/frontend/list_allow_all`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogFrontendParseUrlDirectives = model.NewBool(`catalog/frontend/parse_url_directives`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogPlaceholderPlaceholder = model.NewStr(`catalog/placeholder/placeholder`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSeoTitleSeparator = model.NewStr(`catalog/seo/title_separator`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSeoCategoryCanonicalTag = model.NewBool(`catalog/seo/category_canonical_tag`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogSeoProductCanonicalTag = model.NewBool(`catalog/seo/product_canonical_tag`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogPriceScope = NewConfigPriceScope(`catalog/price/scope`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogNavigationMaxDepth = model.NewStr(`catalog/navigation/max_depth`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogCustomOptionsUseCalendar = model.NewBool(`catalog/custom_options/use_calendar`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogCustomOptionsDateFieldsOrder = model.NewStr(`catalog/custom_options/date_fields_order`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogCustomOptionsTimeFormat = model.NewStr(`catalog/custom_options/time_format`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CatalogCustomOptionsYearRange = model.NewStr(`catalog/custom_options/year_range`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignWatermarkSize = model.NewStr(`design/watermark/size`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignWatermarkImageOpacity = model.NewStr(`design/watermark/imageOpacity`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignWatermarkImage = model.NewStr(`design/watermark/image`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.DesignWatermarkPosition = model.NewStr(`design/watermark/position`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.CmsWysiwygUseStaticUrlsInCatalog = model.NewBool(`cms/wysiwyg/use_static_urls_in_catalog`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.RssCatalogNew = model.NewBool(`rss/catalog/new`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.RssCatalogSpecial = model.NewBool(`rss/catalog/special`, model.WithFieldFromSectionSlice(cfgStruct))
	pp.RssCatalogCategory = model.NewBool(`rss/catalog/category`, model.WithFieldFromSectionSlice(cfgStruct))

	return pp
}
