// +build ignore

package sitemap

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "sitemap",
		Label:     `XML Sitemap`,
		SortOrder: 70,
		Scope:     scope.PermAll,
		Resource:  0, // Otnegam_Sitemap::config_sitemap
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "category",
				Label:     `Categories Options`,
				SortOrder: 1,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sitemap/category/changefreq
						ID:        "changefreq",
						Label:     `Frequency`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `daily`,
						// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: sitemap/category/priority
						ID:        "priority",
						Label:     `Priority`,
						Comment:   element.LongText(`Valid values range from 0.0 to 1.0.`),
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   0.5,
						// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
					},
				),
			},

			&config.Group{
				ID:        "product",
				Label:     `Products Options`,
				SortOrder: 2,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sitemap/product/changefreq
						ID:        "changefreq",
						Label:     `Frequency`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `daily`,
						// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: sitemap/product/priority
						ID:        "priority",
						Label:     `Priority`,
						Comment:   element.LongText(`Valid values range from 0.0 to 1.0.`),
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   1,
						// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
					},

					&config.Field{
						// Path: sitemap/product/image_include
						ID:        "image_include",
						Label:     `Add Images into Sitemap`,
						Type:      config.TypeSelect,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `all`,
						// SourceModel: Otnegam\Sitemap\Model\Source\Product\Image\IncludeImage
					},
				),
			},

			&config.Group{
				ID:        "page",
				Label:     `CMS Pages Options`,
				SortOrder: 3,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sitemap/page/changefreq
						ID:        "changefreq",
						Label:     `Frequency`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   `daily`,
						// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: sitemap/page/priority
						ID:        "priority",
						Label:     `Priority`,
						Comment:   element.LongText(`Valid values range from 0.0 to 1.0.`),
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   0.25,
						// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
					},
				),
			},

			&config.Group{
				ID:        "generate",
				Label:     `Generation Settings`,
				SortOrder: 4,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sitemap/generate/enabled
						ID:        "enabled",
						Label:     `Enabled`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: sitemap/generate/error_email
						ID:        "error_email",
						Label:     `Error Email Recipient`,
						Type:      config.TypeText,
						SortOrder: 5,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},

					&config.Field{
						// Path: sitemap/generate/error_email_identity
						ID:        "error_email_identity",
						Label:     `Error Email Sender`,
						Type:      config.TypeSelect,
						SortOrder: 6,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `general`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
					},

					&config.Field{
						// Path: sitemap/generate/error_email_template
						ID:        "error_email_template",
						Label:     `Error Email Template`,
						Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
						Type:      config.TypeSelect,
						SortOrder: 7,
						Visible:   config.VisibleYes,
						Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
						Default:   `sitemap_generate_error_email_template`,
						// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: sitemap/generate/frequency
						ID:        "frequency",
						Label:     `Frequency`,
						Type:      config.TypeSelect,
						SortOrder: 4,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						// BackendModel: Otnegam\Cron\Model\Config\Backend\Sitemap
						// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
					},

					&config.Field{
						// Path: sitemap/generate/time
						ID:        "time",
						Label:     `Start Time`,
						Type:      config.TypeTime,
						SortOrder: 3,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
					},
				),
			},

			&config.Group{
				ID:        "limit",
				Label:     `Sitemap File Limits`,
				SortOrder: 5,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sitemap/limit/max_lines
						ID:        "max_lines",
						Label:     `Maximum No of URLs Per File`,
						Type:      config.TypeText,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   50000,
					},

					&config.Field{
						// Path: sitemap/limit/max_file_size
						ID:        "max_file_size",
						Label:     `Maximum File Size`,
						Comment:   element.LongText(`File size in bytes.`),
						Type:      config.TypeText,
						SortOrder: 2,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   10485760,
					},
				),
			},

			&config.Group{
				ID:        "search_engines",
				Label:     `Search Engine Submission Settings`,
				SortOrder: 6,
				Scope:     scope.PermAll,
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sitemap/search_engines/submission_robots
						ID:        "submission_robots",
						Label:     `Enable Submission to Robots.txt`,
						Type:      config.TypeSelect,
						SortOrder: 1,
						Visible:   config.VisibleYes,
						Scope:     scope.PermAll,
						Default:   false,
						// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
					},
				),
			},
		),
	},

	// Hidden Configuration, may be visible somewhere else ...
	&config.Section{
		ID: "sitemap",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID: "generate",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sitemap/generate/error_email
						ID:      `error_email`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
					},
				),
			},

			&config.Group{
				ID: "file",
				Fields: config.NewFieldSlice(
					&config.Field{
						// Path: sitemap/file/valid_paths
						ID:      `valid_paths`,
						Type:    config.TypeHidden,
						Visible: config.VisibleNo,
						Default: `{"available":{"any_path":"\/*\/*.xml"}}`,
					},
				),
			},
		),
	},
)
