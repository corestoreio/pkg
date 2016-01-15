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
	"github.com/juju/errgo"
)

// Levels defines how many parts are at least in a path.
// Like a/b/c for 3 parts. And 5 for a fully qualified path.
const Levels int = 3

// Separator used in the database table core_config_data and in config.Service
// to separate the path parts.
var Separator = []byte("/")

const sSeparator = "/"
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

// FQ returns the fully qualified route. Safe for further processing of the
// returned byte slice. If scope is equal to scope.DefaultID and ID is not
// zero then ID gets set to zero.
func (p Path) FQ() (Route, error) {
	if err := p.IsValid(); err != nil {
		return nil, err
	}

	if p.Scope == scope.DefaultID && p.ID > 0 {
		p.ID = 0
	}

	var buf bytes.Buffer
	if _, err := buf.WriteString(p.StrScope()); err != nil {
		return nil, errgo.Mask(err)
	}
	if err := buf.WriteByte(rSeparator); err != nil {
		return nil, errgo.Mask(err)
	}
	bufRaw := buf.Bytes()
	bufRaw = strconv.AppendInt(bufRaw, p.ID, 10)
	buf.Reset()
	if _, err := buf.Write(bufRaw); err != nil {
		return nil, errgo.Mask(err)
	}
	if err := buf.WriteByte(rSeparator); err != nil {
		return nil, errgo.Mask(err)
	}
	if _, err := buf.Write(p.Route); err != nil {
		return nil, errgo.Mask(err)
	}
	return buf.Bytes(), nil
}

// Level joins a configuration path parts by the path separator PS.
// The level argument defines the depth of the path parts to join.
// Level 1 will return the first part like "a", Level 2 returns "a/b"
// Level 3 returns "a/b/c" and so on. Level -1 joins all available path parts.
// Does not generate a fully qualified path.
func (p Path) Level(level int) (Route, error) {
	if err := p.IsValid(); err != nil {
		return nil, err
	}

	lp := len(p.Route)
	if level < 0 || level >= lp {
		return p.Route.Copy(), nil
	}

	if level == 0 {
		return Route(``), nil
	}

	pos := 0
	i := 1
	for pos <= len(p.Route) {
		sc := bytes.IndexRune(p.Route[pos:], rSeparator)
		if sc == -1 {
			break
		}
		pos += sc + 1

		if i == level {
			return p.Route[:pos-1].Copy(), nil
		}
		i++
	}
	return p.Route.Copy(), nil
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
		Route: Route(fqPath[fi+1:]),
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
	if p.Route.IsEmpty() {
		return ErrRouteEmpty
	}

	if false == p.Route.Valid() {
		return ErrRouteInvalidBytes
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
