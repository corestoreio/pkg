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
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/element"
)

// PkgBackend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type PkgBackend struct {
	cfgmodel.PkgBackend
	// CatalogFieldsMasksSku => Mask for SKU.
	// Use {{name}} as Product Name placeholder
	// Path: catalog/fields_masks/sku
	CatalogFieldsMasksSku cfgmodel.Str

	// CatalogFieldsMasksMetaTitle => Mask for Meta Title.
	// Use {{name}} as Product Name placeholder
	// Path: catalog/fields_masks/meta_title
	CatalogFieldsMasksMetaTitle cfgmodel.Str

	// CatalogFieldsMasksMetaKeyword => Mask for Meta Keywords.
	// Use {{name}} as Product Name or {{sku}} as Product SKU placeholders
	// Path: catalog/fields_masks/meta_keyword
	CatalogFieldsMasksMetaKeyword cfgmodel.Str

	// CatalogFieldsMasksMetaDescription => Mask for Meta Description.
	// Use {{name}} and {{description}} as Product Name and Product Description
	// placeholders
	// Path: catalog/fields_masks/meta_description
	CatalogFieldsMasksMetaDescription cfgmodel.Str

	// CatalogFrontendListMode => List Mode.
	// Path: catalog/frontend/list_mode
	// SourceModel: Magento\Catalog\Model\Config\Source\ListMode
	CatalogFrontendListMode cfgmodel.Str

	// CatalogFrontendGridPerPageValues => Products per Page on Grid Allowed Values.
	// Comma-separated.
	// Path: catalog/frontend/grid_per_page_values
	CatalogFrontendGridPerPageValues cfgmodel.Str

	// CatalogFrontendGridPerPage => Products per Page on Grid Default Value.
	// Must be in the allowed values list
	// Path: catalog/frontend/grid_per_page
	CatalogFrontendGridPerPage cfgmodel.Str

	// CatalogFrontendListPerPageValues => Products per Page on List Allowed Values.
	// Comma-separated.
	// Path: catalog/frontend/list_per_page_values
	CatalogFrontendListPerPageValues cfgmodel.Str

	// CatalogFrontendListPerPage => Products per Page on List Default Value.
	// Must be in the allowed values list
	// Path: catalog/frontend/list_per_page
	CatalogFrontendListPerPage cfgmodel.Str

	// CatalogFrontendFlatCatalogCategory => Use Flat Catalog Category.
	// Path: catalog/frontend/flat_catalog_category
	// BackendModel: Magento\Catalog\Model\Indexer\Category\Flat\System\Config\Mode
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogFrontendFlatCatalogCategory cfgmodel.Bool

	// CatalogFrontendFlatCatalogProduct => Use Flat Catalog Product.
	// Path: catalog/frontend/flat_catalog_product
	// BackendModel: Magento\Catalog\Model\Indexer\Product\Flat\System\Config\Mode
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogFrontendFlatCatalogProduct cfgmodel.Bool

	// CatalogFrontendDefaultSortBy => Product Listing Sort by.
	// Path: catalog/frontend/default_sort_by
	// SourceModel: Magento\Catalog\Model\Config\Source\ListSort
	CatalogFrontendDefaultSortBy cfgmodel.Str

	// CatalogFrontendListAllowAll => Allow All Products per Page.
	// Whether to show "All" option in the "Show X Per Page" dropdown
	// Path: catalog/frontend/list_allow_all
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogFrontendListAllowAll cfgmodel.Bool

	// CatalogFrontendParseUrlDirectives => Allow Dynamic Media URLs in Products and Categories.
	// E.g. {{media url="path/to/image.jpg"}} {{skin url="path/to/picture.gif"}}.
	// Dynamic directives parsing impacts catalog performance.
	// Path: catalog/frontend/parse_url_directives
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogFrontendParseUrlDirectives cfgmodel.Bool

	// CatalogPlaceholderPlaceholder => .
	// Path: catalog/placeholder/placeholder
	// BackendModel: Magento\Config\Model\Config\Backend\Image
	CatalogPlaceholderPlaceholder cfgmodel.Str

	// CatalogSeoTitleSeparator => Page Title Separator.
	// Path: catalog/seo/title_separator
	CatalogSeoTitleSeparator cfgmodel.Str

	// CatalogSeoCategoryCanonicalTag => Use Canonical Link Meta Tag For Categories.
	// Path: catalog/seo/category_canonical_tag
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogSeoCategoryCanonicalTag cfgmodel.Bool

	// CatalogSeoProductCanonicalTag => Use Canonical Link Meta Tag For Products.
	// Path: catalog/seo/product_canonical_tag
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogSeoProductCanonicalTag cfgmodel.Bool

	// CatalogPriceScope => Catalog Price Scope.
	// This defines the base currency scope ("Currency Setup" > "Currency Options"
	// > "Base Currency").
	// Path: catalog/price/scope
	// BackendModel: Magento\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope
	// SourceModel: Magento\Catalog\Model\Config\Source\Price\Scope
	CatalogPriceScope PriceScope

	// CatalogNavigationMaxDepth => Maximal Depth.
	// Path: catalog/navigation/max_depth
	CatalogNavigationMaxDepth cfgmodel.Str

	// CatalogCustomOptionsUseCalendar => Use JavaScript Calendar.
	// Path: catalog/custom_options/use_calendar
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CatalogCustomOptionsUseCalendar cfgmodel.Bool

	// CatalogCustomOptionsDateFieldsOrder => Date Fields Order.
	// Path: catalog/custom_options/date_fields_order
	CatalogCustomOptionsDateFieldsOrder cfgmodel.Str

	// CatalogCustomOptionsTimeFormat => Time Format.
	// Path: catalog/custom_options/time_format
	// SourceModel: Magento\Catalog\Model\Config\Source\TimeFormat
	CatalogCustomOptionsTimeFormat cfgmodel.Str

	// CatalogCustomOptionsYearRange => Year Range.
	// Please use a four-digit year format.
	// Path: catalog/custom_options/year_range
	CatalogCustomOptionsYearRange cfgmodel.Str

	// DesignWatermarkSize => Watermark Default Size.
	// Example format: 200x300.
	// Path: design/watermark/size
	DesignWatermarkSize cfgmodel.Str

	// DesignWatermarkImageOpacity => Watermark Opacity, Percent.
	// Path: design/watermark/imageOpacity
	DesignWatermarkImageOpacity cfgmodel.Str

	// DesignWatermarkImage => Watermark.
	// Allowed file types: jpeg, gif, png.
	// Path: design/watermark/image
	// BackendModel: Magento\Config\Model\Config\Backend\Image
	DesignWatermarkImage cfgmodel.Str

	// DesignWatermarkPosition => Watermark Position.
	// Path: design/watermark/position
	// SourceModel: Magento\Catalog\Model\Config\Source\Watermark\Position
	DesignWatermarkPosition cfgmodel.Str

	// CmsWysiwygUseStaticUrlsInCatalog => Use Static URLs for Media Content in WYSIWYG for Catalog.
	// This applies only to catalog products and categories. Media content will be
	// inserted into the editor as a static URL. Media content is not updated if
	// the system configuration base URL changes.
	// Path: cms/wysiwyg/use_static_urls_in_catalog
	// SourceModel: Magento\Config\Model\Config\Source\Yesno
	CmsWysiwygUseStaticUrlsInCatalog cfgmodel.Bool

	// RssCatalogNew => New Products.
	// Path: rss/catalog/new
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssCatalogNew cfgmodel.Bool

	// RssCatalogSpecial => Special Products.
	// Path: rss/catalog/special
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssCatalogSpecial cfgmodel.Bool

	// RssCatalogCategory => Top Level Category.
	// Path: rss/catalog/category
	// SourceModel: Magento\Config\Model\Config\Source\Enabledisable
	RssCatalogCategory cfgmodel.Bool
}

