// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package store

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/configsource"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/corestoreio/csfw/storage/csdb"
)

const (
	// SingleStoreModeEnabled if true then single store mode enabled
	// This flag only shows that admin does not want to show certain
	// UI components in the backend (like store switchers etc)
	// If there is only one store view but it does not check the store view collection. WTF?
	PathSingleStoreModeEnabled = "general/single_store_mode/enabled"
	PathStoreStoreName         = "general/store_information/name"
	PathStoreStorePhone        = "general/store_information/phone"
	PathStoreCodeInURL         = "web/url/use_store"
	PathRedirectToBase         = "web/url/redirect_to_base"
	PathSecureInFrontend       = "web/secure/use_in_frontend"

	PathUnsecureBaseURL = "web/unsecure/base_url"
	PathSecureBaseURL   = "web/secure/base_url"

	//	PathSecureBaseLinkUrl   = "web/secure/base_link_url"
	//	PathUnsecureBaseLinkUrl = "web/unsecure/base_link_url"

	PathSecureBaseStaticURL   = "web/secure/base_static_url"
	PathUnsecureBaseStaticURL = "web/unsecure/base_static_url"

	PathSecureBaseMediaURL   = "web/secure/base_media_url"
	PathUnsecureBaseMediaURL = "web/unsecure/base_media_url"

	// This defines the base currency scope ("Currency Setup" > "Currency Options" > "Base Currency").
	// can be 0 = Global or 1 = Website
	PathPriceScope = "catalog/price/scope"

	PlaceholderBaseURL         = config.LeftDelim + "base_url" + config.RightDelim
	PlaceholderBaseURLSecure   = config.LeftDelim + "secure_base_url" + config.RightDelim
	PlaceholderBaseURLUnSecure = config.LeftDelim + "unsecure_base_url" + config.RightDelim
)

// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
var TableCollection csdb.Manager

// PackageConfiguration contains the main configuration
var PackageConfiguration config.SectionSlice

