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
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/bufferpool"
	"github.com/corestoreio/pkg/util/byteconv"
	"github.com/minio/highwayhash"
)

// PathLevels defines how many parts are at least in a path.
// Like a/b/c for 3 parts. And 5 for a fully qualified path.
const PathLevels = 3

// maxLevels maximum numbers of supported separators. used as array initializer.
const maxLevels = 8 // up to 8. just a guess

// PathSeparator used in the database table core_config_data and in config.Service
// to separate the path parts.
const PathSeparator = '/'

var (
	bSeparator    = []byte(sPathSeparator)
	errRouteEmpty = errors.Empty.Newf("[config] Route is empty")
)

const (
	sPathSeparator          = "/"
	errIncorrectPathTpl     = "[config] Invalid Path %q. Either to short or missing path separator."
	errIncorrectPositionTpl = "[config] Position '%d' does not exists"
	errRouteInvalidBytesTpl = "[config] Route contains invalid bytes %q which are not runes."
)

// Path represents a configuration path bound to a scope. A Path is not safe for
// concurrent use. Bind* method receivers always return a new copy of a path.
type Path struct {
	route string
	// ScopeID a route is bound to this Scope and its ID.
	ScopeID scope.TypeID
	// routeValidated internal flag to avoid running twice the route valid
	// process.
	routeValidated bool
	// UseEnvSuffix if enabled the Path has environment awareness. An
	// environment comes from the *Service type and defines for example
	// PRODUCTION, TEST, CI or STAGING. The final path will use this prefix to
	// distinguish between the environments. Environment awareness should be
	// added to Paths for e.g. payment credentials or order export access data.
	UseEnvSuffix bool
	// envSuffix gets set by the *Service type if the service runs environment
	// aware and a Path has set UseEnvSuffix to true.
	envSuffix string
	//TTL            time.Duration // for later use when the LRU cache becomes time sensitive, like the one from karlseguin
}

// NewPathWithScope creates a new validate Path with a custom scope.
func NewPathWithScope(scp scope.TypeID, route string) (*Path, error) {
	p := &Path{
		route:   route,
		ScopeID: scp,
	}
	if err := p.IsValid(); err != nil {
		return nil, errors.Wrapf(err, "[config] Route %q", p.route)
	}
	return p, nil
}

// MustNewPathWithScope creates a new validate Path with a custom scope but
// panics on error. E.g. invalid path.
func MustNewPathWithScope(scp scope.TypeID, route string) *Path {
	p, err := NewPathWithScope(scp, route)
	if err != nil {
		panic(err)
	}
	return p
}

// NewPath makes a new validated Path. Scope is assigned to Default.
func NewPath(route string) (*Path, error) {
	return NewPathWithScope(scope.DefaultTypeID, route)
}

// MustNewPath same as NewPath but panics on error.
func MustNewPath(route string) *Path {
	return MustNewPathWithScope(scope.DefaultTypeID, route)
}

// Bind binds a path to a new scope with its scope ID. Returns a new Path
// pointer and does not apply the changes to the current Path. Fluent API
// design.
func (p Path) Bind(s scope.TypeID) *Path {
	p.ScopeID = s
	return &p
}

// BindWebsite binds a path to a website scope and its ID. Returns a new Path
// pointer and does not apply the changes to the current Path. Convenience
// helper function. Fluent API design.
func (p Path) BindWebsite(id int64) *Path {
	p.ScopeID = scope.MakeTypeID(scope.Website, id)
	return &p
}

// BindStore binds a path to a store scope and its ID. Returns a new Path
// pointer and does not apply the changes to the current Path. Convenience
// helper function. Fluent API design.
func (p Path) BindStore(id int64) *Path {
	p.ScopeID = scope.MakeTypeID(scope.Store, id)
	return &p
}

// BindDefault binds a path to the default scope. Returns a new Path pointer and
// does not apply the changes to the current Path. Convenience helper function.
// Fluent API design.
func (p Path) BindDefault() *Path {
	p.ScopeID = scope.DefaultTypeID
	return &p
}

// WithEnvSuffix enables that this Path has environment awareness. An
// environment comes from the *Service type and defines for example PRODUCTION,
// TEST, CI or STAGING. The final path will use this prefix to distinguish
// between the environments. Environment awareness should be added to Paths for
// e.g. payment credentials or order export access data.
func (p *Path) WithEnvSuffix() *Path {
	p.UseEnvSuffix = true
	return p
}

func (p *Path) writeEnvSuffix(buf *bytes.Buffer) {
	if p.UseEnvSuffix && p.envSuffix != "" {
		buf.WriteByte(PathSeparator)
		buf.WriteString(p.envSuffix)
	}
}

