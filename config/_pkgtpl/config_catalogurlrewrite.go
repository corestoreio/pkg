// +build ignore

package catalogurlrewrite

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = element.MustNewConfiguration(
	&element.Section{
		ID: "catalog",
		Groups: element.NewGroupSlice(
			&element.Group{
				ID:    "seo",
				Label: `Search Engine Optimization`,
				Fields: element.NewFieldSlice(
					&element.Field{
						// Path: catalog/seo/category_url_suffix
						ID:        "category_url_suffix",
						Label:     `Category URL Suffix`,
						Comment:   element.LongText(`You need to refresh the cache.`),
						Type:      element.TypeText,
						SortOrder: 3,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
					},

					&element.Field{
						// Path: catalog/seo/product_url_suffix
						ID:        "product_url_suffix",
						Label:     `Product URL Suffix`,
						Comment:   element.LongText(`You need to refresh the cache.`),
						Type:      element.TypeText,
						SortOrder: 2,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
					},

					&element.Field{
						// Path: catalog/seo/product_use_categories
						ID:        "product_use_categories",
						Label:     `Use Categories Path for Product URLs`,
						Type:      element.TypeSelect,
						SortOrder: 4,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&element.Field{
						// Path: catalog/seo/save_rewrites_history
						ID:        "save_rewrites_history",
						Label:     `Create Permanent Redirect for URLs if URL Key Changed`,
						Type:      element.TypeSelect,
						SortOrder: 5,
						Visible:   element.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
)
