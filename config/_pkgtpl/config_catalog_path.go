// +build ignore

package catalog

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with PackageConfiguration.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
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
	// SourceModel: Otnegam\Catalog\Model\Config\Source\ListMode
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
	// BackendModel: Otnegam\Catalog\Model\Indexer\Category\Flat\System\Config\Mode
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogFrontendFlatCatalogCategory model.Bool

	// CatalogFrontendFlatCatalogProduct => Use Flat Catalog Product.
	// Path: catalog/frontend/flat_catalog_product
	// BackendModel: Otnegam\Catalog\Model\Indexer\Product\Flat\System\Config\Mode
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogFrontendFlatCatalogProduct model.Bool

	// CatalogFrontendDefaultSortBy => Product Listing Sort by.
	// Path: catalog/frontend/default_sort_by
	// SourceModel: Otnegam\Catalog\Model\Config\Source\ListSort
	CatalogFrontendDefaultSortBy model.Str

	// CatalogFrontendListAllowAll => Allow All Products per Page.
	// Whether to show "All" option in the "Show X Per Page" dropdown
	// Path: catalog/frontend/list_allow_all
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogFrontendListAllowAll model.Bool

	// CatalogFrontendParseUrlDirectives => Allow Dynamic Media URLs in Products and Categories.
	// E.g. {{media url="path/to/image.jpg"}} {{skin url="path/to/picture.gif"}}.
	// Dynamic directives parsing impacts catalog performance.
	// Path: catalog/frontend/parse_url_directives
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogFrontendParseUrlDirectives model.Bool

	// CatalogPlaceholderPlaceholder => .
	// Path: catalog/placeholder/placeholder
	// BackendModel: Otnegam\Config\Model\Config\Backend\Image
	CatalogPlaceholderPlaceholder model.Str

	// CatalogSeoTitleSeparator => Page Title Separator.
	// Path: catalog/seo/title_separator
	CatalogSeoTitleSeparator model.Str

	// CatalogSeoCategoryCanonicalTag => Use Canonical Link Meta Tag For Categories.
	// Path: catalog/seo/category_canonical_tag
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogSeoCategoryCanonicalTag model.Bool

	// CatalogSeoProductCanonicalTag => Use Canonical Link Meta Tag For Products.
	// Path: catalog/seo/product_canonical_tag
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogSeoProductCanonicalTag model.Bool

	// CatalogPriceScope => Catalog Price Scope.
	// This defines the base currency scope ("Currency Setup" > "Currency Options"
	// > "Base Currency").
	// Path: catalog/price/scope
	// BackendModel: Otnegam\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope
	// SourceModel: Otnegam\Catalog\Model\Config\Source\Price\Scope
	CatalogPriceScope model.Str

	// CatalogNavigationMaxDepth => Maximal Depth.
	// Path: catalog/navigation/max_depth
	CatalogNavigationMaxDepth model.Str

	// CatalogCustomOptionsUseCalendar => Use JavaScript Calendar.
	// Path: catalog/custom_options/use_calendar
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CatalogCustomOptionsUseCalendar model.Bool

	// CatalogCustomOptionsDateFieldsOrder => Date Fields Order.
	// Path: catalog/custom_options/date_fields_order
	CatalogCustomOptionsDateFieldsOrder model.Str

	// CatalogCustomOptionsTimeFormat => Time Format.
	// Path: catalog/custom_options/time_format
	// SourceModel: Otnegam\Catalog\Model\Config\Source\TimeFormat
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
	// BackendModel: Otnegam\Config\Model\Config\Backend\Image
	DesignWatermarkImage model.Str

	// DesignWatermarkPosition => Watermark Position.
	// Path: design/watermark/position
	// SourceModel: Otnegam\Catalog\Model\Config\Source\Watermark\Position
	DesignWatermarkPosition model.Str

	// CmsWysiwygUseStaticUrlsInCatalog => Use Static URLs for Media Content in WYSIWYG for Catalog.
	// This applies only to catalog products and categories. Media content will be
	// inserted into the editor as a static URL. Media content is not updated if
	// the system configuration base URL changes.
	// Path: cms/wysiwyg/use_static_urls_in_catalog
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	CmsWysiwygUseStaticUrlsInCatalog model.Bool

	// RssCatalogNew => New Products.
	// Path: rss/catalog/new
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	RssCatalogNew model.Bool

	// RssCatalogSpecial => Special Products.
	// Path: rss/catalog/special
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	RssCatalogSpecial model.Bool

	// RssCatalogCategory => Top Level Category.
	// Path: rss/catalog/category
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	RssCatalogCategory model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(pkgCfg element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(pkgCfg)
}