func (p *Path) stripEnvSuffixStr(r string) string {
	if p.envSuffix != "" && strings.HasSuffix(r, p.envSuffix) {
		r = r[:len(r)-len(p.envSuffix)-1] // 1 == PathSeparator length
	}
	return r
}

func (p *Path) stripEnvSuffixByte(r []byte) []byte {
	if p.envSuffix != "" && bytes.HasSuffix(r, []byte(p.envSuffix)) {
		r = r[:len(r)-len(p.envSuffix)-1] // 1 == PathSeparator length
	}
	return r
}

// String returns a fully qualified path. Errors get logged if debug mode
// is enabled. String starts with `[config] Error:` on error.
// Error behaviour: NotValid, Empty or WriteFailed
func (p *Path) String() string {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := p.AppendFQ(buf); err != nil {
		return fmt.Sprintf("[config] Error: %+v", err)
	}
	return buf.String()
}

// FQ returns the fully qualified route. Safe for further processing of the
// returned byte slice. If scope is equal to scope.DefaultID and ID is not
// zero then ID gets set to zero.
// Error behaviour: NotValid, Empty or WriteFailed
func (p *Path) FQ() (string, error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	if err := p.AppendFQ(buf); err != nil {
		return "", errors.Wrapf(err, "[config] Scope %d Path %q", p.ScopeID, p.route)
	}
	return buf.String(), nil
}

// AppendFQ validates the Path and appends it to the buffer.
func (p *Path) AppendFQ(buf *bytes.Buffer) error {
	if err := p.IsValid(); err != nil {
		return err
	}

	scp, id := p.ScopeID.Unpack()
	if scp != scope.Website && scp != scope.Store {
		scp = scope.Default
		id = 0
	}

	buf.Write(scp.StrBytes())
	buf.WriteByte(PathSeparator)

	bufRaw := buf.Bytes()
	bufRaw = strconv.AppendInt(bufRaw, id, 10)
	buf.Reset()
	buf.Write(bufRaw)
	buf.WriteByte(PathSeparator)
	buf.WriteString(p.route)
	p.writeEnvSuffix(buf)
	return nil
}

