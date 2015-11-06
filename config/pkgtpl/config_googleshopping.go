// +build ignore

package googleshopping

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID:        "google",
		Label:     "",
		SortOrder: 0,
		Scope:     scope.NewPerm(),
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "googleshopping",
				Label:     `Google Shopping`,
				Comment:   ``,
				SortOrder: 20,
				Scope:     scope.PermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `google/googleshopping/account_id`,
						ID:           "account_id",
						Label:        `Account ID`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `google/googleshopping/login`,
						ID:           "login",
						Label:        `Account Login`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `google/googleshopping/password`,
						ID:           "password",
						Label:        `Account Password`,
						Comment:      ``,
						Type:         config.TypeObscure,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil, // Magento\Config\Model\Config\Backend\Encrypted // @todo Magento\Config\Model\Config\Backend\Encrypted
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `google/googleshopping/account_type`,
						ID:           "account_type",
						Label:        `Account Type`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `HOSTED_OR_GOOGLE`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\GoogleShopping\Model\Source\Accounttype
					},

					&config.Field{
						// Path: `google/googleshopping/target_country`,
						ID:           "target_country",
						Label:        `Target Country`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    90,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `US`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\GoogleShopping\Model\Source\Country
					},

					&config.Field{
						// Path: `google/googleshopping/observed`,
						ID:           "observed",
						Label:        `Update Google Shopping Item when Product is Updated`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    100,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `google/googleshopping/verify_meta_tag`,
						ID:           "verify_meta_tag",
						Label:        `Verifying Meta Tag`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    110,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `google/googleshopping/debug`,
						ID:           "debug",
						Label:        `Debug`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    140,
						Visible:      config.VisibleYes,
						Scope:        scope.NewPerm(config.IDScopeDefault, config.IDScopeWebsite),
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `google/googleshopping/destinations`,
						ID:           "destinations",
						Label:        `Destinations`,
						Comment:      ``,
						Type:         config.Type,
						SortOrder:    150,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      `{"product_search":"ProductSearch","product_ads":"ProductAds","commerce_search":"CommerceSearch"}`,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `google/googleshopping/product_search`,
						ID:           "product_search",
						Label:        `Google Product Search`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    151,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\GoogleShopping\Model\Source\Destinationstates
					},

					&config.Field{
						// Path: `google/googleshopping/product_ads`,
						ID:           "product_ads",
						Label:        `Product Ads`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    152,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\GoogleShopping\Model\Source\Destinationstates
					},

					&config.Field{
						// Path: `google/googleshopping/commerce_search`,
						ID:           "commerce_search",
						Label:        `Commerce Search`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    153,
						Visible:      config.VisibleYes,
						Scope:        scope.PermAll,
						Default:      nil,
						BackendModel: nil,
						SourceModel:  nil, // Magento\GoogleShopping\Model\Source\Destinationstates
					},
				},
			},
		},
	},

	// Hidden Configuration
	&config.Section{
		ID: "google",
		Groups: config.GroupSlice{
			&config.Group{
				ID: "googleshopping",
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `google/googleshopping/allowed_countries`,
						ID:      "allowed_countries",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(config.IDScopeDefault), // @todo search for that
						Default: `{"AU":{"_value":{"name":"Australia","language":"en","currency":"AUD","currency_name":"Australian Dollar"},"_attribute":{"translate":"name currency_name"}},"BR":{"_value":{"name":"Brazil","language":"pt","locale":"pt_BR","currency":"BRL","currency_name":"Brazilian Real"},"_attribute":{"translate":"name currency_name"}},"CN":{"_value":{"name":"China","language":"zh_CN","currency":"CNY","currency_name":"Chinese Yuan Renminbi"},"_attribute":{"translate":"name currency_name"}},"FR":{"_value":{"name":"France","language":"fr","currency":"EUR","currency_name":"Euro"},"_attribute":{"translate":"name currency_name"}},"DE":{"_value":{"name":"Germany","language":"de","locale":"de_DE","currency":"EUR","currency_name":"Euro"},"_attribute":{"translate":"name currency_name"}},"IT":{"_value":{"name":"Italy","language":"it","currency":"EUR","currency_name":"Euro"},"_attribute":{"translate":"name currency_name"}},"JP":{"_value":{"name":"Japan","language":"ja","currency":"JPY","currency_name":"Japanese Yen"},"_attribute":{"translate":"name currency_name"}},"NL":{"_value":{"name":"Netherlands","language":"nl","currency":"EUR","currency_name":"Euro"},"_attribute":{"translate":"name currency_name"}},"ES":{"_value":{"name":"Spain","language":"es","currency":"EUR","currency_name":"Euro"},"_attribute":{"translate":"name currency_name"}},"CH":{"_value":{"name":"Switzerland","language":"de","locale":"de_CH","currency":"CHF","currency_name":"Swiss Franc"},"_attribute":{"translate":"name currency_name"}},"GB":{"_value":{"name":"United Kingdom","language":"en","locale":"en_GB","currency":"GBP","currency_name":"British Pound Sterling"},"_attribute":{"translate":"name currency_name"}},"US":{"_value":{"name":"United States","language":"en","locale":"en_US","currency":"USD","currency_name":"US Dollar"},"_attribute":{"translate":"name currency_name"}}}`,
					},

					&config.Field{
						// Path: `google/googleshopping/attributes`,
						ID:      "attributes",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(config.IDScopeDefault), // @todo search for that
						Default: `{"Item":{"title":{"_value":{"name":"Title","required":"1"},"_attribute":{"translate":"name"}},"content":{"_value":{"name":"Description","required":"1"},"_attribute":{"translate":"name"}},"expiration_date":{"_value":{"name":"Expiration date","required":"0"},"_attribute":{"translate":"name"}},"adult":{"_value":{"name":"Adult","required":"0"},"_attribute":{"translate":"name"}}},"ProductSearch":{"condition":{"_value":{"name":"Condition","required":"1"},"_attribute":{"translate":"name"}},"price":{"_value":{"name":"Price","required":"1"},"_attribute":{"translate":"name"}},"sale_price":{"_value":{"name":"Sale Price","required":"0","country":"US"},"_attribute":{"translate":"name"}},"sale_price_effective_date_from":{"_value":{"name":"Sale Price Effective From Date","required":"0","country":"US"},"_attribute":{"translate":"name"}},"sale_price_effective_date_to":{"_value":{"name":"Sale Price Effective To Date","required":"0","country":"US"},"_attribute":{"translate":"name"}},"age_group":{"_value":{"name":"Age Group","required":"1"},"_attribute":{"translate":"name"}},"brand":{"_value":{"name":"Brand","required":"1"},"_attribute":{"translate":"name"}},"color":{"_value":{"name":"Color","required":"1"},"_attribute":{"translate":"name"}},"gender":{"_value":{"name":"Gender","required":"1"},"_attribute":{"translate":"name"}},"mpn":{"_value":{"name":"Manufacturer\\'s Part Number (MPN)","required":"1"},"_attribute":{"translate":"name"}},"online_only":{"_value":{"name":"Online Only","required":"0"},"_attribute":{"translate":"name"}},"gtin":{"_value":{"name":"GTIN","required":"1"},"_attribute":{"translate":"name"}},"product_type":{"_value":{"name":"Product Type (Category)","required":"0"},"_attribute":{"translate":"name"}},"product_review_average":{"_value":{"name":"Product Review Average","required":"0"},"_attribute":{"translate":"name"}},"product_review_count":{"_value":{"name":"Product Review Count","required":"0"},"_attribute":{"translate":"name"}},"shipping_weight":{"_value":{"name":"Shipping Weight","required":"0"},"_attribute":{"translate":"name"}},"size":{"_value":{"name":"Size","required":"1"},"_attribute":{"translate":"name"}},"material":{"_value":{"name":"Material","required":"1"},"_attribute":{"translate":"name"}},"pattern":{"_value":{"name":"Pattern\/Graphic","required":"1"},"_attribute":{"translate":"name"}}},"ProductAds":{"adwords_grouping":{"_value":{"name":"Grouping","required":"0"},"_attribute":{"translate":"name"}},"adwords_labels":{"_value":{"name":"Labels","required":"0"},"_attribute":{"translate":"name"}},"adwords_redirect":{"_value":{"name":"Redirect","required":"0"},"_attribute":{"translate":"name"}},"adwords_queryparam":{"_value":{"name":"Query Param","required":"0"},"_attribute":{"translate":"name"}}}}`,
					},

					&config.Field{
						// Path: `google/googleshopping/attribute_groups`,
						ID:      "attribute_groups",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(config.IDScopeDefault), // @todo search for that
						Default: `{"price":{"sale_price":null,"tax":null,"sale_price_effective_date":null,"sale_price_effective_date_from":null,"sale_price_effective_date_to":null},"shipping_weight":{"weight":null},"title":{"name":null},"content":{"description":null}}`,
					},

					&config.Field{
						// Path: `google/googleshopping/base_attributes`,
						ID:      "base_attributes",
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Scope:   scope.NewPerm(config.IDScopeDefault), // @todo search for that
						Default: `{"id":null,"title":null,"link":null,"content":null,"price":null,"image_link":null,"condition":null,"target_country":null,"content_language":null,"destinations":null,"availability":null,"google_product_category":null,"product_type":null}`,
					},
				},
			},
		},
	},
)
