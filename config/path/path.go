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

	"github.com/corestoreio/csfw/store/scope"
	"github.com/juju/errgo"
)

// Levels defines how many parts are at least in a path.
// Like a/b/c for 3 parts. And 5 for a fully qualified path.
const Levels int = 3

// realLevels maximum numbers of supported separators. used as array initializer.
const maxLevels int = 8 // up to 8. just a guess

// Separator used in the database table core_config_data and in config.Service
// to separate the path parts.
const Separator byte = '/'

const sSeparator = "/"

var bSeparator = []byte(sSeparator)

// ErrRouteEmpty path parts are empty
var ErrRouteEmpty = errors.New("Route is empty")

// ErrIncorrectPath a path is missing a path separator or is too short
var ErrIncorrectPath = errors.New("Incorrect Path. Either to short or missing path separator.")

// ErrIncorrectPosition returned by function Part() whenever an invalid input
// position has been applied.
var ErrIncorrectPosition = errors.New("Position does not exists")

// Path represents a configuration path bound to a scope.
type Path struct {
	Route
	Scope scope.Scope
	// ID represents a website, group or store ID
	ID int64
	// RouteLevelValid allows to bypass validation of separators in a Route
	// in cases where only a partial Route has been provided.
	RouteLevelValid bool
	// routeValidated internal flag to avoid running twice the route valid process
	routeValidated bool
}

// New creates a new validated Path. Scope is assigned to Default.
func New(rs ...Route) (Path, error) {
	p := Path{
		Scope: scope.DefaultID,
	}
	if len(rs) == 1 {
		p.Route = rs[0]
	} else {
		p.Route.Append(rs...)
	}
	if err := p.IsValid(); err != nil {
		return Path{}, err
	}
	return p, nil
}

// MustNew same as New but panics on error.
func MustNew(rs ...Route) Path {
	p, err := New(rs...)
	if err != nil {
		panic(err)
	}
	return p
}

// NewByParts creates a new Path from path part strings.
// Parts gets merged via Separator.
//		p := NewByParts("catalog","product",")
func NewByParts(parts ...string) (Path, error) {
	return New(NewRoute(parts...))
}