func isDigitOnly(str string) bool {
	for _, r := range str {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Parse takes a route or a fully qualified path and splits it into its parts
// with final validation. Input: stores/5/catalog/frontend/list_allow_all or
// just catalog/frontend/list_allow_all to use default scope.
//	=>
//		scope: 		stores
//		scopeID: 	5
//		route: 		catalog/frontend/list_allow_all
// Zero allocations to memory. Useful to reduce allocations by reusing Path
// pointer because it calls internally Reset.
func (p *Path) Parse(routeOrFQPath string) (err error) {
	p.Reset()
	routeOrFQPath = p.stripEnvSuffixStr(routeOrFQPath)
	// this is the most fast version I come up with.
	// moving from strings to bytes was even slower despite inline
	// th parse int64 function
	if strings.Count(routeOrFQPath, sPathSeparator) < PathLevels-1 {
		return errors.NotValid.Newf("[config] Expecting: `aa/bb/cc` or `strScope/ID/aa/bb/cc` but got %q`", routeOrFQPath)
	}

	fi1 := strings.Index(routeOrFQPath, sPathSeparator)
	scopeStr := routeOrFQPath[:fi1]

	fi2 := strings.Index(routeOrFQPath[fi1+1:], sPathSeparator)
	scopeIDStr := routeOrFQPath[fi1+1 : fi1+1+fi2]

	var scopeID int64
	p.route = routeOrFQPath
	p.ScopeID = scope.DefaultTypeID

	if isDigitOnly(scopeIDStr) {
		scopeID, err = strconv.ParseInt(scopeIDStr, 10, 64)
		if err != nil {
			return errors.NotValid.New(err, "[config] ParseInt with value: %q", scopeIDStr)
		}
		if !scope.Valid(scopeStr) {
			// if scope is not valid, the next part MUST no be an integer
			return errors.NotSupported.Newf("[config] Unknown Scope: %q", scopeStr)
		}
		p.route = routeOrFQPath[fi1+1+fi2+1:]
		p.ScopeID = scope.MakeTypeID(scope.FromString(scopeStr), scopeID)
	}

	return p.IsValid()
}

// ParseStrings parses the arguments into a valid path. scp must be a valid
// scope string as defined in stores/scope package. id must be a stringified
// uint value.
func (p *Path) ParseStrings(scp, id, route string) error {
	p.Reset()
	route = p.stripEnvSuffixStr(route)
	if !scope.Valid(scp) {
		return errors.NotValid.Newf("[config] %q Invalid scope: %q", route, scp)
	}
	id2, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.CorruptData.New(err, "[config] %q failed to parse %q to uint", route, id)
	}
	p.route = route
	p.ScopeID = scope.FromString(scp).WithID(int64(id2))
	return p.IsValid()
}

// IsValid checks for valid configuration path. Returns nil on success.
// Configuration path attribute can have only three groups of [a-zA-Z0-9_] characters split by '/'.
// Minimal length per part 2 characters. Case sensitive.
//
// Error behaviour: NotValid or Empty
func (p *Path) IsValid() error {
	if p == nil {
		return errors.Empty.Newf("[config] Path is nil")
	}
	seps := p.Separators()
	if !p.routeValidated {
		if "" == p.route {
			return errRouteEmpty
		}
		rLen := len(p.route)

		if seps == rLen {
			return errors.NotValid.Newf(errIncorrectPathTpl, p.route)
		}

		if !utf8.ValidString(p.route) {
			return errors.NotValid.Newf(errRouteInvalidBytesTpl, p.route)
		}
		for _, rn := range p.route {
			switch {
			case rn == rune(PathSeparator), rn == '_':
				// ok
			case unicode.IsDigit(rn), unicode.IsLetter(rn), unicode.IsNumber(rn):
				// ok
			default:
				return errors.NotValid.Newf("[config] Route %q contains invalid character: %q", p.route, rn)
			}
		}
		if err := p.ScopeID.IsValid(); err != nil {
			return errors.NotValid.Newf("[config] Route %q contains invalid ScopeID: %q", p.route, p.ScopeID.String())
		}
		if idx := strings.Index(p.route, scope.StrDefault.String()); idx >= 0 && rLen > idx+7 && p.route[:idx+7] == scope.StrDefault.String() {
			return errors.NotValid.Newf("[config] Route cannot start with: %q", scope.StrDefault.String())
		}
		if idx := strings.Index(p.route, scope.StrWebsites.String()); idx >= 0 && rLen > idx+8 && p.route[:idx+8] == scope.StrWebsites.String() {
			return errors.NotValid.Newf("[config] Route cannot start with: %q", scope.StrWebsites.String())
		}
		if idx := strings.Index(p.route, scope.StrStores.String()); idx >= 0 && rLen > idx+6 && p.route[:idx+6] == scope.StrStores.String() {
			return errors.NotValid.Newf("[config] Route cannot start with: %q", scope.StrStores.String())
		}
	}
	if seps < PathLevels-1 || utf8.RuneCountInString(p.route) < 8 /*aa/bb/cc*/ {
		return errors.NotValid.Newf(errIncorrectPathTpl, p.route)
	}
	return nil
}

// IsEmpty returns true if the underlying route is empty.
func (p *Path) IsEmpty() bool {
	return p == nil || p.route == ""
}

// Equal compares the scope and the route
func (p *Path) Equal(b *Path) bool {
	return p != nil && b != nil && p.ScopeID == b.ScopeID && p.route == b.route
}

// EqualRoute compares only the route.
func (p *Path) EqualRoute(b *Path) bool {
	return p != nil && b != nil && p.route == b.route
}

// Reset sets all fields to the zero value for pointer reuse.
func (p *Path) Reset() *Path {
	p.route = ""
	p.ScopeID = 0
	p.routeValidated = false
	p.UseEnvSuffix = false
	return p
}

// MarshalText implements interface encoding.TextMarshaler.
func (p *Path) MarshalText() (text []byte, err error) {
	var buf bytes.Buffer
	if err := p.AppendFQ(&buf); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}

// UnmarshalText transforms the text into a route with performed validation
// checks. Implements encoding.TextUnmarshaler.
// Error behaviour: NotValid, Empty.
func (p *Path) UnmarshalText(txt []byte) error {
	p.Reset()
	txt = p.stripEnvSuffixByte(txt)
	if !(bytes.Count(txt, bSeparator) >= PathLevels+1) {
		return errors.NotValid.Newf("[config] Incorrect fully qualified path: %q. Expecting: strScope/ID/%s", txt, txt)
	}

	fi := bytes.Index(txt, bSeparator)
	scopeStr := txt[:fi]

	if !scope.ValidBytes(scopeStr) {
		return errors.NotSupported.Newf("[config] Unknown Scope: %q", scopeStr)
	}

	txt = txt[fi+1:]
	fi = bytes.Index(txt, bSeparator)
	scopeID, _, err := byteconv.ParseInt(txt[:fi])
	if err != nil {
		return errors.NotValid.New(err, "[config] ParseInt")
	}

	p.route = string(txt[fi+1:])
	p.ScopeID = scope.MakeTypeID(scope.FromBytes(scopeStr), scopeID)
	return errors.NotValid.New(p.IsValid(), "[config] ParseInt")
}

// MarshalBinary implements interface encoding.BinaryMarshaler.
func (p *Path) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	buf.Grow(8)
	sBuf := buf.Bytes()[:8]
	binary.LittleEndian.PutUint64(sBuf[:], p.ScopeID.ToUint64())
	buf.Reset()
	buf.Write(sBuf[:])
	buf.WriteString(p.route)
	p.writeEnvSuffix(&buf)
	return buf.Bytes(), nil
}

