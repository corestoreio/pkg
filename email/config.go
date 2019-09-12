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

package email

import (
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
)

var PackageConfiguration = config.MustNewConfiguration(
	&config.Section{
		ID: "design",
		Groups: config.Groups{
			&config.Group{
				ID:        "email",
				Label:     `Transactional Emails`,
				Comment:   ``,
				SortOrder: 510,
				Scopes:    scope.PermStore,
				Fields: config.Fields{
					&config.Field{
						// Path: `design/email/logo`,
						ID:        "logo",
						Label:     `Logo Image`,
						Comment:   `Allowed file types: jpg, jpeg, gif, png. To optimize logo for high-resolution displays, upload an image that is 3x normal size and then specify 1x dimensions in width/height fields below.`,
						Type:      config.TypeImage,
						SortOrder: 10,
						Visible:   true,
						Scopes:    scope.PermStore,
						Default:   nil,
						// // BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Logo

					},
					&config.Field{
						// Path: `design/email/logo_alt`,
						ID:        "logo_alt",
						Label:     `Logo Image Alt`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   true,
						Scopes:    scope.PermStore,
					},
					&config.Field{
						// Path: `design/email/logo_width`,
						ID:        "logo_width",
						Label:     `Logo Width`,
						Comment:   `Only necessary if image has been uploaded above. Enter number of pixels, without appending "px".`,
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   true,
						Scopes:    scope.PermStore,
					},
					&config.Field{
						// Path: `design/email/logo_height`,
						ID:        "logo_height",
						Label:     `Logo Height`,
						Comment:   `Only necessary if image has been uploaded above. Enter number of pixels, without appending "px".`,
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   true,
						Scopes:    scope.PermStore,
					},
					&config.Field{
						// Path: `design/email/header_template`,
						ID:        "header_template",
						Label:     `Header Template`,
						Comment:   `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:      config.TypeSelect,
						SortOrder: 50,
						Visible:   true,
						Scopes:    scope.PermStore,
						Default:   `design_email_header_template`,
						// Magento\Config\Model\Config\Source\Email\Template
					},
					&config.Field{
						// Path: `design/email/footer_template`,
						ID:        "footer_template",
						Label:     `Footer Template`,
						Comment:   `Email template chosen based on theme fallback when "Default" option is selected.`,
						Type:      config.TypeSelect,
						SortOrder: 60,
						Visible:   true,
						Scopes:    scope.PermStore,
						Default:   `design_email_footer_template`,
						// Magento\Config\Model\Config\Source\Email\Template
					},
				},
			},
		},
	},
	&config.Section{
		ID:        "trans_email",
		Label:     "Store Email Addresses",
		SortOrder: 90,
		Scopes:    scope.PermStore,
		Groups: config.Groups{
			&config.Group{
				ID:        "ident_custom1",
				Label:     `Custom Email 1`,
				Comment:   ``,
				SortOrder: 4,
				Scopes:    scope.PermStore,
				Fields: config.Fields{
					&config.Field{
						// Path: `trans_email/ident_custom1/email`,
						ID:        "email",
						Label:     `Sender Email`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   true,
						Scopes:    scope.PermStore,
						//// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address => validation for correct mail address
					},

					&config.Field{
						// Path: `trans_email/ident_custom1/name`,
						ID:        "name",
						Label:     `Sender Name`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   true,
						Scopes:    scope.PermStore,
						//// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender => validation for a name to use only visible characters & is max 255 long
					},
				},
			},

			&config.Group{
				ID:        "ident_custom2",
				Label:     `Custom Email 2`,
				Comment:   ``,
				SortOrder: 5,
				Scopes:    scope.PermStore,
				Fields: config.Fields{
					&config.Field{
						// Path: `trans_email/ident_custom2/email`,
						ID:        "email",
						Label:     `Sender Email`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   true,
						Scopes:    scope.PermStore,
						//// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address
					},

					&config.Field{
						// Path: `trans_email/ident_custom2/name`,
						ID:        "name",
						Label:     `Sender Name`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   true,
						Scopes:    scope.PermStore,
						//// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender
					},
				},
			},

			&config.Group{
				ID:        "ident_general",
				Label:     `General Contact`,
				Comment:   ``,
				SortOrder: 1,
				Scopes:    scope.PermStore,
				Fields: config.Fields{
					&config.Field{
						// Path: `trans_email/ident_general/email`,
						ID:        "email",
						Label:     `Sender Email`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address

					},

					&config.Field{
						// Path: `trans_email/ident_general/name`,
						ID:        "name",
						Label:     `Sender Name`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender

					},
				},
			},

			&config.Group{
				ID:        "ident_sales",
				Label:     `Sales Representative`,
				Comment:   ``,
				SortOrder: 2,
				Scopes:    scope.PermStore,
				Fields: config.Fields{
					&config.Field{
						// Path: `trans_email/ident_sales/email`,
						ID:        "email",
						Label:     `Sender Email`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address

					},

					&config.Field{
						// Path: `trans_email/ident_sales/name`,
						ID:        "name",
						Label:     `Sender Name`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender

					},
				},
			},

			&config.Group{
				ID:        "ident_support",
				Label:     `Customer Support`,
				Comment:   ``,
				SortOrder: 3,
				Scopes:    scope.PermStore,
				Fields: config.Fields{
					&config.Field{
						// Path: `trans_email/ident_support/email`,
						ID:        "email",
						Label:     `Sender Email`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address

					},

					&config.Field{
						// Path: `trans_email/ident_support/name`,
						ID:        "name",
						Label:     `Sender Name`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Sender

					},
				},
			},
		},
	},

	&config.Section{
		ID: "system",
		Groups: config.Groups{
			&config.Group{
				ID:        "smtp",
				Label:     `Mail Sending Settings`,
				Comment:   ``,
				SortOrder: 20,
				Scopes:    scope.PermStore,
				Fields: config.Fields{
					&config.Field{
						// Path: `system/smtp/disable`,
						ID:        "disable",
						Label:     `Disable Email Communications. Output will be logged if disabled.`,
						Comment:   ``,
						Type:      config.TypeSelect,
						SortOrder: 10,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil,
						// Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `system/smtp/host`,
						ID:        "host",
						Label:     `Host`,
						Comment:   `SMTP Host`,
						Type:      config.TypeText,
						SortOrder: 20,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil,

					},

					&config.Field{
						// Path: `system/smtp/port`,
						ID:        "port",
						Label:     `Port (25)`,
						Comment:   `SMTP Port`,
						Type:      config.TypeText,
						SortOrder: 30,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil,

					},

					&config.Field{ // CS feature, not available in Magento
						// Path: `system/smtp/username`,
						ID:        "username",
						Label:     `Username`,
						Comment:   `SMTP Username`,
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil,

					},

					&config.Field{ // CS feature, not available in Magento
						// Path: `system/smtp/password`,
						ID:        "password",
						Label:     `Password`,
						Comment:   `SMTP Passowrd`,
						Type:      config.TypeText,
						SortOrder: 40,
						Visible:   true,
						Scopes:    scope.PermStore,

						// BackendModel: nil, // @todo encryption
						// @todo encryption
					},

					&config.Field{
						// Path: `system/smtp/set_return_path`,
						ID:        "set_return_path",
						Label:     `Set Return-Path`,
						Comment:   ``,
						Type:      config.TypeSelect,
						SortOrder: 70,
						Visible:   true,
						Scopes:    scope.PermDefault,

						// BackendModel: nil,
						// Magento\Config\Model\Config\Source\Yesnocustom
					},

					&config.Field{
						// Path: `system/smtp/return_path_email`,
						ID:        "return_path_email",
						Label:     `Return-Path Email`,
						Comment:   ``,
						Type:      config.TypeText,
						SortOrder: 80,
						Visible:   true,
						Scopes:    scope.PermDefault,

						// BackendModel: nil, // Magento\Config\Model\Config\Backend\Email\Address

					},
				},
			},
		},
	},
)
