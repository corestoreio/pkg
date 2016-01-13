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

package path

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"unicode/utf8"

	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/bufferpool"
)

// Levels defines how many parts are at least in a path.
// Like a/b/c for 3 parts. And 5 for a fully qualified path.
const Levels int = 3

// Separator used in the database table core_config_data and in config.Service
// to separate the path parts.
var Separator = []byte("/")

const sSeparator = "s"
const rSeparator = '/'

// ErrRouteEmpty path parts are empty
var ErrRouteEmpty = errors.New("Route is empty")

// ErrIncorrectPath a path is missing a path separator or is too short
var ErrIncorrectPath = errors.New("Incorrect Path. Either to short or missing path separator.")

// Path represents a configuration path bound to a scope.
type Path struct {
	Route
	Scope scope.Scope
	// ID represents a website, group or store ID
	ID int64
}

// New creates a new validated Path. Scope is assigned to Default.
func New(r Route) (Path, error) {
	p := Path{
		Route: r,
		Scope: scope.DefaultID,
	}
	if err := p.IsValid(); err != nil {
		return Path{}, err
	}
	return p, nil
}

//// NewSplit takes a path argument like a/b/c or path parts like "a","b","c".
//// If a path has been provided it gets split into its parts.
//// Scope is assigned to Default.
//func NewSplit(paths ...string) (Path, error) {
//	p := Path{
//		Scope: scope.DefaultID,
//	}
//	switch {
//	case len(paths) >= Levels:
//		p.Route = paths
//	case len(paths) == 1 && paths[0] != "":
//		p.Route = Split(paths[0])
//	default:
//		return Path{}, fmt.Errorf("Incorrect number of paths elements: want %d, have %d, Path: %v", Levels, len(paths), paths)
//	}
//
//	if err := p.IsValid(); err != nil {
//		return Path{}, err
//	}
//	return p, nil
//}

// MustNew same as New but panics on error.
func MustNew(r Route) Path {
	p, err := New(r)
	if err != nil {
		panic(err)
	}
	return p
}

// BindStr binds a path to a new scope with its scope ID.
// The scope gets extracted from the StrScope.
func (p Path) BindStr(s scope.StrScope, id int64) Path {
	p.Scope = s.Scope()
	p.ID = id
	return p
}

// Bind binds a path to a new scope with its scope ID.
// Group Scope is not supported and falls back to default.
func (p Path) Bind(s scope.Scope, id int64) Path {
	p.Scope = s
	p.ID = id
	return p
}

// StrScope wrapper function. Converts the Path.Scope to a StrScope.
func (p Path) StrScope() string {
	return scope.FromScope(p.Scope).String()
}

// String returns a fully qualified path. Errors get logged if debug mode
// is enabled.
func (p Path) String() string {
	s, err := p.FQ()
	if PkgLog.IsDebug() {
		PkgLog.Debug("path.Path.FQ.String", "err", err, "path", p)
	}
	return string(s)
}

// FQ returns the fully qualified path.
func (p Path) FQ() ([]byte, error) {
	if err := p.IsValid(); err != nil {
		return nil, err
	}

	buf := bufferpool.Get()
	defer bufferpool.Put(buf)

	idStr := "0"
	if p.ID > 0 {
		if p.ID <= int64CacheLen {
			idStr = int64Cache[p.ID]
		} else {
			idStr = strconv.FormatInt(p.ID, 10)
		}
	} else {
		buf.WriteRune('0')
	}

	scopeStr := scope.FromScope(p.Scope)

	buf.WriteString(scopeStr.String())
	buf.WriteByte(rSeparator)
	buf.WriteString(idStr)
	buf.WriteByte(rSeparator)
	buf.Write(p.Route)
	return buf.Bytes(), nil // can be a bug because of defer and put back
}

// this "cache" should cover ~80% of all store setups
var int64Cache = [...]string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
}
var int64CacheLen = int64(len(int64Cache))

// Level joins a configuration path parts by the path separator PS.
// The level argument defines the depth of the path parts to join.
// Level 1 will return the first part like "a", Level 2 returns "a/b"
// Level 3 returns "a/b/c" and so on. Level -1 joins all available path parts.
// Does not generate a fully qualified path.
func (p Path) Level(level int) Route {
	lp := len(p.Route)
	if level <= 0 || level >= lp {
		level = lp
	}

	sc := bytes.Count(p.Route, Separator)
	idx := 0
	for i := 0; i < sc; i++ {
		idx += bytes.IndexRune(p.Route[idx:], rSeparator)
		if level <= i {
			return Route(p.Route[idx:])
		}
	}
	return p.Route
}

// SplitFQPath takes a fully qualified path and splits it into its parts.
// 	Input: stores/5/catalog/frontend/list_allow_all
//	=>
//		scope: 		stores
//		scopeID: 	5
//		path: 		catalog/frontend/list_allow_all
// Zero allocations to memory. Err may contain an ErrUnsupportedScope or
// failed to parse a string into an int64 or invalid fqPath.
func SplitFQ(fqPath string) (Path, error) {
	if false == isFQ(fqPath) {
		return Path{}, fmt.Errorf("Incorrect fully qualified path: %q", fqPath)
	}

	fi := strings.Index(fqPath, sSeparator)
	scopeStr := fqPath[:fi]

	if false == scope.Valid(scopeStr) {
		return Path{}, scope.ErrUnsupportedScope
	}

	fqPath = fqPath[fi+1:]

	fi = strings.Index(fqPath, sSeparator)
	scopeID, err := strconv.ParseInt(fqPath[:fi], 10, 64)
	path := fqPath[fi+1:]
	return Path{
		Route: Route(path),
		Scope: scope.FromString(scopeStr),
		ID:    scopeID,
	}, err
}

func isFQ(fqPath string) bool {
	return strings.Count(fqPath, sSeparator) >= Levels+1 // like stores/1/a/b/c
}

// IsValid checks for valid configuration path. Returns nil on success.
// Configuration path attribute can have only three groups of [a-zA-Z0-9_] characters split by '/'.
// Minimal length per part 2 characters. Case sensitive.
//
// IsValid can return ErrRouteEmpty or ErrIncorrectPath or a custom error.
func (p Path) IsValid() error {
	if p.Route.IsEmpty() {
		return ErrRouteEmpty
	}

	if false == utf8.Valid(p.Route) {
		return ErrRouteInvalidBytes
	}

	if p.Scope == scope.DefaultID && p.ID != 0 {
		return fmt.Errorf("Default Scope can only have an ID of 0. Yours is: %d", p.ID)
	}

	var sepCount, length int
	i := 0
	for i < len(p.Route) {
		var r rune
		if p.Route[i] < utf8.RuneSelf {
			r = rune(p.Route[i])
			i++
		} else {
			dr, _ := utf8.DecodeRune(p.Route[i:])
			return fmt.Errorf("This character %q is not allowed in Route %s", string(dr), p.Route)
		}
		ok := false
		switch {
		case '0' <= r && r <= '9':
			ok = true
		case 'a' <= r && r <= 'z':
			ok = true
		case 'A' <= r && r <= 'Z':
			ok = true
		case r == '_':
			ok = true
		case r == rSeparator:
			sepCount++
			ok = true
		}
		if !ok {
			return fmt.Errorf("This character %q is not allowed in Route %s", string(r), p.Route)
		}
		length++
	}
	if sepCount != Levels-1 || length < 8 {
		return ErrIncorrectPath
	}

	return nil
}
