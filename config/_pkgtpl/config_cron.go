// +build ignore

package cron

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/store/scope"
)

// PackageConfiguration global configuration options for this package. Used in
// Frontend and Backend.
var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID: "system",
		Groups: config.NewGroupSlice(
			&config.Group{
				ID:        "cron",
				Label:     `Cron (Scheduled Tasks) - all the times are in minutes`,
				Comment:   element.LongText(`For correct URLs generated during cron runs please make sure that Web > Secure and Unsecure Base URLs are explicitly set.`),
				SortOrder: 15,
				Scope:     scope.NewPerm(scope.DefaultID),
				Fields:    config.NewFieldSlice(),
			},
		),
	},
)
