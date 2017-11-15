// FilterThis filters the current slice by predicate f without memory allocation.
// Auto generated via dmlgen.
func (s slice{{.Entity}}) FilterThis(f func(*{{.Entity}}) bool) slice{{.Entity}} {
	b := s[:0]
	for _, x := range s {
		if f(x) {
			b = append(b, x)
		}
	}
	return b
}

// Filter returns a new slice filtered by predicate f.
// Auto generated via dmlgen.
func (s slice{{.Entity}}) Filter(f func(*{{.Entity}}) bool) slice{{.Entity}} {
	sl := make(slice{{.Entity}}, 0, len(s))
	for _, e := range s {
		if f(e) {
			sl = append(sl, e)
		}
	}
	return sl
}

// Each will run function f on all items in []*{{.Entity}}.
// Auto generated via dmlgen.
func (s slice{{.Entity}}) Each(f func(*{{.Entity}})) slice{{.Entity}} {
	for i := range s {
		f(s[i])
	}
	return s
}

// Swap will satisfy the sort.Interface.
// Auto generated via dmlgen.
func (s slice{{.Entity}}) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Cut will remove items i through j-1.
// Auto generated via dmlgen.
func (s *slice{{.Entity}}) Cut(i, j int) {
	z := *s // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this should avoid the memory leak
	}
	z = z[:len(z)-j+i]
	*s = z
}

// Delete will remove an item from the slice.
// Auto generated via dmlgen.
func (s *slice{{.Entity}}) Delete(i int) {
	z := *s // copy the slice header
	end := len(z) - 1
	s.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	*s = z
}

// Insert will place a new item at position i.
// Auto generated via dmlgen.
func (s *slice{{.Entity}}) Insert(n *{{.Entity}}, i int) {
	z := *s // copy the slice header
	z = append(z, &{{.Entity}}{})
	copy(z[i+1:], z[i:])
	z[i] = n
	*s = z
}

// Append will add a new item at the end of slice{{.Entity}}.
// Auto generated via dmlgen.
func (s *slice{{.Entity}}) Append(n ...*{{.Entity}}) {
	*s = append(*s, n...)
}

// Prepend will add a new item at the beginning of slice{{.Entity}}.
// Auto generated via dmlgen.
func (s *slice{{.Entity}}) Prepend(n *{{.Entity}}) {
	s.Insert(n, 0)
}