func init() {
	PackageConfiguration = config.NewConfiguration(
		&config.Section{
			ID:        "general",
			Label:     "General",
			SortOrder: 10,
			Scope:     scope.PermAll,

			Groups: config.GroupSlice{
				&config.Group{
					ID:        "single_store_mode",
					Label:     `Single-Store Mode`,
					Comment:   ``,
					SortOrder: 150,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `general/single_store_mode/enabled`,
							ID:           "enabled",
							Label:        `Enable Single-Store Mode`,
							Comment:      `This setting will not be taken into account if system has more than one store view.`,
							Type:         config.TypeSelect,
							SortOrder:    10,
							Visible:      config.VisibleYes,
							Scope:        scope.NewPerm(scope.DefaultID),
							Default:      nil,
							BackendModel: nil,
							SourceModel:  configsource.YesNo, // Magento\Config\Model\Config\Source\Yesno
						},
					},
				},

				&config.Group{
					ID:        "store_information",
					Label:     `Store Information`,
					Comment:   ``,
					SortOrder: 100,
					Scope:     scope.PermAll,
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `general/store_information/name`,
							ID:           "name",
							Label:        `Store Name`,
							Comment:      ``,
							Type:         config.TypeText,
							SortOrder:    10,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil,
							SourceModel:  nil,
						},

						&config.Field{
							// Path: `general/store_information/phone`,
							ID:           "phone",
							Label:        `Store Phone Number`,
							Comment:      ``,
							Type:         config.TypeText,
							SortOrder:    20,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil,
							SourceModel:  nil,
						},
					},
				},
			},
		},

		&config.Section{
			ID:        "web",
			Label:     "Web",
			SortOrder: 20,
			Scope:     scope.PermAll,
			Groups: config.GroupSlice{
				&config.Group{
					ID:        "url",
					Label:     `Url Options`,
					Comment:   ``,
					SortOrder: 3,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `web/url/use_store`,
							ID:           "use_store",
							Label:        `Add Store Code to Urls`,
							Comment:      `<strong style="color:red">Warning!</strong> When using Store Code in URLs, in some cases system may not work properly if URLs without Store Codes are specified in the third party services (e.g. PayPal etc.).`,
							Type:         config.TypeSelect,
							SortOrder:    10,
							Visible:      config.VisibleYes,
							Scope:        scope.NewPerm(scope.DefaultID),
							Default:      nil,
							BackendModel: nil,                // Magento\Config\Model\Config\Backend\Store
							SourceModel:  configsource.YesNo, // Magento\Config\Model\Config\Source\Yesno
						},

						&config.Field{
							// Path: `web/url/redirect_to_base`,
							ID:           "redirect_to_base",
							Label:        `Auto-redirect to Base URL`,
							Comment:      `I.e. redirect from http://example.com/store/ to http://www.example.com/store/`,
							Type:         config.TypeSelect,
							SortOrder:    20,
							Visible:      config.VisibleYes,
							Scope:        scope.NewPerm(scope.DefaultID),
							Default:      nil,
							BackendModel: nil,
							SourceModel:  configsource.Redirect, // Magento\Config\Model\Config\Source\Web\Redirect
						},
					},
				},

				&config.Group{
					ID:        "unsecure",
					Label:     `Base URLs`,
					Comment:   `Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`,
					SortOrder: 10,
					Scope:     scope.PermAll,
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `web/unsecure/base_url`,
							ID:           "base_url",
							Label:        `Base URL`,
							Comment:      `Specify URL or {{base_url}} placeholder.`,
							Type:         config.TypeText,
							SortOrder:    10,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							SourceModel:  nil,
						},

						&config.Field{
							// Path: `web/unsecure/base_link_url`,
							ID:           "base_link_url",
							Label:        `Base Link URL`,
							Comment:      `May start with {{unsecure_base_url}} placeholder.`,
							Type:         config.TypeText,
							SortOrder:    20,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							SourceModel:  nil,
						},

						&config.Field{
							// Path: `web/unsecure/base_static_url`,
							ID:           "base_static_url",
							Label:        `Base URL for Static View Files`,
							Comment:      `May be empty or start with {{unsecure_base_url}} placeholder.`,
							Type:         config.TypeText,
							SortOrder:    25,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							SourceModel:  nil,
						},

						&config.Field{
							// Path: `web/unsecure/base_media_url`,
							ID:           "base_media_url",
							Label:        `Base URL for User Media Files`,
							Comment:      `May be empty or start with {{unsecure_base_url}} placeholder.`,
							Type:         config.TypeText,
							SortOrder:    40,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							SourceModel:  nil,
						},
					},
				},

				&config.Group{
					ID:        "secure",
					Label:     `Base URLs (Secure)`,
					Comment:   `Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. https://example.com/magento/`,
					SortOrder: 20,
					Scope:     scope.PermAll,
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `web/secure/base_url`,
							ID:           "base_url",
							Label:        `Secure Base URL`,
							Comment:      `Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.`,
							Type:         config.TypeText,
							SortOrder:    10,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							SourceModel:  nil,
						},

						&config.Field{
							// Path: `web/secure/base_link_url`,
							ID:           "base_link_url",
							Label:        `Secure Base Link URL`,
							Comment:      `May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.`,
							Type:         config.TypeText,
							SortOrder:    20,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							SourceModel:  nil,
						},

						&config.Field{
							// Path: `web/secure/base_static_url`,
							ID:           "base_static_url",
							Label:        `Secure Base URL for Static View Files`,
							Comment:      `May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`,
							Type:         config.TypeText,
							SortOrder:    25,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							SourceModel:  nil,
						},

						&config.Field{
							// Path: `web/secure/base_media_url`,
							ID:           "base_media_url",
							Label:        `Secure Base URL for User Media Files`,
							Comment:      `May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`,
							Type:         config.TypeText,
							SortOrder:    40,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							SourceModel:  nil,
						},

						&config.Field{
							// Path: `web/secure/use_in_frontend`,
							ID:           "use_in_frontend",
							Label:        `Use Secure URLs on Storefront`,
							Comment:      `Enter https protocol to use Secure URLs on Storefront.`,
							Type:         config.TypeSelect,
							SortOrder:    50,
							Visible:      config.VisibleYes,
							Scope:        scope.PermAll,
							Default:      false,
							BackendModel: nil,                // Magento\Config\Model\Config\Backend\Secure
							SourceModel:  configsource.YesNo, // Magento\Config\Model\Config\Source\Yesno
						},

						&config.Field{
							// Path: `web/secure/use_in_adminhtml`,
							ID:           "use_in_adminhtml",
							Label:        `Use Secure URLs in Admin`,
							Comment:      `Enter https protocol to use Secure URLs in Admin.`,
							Type:         config.TypeSelect,
							SortOrder:    60,
							Visible:      config.VisibleYes,
							Scope:        scope.NewPerm(scope.DefaultID),
							Default:      false,
							BackendModel: nil,                // Magento\Config\Model\Config\Backend\Secure
							SourceModel:  configsource.YesNo, // Magento\Config\Model\Config\Source\Yesno
						},

						&config.Field{
							// Path: `web/secure/offloader_header`,
							ID:           "offloader_header",
							Label:        `Offloader header`,
							Comment:      ``,
							Type:         config.TypeText,
							SortOrder:    70,
							Visible:      config.VisibleYes,
							Scope:        scope.NewPerm(scope.DefaultID),
							Default:      "SSL_OFFLOADED",
							BackendModel: nil,
							SourceModel:  nil,
						},
					},
				},
			},
		},
		&config.Section{
			ID: "catalog",
			Groups: config.GroupSlice{
				&config.Group{
					ID: "price",
					Fields: config.FieldSlice{
						&config.Field{
							// Path: `catalog/price/scope`,
							ID:      "scope",
							Default: PriceScopeGlobal,
						},
					},
				},
			},
		},
	)
}