func (pp *PkgPath) init(pkgCfg element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.CatalogFieldsMasksSku = model.NewStr(`catalog/fields_masks/sku`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFieldsMasksMetaTitle = model.NewStr(`catalog/fields_masks/meta_title`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFieldsMasksMetaKeyword = model.NewStr(`catalog/fields_masks/meta_keyword`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFieldsMasksMetaDescription = model.NewStr(`catalog/fields_masks/meta_description`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendListMode = model.NewStr(`catalog/frontend/list_mode`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendGridPerPageValues = model.NewStr(`catalog/frontend/grid_per_page_values`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendGridPerPage = model.NewStr(`catalog/frontend/grid_per_page`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendListPerPageValues = model.NewStr(`catalog/frontend/list_per_page_values`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendListPerPage = model.NewStr(`catalog/frontend/list_per_page`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendFlatCatalogCategory = model.NewBool(`catalog/frontend/flat_catalog_category`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendFlatCatalogProduct = model.NewBool(`catalog/frontend/flat_catalog_product`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendDefaultSortBy = model.NewStr(`catalog/frontend/default_sort_by`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendListAllowAll = model.NewBool(`catalog/frontend/list_allow_all`, model.WithPkgCfg(pkgCfg))
	pp.CatalogFrontendParseUrlDirectives = model.NewBool(`catalog/frontend/parse_url_directives`, model.WithPkgCfg(pkgCfg))
	pp.CatalogPlaceholderPlaceholder = model.NewStr(`catalog/placeholder/placeholder`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSeoTitleSeparator = model.NewStr(`catalog/seo/title_separator`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSeoCategoryCanonicalTag = model.NewBool(`catalog/seo/category_canonical_tag`, model.WithPkgCfg(pkgCfg))
	pp.CatalogSeoProductCanonicalTag = model.NewBool(`catalog/seo/product_canonical_tag`, model.WithPkgCfg(pkgCfg))
	pp.CatalogPriceScope = model.NewStr(`catalog/price/scope`, model.WithPkgCfg(pkgCfg))
	pp.CatalogNavigationMaxDepth = model.NewStr(`catalog/navigation/max_depth`, model.WithPkgCfg(pkgCfg))
	pp.CatalogCustomOptionsUseCalendar = model.NewBool(`catalog/custom_options/use_calendar`, model.WithPkgCfg(pkgCfg))
	pp.CatalogCustomOptionsDateFieldsOrder = model.NewStr(`catalog/custom_options/date_fields_order`, model.WithPkgCfg(pkgCfg))
	pp.CatalogCustomOptionsTimeFormat = model.NewStr(`catalog/custom_options/time_format`, model.WithPkgCfg(pkgCfg))
	pp.CatalogCustomOptionsYearRange = model.NewStr(`catalog/custom_options/year_range`, model.WithPkgCfg(pkgCfg))
	pp.DesignWatermarkSize = model.NewStr(`design/watermark/size`, model.WithPkgCfg(pkgCfg))
	pp.DesignWatermarkImageOpacity = model.NewStr(`design/watermark/imageOpacity`, model.WithPkgCfg(pkgCfg))
	pp.DesignWatermarkImage = model.NewStr(`design/watermark/image`, model.WithPkgCfg(pkgCfg))
	pp.DesignWatermarkPosition = model.NewStr(`design/watermark/position`, model.WithPkgCfg(pkgCfg))
	pp.CmsWysiwygUseStaticUrlsInCatalog = model.NewBool(`cms/wysiwyg/use_static_urls_in_catalog`, model.WithPkgCfg(pkgCfg))
	pp.RssCatalogNew = model.NewBool(`rss/catalog/new`, model.WithPkgCfg(pkgCfg))
	pp.RssCatalogSpecial = model.NewBool(`rss/catalog/special`, model.WithPkgCfg(pkgCfg))
	pp.RssCatalogCategory = model.NewBool(`rss/catalog/category`, model.WithPkgCfg(pkgCfg))

	return pp
}
