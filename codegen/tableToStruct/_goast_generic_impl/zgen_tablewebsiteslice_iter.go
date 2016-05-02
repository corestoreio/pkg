package store

import "sort"

type sortTableWebsiteSlice struct {
	TableWebsiteSlice
	LessFunc func(*TableWebsite, *TableWebsite) bool
}

func (s sortTableWebsiteSlice) Less(i, j int) bool {
	return s.LessFunc(s.TableWebsiteSlice[i], s.TableWebsiteSlice[j])
}
func (s TableWebsiteSlice) Len() int {
	return len(s)
}
func (s TableWebsiteSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s TableWebsiteSlice) Sort(less func(*TableWebsite, *TableWebsite) bool) {
	sort.Sort(sortTableWebsiteSlice{TableWebsiteSlice: s, LessFunc: less})
}
func (s TableWebsiteSlice) Find(f func(*TableWebsite) bool) (match *TableWebsite, found bool) {
	for _, u := range s {
		if f(u) {
			match = u
			found = true
			return
		}
	}
	return
}
func (s TableWebsiteSlice) FilterThis(f func(*TableWebsite) bool) TableWebsiteSlice {
	b := s[:0]
	for _, x := range s {
		if f(x) {
			b = append(b, x)
		}
	}
	return b
}
func (s TableWebsiteSlice) Filter(f func(*TableWebsite) bool) TableWebsiteSlice {
	sl := make(TableWebsiteSlice, 0, len(s))
	for _, w := range s {
		if f(w) {
			sl = append(sl, w)
		}
	}
	return sl
}
func (s TableWebsiteSlice) FilterNot(f func(*TableWebsite) bool) TableWebsiteSlice {
	sl := make(TableWebsiteSlice, 0, len(s))
	for _, v := range s {
		if f(v) == false {
			sl = append(sl, v)
		}
	}
	return sl
}
func (s TableWebsiteSlice) Each(f func(int, *TableWebsite)) TableWebsiteSlice {
	for i := range s {
		f(i, s[i])
	}
	return s
}
func (s *TableWebsiteSlice) Cut(i, j int) {
	z := *s
	copy(z[i:], z[j:])
	for k, n := len(z)-j+i, len(z); k < n; k++ {
		z[k] = nil
	}
	z = z[:len(z)-j+i]
	*s = z
}
func (s *TableWebsiteSlice) Delete(i int) {
	z := *s
	end := len(z) - 1
	s.Swap(i, end)
	copy(z[i:], z[i+1:])
	z[end] = nil
	z = z[:end]
	*s = z
}
func (s *TableWebsiteSlice) Insert(n *TableWebsite, i int) {
	z := *s
	z = append(z, new(TableWebsite))
	copy(z[i+1:], z[i:])
	z[i] = n
	*s = z
}
func (s *TableWebsiteSlice) Append(n ...*TableWebsite) {
	*s = append(*s, n...)
}
func (s *TableWebsiteSlice) Prepend(n *TableWebsite) {
	s.Insert(n, 0)
}
