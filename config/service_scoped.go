package config

import "time"

// ScopedGetter is equal to Getter but the underlying implementation takes
// care of providing the correct scope: default, website or store and bubbling
// up the scope chain from store -> website -> default.
//
// This interface is mainly implemented in the store package. The functions
// should be the same as in Getter but only the different is the paths
// argument. A path can be either one string containing a valid path like a/b/c
// or it can consists of 3 path parts like "a", "b", "c". All other arguments
// are invalid, should log an error if debug is enabled and return the empty type value.
type ScopedGetter interface {
	String(paths ...string) string
	Bool(paths ...string) bool
	Float64(paths ...string) float64
	Int(paths ...string) int
	DateTime(paths ...string) time.Time
}

type scopedService struct {
	root      Getter
	websiteID int64
	groupID   int64
	storeID   int64
}

var _ ScopedGetter = (*scopedService)(nil)

func newScopedService(r Getter, websiteID, groupID, storeID int64) *scopedService {
	return &scopedService{
		root:      r,
		websiteID: websiteID,
		groupID:   groupID,
		storeID:   storeID,
	}
}

// TODO some internal refactoring to keep it DRY ;-)

// String returns a string. Enable debug logging to see possible errors.
func (ss *scopedService) String(paths ...string) string {
	argP := Path(paths...)
	if ss.storeID > 0 {
		if v, err := ss.root.String(ScopeStore(ss.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.String.store.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.groupID > 0 {
		if v, err := ss.root.String(ScopeGroup(ss.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.String.group.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.websiteID > 0 {
		if v, err := ss.root.String(ScopeWebsite(ss.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.String.website.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	v, err := ss.root.String(argP)
	if err != nil && err != ErrKeyNotFound && PkgLog.IsDebug() {
		PkgLog.Debug("config.scopedService.String.default.error", "err", err, "scopedService", ss, "paths", paths)
	}
	return v
}

// Bool returns a bool value. Enable debug logging to see possible errors.
func (ss *scopedService) Bool(paths ...string) bool {
	argP := Path(paths...)
	if ss.storeID > 0 {
		if v, err := ss.root.Bool(ScopeStore(ss.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.Bool.store.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.groupID > 0 {
		if v, err := ss.root.Bool(ScopeGroup(ss.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.Bool.group.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.websiteID > 0 {
		if v, err := ss.root.Bool(ScopeWebsite(ss.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.Bool.website.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	v, err := ss.root.Bool(argP)
	if err != nil && err != ErrKeyNotFound && PkgLog.IsDebug() {
		PkgLog.Debug("config.scopedService.Bool.default.error", "err", err, "scopedService", ss, "paths", paths)
	}
	return v
}

// Float64 returns a float number. Enable debug logging for possible errors.
func (ss *scopedService) Float64(paths ...string) float64 {
	argP := Path(paths...)
	if ss.storeID > 0 {
		if v, err := ss.root.Float64(ScopeStore(ss.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.Float64.store.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.groupID > 0 {
		if v, err := ss.root.Float64(ScopeGroup(ss.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.Float64.group.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.websiteID > 0 {
		if v, err := ss.root.Float64(ScopeWebsite(ss.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.Float64.website.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	v, err := ss.root.Float64(argP)
	if err != nil && err != ErrKeyNotFound && PkgLog.IsDebug() {
		PkgLog.Debug("config.scopedService.Float64.default.error", "err", err, "scopedService", ss, "paths", paths)
	}
	return v
}

// Int returns an int. Enable debug logging for possible errors.
func (ss *scopedService) Int(paths ...string) int {
	argP := Path(paths...)
	if ss.storeID > 0 {
		if v, err := ss.root.Int(ScopeStore(ss.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.Int.store.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.groupID > 0 {
		if v, err := ss.root.Int(ScopeGroup(ss.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.Int.group.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.websiteID > 0 {
		if v, err := ss.root.Int(ScopeWebsite(ss.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.Int.website.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	v, err := ss.root.Int(argP)
	if err != nil && err != ErrKeyNotFound && PkgLog.IsDebug() {
		PkgLog.Debug("config.scopedService.Int.default.error", "err", err, "scopedService", ss, "paths", paths)
	}
	return v
}

// DateTime returns a time. Enable debug logging for possible errors.
func (ss *scopedService) DateTime(paths ...string) time.Time {
	argP := Path(paths...)
	if ss.storeID > 0 {
		if v, err := ss.root.DateTime(ScopeStore(ss.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.DateTime.store.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.groupID > 0 {
		if v, err := ss.root.DateTime(ScopeGroup(ss.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.DateTime.group.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	if ss.websiteID > 0 {
		if v, err := ss.root.DateTime(ScopeWebsite(ss.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && PkgLog.IsDebug() {
			PkgLog.Debug("config.scopedService.DateTime.website.error", "err", err, "scopedService", ss, "paths", paths)
		}
	}
	v, err := ss.root.DateTime(argP)
	if err != nil && err != ErrKeyNotFound && PkgLog.IsDebug() {
		PkgLog.Debug("config.scopedService.DateTime.default.error", "err", err, "scopedService", ss, "paths", paths)
	}
	return v
}
