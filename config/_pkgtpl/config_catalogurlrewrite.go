// +build ignore

package catalogurlrewrite

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "catalog",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:    "seo",
				Label: `Search Engine Optimization`,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: catalog/seo/category_url_suffix
						ID:        "category_url_suffix",
						Label:     `Category URL Suffix`,
						Comment:   element.LongText(`You need to refresh the cache.`),
						Type:      config.TypeText,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
					},

					&config.Field{
						// Path: catalog/seo/product_url_suffix
						ID:        "product_url_suffix",
						Label:     `Product URL Suffix`,
						Comment:   element.LongText(`You need to refresh the cache.`),
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
					},

					&config.Field{
						// Path: catalog/seo/product_use_categories
						ID:        "product_use_categories",
						Label:     `Use Categories Path for Product URLs`,
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: catalog/seo/save_rewrites_history
						ID:        "save_rewrites_history",
						Label:     `Create Permanent Redirect for URLs if URL Key Changed`,
						Type:      config.TypeSelect,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},
)
