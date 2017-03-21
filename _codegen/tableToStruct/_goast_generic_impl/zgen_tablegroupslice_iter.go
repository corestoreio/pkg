package store

import "sort"

type sortTableGroupSlice struct {
	TableGroupSlice
	LessFunc func(*TableGroup, *TableGroup) bool
}

func (s sortTableGroupSlice) Less(i, j int) bool {
	return s.LessFunc(s.TableGroupSlice[i], s.TableGroupSlice[j])
}
func (s TableGroupSlice) Len() int {
	return len(s)
}
func (s TableGroupSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s TableGroupSlice) Sort(less func(*TableGroup, *TableGroup) bool) {
	sort.Sort(sortTableGroupSlice{TableGroupSlice: s, LessFunc: less})
}
func (s TableGroupSlice) Find(f func(*TableGroup) bool) (match *TableGroup, found bool) {
	for _, u := range s {
		if f(u) {
			match = u
			found = true
			return
		}
	}
	return
}
func (s TableGroupSlice) FilterThis(f func(*TableGroup) bool) TableGroupSlice {
	b := s[:0]
	for _, x := range s {
		if f(x) {
			b = append(b, x)
		}
	}
	return b
}
func (s TableGroupSlice) Filter(f func(*TableGroup) bool) TableGroupSlice {
	sl := make(TableGroupSlice, 0, len(s))
	for _, w := range s {
		if f(w) {
			sl = append(sl, w)
		}
	}
	return sl
}
func (s TableGroupSlice) FilterNot(f func(*TableGroup) bool) TableGroupSlice {
	sl := make(TableGroupSlice, 0, len(s))
	for _, v := range s {
		if f(v) == false {
			sl = append(sl, v)
		}
	}
	return sl
}
func (s TableGroupSlice) Each(f func(int, *TableGroup)) TableGroupSlice {
	for i := range s {
		f(i, s[i])
	}
	return s
}
func (s *TableGroupSlice) Cut(i, j int) {
	z := *s
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil
	}
	z = z[:len(z)-j+i]
	*s = z
}
func (s *TableGroupSlice) Delete(i int) {
	z := *s
	end := len(z) - 1
	s.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil
	z = z[:end]
	*s = z
}
func (s *TableGroupSlice) Insert(n *TableGroup, i int) {
	z := *s
	z = append(z, new(TableGroup))
	copy(z[i+1:], z[i:])
	z[i] = n
	*s = z
}
func (s *TableGroupSlice) Append(n ...*TableGroup) {
	*s = append(*s, n...)
}
func (s *TableGroupSlice) Prepend(n *TableGroup) {
	s.Insert(n, 0)
}
