package config

import (
	"time"

	"github.com/corestoreio/csfw/utils/log"
)

// ScopedReader is equal to Reader but the underlying implementation takes
// care of providing the correct scope: default, website or store and bubbling
// up the scope chain from store -> website -> default.
//
// This interface is mainly implemented in the store package. The functions
// should be the same as in Reader but only the different is the paths
// argument. A path can be either one string containing a valid path like a/b/c
// or it can consists of 3 path parts like "a", "b", "c". All other arguments
// are invalid, should log an error if debug is enabled and return the empty type value.
type ScopedReader interface {
	GetString(paths ...string) string
	GetBool(paths ...string) bool
	GetFloat64(paths ...string) float64
	GetInt(paths ...string) int
	GetDateTime(paths ...string) time.Time
}

type scopedManager struct {
	root      Reader
	websiteID int64
	groupID   int64
	storeID   int64
}

var _ ScopedReader = (*scopedManager)(nil)

func newScopedManager(r Reader, websiteID, groupID, storeID int64) *scopedManager {
	return &scopedManager{
		root:      r,
		websiteID: websiteID,
		groupID:   groupID,
		storeID:   storeID,
	}
}

// TODO some internal refactoring to avoid repetition ;-)

func (m *scopedManager) GetString(paths ...string) string {
	argP := Path(paths...)
	if m.storeID > 0 {
		if v, err := m.root.GetString(ScopeStore(m.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetString.store.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.groupID > 0 {
		if v, err := m.root.GetString(ScopeGroup(m.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetString.group.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.websiteID > 0 {
		if v, err := m.root.GetString(ScopeWebsite(m.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetString.website.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	v, err := m.root.GetString(argP)
	if err != nil && err != ErrKeyNotFound && log.IsDebug() {
		log.Debug("config.scopedManager.GetString.default.error", "err", err, "scopedManager", m, "paths", paths)
	}
	return v
}

func (m *scopedManager) GetBool(paths ...string) bool {
	argP := Path(paths...)
	if m.storeID > 0 {
		if v, err := m.root.GetBool(ScopeStore(m.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetBool.store.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.groupID > 0 {
		if v, err := m.root.GetBool(ScopeGroup(m.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetBool.group.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.websiteID > 0 {
		if v, err := m.root.GetBool(ScopeWebsite(m.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetBool.website.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	v, err := m.root.GetBool(argP)
	if err != nil && err != ErrKeyNotFound && log.IsDebug() {
		log.Debug("config.scopedManager.GetBool.default.error", "err", err, "scopedManager", m, "paths", paths)
	}
	return v
}

func (m *scopedManager) GetFloat64(paths ...string) float64 {
	argP := Path(paths...)
	if m.storeID > 0 {
		if v, err := m.root.GetFloat64(ScopeStore(m.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetFloat64.store.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.groupID > 0 {
		if v, err := m.root.GetFloat64(ScopeGroup(m.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetFloat64.group.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.websiteID > 0 {
		if v, err := m.root.GetFloat64(ScopeWebsite(m.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetFloat64.website.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	v, err := m.root.GetFloat64(argP)
	if err != nil && err != ErrKeyNotFound && log.IsDebug() {
		log.Debug("config.scopedManager.GetFloat64.default.error", "err", err, "scopedManager", m, "paths", paths)
	}
	return v
}

func (m *scopedManager) GetInt(paths ...string) int {
	argP := Path(paths...)
	if m.storeID > 0 {
		if v, err := m.root.GetInt(ScopeStore(m.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetInt.store.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.groupID > 0 {
		if v, err := m.root.GetInt(ScopeGroup(m.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetInt.group.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.websiteID > 0 {
		if v, err := m.root.GetInt(ScopeWebsite(m.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetInt.website.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	v, err := m.root.GetInt(argP)
	if err != nil && err != ErrKeyNotFound && log.IsDebug() {
		log.Debug("config.scopedManager.GetInt.default.error", "err", err, "scopedManager", m, "paths", paths)
	}
	return v
}

func (m *scopedManager) GetDateTime(paths ...string) time.Time {
	argP := Path(paths...)
	if m.storeID > 0 {
		if v, err := m.root.GetDateTime(ScopeStore(m.storeID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetDateTime.store.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.groupID > 0 {
		if v, err := m.root.GetDateTime(ScopeGroup(m.groupID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetDateTime.group.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	if m.websiteID > 0 {
		if v, err := m.root.GetDateTime(ScopeWebsite(m.websiteID), argP); err == nil {
			return v
		} else if err != ErrKeyNotFound && log.IsDebug() {
			log.Debug("config.scopedManager.GetDateTime.website.error", "err", err, "scopedManager", m, "paths", paths)
		}
	}
	v, err := m.root.GetDateTime(argP)
	if err != nil && err != ErrKeyNotFound && log.IsDebug() {
		log.Debug("config.scopedManager.GetDateTime.default.error", "err", err, "scopedManager", m, "paths", paths)
	}
	return v
}