// UnmarshalBinary decodes input bytes into a valid Path. Implements
// encoding.BinaryUnmarshaler.
func (p *Path) UnmarshalBinary(data []byte) error {
	p.Reset()
	if len(data) < 8+5 { // 8 for the uint and min 5 bytes for a/b/c
		return errors.TooShort.Newf("[config] UnmarshalBinary: input data too short")
	}
	p.ScopeID = scope.TypeID(binary.LittleEndian.Uint64(data[:8]))
	p.route = string(p.stripEnvSuffixByte(data[8:]))
	return errors.WithStack(p.IsValid())
}

// Level returns a hierarchical based route depending on the depth. The depth
// argument defines the depth of levels to be returned. Depth 1 will return the
// first part like "a", Depth 2 returns "a/b" Depth 3 returns "a/b/c" and so on.
// Level -1 gives you all available levels. Does generate a fully qualified
// path.
func (p *Path) Level(depth int) (string, error) {
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
		sc := strings.IndexByte(fq[pos:], PathSeparator)
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

// Hash64ByLevel same as Level() but returns a HighwayHash-64 checksum of data.
// Usage as map key.
func (p *Path) Hash64ByLevel(depth int) uint64 {
	r2, err := p.Level(depth)
	if err != nil {
		return 0
	}
	return highwayhash.Sum64([]byte(r2), highwayHashKey[:])
}

// for now a math.rand.Read random key.
var highwayHashKey = [highwayhash.Size]byte{0x52, 0xfd, 0xfc, 0x7, 0x21, 0x82, 0x65, 0x4f, 0x16, 0x3f, 0x5f, 0xf, 0x9a, 0x62, 0x1d, 0x72, 0x95, 0x66, 0xc7, 0x4d, 0x10, 0x3, 0x7c, 0x4d, 0x7b, 0xbb, 0x4, 0x7, 0xd1, 0xe2, 0xc6, 0x49}

// Hash64 computes the HighwayHash-64 checksum of data.
// Returns zero in case of an error.
// Usage as map key.
func (p *Path) Hash64() uint64 {
	buf := bufferpool.Get()
	if err := p.AppendFQ(buf); err != nil {
		bufferpool.Put(buf)
		return 0
	}
	s := highwayhash.Sum64(buf.Bytes(), highwayHashKey[:])
	bufferpool.Put(buf)
	return s
}

// Separators returns the number of separators
func (p *Path) Separators() (count int) {
	for _, b := range p.route {
		if b == rune(PathSeparator) {
			count++
		}
	}
	return
}

// ScopeRoute returns the assigned scope and its ID and the route.
func (p *Path) ScopeRoute() (scope.TypeID, string) {
	if p.UseEnvSuffix && p.envSuffix != "" {
		return p.ScopeID, p.route + sPathSeparator + p.envSuffix
	}
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
func (p *Path) Part(pos int) (string, error) {
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
		if b == rune(PathSeparator) && sp < maxLevels {
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
func (p *Path) Split(ret ...string) (_ []string, err error) {

	const sepCount = PathLevels - 1 // only two separators supported
	var sepPos [sepCount]int
	sp := 0
	for i, b := range p.route {
		if b == rune(PathSeparator) && sp < sepCount {
			sepPos[sp] = i // positions of the separators in the slice
			sp++
		}
	}
	if sp < 1 {
		return nil, errors.NotValid.Newf(errIncorrectPathTpl, p.route)
	}
	if ret == nil {
		ret = make([]string, 0, sp+1)
	}

	min := 0
	for i := 0; i < PathLevels; i++ {
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

// NewValue creates a new value with an assigned path. Guarantees not to return
// nil.
func (p Path) NewValue(data []byte) *Value {
	v := NewValue(data)
	v.Path = p
	return v
}

// RouteHasPrefix returns true if the Paths' route starts with the argument route
func (p *Path) RouteHasPrefix(route string) bool {
	lr := len(route)
	return p != nil && len(p.route) >= lr && lr > 0 && p.route[0:lr] == route
}

// PathSlice represents a collection of Paths
type PathSlice []*Path

// add more functions if needed

// Contains return true if the Path p can be found within the slice.
// It must match ID, Scope and Route.
func (ps PathSlice) Contains(p *Path) bool {
	for _, pps := range ps {
		if p != nil && pps.ScopeID == p.ScopeID && pps.route == p.route {
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
