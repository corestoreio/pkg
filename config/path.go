// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package config

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/bufferpool"
)

// Levels defines how many parts are at least in a path.
// Like a/b/c for 3 parts. And 5 for a fully qualified path.
const Levels = 3

// realLevels maximum numbers of supported separators. used as array initializer.
const maxLevels = 8 // up to 8. just a guess

// Separator used in the database table core_config_data and in config.Service
// to separate the path parts.
const Separator byte = '/'

const sSeparator = "/"

const errIncorrectPathTpl = "[config] Invalid Path %q. Either to short or missing path separator."

const errIncorrectPositionTpl = "[config] Position '%d' does not exists"

const errRouteInvalidBytesTpl = "[config] Route contains invalid bytes %q which are not runes."

var errRouteEmpty = errors.Empty.Newf("[config] Route is empty")

// Path represents a configuration path bound to a scope.
type Path struct {
	route string
	// ScopeID a route is bound to this Scope and its ID.
	ScopeID scope.TypeID
	// routeValidated internal flag to avoid running twice the route valid process
	// TODO this flag makes only then sense when field `Route` is private and there a public functions to set/modify it.
	routeValidated bool
}

func ParsePathBytes(path []byte) (Path, error) {
	return Path{}, errors.NotImplemented.Newf("TODO")
}

// MakePath creates a new validated Path. Scope is assigned to Default.
func MakePath(route string) (Path, error) {
	p := Path{
		route:   route,
		ScopeID: scope.DefaultTypeID,
	}
	if err := p.IsValid(); err != nil {
		return Path{}, errors.Wrapf(err, "[config] Route %q", p.route)
	}
	return p, nil
}

// MustMakePath same as MakePath but panics on error.
func MustMakePath(route string) Path {
	p, err := MakePath(route)
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

// String returns a fully qualified path. Errors get logged if debug mode
// is enabled. String starts with `[config] Error:` on error.
// Error behaviour: NotValid, Empty or WriteFailed
func (p Path) String() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := p.fq(buf); err != nil {
		return fmt.Sprintf("[config] Error: %+v", err)
	}
	return buf.String()
}

// GoString returns the internal representation of Path
func (p Path) GoString() string {
	return fmt.Sprintf("cfgpath.Path{ Route:cfgpath.MakeRoute(%q), ScopeHash: %d }", p.route, p.ScopeID)
}

// FQ returns the fully qualified route. Safe for further processing of the
// returned byte slice. If scope is equal to scope.DefaultID and ID is not
// zero then ID gets set to zero.
// Error behaviour: NotValid, Empty or WriteFailed
func (p Path) FQ() (string, error) {
	// bufPool not possible because we're returning bytes, which can be modified
	// and bufPool truncates the slice, so return would a zero slice.
	var buf bytes.Buffer
	if err := p.fq(&buf); err != nil {
		return "", errors.Wrapf(err, "[config] Scope %d Path %q", p.ScopeID, p.route)
	}
	return buf.String(), nil
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

	buf.Write(scp.StrBytes())
	buf.WriteByte(Separator)

	bufRaw := buf.Bytes()
	bufRaw = strconv.AppendInt(bufRaw, id, 10)
	buf.Reset()
	buf.Write(bufRaw)
	buf.WriteByte(Separator)
	buf.WriteString(p.route)
	return nil
}

// SplitFQ takes a fully qualified path and splits it into its parts.
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
	if !(strings.Count(fqPath, sSeparator) >= Levels+1) {
		return Path{}, errors.NotValid.Newf("[config] Incorrect fully qualified path: %q. Expecting: strScope/ID/%s", fqPath, fqPath)
	}

	fi := strings.Index(fqPath, sSeparator)
	scopeStr := fqPath[:fi]

	if !scope.Valid(scopeStr) {
		return Path{}, errors.NotSupported.Newf("[config] Unknown Scope: %q", scopeStr)
	}

	fqPath = fqPath[fi+1:]
	fi = strings.Index(fqPath, sSeparator)
	scopeID, err := strconv.ParseInt(fqPath[:fi], 10, 64)

	return Path{
		route:   fqPath[fi+1:],
		ScopeID: scope.MakeTypeID(scope.FromString(scopeStr), scopeID),
	}, errors.NotValid.New(err, "[config] ParseInt")
}

