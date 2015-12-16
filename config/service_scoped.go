package config

import (
	"github.com/corestoreio/csfw/config/scope"
	"time"
)

// ScopedGetter is equal to Getter but the underlying implementation takes
// care of providing the correct scope: default, website or store and bubbling
// up the scope chain from store -> website -> default.
//
// This interface is mainly implemented in the store package. The functions
// should be the same as in Getter but only the different is the paths
// argument. A path can be either one string containing a valid path like a/b/c
// or it can consists of 3 path parts like "a", "b", "c". All other arguments
// are invalid. Returned error is mostly of ErrKeyNotFound.
type ScopedGetter interface {
	// Scope tells you the current underlying scope and its website, group or store ID
	Scope() (scope.Scope, int64)
	String(paths ...string) (string, error)
	Bool(paths ...string) (bool, error)
	Float64(paths ...string) (float64, error)
	Int(paths ...string) (int, error)
	DateTime(paths ...string) (time.Time, error)
}

type scopedService struct {
	root      Getter
	websiteID int64
	groupID   int64
	storeID   int64
}

var _ ScopedGetter = (*scopedService)(nil)

func newScopedService(r Getter, websiteID, groupID, storeID int64) scopedService {
	return scopedService{
		root:      r,
		websiteID: websiteID,
		groupID:   groupID,
		storeID:   storeID,
	}
}

// Scope tells you the current underlying scope and its website, group or store ID
func (ss scopedService) Scope() (scope.Scope, int64) {
	switch {
	case ss.storeID > 0:
		return scope.StoreID, ss.storeID
	case ss.groupID > 0:
		return scope.GroupID, ss.groupID
	case ss.websiteID > 0:
		return scope.WebsiteID, ss.websiteID
	default:
		return scope.DefaultID, 0
	}
}

// String returns a string. Enable debug logging to see possible errors.
func (ss scopedService) String(paths ...string) (string, error) {
	return ss.root.String(Scope(ss.Scope()), Path(paths...))
}

// Bool returns a bool value. Enable debug logging to see possible errors.
func (ss scopedService) Bool(paths ...string) (bool, error) {
	return ss.root.Bool(Scope(ss.Scope()), Path(paths...))
}

// Float64 returns a float number. Enable debug logging for possible errors.
func (ss scopedService) Float64(paths ...string) (float64, error) {
	return ss.root.Float64(Scope(ss.Scope()), Path(paths...))
}

// Int returns an int. Enable debug logging for possible errors.
func (ss scopedService) Int(paths ...string) (int, error) {
	return ss.root.Int(Scope(ss.Scope()), Path(paths...))
}

// DateTime returns a time. Enable debug logging for possible errors.
func (ss scopedService) DateTime(paths ...string) (time.Time, error) {
	return ss.root.DateTime(Scope(ss.Scope()), Path(paths...))
}
