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

package store

import (
	"net/http"

	"github.com/corestoreio/csfw/catalog/catconfig"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/configsource"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/storage/csdb"
	"github.com/corestoreio/csfw/store/scope"
)

// PathStoreCodeInURL if yes the ___store variable will be added to the URL. TODO.
//var PathStoreCodeInURL = model.NewBool("web/url/use_store", )

var PathUnsecureBaseURL = model.NewBaseURL("web/unsecure/base_url")

var PathSecureBaseURL = model.NewBaseURL("web/secure/base_url")

//	PathSecureBaseLinkUrl   = "web/secure/base_link_url"
//	PathUnsecureBaseLinkUrl = "web/unsecure/base_link_url"

var PathSecureBaseStaticURL = model.NewBaseURL("web/secure/base_static_url")
var PathUnsecureBaseStaticURL = model.NewBaseURL("web/unsecure/base_static_url")

var PathSecureBaseMediaURL = model.NewBaseURL("web/secure/base_media_url")
var PathUnsecureBaseMediaURL = model.NewBaseURL("web/unsecure/base_media_url")

// PathPriceScope defines the base currency scope
// ("Currency Setup" > "Currency Options" > "Base Currency").
// can be 0 = Global or 1 = Website
// See constants PriceScopeGlobal and PriceScopeWebsite.
var PathPriceScope = catconfig.NewConfigPriceScope("catalog/price/scope")

// Placeholder constants and their values can occur in the table core_config_data.
// These placeholder must be replaced with the current values.
const (
	PlaceholderBaseURL         = config.LeftDelim + "base_url" + config.RightDelim
	PlaceholderBaseURLSecure   = config.LeftDelim + "secure_base_url" + config.RightDelim
	PlaceholderBaseURLUnSecure = config.LeftDelim + "unsecure_base_url" + config.RightDelim
)

// TableCollection handles all tables and its columns. init() in generated Go file will set the value.
var TableCollection csdb.Manager

// PackageConfiguration contains the main configuration
var PackageConfiguration element.SectionSlice

