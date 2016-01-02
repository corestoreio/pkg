// +build ignore

package catalog

import (
	"github.com/corestoreio/csfw/config/model"
)

// PathCatalogFieldsMasksSku => Mask for SKU.
// Use {{name}} as Product Name placeholder
var PathCatalogFieldsMasksSku = model.NewStr(`catalog/fields_masks/sku`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFieldsMasksMetaTitle => Mask for Meta Title.
// Use {{name}} as Product Name placeholder
var PathCatalogFieldsMasksMetaTitle = model.NewStr(`catalog/fields_masks/meta_title`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFieldsMasksMetaKeyword => Mask for Meta Keywords.
// Use {{name}} as Product Name or {{sku}} as Product SKU placeholders
var PathCatalogFieldsMasksMetaKeyword = model.NewStr(`catalog/fields_masks/meta_keyword`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFieldsMasksMetaDescription => Mask for Meta Description.
// Use {{name}} and {{description}} as Product Name and Product Description
// placeholders
var PathCatalogFieldsMasksMetaDescription = model.NewStr(`catalog/fields_masks/meta_description`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendListMode => List Mode.
// SourceModel: Otnegam\Catalog\Model\Config\Source\ListMode
var PathCatalogFrontendListMode = model.NewStr(`catalog/frontend/list_mode`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendGridPerPageValues => Products per Page on Grid Allowed Values.
// Comma-separated.
var PathCatalogFrontendGridPerPageValues = model.NewStr(`catalog/frontend/grid_per_page_values`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendGridPerPage => Products per Page on Grid Default Value.
// Must be in the allowed values list
var PathCatalogFrontendGridPerPage = model.NewStr(`catalog/frontend/grid_per_page`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendListPerPageValues => Products per Page on List Allowed Values.
// Comma-separated.
var PathCatalogFrontendListPerPageValues = model.NewStr(`catalog/frontend/list_per_page_values`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendListPerPage => Products per Page on List Default Value.
// Must be in the allowed values list
var PathCatalogFrontendListPerPage = model.NewStr(`catalog/frontend/list_per_page`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendFlatCatalogCategory => Use Flat Catalog Category.
// BackendModel: Otnegam\Catalog\Model\Indexer\Category\Flat\System\Config\Mode
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogFrontendFlatCatalogCategory = model.NewBool(`catalog/frontend/flat_catalog_category`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendFlatCatalogProduct => Use Flat Catalog Product.
// BackendModel: Otnegam\Catalog\Model\Indexer\Product\Flat\System\Config\Mode
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogFrontendFlatCatalogProduct = model.NewBool(`catalog/frontend/flat_catalog_product`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendDefaultSortBy => Product Listing Sort by.
// SourceModel: Otnegam\Catalog\Model\Config\Source\ListSort
var PathCatalogFrontendDefaultSortBy = model.NewStr(`catalog/frontend/default_sort_by`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendListAllowAll => Allow All Products per Page.
// Whether to show "All" option in the "Show X Per Page" dropdown
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogFrontendListAllowAll = model.NewBool(`catalog/frontend/list_allow_all`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogFrontendParseUrlDirectives => Allow Dynamic Media URLs in Products and Categories.
// E.g. {{media url="path/to/image.jpg"}} {{skin url="path/to/picture.gif"}}.
// Dynamic directives parsing impacts catalog performance.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogFrontendParseUrlDirectives = model.NewBool(`catalog/frontend/parse_url_directives`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogPlaceholderPlaceholder => .
// BackendModel: Otnegam\Config\Model\Config\Backend\Image
var PathCatalogPlaceholderPlaceholder = model.NewStr(`catalog/placeholder/placeholder`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogSeoTitleSeparator => Page Title Separator.
var PathCatalogSeoTitleSeparator = model.NewStr(`catalog/seo/title_separator`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogSeoCategoryCanonicalTag => Use Canonical Link Meta Tag For Categories.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogSeoCategoryCanonicalTag = model.NewBool(`catalog/seo/category_canonical_tag`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogSeoProductCanonicalTag => Use Canonical Link Meta Tag For Products.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogSeoProductCanonicalTag = model.NewBool(`catalog/seo/product_canonical_tag`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogPriceScope => Catalog Price Scope.
// This defines the base currency scope ("Currency Setup" > "Currency Options"
// > "Base Currency").
// BackendModel: Otnegam\Catalog\Model\Indexer\Product\Price\System\Config\PriceScope
// SourceModel: Otnegam\Catalog\Model\Config\Source\Price\Scope
var PathCatalogPriceScope = model.NewStr(`catalog/price/scope`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogNavigationMaxDepth => Maximal Depth.
var PathCatalogNavigationMaxDepth = model.NewStr(`catalog/navigation/max_depth`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogCustomOptionsUseCalendar => Use JavaScript Calendar.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCatalogCustomOptionsUseCalendar = model.NewBool(`catalog/custom_options/use_calendar`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogCustomOptionsDateFieldsOrder => Date Fields Order.
var PathCatalogCustomOptionsDateFieldsOrder = model.NewStr(`catalog/custom_options/date_fields_order`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogCustomOptionsTimeFormat => Time Format.
// SourceModel: Otnegam\Catalog\Model\Config\Source\TimeFormat
var PathCatalogCustomOptionsTimeFormat = model.NewStr(`catalog/custom_options/time_format`, model.WithPkgCfg(PackageConfiguration))

// PathCatalogCustomOptionsYearRange => Year Range.
// Please use a four-digit year format.
var PathCatalogCustomOptionsYearRange = model.NewStr(`catalog/custom_options/year_range`, model.WithPkgCfg(PackageConfiguration))

// PathDesignWatermarkSize => Watermark Default Size.
// Example format: 200x300.
var PathDesignWatermarkSize = model.NewStr(`design/watermark/size`, model.WithPkgCfg(PackageConfiguration))

// PathDesignWatermarkImageOpacity => Watermark Opacity, Percent.
var PathDesignWatermarkImageOpacity = model.NewStr(`design/watermark/imageOpacity`, model.WithPkgCfg(PackageConfiguration))

// PathDesignWatermarkImage => Watermark.
// Allowed file types: jpeg, gif, png.
// BackendModel: Otnegam\Config\Model\Config\Backend\Image
var PathDesignWatermarkImage = model.NewStr(`design/watermark/image`, model.WithPkgCfg(PackageConfiguration))

// PathDesignWatermarkPosition => Watermark Position.
// SourceModel: Otnegam\Catalog\Model\Config\Source\Watermark\Position
var PathDesignWatermarkPosition = model.NewStr(`design/watermark/position`, model.WithPkgCfg(PackageConfiguration))

// PathCmsWysiwygUseStaticUrlsInCatalog => Use Static URLs for Media Content in WYSIWYG for Catalog.
// This applies only to catalog products and categories. Media content will be
// inserted into the editor as a static URL. Media content is not updated if
// the system configuration base URL changes.
// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
var PathCmsWysiwygUseStaticUrlsInCatalog = model.NewBool(`cms/wysiwyg/use_static_urls_in_catalog`, model.WithPkgCfg(PackageConfiguration))

// PathRssCatalogNew => New Products.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathRssCatalogNew = model.NewBool(`rss/catalog/new`, model.WithPkgCfg(PackageConfiguration))

// PathRssCatalogSpecial => Special Products.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathRssCatalogSpecial = model.NewBool(`rss/catalog/special`, model.WithPkgCfg(PackageConfiguration))

// PathRssCatalogCategory => Top Level Category.
// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
var PathRssCatalogCategory = model.NewBool(`rss/catalog/category`, model.WithPkgCfg(PackageConfiguration))
