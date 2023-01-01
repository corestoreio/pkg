// Package generic contains all generic functions.
// Please use https://github.com/samber/lo
package generic

func Contains[T comparable](s []T, v T) bool {
	for _, vs := range s {
		if v == vs {
			return true
		}
	}
	return false
}

func Reverse[T any](s []T) []T {
	result := make([]T, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		result = append(result, s[i])
	}
	return result
}

// MapNew maps function f to each slice element and returns a new slice, leaving
// input slice "s" untouched.
func MapNew[T any](s []T, f func(T) T) []T {
	result := make([]T, len(s))
	for i := range s {
		result[i] = f(s[i])
	}
	return result
}

// Map maps function f to each element of slice "s". It changes slice "s".
func Map[T any](s []T, f func(T) T) []T {
	for i := range s {
		s[i] = f(s[i])
	}
	return s
}

// FilterNew filters slice s with predicate f. Returns a new slice, leaving s
// untouched. Might return a nil slice.
func FilterNew[T any](s []T, f func(T) bool) []T {
	var result []T
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Filter filters slice s with predicate f and re-uses the backing array of "s".
func Filter[T any](s []T, f func(T) bool) []T {
	result := s[:0]
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func Reduce[T any](s []T, init T, f func(T, T) T) T {
	cur := init
	for _, v := range s {
		cur = f(cur, v)
	}
	return cur
}

func ToAnySlice[T any](slice []T) []any {
	ret := make([]any, len(slice))
	for i, item := range slice {
		ret[i] = item
	}
	return ret
}

func RemoveEmpty[T comparable](slIN []T) []T {
	if len(slIN) == 0 {
		return nil
	}
	sl := slIN[:0]
	for _, d := range slIN {
		if d != "" {
			sl = append(sl, d)
		}
	}
	return sl
}

func ToStringSlice[T ~string](strs ...T) []string {
	ret := make([]string, len(strs))
	for i, t := range strs {
		ret[i] = string(t)
	}
	return ret
}

// Intersection returns a list of common elements present in both arrays.
// The elements in the output can be in any order.
func Intersection[A comparable](a []A, b []A) (intersection []A) {
	for _, aa := range a {
		for _, bb := range b {
			if aa == bb {
				intersection = append(intersection, aa)
			}
		}
	}
	return intersection
}

func Window[S any](elements []S, size int) [][]S {
	batchSize := (len(elements) + size - 1) / size
	var batches [][]S
	for batchSize < len(elements) {
		elements, batches = elements[batchSize:], append(batches, elements[0:batchSize:batchSize])
	}
	batches = append(batches, elements)
	return batches
}
