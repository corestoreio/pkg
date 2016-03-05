// +build ignore

package catalogurlrewrite

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
							Comment:   text.Long(`You need to refresh the cache.`),
							Type:      element.TypeText,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
						},

						&element.Field{
							// Path: catalog/seo/product_url_suffix
							ID:        "product_url_suffix",
							Label:     `Product URL Suffix`,
							Comment:   text.Long(`You need to refresh the cache.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// BackendModel: Magento\Catalog\Model\System\Config\Backend\Catalog\Url\Rewrite\Suffix
						},

						&element.Field{
							// Path: catalog/seo/product_use_categories
							ID:        "product_use_categories",
							Label:     `Use Categories Path for Product URLs`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: catalog/seo/save_rewrites_history
							ID:        "save_rewrites_history",
							Label:     `Create Permanent Redirect for URLs if URL Key Changed`,
							Type:      element.TypeSelect,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scopes:    scope.PermStore,
							// SourceModel: Magento\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},
	)
	Backend = NewBackend(ConfigStructure)
}
