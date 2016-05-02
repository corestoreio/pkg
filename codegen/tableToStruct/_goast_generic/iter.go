package main

import "sort"

//Impl Types: Requires a slice of any type (interface{})
type I interface{}
type Slice []*I

type sort_ struct {
	Slice
	LessFunc func(*I, *I) bool
}

// Less will satisfy the sort.Interface and compares via
// the primary key.
// Generated via tableToStruct.
func (s sort_) Less(i, j int) bool {
	return s.LessFunc(s.Slice[i], s.Slice[j])
}

// Len returns the length and  will satisfy the sort.Interface.
// Generated via tableToStruct.
func (s Slice) Len() int { return len(s) }

// Swap will satisfy the sort.Interface.
// Generated via tableToStruct.
func (s Slice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Sort will sort Slice.
// Generated via tableToStruct.
func (s Slice) Sort(less func(*I, *I) bool) {
	sort.Sort(sort_{Slice: s, LessFunc: less})
}

// Find takes predicate f and searches through the slice to return the first match.
// Generated via tableToStruct.
func (s Slice) Find(f func(*I) bool) (match *I, found bool) {
	for _, u := range s {
		if f(u) {
			match = u
			found = true
			return
		}
	}
	return
}

// FilterThis filters the current slice by predicate f without memory allocation.
// Generated via tableToStruct.
func (s Slice) FilterThis(f func(*I) bool) Slice {
	b := s[:0]
	for _, x := range s {
		if f(x) {
			b = append(b, x)
		}
	}
	return b
}

// Filter returns a new slice filtered by predicate f.
// Generated via tableToStruct.
func (s Slice) Filter(f func(*I) bool) Slice {
	sl := make(Slice, 0, len(s))
	for _, w := range s {
		if f(w) {
			sl = append(sl, w)
		}
	}
	return sl
}

// FilterNot will return a new Slice that does not match
// by calling the function f
// Generated via tableToStruct.
func (s Slice) FilterNot(f func(*I) bool) Slice {
	sl := make(Slice, 0, len(s))
	for _, v := range s {
		if f(v) == false {
			sl = append(sl, v)
		}
	}
	return sl
}

// Each will run function f on all items in Slice.
// Generated via tableToStruct.
func (s Slice) Each(f func(int, *I)) Slice {
	for i := range s {
		f(i, s[i])
	}
	return s
}

// Cut will remove items i through j-1.
// Generated via tableToStruct.
func (s *Slice) Cut(i, j int) {
	z := *s // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this should avoid the memory leak
	}
	z = z[:len(z)-j+i]
	*s = z
}

// Delete will remove an item from the slice.
// Generated via tableToStruct.
func (s *Slice) Delete(i int) {
	z := *s // copy the slice header
	end := len(z) - 1
	s.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	*s = z
}

// Insert will place a new item at position i.
// Generated via tableToStruct.
func (s *Slice) Insert(n *I, i int) {
	z := *s // copy the slice header
	z = append(z, new(I))
	copy(z[i+1:], z[i:])
	z[i] = n
	*s = z
}

// Append will add a new item at the end of Slice.
// Generated via tableToStruct.
func (s *Slice) Append(n ...*I) {
	*s = append(*s, n...)
}

// Prepend will add a new item at the beginning of Slice.
// Generated via tableToStruct.
func (s *Slice) Prepend(n *I) {
	s.Insert(n, 0)
}
