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
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/juju/errgo"
)

// ErrIncorrect gets returned whenever the path consists of less than
// Levels parts.
var ErrIncorrect = errors.New("Incorrect Path. Expecting at least three path parts like a/b/c")

// Levels defines how many parts are at least in a path.
// Like a/b/c for 3 parts. And 5 for a fully qualified path.
const Levels int = 3

// PS path separator used in the database table core_config_data and in config.Service
const PS = "/"

const rPS = '/'
const strDefaultID = "0"

// Path represents a configuration path.
type Path struct {
	// Parts either one short path or three path parts
	Parts []string
	Scope scope.Scope
	// ID represents a website, group or store ID
	ID int64
	// NoValidation disables validation in FQ() function
	NoValidation bool
}

// New creates a new validated Path. Argument can either be a path like
// a/b/c or path parts like "a","b","c".
func New(paths ...string) (Path, error) {
	p := Path{
		Parts: paths,
		Scope: scope.DefaultID,
	}
	if err := p.IsValid(); err != nil {
		return Path{}, err
	}
	return p, nil
}

// NewSplit takes a path argument like a/b/c or path parts like "a","b","c".
// If a path have been provided it gets split into its parts.
func NewSplit(paths ...string) (Path, error) {
	p := Path{
		Scope: scope.DefaultID,
	}
	switch {
	case len(paths) >= Levels:
		p.Parts = paths
	case len(paths) == 1:
		p.Parts = Split(paths[0])
	default:
		return Path{}, fmt.Errorf("Incorrect number of paths elements: want %d, have %d, Path: %v", Levels, len(paths), paths)
	}

	if err := p.IsValid(); err != nil {
		return Path{}, err
	}
	return p, nil
}

// MustNew same as New but panics on error.
func MustNew(paths ...string) Path {
	p, err := New(paths...)
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
	return s
}

// FQ returns the fully qualified path. Validation can be disabled by setting
// NoValidation to true.
func (p Path) FQ() (string, error) {
	if err := p.IsValid(); err != nil {
		return "", err
	}

	idStr := "0"
	if p.ID > 0 {
		if p.ID <= int64CacheLen {
			idStr = int64Cache[p.ID]
		} else {
			idStr = strconv.FormatInt(p.ID, 10)
		}
	}

	scopeStr := scope.FromScope(p.Scope)
	if scopeStr == scope.StrDefault && idStr != strDefaultID {
		idStr = strDefaultID // default scope is always 0
	}
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	buf.WriteString(scopeStr.String())
	buf.WriteString(PS)
	buf.WriteString(idStr)
	buf.WriteString(PS)
	join(buf, p.Parts)
	return buf.String(), nil
}

// this "cache" should cover ~80% of all store setups
var int64Cache = [...]string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
}
var int64CacheLen = int64(len(int64Cache))

// Split splits a configuration path by the path separator PS.
func Split(path string) []string {
	if path[:1] == PS {
		path = path[1:] // trim first PS
	}
	return strings.Split(path, PS)
}

func join(buf *bytes.Buffer, paths []string) {
	for i, p := range paths {
		buf.WriteString(p)
		if i < (len(paths) - 1) {
			buf.WriteString(PS)
		}
	}
}

// Level joins a configuration path parts by the path separator PS.
// The level argument defines the depth of the path parts to join.
// Level 1 will return the first part like "a", Level 2 returns "a/b"
// Level 3 returns "a/b/c" and so on. Level -1 joins all available path parts.
// Does not generate a fully qualified path.
func (p Path) Level(level int) string {
	if lp := len(p.Parts); level <= 0 || level >= lp {
		level = lp
	}
	buf := bufferpool.Get()
	join(buf, p.Parts[:level])
	s := buf.String()
	bufferpool.Put(buf)
	return s
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

	fi := strings.Index(fqPath, PS)
	scopeStr := fqPath[:fi]

	if false == scope.Valid(scopeStr) {
		return Path{}, scope.ErrUnsupportedScope
	}

	fqPath = fqPath[fi+1:]

	fi = strings.Index(fqPath, PS)
	scopeID, err := strconv.ParseInt(fqPath[:fi], 10, 64)
	path := fqPath[fi+1:]
	return Path{
		Parts: []string{path},
		Scope: scope.FromString(scopeStr),
		ID:    scopeID,
	}, err
}

func isFQ(fqPath string) bool {
	return strings.Count(fqPath, PS) >= Levels+1 // like stores/1/a/b/c
}

// IsValid checks for valid configuration path. Returns nil on success.
// Configuration path attribute can have only three groups of [a-zA-Z0-9_] characters split by '/'.
// Minimal length per part 2 characters. Case sensitive.
func (p Path) IsValid() error {
	lp := len(p.Parts)
	if lp < 1 {
		return errgo.New("Parts are empty")
	}

	// first argument only without a slash
	if lp == 1 && (strings.Count(p.Parts[0], PS) != Levels-1 || len(p.Parts[0]) < 8) { // must contain at least two slashes
		return errgo.Newf("Incorrect Parts: %#v. Either to short or missing path separator.", p.Parts)
	}

	valid := 0
	for _, part := range p.Parts {
		if len(part) < 2 {
			return errgo.Newf("This path part %q is too short. Parts: %#v", part, p.Parts)
		}

		for _, r := range part {
			ok := false
			switch {
			case '0' <= r && r <= '9':
				ok = true
			case 'a' <= r && r <= 'z':
				ok = true
			case 'A' <= r && r <= 'Z':
				ok = true
			case r == '_', r == rPS:
				ok = true
			}
			if !ok {
				return errgo.Newf("This character %q is not allowed in Parts %#v", string(r), p.Parts)
			}
		}
		valid++
	}

	if lp > 1 && valid < Levels { // if more than one arg has been provided all 3 must be valid
		return errgo.Newf("All arguments must be valid! Min want: %d. Have: %d. Parts %#v", Levels, valid, p.Parts)
	}

	return nil
}
