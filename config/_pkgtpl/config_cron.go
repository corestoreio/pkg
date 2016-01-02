// +build ignore

package cron

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
			ID: "system",
			Groups: element.NewGroupSlice(
				&element.Group{
					ID:        "cron",
					Label:     `Cron (Scheduled Tasks) - all the times are in minutes`,
					Comment:   element.LongText(`For correct URLs generated during cron runs please make sure that Web > Secure and Unsecure Base URLs are explicitly set.`),
					SortOrder: 15,
					Scope:     scope.NewPerm(scope.DefaultID),
					Fields:    element.NewFieldSlice(),
				},
			),
		},
	)
	Path = NewPath(PackageConfiguration)
}