// MustNewByParts same as NewByParts but panics on error.
func MustNewByParts(parts ...string) Path {
	p, err := NewByParts(parts...)
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
// is enabled. String is empty on error.
func (p Path) String() string {
	s, err := p.FQ()
	if PkgLog.IsDebug() {
		PkgLog.Debug("path.Path.FQ.String", "err", err, "path", p)
	}
	return s.String()
}

// GoString returns the internal representation of Path
func (p Path) GoString() string {
	return fmt.Sprintf("path.Path{ Route:path.NewRoute(`%s`), Scope: %d, ID: %d }", p.Route, p.Scope, p.ID)
}

// FQ returns the fully qualified route. Safe for further processing of the
// returned byte slice. If scope is equal to scope.DefaultID and ID is not
// zero then ID gets set to zero.
// The returned Route slice is owned by Path. For further modifications you must
// copy it via Route.Copy().
func (p Path) FQ() (Route, error) {
	if err := p.IsValid(); err != nil {
		return Route{}, err
	}

	if (p.Scope == scope.DefaultID || p.Scope == scope.GroupID) && p.ID > 0 {
		p.Scope = scope.DefaultID
		p.ID = 0
	}

	var buf bytes.Buffer
	if _, err := buf.WriteString(p.StrScope()); err != nil {
		return Route{}, errgo.Mask(err)
	}
	if err := buf.WriteByte(Separator); err != nil {
		return Route{}, errgo.Mask(err)
	}
	bufRaw := buf.Bytes()
	bufRaw = strconv.AppendInt(bufRaw, p.ID, 10)
	buf.Reset()
	if _, err := buf.Write(bufRaw); err != nil {
		return Route{}, errgo.Mask(err)
	}
	if err := buf.WriteByte(Separator); err != nil {
		return Route{}, errgo.Mask(err)
	}
	if _, err := buf.Write(p.Route.Chars); err != nil {
		return Route{}, errgo.Mask(err)
	}
	return newRoute(buf.Bytes()), nil
}

// Level joins a configuration path parts by the path separator PS.
// The level argument defines the depth of the path parts to join.
// Level 1 will return the first part like "a", Level 2 returns "a/b"
// Level 3 returns "a/b/c" and so on. Level -1 joins all available path parts.
// Does not generate a fully qualified path.
// The returned Route slice is owned by Path. For further modifications you must
// copy it via Route.Copy().
func (p Path) Level(level int) (r Route, err error) {
	p.routeValidated = true
	if err = p.IsValid(); err != nil {
		return
	}
	return p.Route.Level(level)
}

// Hash same as Level() but returns a fnv32a value or an error if the route is
// invalid.
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Hash implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// http://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function.
func (p Path) Hash(level int) (uint32, error) {
	p.routeValidated = true
	if err := p.IsValid(); err != nil {
		return 0, err
	}
	return p.Route.Hash(level)
}

// Part returns the route part on the desired position. The Route gets validated
// before extracting the part.
//		Have Route: general/single_store_mode/enabled
//		Pos<1 => ErrIncorrectPosition
//		Pos=1 => general
//		Pos=2 => single_store_mode
//		Pos=3 => enabled
//		Pos>3 => ErrIncorrectPosition
// The returned Route slice is owned by Path. For further modifications you must
// copy it via Route.Copy().
func (p Path) Part(pos int) (r Route, err error) {
	p.routeValidated = true
	if err = p.IsValid(); err != nil {
		return
	}
	return p.Route.Part(pos)
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
	// this is the most fast version I come up with.
	// moving from strings to bytes was even slower despite inline
	// th parse int64 function
	if false == (strings.Count(fqPath, sSeparator) >= Levels+1) {
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

	return Path{
		Route: Route{Chars: []byte(fqPath[fi+1:])},
		Scope: scope.FromString(scopeStr),
		ID:    scopeID,
	}, err
}

// BenchmarkSplitFQ-4  	 2000000	       761 ns/op	      32 B/op	       1 allocs/op
// slower than the string version above. this commented out will be kept for historical
// reasons. maybe some one can speed it more up than the above string version.
//
// ErrInvalidScopeID when parsing the scope ID fails.
// var ErrInvalidScopeID = errors.New("Scope ID contains invalid bytes. Cannot extract an integer value.")
//func SplitFQ(fqPath Route) (Path, error) {
//	if false == (bytes.Count(fqPath, Separator) >= Levels+1) || false == fqPath.Valid() {
//		return Path{}, fmt.Errorf("Incorrect fully qualified path: %q", fqPath)
//	}
//
//	fi := bytes.IndexRune(fqPath, rSeparator)
//	scopeBytes := fqPath[:fi]
//
//	if false == scope.ValidBytes(scopeBytes) {
//		return Path{}, scope.ErrUnsupportedScope
//	}
//
//	fqPath = fqPath[fi+1:]                   // remove scope string
//	fi = bytes.IndexRune(fqPath, rSeparator) // find scope id
//
//	scopeIDBytes := fqPath[:fi]
//	if len(scopeIDBytes) > 5 { // i have never seen more than 10k stores, websites or groups
//		return Path{}, ErrInvalidScopeID
//	}
//
//	const maxUint64 = (1<<64 - 1)
//	const cutoff = maxUint64/10 + 1
//	var n uint64
//	base := 10
//	for i := 0; i < len(scopeIDBytes); i++ {
//		var v byte
//		d := scopeIDBytes[i]
//		switch {
//		case '0' <= d && d <= '9':
//			v = d - '0'
//		case 'a' <= d && d <= 'z':
//			v = d - 'a' + 10
//		case 'A' <= d && d <= 'Z':
//			v = d - 'A' + 10
//		default:
//			n = 0
//			return Path{}, ErrInvalidScopeID
//		}
//		if v >= byte(base) {
//			n = 0
//			return Path{}, ErrInvalidScopeID
//		}
//
//		if n >= cutoff {
//			// n*base overflows
//			n = maxUint64
//			return Path{}, ErrInvalidScopeID
//		}
//		n *= uint64(base)
//
//		n1 := n + uint64(v)
//		if n1 < n || n1 > (1<<uint(64)-1) { // 64 bits
//			// n+v overflows
//			n = maxUint64
//			return Path{}, ErrInvalidScopeID
//		}
//		n = n1
//	}
//
//	return Path{
//		Route: Route(fqPath[fi+1:].Copy()),
//		Scope: scope.FromBytes(scopeBytes),
//		ID:    int64(n),
//	}, nil
//}

// IsValid checks for valid configuration path. Returns nil on success.
// Configuration path attribute can have only three groups of [a-zA-Z0-9_] characters split by '/'.
// Minimal length per part 2 characters. Case sensitive.
//
// IsValid can return ErrRouteEmpty or ErrIncorrectPath or a custom error.
func (p Path) IsValid() error {
	if !p.routeValidated {
		// only validate the route when it has not yet been done
		if err := p.Route.Validate(); err != nil {
			return err
		}
	}
	if p.RouteLevelValid {
		return nil
	}
	if p.Route.Separators() < Levels-1 || p.Route.RuneCount() < 8 /*aa/bb/cc*/ {
		return ErrIncorrectPath
	}
	return nil
}