// NewBackend initializes the global configuration models containing the
// cfgpath.Route variable to the appropriate entry.
// The function Load() will be executed to apply the Sections
// to all models. See Load() for more details.
func NewBackend(cfgStruct element.Sections) *PkgBackend {
	return (&PkgBackend{}).Load(cfgStruct)
}

// Load creates the configuration models for each PkgBackend field.
// Internal mutex will protect the fields during loading.
// The argument Sections will be applied to all models.
func (pp *PkgBackend) Load(cfgStruct element.Sections) *PkgBackend {
	pp.Lock()
	defer pp.Unlock()

	opt := cfgmodel.WithFieldFromSectionSlice(cfgStruct)

	pp.CatalogFieldsMasksSku = cfgmodel.NewStr(`catalog/fields_masks/sku`, opt)
	pp.CatalogFieldsMasksMetaTitle = cfgmodel.NewStr(`catalog/fields_masks/meta_title`, opt)
	pp.CatalogFieldsMasksMetaKeyword = cfgmodel.NewStr(`catalog/fields_masks/meta_keyword`, opt)
	pp.CatalogFieldsMasksMetaDescription = cfgmodel.NewStr(`catalog/fields_masks/meta_description`, opt)
	pp.CatalogFrontendListMode = cfgmodel.NewStr(`catalog/frontend/list_mode`, opt)
	pp.CatalogFrontendGridPerPageValues = cfgmodel.NewStr(`catalog/frontend/grid_per_page_values`, opt)
	pp.CatalogFrontendGridPerPage = cfgmodel.NewStr(`catalog/frontend/grid_per_page`, opt)
	pp.CatalogFrontendListPerPageValues = cfgmodel.NewStr(`catalog/frontend/list_per_page_values`, opt)
	pp.CatalogFrontendListPerPage = cfgmodel.NewStr(`catalog/frontend/list_per_page`, opt)
	pp.CatalogFrontendFlatCatalogCategory = cfgmodel.NewBool(`catalog/frontend/flat_catalog_category`, opt)
	pp.CatalogFrontendFlatCatalogProduct = cfgmodel.NewBool(`catalog/frontend/flat_catalog_product`, opt)
	pp.CatalogFrontendDefaultSortBy = cfgmodel.NewStr(`catalog/frontend/default_sort_by`, opt)
	pp.CatalogFrontendListAllowAll = cfgmodel.NewBool(`catalog/frontend/list_allow_all`, opt)
	pp.CatalogFrontendParseUrlDirectives = cfgmodel.NewBool(`catalog/frontend/parse_url_directives`, opt)
	pp.CatalogPlaceholderPlaceholder = cfgmodel.NewStr(`catalog/placeholder/placeholder`, opt)
	pp.CatalogSeoTitleSeparator = cfgmodel.NewStr(`catalog/seo/title_separator`, opt)
	pp.CatalogSeoCategoryCanonicalTag = cfgmodel.NewBool(`catalog/seo/category_canonical_tag`, opt)
	pp.CatalogSeoProductCanonicalTag = cfgmodel.NewBool(`catalog/seo/product_canonical_tag`, opt)
	pp.CatalogPriceScope = NewPriceScope(`catalog/price/scope`, opt)
	pp.CatalogNavigationMaxDepth = cfgmodel.NewStr(`catalog/navigation/max_depth`, opt)
	pp.CatalogCustomOptionsUseCalendar = cfgmodel.NewBool(`catalog/custom_options/use_calendar`, opt)
	pp.CatalogCustomOptionsDateFieldsOrder = cfgmodel.NewStr(`catalog/custom_options/date_fields_order`, opt)
	pp.CatalogCustomOptionsTimeFormat = cfgmodel.NewStr(`catalog/custom_options/time_format`, opt)
	pp.CatalogCustomOptionsYearRange = cfgmodel.NewStr(`catalog/custom_options/year_range`, opt)
	pp.DesignWatermarkSize = cfgmodel.NewStr(`design/watermark/size`, opt)
	pp.DesignWatermarkImageOpacity = cfgmodel.NewStr(`design/watermark/imageOpacity`, opt)
	pp.DesignWatermarkImage = cfgmodel.NewStr(`design/watermark/image`, opt)
	pp.DesignWatermarkPosition = cfgmodel.NewStr(`design/watermark/position`, opt)
	pp.CmsWysiwygUseStaticUrlsInCatalog = cfgmodel.NewBool(`cms/wysiwyg/use_static_urls_in_catalog`, opt)
	pp.RssCatalogNew = cfgmodel.NewBool(`rss/catalog/new`, opt)
	pp.RssCatalogSpecial = cfgmodel.NewBool(`rss/catalog/special`, opt)
	pp.RssCatalogCategory = cfgmodel.NewBool(`rss/catalog/category`, opt)

	return pp
}
