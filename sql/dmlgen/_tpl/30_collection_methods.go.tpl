// Filter filters the current slice by predicate f without memory allocation.
// Auto generated via dmlgen.
func (cc *{{.Collection}}) Filter(f func(*{{.Entity}}) bool) *{{.Collection}} {
	b,i := cc.Data[:0],0
	for _, e := range cc.Data {
		if f(e) {
			b = append(b, e)
			cc.Data[i] = nil // this avoids the memory leak
		}
		i++
	}
	cc.Data = b
	return cc
}

// Each will run function f on all items in []*{{.Entity}}.
// Auto generated via dmlgen.
func (cc *{{.Collection}}) Each(f func(*{{.Entity}})) *{{.Collection}} {
	for i := range cc.Data {
		f(cc.Data[i])
	}
	return cc
}

// Cut will remove items i through j-1.
// Auto generated via dmlgen.
func (cc *{{.Collection}}) Cut(i, j int) *{{.Collection}} {
	z := cc.Data // copy slice header
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil // this avoids the memory leak
	}
	z = z[:len(z)-j+i]
	cc.Data = z
	return cc
}

// Swap will satisfy the sort.Interface.
// Auto generated via dmlgen.
func (cc *{{.Collection}}) Swap(i, j int) { cc.Data[i], cc.Data[j] = cc.Data[j], cc.Data[i] }

// Delete will remove an item from the slice.
// Auto generated via dmlgen.
func (cc *{{.Collection}}) Delete(i int) *{{.Collection}} {
	z := cc.Data // copy the slice header
	end := len(z) - 1
	cc.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil // this should avoid the memory leak
	z = z[:end]
	cc.Data = z
	return cc
}

// Insert will place a new item at position i.
// Auto generated via dmlgen.
func (cc *{{.Collection}}) Insert(n *{{.Entity}}, i int) *{{.Collection}} {
	z := cc.Data // copy the slice header
	z = append(z, &{{.Entity}}{})
	copy(z[i+1:], z[i:])
	z[i] = n
	cc.Data = z
	return cc
}

// Append will add a new item at the end of *{{.Collection}}.
// Auto generated via dmlgen.
func (cc *{{.Collection}}) Append(n ...*{{.Entity}}) *{{.Collection}} {
	cc.Data = append(cc.Data, n...)
	return cc
}
