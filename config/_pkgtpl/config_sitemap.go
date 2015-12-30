// +build ignore

package sitemap

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package.
// Used in frontend and backend. See init() for details.
var PackageConfiguration element.SectionSlice

func init() {
	PackageConfiguration = element.MustNewConfiguration(
		&element.Section{
			ID:        "sitemap",
			Label:     `XML Sitemap`,
			SortOrder: 70,
			Scope:     scope.PermAll,
			Resource:  0, // Otnegam_Sitemap::config_sitemap
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "category",
					Label:     `Categories Options`,
					SortOrder: 1,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sitemap/category/changefreq
							ID:        "changefreq",
							Label:     `Frequency`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `daily`,
							// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
						},

						&element.Field{
							// Path: sitemap/category/priority
							ID:        "priority",
							Label:     `Priority`,
							Comment:   element.LongText(`Valid values range from 0.0 to 1.0.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   0.5,
							// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
						},
					),
				},

				&element.Group{
					ID:        "product",
					Label:     `Products Options`,
					SortOrder: 2,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sitemap/product/changefreq
							ID:        "changefreq",
							Label:     `Frequency`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `daily`,
							// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
						},

						&element.Field{
							// Path: sitemap/product/priority
							ID:        "priority",
							Label:     `Priority`,
							Comment:   element.LongText(`Valid values range from 0.0 to 1.0.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   1,
							// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
						},

						&element.Field{
							// Path: sitemap/product/image_include
							ID:        "image_include",
							Label:     `Add Images into Sitemap`,
							Type:      element.TypeSelect,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `all`,
							// SourceModel: Otnegam\Sitemap\Model\Source\Product\Image\IncludeImage
						},
					),
				},

				&element.Group{
					ID:        "page",
					Label:     `CMS Pages Options`,
					SortOrder: 3,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sitemap/page/changefreq
							ID:        "changefreq",
							Label:     `Frequency`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   `daily`,
							// SourceModel: Otnegam\Sitemap\Model\Config\Source\Frequency
						},

						&element.Field{
							// Path: sitemap/page/priority
							ID:        "priority",
							Label:     `Priority`,
							Comment:   element.LongText(`Valid values range from 0.0 to 1.0.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   0.25,
							// BackendModel: Otnegam\Sitemap\Model\Config\Backend\Priority
						},
					),
				},

				&element.Group{
					ID:        "generate",
					Label:     `Generation Settings`,
					SortOrder: 4,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sitemap/generate/enabled
							ID:        "enabled",
							Label:     `Enabled`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},

						&element.Field{
							// Path: sitemap/generate/error_email
							ID:        "error_email",
							Label:     `Error Email Recipient`,
							Type:      element.TypeText,
							SortOrder: 5,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},

						&element.Field{
							// Path: sitemap/generate/error_email_identity
							ID:        "error_email_identity",
							Label:     `Error Email Sender`,
							Type:      element.TypeSelect,
							SortOrder: 6,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `general`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
						},

						&element.Field{
							// Path: sitemap/generate/error_email_template
							ID:        "error_email_template",
							Label:     `Error Email Template`,
							Comment:   element.LongText(`Email template chosen based on theme fallback when "Default" option is selected.`),
							Type:      element.TypeSelect,
							SortOrder: 7,
							Visible:   element.VisibleYes,
							Scope:     scope.NewPerm(scope.DefaultID, scope.WebsiteID),
							Default:   `sitemap_generate_error_email_template`,
							// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
						},

						&element.Field{
							// Path: sitemap/generate/frequency
							ID:        "frequency",
							Label:     `Frequency`,
							Type:      element.TypeSelect,
							SortOrder: 4,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							// BackendModel: Otnegam\Cron\Model\Config\Backend\Sitemap
							// SourceModel: Otnegam\Cron\Model\Config\Source\Frequency
						},

						&element.Field{
							// Path: sitemap/generate/time
							ID:        "time",
							Label:     `Start Time`,
							Type:      element.TypeTime,
							SortOrder: 3,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
						},
					),
				},

				&element.Group{
					ID:        "limit",
					Label:     `Sitemap File Limits`,
					SortOrder: 5,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sitemap/limit/max_lines
							ID:        "max_lines",
							Label:     `Maximum No of URLs Per File`,
							Type:      element.TypeText,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   50000,
						},

						&element.Field{
							// Path: sitemap/limit/max_file_size
							ID:        "max_file_size",
							Label:     `Maximum File Size`,
							Comment:   element.LongText(`File size in bytes.`),
							Type:      element.TypeText,
							SortOrder: 2,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   10485760,
						},
					),
				},

				&element.Group{
					ID:        "search_engines",
					Label:     `Search Engine Submission Settings`,
					SortOrder: 6,
					Scope:     scope.PermAll,
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sitemap/search_engines/submission_robots
							ID:        "submission_robots",
							Label:     `Enable Submission to Robots.txt`,
							Type:      element.TypeSelect,
							SortOrder: 1,
							Visible:   element.VisibleYes,
							Scope:     scope.PermAll,
							Default:   false,
							// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
						},
					),
				},
			),
		},

		// Hidden Configuration, may be visible somewhere else ...
		&element.Section{
			ID: "sitemap",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID: "generate",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sitemap/generate/error_email
							ID:      `error_email`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
						},
					),
				},

				&element.Group{
					ID: "file",
					Fields: element.NewFieldSlice(
						&element.Field{
							// Path: sitemap/file/valid_paths
							ID:      `valid_paths`,
							Type:    element.TypeHidden,
							Visible: element.VisibleNo,
							Default: `{"available":{"any_path":"\/*\/*.xml"}}`,
						},
					),
				},
			),
		},
	)
}