// BenchmarkSplitFQ-4  	 2000000	       761 ns/op	      32 B/op	       1 allocs/op
// slower than the string version above. this commented out will be kept for historical
// reasons. maybe some one can speed it more up than the above string version.
//
// ErrInvalidScopeID when parsing the scope ID fails.
// var ErrInvalidScopeID = errors.Make("Scope ID contains invalid bytes. Cannot extract an integer value.")
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
	seps := p.Separators()
	if !p.routeValidated {
		if "" == p.route {
			return errRouteEmpty
		}

		if seps == len(p.route) {
			return errors.NotValid.Newf(errIncorrectPathTpl, p.route)
		}

		if !utf8.ValidString(p.route) {
			return errors.NotValid.Newf(errRouteInvalidBytesTpl, p.route)
		}
		for _, rn := range p.route {
			switch {
			case rn == rune(Separator), rn == '_':
				// ok
			case unicode.IsDigit(rn), unicode.IsLetter(rn), unicode.IsNumber(rn):
				// ok
			default:
				return errors.NotValid.Newf("[config] Route %q contains invalid character %q.", p.route, rn)
			}
		}

	}
	if seps < Levels-1 || utf8.RuneCountInString(p.route) < 8 /*aa/bb/cc*/ {
		return errors.NotValid.Newf(errIncorrectPathTpl, p.route)
	}
	return nil
}

func (p Path) IsEmpty() bool {
	return p.route == ""
}

// Equal compares the Sum32 of both routes
func (p Path) Equal(b Path) bool {
	return p.ScopeID == b.ScopeID && p.route == b.route
}

func (p Path) EqualRoute(b Path) bool {
	return p.route == b.route
}

// Append adds other partial routes with a Separator between. After the partial
// routes have been added a validation check will NOT be done.
// Internally it creates a new byte slice.
//
//		a := cfgpath.Route(`catalog/product`)
//		b := cfgpath.Route(`enable_flat_tables`)
//		if err := a.Append(b); err != nil {
//			panic(err)
//		}
//		println(a.String())
//		// Should print: catalog/product/enable_flat_tables
func (p Path) Append(routes ...string) Path {

	if i1, i2 := strings.LastIndexByte(p.route, Separator), len(p.route)-1; i1 > 0 && i2 > 0 && i1 == i2 {
		p.route = p.route[:len(p.route)-1] // strip last Separator
	}

	var buf strings.Builder
	buf.WriteString(p.route)

	i := 0
	for _, r := range routes {
		if r != "" {
			if len(p.route) > 0 && len(r) > 0 && r[0] != Separator {
				buf.WriteByte(Separator)
			}
			buf.WriteString(r)
			i++
		}
	}
	return Path{
		route: buf.String(),
	}
}

// MarshalText implements interface encoding.TextMarshaler.
func (p Path) MarshalText() (text []byte, err error) {
	return []byte(p.route), nil
}

// UnmarshalText transforms the text into a route with performed validation
// checks. Implements encoding.TextUnmarshaler.
// Error behaviour: NotValid, Empty.
func (p *Path) UnmarshalText(txt []byte) error {
	p.route = string(txt)
	return errors.WithStack(p.IsValid())
}

// Level returns a hierarchical based route depending on the depth. The depth
// argument defines the depth of levels to be returned. Depth 1 will return the
// first part like "a", Depth 2 returns "a/b" Depth 3 returns "a/b/c" and so on.
// Level -1 gives you all available levels. Does generate a fully qualified
// path.
func (p Path) Level(depth int) (string, error) {
	if err := p.IsValid(); err != nil {
		return "", errors.WithStack(err)
	}

	fq, err := p.FQ()
	if err != nil {
		return "", errors.WithStack(err)
	}

	lp := len(fq)
	switch {
	case depth < 0:
		return fq, nil
	case depth == 0:
		return "", nil
	case depth >= lp:
		return fq, nil
	}

	pos := 0
	i := 1
	for pos <= lp {
		sc := strings.IndexByte(fq[pos:], Separator)
		if sc == -1 {
			return fq, nil
		}
		pos += sc + 1

		if i == depth {
			return fq[:pos-1], nil
		}
		i++
	}
	return fq, nil
}

const (
	offset32 = 2166136261
	prime32  = 16777619
)

// Hash same as Level() but returns a fnv32a value or an error if the route is
// invalid. 32 has been chosen because routes consume not that much space as
// a []byte.
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
// Usage as map key.
func (p Path) Hash(depth int) (uint32, error) {
	r2, err := p.Level(depth)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	var hash uint32 = offset32
	for _, c := range r2 {
		hash ^= uint32(c)
		hash *= prime32
	}
	return hash, nil
}

