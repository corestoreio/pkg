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

package scope

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
)

// MaxID maximum allowed ID which can be packed into a TypeID. The ID relates to
// an auto_increment column in the database. Doesn't matter whether we have a
// website, group or store scope. int24 (8388607) size at the moment.
const MaxID int64 = 1<<23 - 1

// DefaultTypeID default Hash value for Default Scope and ID 0. Avoids typing
// 		scope.NewHash(DefaultID,0)
const DefaultTypeID TypeID = TypeID(Default)<<24 | 0

// Hash defines a merged Scope with its ID. The first 8 bit represents the
// Type: Default, Website, Group or Store. The last 24 bit represents the
// assigned ID. This ID relates to the database table in M2 to `website`,
// `store` or `store_group` and for M1 to `core_website`, `core_store` and
// `core_store_group`. The maximum ID which can be used is defined in constant
// MaxID.
type TypeID uint32

// If we have need for more store IDs then we can change the underlying types here.

// String human readable output
func (t TypeID) String() string {
	scp, id := t.Unpack()
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	_, _ = buf.WriteString("Type(")
	_, _ = buf.WriteString(scp.String())
	_, _ = buf.WriteString(") ID(")
	nb := strconv.AppendInt(buf.Bytes(), id, 10)
	buf.Reset()
	_, _ = buf.Write(nb)
	_ = buf.WriteByte(')')
	return buf.String()
}

// GoString compilable representation of a hash.
func (t TypeID) GoString() string {
	scp, id := t.Unpack()
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	_, _ = buf.WriteString("scope.MakeTypeID(scope.")
	_, _ = buf.WriteString(scp.String())
	_, _ = buf.WriteString(", ")
	nb := strconv.AppendInt(buf.Bytes(), id, 10)
	buf.Reset()
	_, _ = buf.Write(nb)
	_ = buf.WriteByte(')')
	return buf.String()
}

// ToUint64 converts the hash
func (t TypeID) ToUint64() uint64 {
	return uint64(t)
}

// Unpack extracts a Scope and its ID from a hash. Returned ID can be -1 when
// the Hash contains invalid data. An ID of -1 is considered an error.
func (t TypeID) Unpack() (s Type, id int64) {

	prospectS := t >> 24
	if prospectS > maxUint8 || prospectS < 0 {
		return Absent, -1
	}
	s = Type(prospectS)

	h64 := int64(t)
	prospectID := h64 ^ (h64>>24)<<24
	if prospectID > MaxID || prospectID < 0 {
		return Absent, -1
	}

	id = prospectID
	return
}

// EqualTypes compares the type of two TypeIDs and returns true if their type
// matches. This functions checks overflows, would then return false. Two
// TypeIDs with an Absent type are never equal.
func (t TypeID) EqualTypes(other TypeID) bool {
	hScope := t >> 24
	if hScope > maxUint8 || hScope <= 0 {
		return false
	}
	oScope := other >> 24
	if oScope > maxUint8 || oScope <= 0 {
		return false
	}
	return hScope == oScope
}

// Type returns the underlying assigned type.
func (t TypeID) Type() Type {
	hScope := t >> 24
	if hScope > maxUint8 || hScope < 0 {
		return Absent
	}
	return Type(hScope)
}

// ID returns the underlying assigned ID. If the ID overflows the MaxID or
// is smaller than zero then it returns -1.
func (t TypeID) ID() int64 {
	h64 := int64(t)
	prospectID := h64 ^ (h64>>24)<<24
	if prospectID > MaxID || prospectID < 0 {
		return -1
	}
	return prospectID
}

// ValidParent validates if the parent Type is within the hierarchical chain:
// default -> website -> store.
func (t TypeID) ValidParent(parent TypeID) bool {
	p, pID := parent.Unpack()
	c, cID := t.Unpack()
	return (p == Default && pID == 0 && c == Default && cID == 0) ||
		(p == Default && pID == 0 && c == Website && cID >= 0) ||
		(p == Website && pID >= 0 && c == Store && cID >= 0)
}

