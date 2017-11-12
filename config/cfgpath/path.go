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

package cfgpath

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/errors"
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

// Path represents a configuration path bound to a scope.
type Path struct {
	Route
	// ScopeID a path is bound to this Scope and its ID
	ScopeID scope.TypeID
	// RouteLevelValid allows to bypass validation of separators in a Route
	// in cases where only a partial Route has been provided.
	RouteLevelValid bool
	// routeValidated internal flag to avoid running twice the route valid process
	routeValidated bool
}

// New creates a new validated Path. Scope is assigned to Default.
func New(rs ...Route) (Path, error) {
	p := Path{
		ScopeID: scope.DefaultTypeID,
	}
	if len(rs) == 1 {
		p.Route = rs[0]
	} else {
		p.Route.Append(rs...)
	}
	if p.Route.Sum32 == 0 {
		p.Route.Sum32 = p.Route.Hash32()
	}
	if err := p.IsValid(); err != nil {
		return Path{}, errors.Wrapf(err, "[cfgpath] Route %q", p.Route)
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

// Bind binds a path to a new scope with its scope ID. Group Scope is not
// supported and falls back to default. Fluent API design.
func (p Path) Bind(s scope.TypeID) Path {
	p.ScopeID = s
	return p
}

// BindWebsite binds a path to a website scope and its ID. Convenience helper
// function. Fluent API design.
func (p Path) BindWebsite(id int64) Path {
	p.ScopeID = scope.MakeTypeID(scope.Website, id)
	return p
}

// BindStore binds a path to a store scope and its ID. Convenience helper
// function. Fluent API design.
func (p Path) BindStore(id int64) Path {
	p.ScopeID = scope.MakeTypeID(scope.Store, id)
	return p
}

// Clone returns a new allocated Path with copied data. Clone is not needed if
// you before or after the assignment to a new variable use Path.Append() or the
// classic append(Path.Chars,[]byte() ....) to also allocate a new slice.
func (p Path) Clone() Path {
	p2 := p
	p2.Route = p.Route.Clone()
	return p2
}

// String returns a fully qualified path. Errors get logged if debug mode
// is enabled. String starts with `[cfgpath] Error:` on error.
// Error behaviour: NotValid, Empty or WriteFailed
func (p Path) String() string {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	if err := p.fq(buf); err != nil {
		return fmt.Sprintf("[cfgpath] Error: %+v", err)
	}
	return buf.String()
}

// GoString returns the internal representation of Path
func (p Path) GoString() string {
	return fmt.Sprintf("cfgpath.Path{ Route:cfgpath.NewRoute(`%s`), ScopeHash: %d }", p.Route, p.ScopeID)
}

// FQ returns the fully qualified route. Safe for further processing of the
// returned byte slice. If scope is equal to scope.DefaultID and ID is not
// zero then ID gets set to zero.
// Error behaviour: NotValid, Empty or WriteFailed
func (p Path) FQ() (Route, error) {
	// bufPool not possible because we're returning bytes, which can be modified
	// and bufPool truncates the slice, so return would a zero slice.
	var buf bytes.Buffer
	err := p.fq(&buf)
	return newRoute(buf.Bytes()), err
}

// Level returns a hierarchical based route depending on the depth.
// The depth argument defines the depth of levels to be returned.
// Depth 1 will return the first part like "a", Depth 2 returns "a/b"
// Depth 3 returns "a/b/c" and so on. Level -1 gives you all available levels.
// Does not generate a fully qualified path.
// The returned Route slice is owned by Path.Route. For further modifications you must
// copy it via Route.Copy().
// Error behaviour: NotValid or Empty
func (p Path) Level(depth int) (_ Route, err error) {
	p.routeValidated = true
	if err = p.IsValid(); err != nil {
		return
	}
	return p.Route.Level(depth)
}

// Hash same as Level() but returns a fnv32a value or an error if the route is
// invalid. The hash value contains the scope, scopeID and route.
// The returned Hash is equal to FQ().Hash32, but this function has less allocs.
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Hash implements FNV-1 and FNV-1a, non-cryptographic hash functions
// created by Glenn Fowler, Landon Curt Noll, and Phong Vo.
// See
// http://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function.
//
// Error behaviour: NotValid, Empty or WriteFailed
func (p Path) Hash(depth int) (uint32, error) {
	p.routeValidated = true

	buf := bufPool.Get()
	defer bufPool.Put(buf)
	err := p.fq(buf)
	if err != nil {
		return 0, err
	}
	r := newRoute(buf.Bytes())
	if depth < 0 {
		return r.Sum32, nil
	}
	return r.Hash(depth + 2)
}

// Error behaviour: NotValid or WriteFailed
func (p Path) fq(buf *bytes.Buffer) error {
	if err := p.IsValid(); err != nil {
		return err
	}

	scp, id := p.ScopeID.Unpack()
	if scp != scope.Website && scp != scope.Store {
		scp = scope.Default
		id = 0
	}

	if _, err := buf.Write(scp.StrBytes()); err != nil {
		return errors.NewWriteFailed(err, "[cfgpath] buf.Write")
	}
	if err := buf.WriteByte(Separator); err != nil {
		return errors.NewWriteFailed(err, "[cfgpath] buf.Write")
	}
	bufRaw := buf.Bytes()
	bufRaw = strconv.AppendInt(bufRaw, id, 10)
	buf.Reset()
	if _, err := buf.Write(bufRaw); err != nil {
		return errors.NewWriteFailed(err, "[cfgpath] buf.Write")
	}
	if err := buf.WriteByte(Separator); err != nil {
		return errors.NewWriteFailed(err, "[cfgpath] buf.Write")
	}
	if _, err := buf.Write(p.Route.Chars); err != nil {
		return errors.NewWriteFailed(err, "[cfgpath] buf.Write")
	}
	return nil
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
		return Path{}, errors.NewNotValidf("[cfgpath] Incorrect fully qualified path: %q. Expecting: strScope/ID/%s", fqPath, fqPath)
	}

	fi := strings.Index(fqPath, sSeparator)
	scopeStr := fqPath[:fi]

	if false == scope.Valid(scopeStr) {
		return Path{}, errors.NewNotSupportedf("[cfgpath] Unknown Scope: %q", scopeStr)
	}

	fqPath = fqPath[fi+1:]
	fi = strings.Index(fqPath, sSeparator)
	scopeID, err := strconv.ParseInt(fqPath[:fi], 10, 64)

	return Path{
		Route:   NewRoute(fqPath[fi+1:]),
		ScopeID: scope.MakeTypeID(scope.FromString(scopeStr), scopeID),
	}, errors.NewNotValid(err, "[cfgpath] ParseInt")
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
// Error behaviour: NotValid or Empty
func (p Path) IsValid() error {
	if !p.routeValidated {
		// only validate the route when it has not yet been done
		if err := p.Route.Validate(); err != nil {
			return errors.Wrap(err, "[cfgpath] Route.Validate")
		}
	}
	if p.RouteLevelValid {
		return nil
	}
	if p.Route.Separators() < Levels-1 || p.Route.RuneCount() < 8 /*aa/bb/cc*/ {
		return errors.NewNotValidf(errIncorrectPathTpl, p.Route.String())
	}
	return nil
}
