package z

import "fmt"

func SliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func PointerOf[T any](v T) *T { return &v }

func ValueOf[T any](v *T) (o T) { return When(v == nil, o, *v) }

func When[T any](cond bool, Tval, Fval T) T {
	if cond {
		return Tval
	}
	return Fval
}

func Yield[T any](v T) func() T { return func() T { return v } }

func String[T any](v T) string { return fmt.Sprint(v) }

func CloneMap[T1 comparable, T2 any](m map[T1]T2) (o map[T1]T2) {
	o = make(map[T1]T2)
	for k, v := range m {
		o[k] = v
	}
	return o
}