// Hash32 returns a fnv32a value of the route. 32 has been chosen because
// routes consume not that much space as a []byte.
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
// Usage as map key.
func (p Path) Hash32() uint32 {
	fq, _ := p.FQ()
	if fq == "" {
		return 0
	}
	var hash uint32 = offset32
	for _, c := range fq {
		hash ^= uint32(c)
		hash *= prime32
	}
	return hash
}

// Separators returns the number of separators
func (p Path) Separators() (count int) {
	for _, b := range p.route {
		if b == rune(Separator) {
			count++
		}
	}
	return
}

// ScopeRoute returns the assigned scope and its ID and the route.
func (p Path) ScopeRoute() (scope.TypeID, string) {
	return p.ScopeID, p.route
}

// Part returns the route part on the desired position. The Route gets validated
// before extracting the part. Does not run Validate()
//		Have Route: general/single_store_mode/enabled
//		Pos<1 => ErrIncorrectPosition
//		Pos=1 => general
//		Pos=2 => single_store_mode
//		Pos=3 => enabled
//		Pos>3 => ErrIncorrectPosition
// The returned Route slice is owned by Path. For further modifications you must
// copy it via Route.Copy().
func (p Path) Part(pos int) (string, error) {
	p.routeValidated = true
	if err := p.IsValid(); err != nil {
		return "", err
	}

	if pos < 1 {
		return "", errors.NotValid.Newf(errIncorrectPositionTpl, pos)
	}

	sepCount := p.Separators()
	if sepCount < 1 { // no separator found
		return p.route, nil
	}
	if pos > sepCount+1 {
		return "", errors.NotValid.Newf(errIncorrectPositionTpl, pos)
	}

	var sepPos [maxLevels]int
	sp := 0
	for i, b := range p.route {
		if b == rune(Separator) && sp < maxLevels {
			sepPos[sp] = i + 1 // positions of the separators in the slice
			sp++
		}
	}

	pos--
	min := 0
	for i := 0; i < maxLevels; i++ {
		if sepPos[i] == 0 { // no more separators found
			return p.route[min:], nil
		}
		max := sepPos[i]
		if i == pos {
			return p.route[min : max-1], nil
		}
		min = max
	}
	return p.route[min:], nil
}

// Split splits the route into its three parts and appends it to the slice
// `ret`. Does not run Validate() Example:
// 		routes := cfgpath.MakeRoute("aa/bb/cc")
//		rs, err := routes.Split()
//		rs[0].String() == "aa"
//		rs[1].String() == "bb"
//		rs[2].String() == "cc"
//
// Error behaviour: NotValid
func (p Path) Split(ret ...string) (_ []string, err error) {

	const sepCount = Levels - 1 // only two separators supported
	var sepPos [sepCount]int
	sp := 0
	for i, b := range p.route {
		if b == rune(Separator) && sp < sepCount {
			sepPos[sp] = i // positions of the separators in the slice
			sp++
		}
	}
	if sp < 1 {
		return nil, errors.NotValid.Newf(errIncorrectPathTpl, p.route)
	}

	min := 0
	for i := 0; i < Levels; i++ {
		var max int
		if i < sepCount && sepPos[i] > 0 {
			max = sepPos[i]
		} else {
			max = len(p.route)
		}
		ret = append(ret, p.route[min:max])
		if i < sepCount && sepPos[i] == 0 {
			return
		}
		min = max + 1
	}
	return ret, err
}

// PathSlice represents a collection of Paths
type PathSlice []Path

// add more functions if needed

// Contains return true if the Path p can be found within the slice.
// It must match ID, Scope and Route.
func (ps PathSlice) Contains(p Path) bool {
	for _, pps := range ps {
		if pps.ScopeID == p.ScopeID && pps.route == p.route {
			return true
		}
	}
	return false
}

func (ps PathSlice) Len() int { return len(ps) }

// Less sorts by scope, id and route
func (ps PathSlice) Less(i, j int) bool {
	p1 := ps[i]
	p2 := ps[j]
	if p1.ScopeID != p2.ScopeID && p1.route == p2.route {
		return p1.ScopeID < p2.ScopeID
	}
	return p1.route < p2.route
}

func (ps PathSlice) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

// Sort is a convenience method to sort stable.
func (ps PathSlice) Sort() { sort.Stable(ps) }
