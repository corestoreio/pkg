package store

import "sort"

type sortTableStoreSlice struct {
	TableStoreSlice
	LessFunc func(*TableStore, *TableStore) bool
}

func (s sortTableStoreSlice) Less(i, j int) bool {
	return s.LessFunc(s.TableStoreSlice[i], s.TableStoreSlice[j])
}
func (s TableStoreSlice) Len() int {
	return len(s)
}
func (s TableStoreSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s TableStoreSlice) Sort(less func(*TableStore, *TableStore) bool) {
	sort.Sort(sortTableStoreSlice{TableStoreSlice: s, LessFunc: less})
}
func (s TableStoreSlice) Find(f func(*TableStore) bool) (match *TableStore, found bool) {
	for _, u := range s {
		if f(u) {
			match = u
			found = true
			return
		}
	}
	return
}
func (s TableStoreSlice) FilterThis(f func(*TableStore) bool) TableStoreSlice {
	b := s[:0]
	for _, x := range s {
		if f(x) {
			b = append(b, x)
		}
	}
	return b
}
func (s TableStoreSlice) Filter(f func(*TableStore) bool) TableStoreSlice {
	sl := make(TableStoreSlice, 0, len(s))
	for _, w := range s {
		if f(w) {
			sl = append(sl, w)
		}
	}
	return sl
}
func (s TableStoreSlice) FilterNot(f func(*TableStore) bool) TableStoreSlice {
	sl := make(TableStoreSlice, 0, len(s))
	for _, v := range s {
		if f(v) == false {
			sl = append(sl, v)
		}
	}
	return sl
}
func (s TableStoreSlice) Each(f func(int, *TableStore)) TableStoreSlice {
	for i := range s {
		f(i, s[i])
	}
	return s
}
func (s *TableStoreSlice) Cut(i, j int) {
	z := *s
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil
	}
	z = z[:len(z)-j+i]
	*s = z
}
func (s *TableStoreSlice) Delete(i int) {
	z := *s
	end := len(z) - 1
	s.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil
	z = z[:end]
	*s = z
}
func (s *TableStoreSlice) Insert(n *TableStore, i int) {
	z := *s
	z = append(z, new(TableStore))
	copy(z[i+1:], z[i:])
	z[i] = n
	*s = z
}
func (s *TableStoreSlice) Append(n ...*TableStore) {
	*s = append(*s, n...)
}
func (s *TableStoreSlice) Prepend(n *TableStore) {
	s.Insert(n, 0)
}