// CalculateRunMode transforms the Hash into a runMode. On an invalid Hash (the
// Type is < Website or Type > Store) it falls back to the default run mode,
// which is a zero Hash. Implements interface RunModeCalculater.
func (t TypeID) CalculateRunMode(_ *http.Request) TypeID {
	if s := t.Type(); s < Website || s > Store {
		// fall back to default because only Website, Group and Store are allowed.
		t = DefaultRunMode
	}
	return t
}

// TypeIDMaxSegments maximum supported segments or also known as shards. This
// constant can be used to create the segmented array in other packages.
const TypeIDMaxSegments uint16 = 256

const hashBitAnd TypeID = TypeID(TypeIDMaxSegments) - 1

// Segment generates an 0 < ID <= 255 from a hash. Only used within an array
// index to optimize map[] usage in high concurrent situations. Also known as
// shard. An array of N shards is created, each shard contains its own instance
// of the cache with a lock. When an item with unique key needs to be cached a
// shard for it is chosen at first by the function Segment(). After that the
// cache lock is acquired and a write to the cache takes place. Reads are
// analogue.
func (h TypeID) Segment() uint8 {
	return uint8(h & hashBitAnd)
}

// MakeTypeID creates a new merged value of a Type and its ID. An error is equal
// to returning 0. An error occurs when id is greater than MaxStoreID or smaller
// 0. An errors occurs when the Scope is Default and ID anything else than 0.
func MakeTypeID(s Type, id int64) TypeID {
	if id > MaxID || (s > Default && id < 0) {
		return 0
	}
	if s < Website {
		id = 0
	}
	return TypeID(s)<<24 | TypeID(id)
}

// TypeIDs collection of multiple TypeID values.
type TypeIDs []TypeID

// Len is part of sort.Interface.
func (t TypeIDs) Len() int { return len(t) }

// Swap is part of sort.Interface.
func (t TypeIDs) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Less is part of sort.Interface.
func (t TypeIDs) Less(i, j int) bool { return t[i] < t[j] }

// Lowest finds from hashes the common lowest Type. All Types must have within
// their Type the same ID otherwise an error will be returned. This functions
// gets mainly used in backend* packages if several configuration paths must be
// applied to one functional option. Eg. config path A has Type Website(1) but
// config path B has Type Store(2) and config path C has Type Website(1) so
// the common valid hash resolves to Store(2). If there would be a config path
// with Type Store(3) then a NotValid error gets returned.
func (t TypeIDs) Lowest() (TypeID, error) {
	sort.Stable(t)
	var pick = DefaultTypeID
	wIDs, gIDs, sIDs := float64(0), float64(0), float64(0)
	wC, gC, sC := float64(0), float64(0), float64(0)
	for _, v := range t {

		if v.Type() > pick.Type() {
			pick = v
		}

		switch v.Type() {
		case Website:
			wC++
			wIDs += float64(v.ID())
		case Group:
			gC++
			gIDs += float64(v.ID())
		case Store:
			sC++
			sIDs += float64(v.ID())
		}
	}

	switch pick.Type() {
	case Website:
		if float64(pick.ID()) != wIDs/wC {
			return 0, errors.NewNotValidf("[scope] Invalid hash: %s in slice.", pick)
		}
	case Group:
		if float64(pick.ID()) != gIDs/gC {
			return 0, errors.NewNotValidf("[scope] Invalid hash: %s in slice.", pick)
		}
	case Store:
		if float64(pick.ID()) != sIDs/sC {
			return 0, errors.NewNotValidf("[scope] Invalid hash: %s in slice.", pick)
		}
	case Default, Absent:
		// do nothing
	default:
		// todo implement scope independent solution ...
		return 0, errors.NewNotValidf("[scope] Invalid hash: %s in slice.", pick)

	}

	return pick, nil
}