func init() {
	PackageConfiguration = element.MustNewConfiguration(

		&element.Section{
			ID:        "web",
			Label:     "Web",
			SortOrder: 20,
			Scope:     scope.PermAll,
			Groups: element.GroupSlice{
				&element.Group{
					ID:        "url",
					Label:     `Url Options`,
					Comment:   ``,
					SortOrder: 3,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.FieldSlice{
						&element.Field{
							// Path: `web/url/use_store`,
							ID:           "use_store",
							Label:        `Add Store Code to Urls`,
							Comment:      `<strong style="color:red">Warning!</strong> When using Store Code in URLs, in some cases system may not work properly if URLs without Store Codes are specified in the third party services (e.g. PayPal etc.).`,
							Type:         element.TypeSelect,
							SortOrder:    10,
							Visible:      element.VisibleYes,
							Scope:        scope.NewPerm(scope.DefaultID),
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Store
							// SourceModel:  configsource.YesNo, // Magento\Config\Model\Config\Source\Yesno
						},
					},
				},

				&element.Group{
					ID:        "unsecure",
					Label:     `Base URLs`,
					Comment:   `Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. http://example.com/magento/`,
					SortOrder: 10,
					Scope:     scope.PermAll,
					Fields: element.FieldSlice{
						&element.Field{
							// Path: `web/unsecure/base_url`,
							ID:           "base_url",
							Label:        `Base URL`,
							Comment:      `Specify URL or {{base_url}} placeholder.`,
							Type:         element.TypeText,
							SortOrder:    10,
							Visible:      element.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							// SourceModel:  nil,
						},

						&element.Field{
							// Path: `web/unsecure/base_link_url`,
							ID:           "base_link_url",
							Label:        `Base Link URL`,
							Comment:      `May start with {{unsecure_base_url}} placeholder.`,
							Type:         element.TypeText,
							SortOrder:    20,
							Visible:      element.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							// SourceModel:  nil,
						},

						&element.Field{
							// Path: `web/unsecure/base_static_url`,
							ID:           "base_static_url",
							Label:        `Base URL for Static View Files`,
							Comment:      `May be empty or start with {{unsecure_base_url}} placeholder.`,
							Type:         element.TypeText,
							SortOrder:    25,
							Visible:      element.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							// SourceModel:  nil,
						},

						&element.Field{
							// Path: `web/unsecure/base_media_url`,
							ID:           "base_media_url",
							Label:        `Base URL for User Media Files`,
							Comment:      `May be empty or start with {{unsecure_base_url}} placeholder.`,
							Type:         element.TypeText,
							SortOrder:    40,
							Visible:      element.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							// SourceModel:  nil,
						},
					},
				},

				&element.Group{
					ID:        "secure",
					Label:     `Base URLs (Secure)`,
					Comment:   `Any of the fields allow fully qualified URLs that end with '/' (slash) e.g. https://example.com/magento/`,
					SortOrder: 20,
					Scope:     scope.PermAll,
					Fields: element.FieldSlice{
						&element.Field{
							// Path: `web/secure/base_url`,
							ID:           "base_url",
							Label:        `Secure Base URL`,
							Comment:      `Specify URL or {{base_url}}, or {{unsecure_base_url}} placeholder.`,
							Type:         element.TypeText,
							SortOrder:    10,
							Visible:      element.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							// SourceModel:  nil,
						},

						&element.Field{
							// Path: `web/secure/base_link_url`,
							ID:           "base_link_url",
							Label:        `Secure Base Link URL`,
							Comment:      `May start with {{secure_base_url}} or {{unsecure_base_url}} placeholder.`,
							Type:         element.TypeText,
							SortOrder:    20,
							Visible:      element.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							// SourceModel:  nil,
						},

						&element.Field{
							// Path: `web/secure/base_static_url`,
							ID:           "base_static_url",
							Label:        `Secure Base URL for Static View Files`,
							Comment:      `May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`,
							Type:         element.TypeText,
							SortOrder:    25,
							Visible:      element.VisibleYes,
							Scope:        scope.PermAll,
							Default:      nil,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							// SourceModel:  nil,
						},

						&element.Field{
							// Path: `web/secure/base_media_url`,
							ID:        "base_media_url",
							Label:     `Secure Base URL for User Media Files`,
							Comment:   `May be empty or start with {{secure_base_url}}, or {{unsecure_base_url}} placeholder.`,
							Type:      element.TypeText,
							SortOrder: 40,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   nil,
							//BackendModel: nil, // Magento\Config\Model\Config\Backend\Baseurl
							//// SourceModel:  nil,
						},

						&element.Field{
							// Path: `web/secure/use_in_adminhtml`,
							ID:           "use_in_adminhtml",
							Label:        `Use Secure URLs in Admin`,
							Comment:      `Enter https protocol to use Secure URLs in Admin.`,
							Type:         element.TypeSelect,
							SortOrder:    60,
							Visible:      element.VisibleYes,
							Scope:        scope.NewPerm(scope.DefaultID),
							Default:      false,
							BackendModel: nil, // Magento\Config\Model\Config\Backend\Secure
							// SourceModel:  configsource.YesNo, // Magento\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: `web/secure/offloader_header`,
							ID:           "offloader_header",
							Label:        `Offloader header`,
							Comment:      ``,
							Type:         element.TypeText,
							SortOrder:    70,
							Visible:      element.VisibleYes,
							Scope:        scope.NewPerm(scope.DefaultID),
							Default:      "SSL_OFFLOADED",
							BackendModel: nil,
							// SourceModel:  nil,
						},
					},
				},
			},
		},
		&element.Section{
			ID: "catalog",
			Groups: element.GroupSlice{
				&element.Group{
					ID: "price",
					Fields: element.FieldSlice{
						&element.Field{
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
