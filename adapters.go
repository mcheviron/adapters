package adapters

import (
	"iter"
)

func Filter[T any](s iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if pred(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Filter2[K, V any](s iter.Seq2[K, V], pred func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if pred(k, v) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

func Map[T, R any](s iter.Seq[T], transform func(T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			result := transform(v)
			if !yield(result) {
				return
			}
		}
	}
}

func Map2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], transform func(K1, V1) (K2, V2)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s {
			k2, v2 := transform(k, v)
			if !yield(k2, v2) {
				return
			}
		}
	}
}

func Reduce[T, R any](s iter.Seq[T], initial R, reducer func(R, T) R) R {
	result := initial
	for v := range s {
		result = reducer(result, v)
	}
	return result
}

func Take[T any](s iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		for v := range s {
			if count >= n {
				return
			}
			if !yield(v) {
				return
			}
			count++
		}
	}
}

func Take2[K, V any](s iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		count := 0
		for k, v := range s {
			if count >= n {
				return
			}
			if !yield(k, v) {
				return
			}
			count++
		}
	}
}

func Skip[T any](s iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		for v := range s {
			if count >= n {
				if !yield(v) {
					return
				}
			}
			count++
		}
	}
}

func Skip2[K, V any](s iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		count := 0
		for k, v := range s {
			if count >= n {
				if !yield(k, v) {
					return
				}
			}
			count++
		}
	}
}

func Zip[T, U any](s1 iter.Seq[T], s2 iter.Seq[U]) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		next1, stop1 := iter.Pull(s1)
		next2, stop2 := iter.Pull(s2)
		defer stop1()
		defer stop2()

		for {
			v1, ok1 := next1()
			if !ok1 {
				return
			}
			v2, ok2 := next2()
			if !ok2 {
				return
			}
			if !yield(v1, v2) {
				return
			}
		}
	}
}

func FlatMap[T, R any](s iter.Seq[T], transform func(T) iter.Seq[R]) iter.Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			innerSeq := transform(v)
			for innerV := range innerSeq {
				if !yield(innerV) {
					return
				}
			}
		}
	}
}

func FlatMap2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], transform func(K1, V1) iter.Seq2[K2, V2]) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range s {
			innerSeq := transform(k1, v1)
			for k2, v2 := range innerSeq {
				if !yield(k2, v2) {
					return
				}
			}
		}
	}
}

func Flatten[T any](s iter.Seq[any]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range s {
			switch v := item.(type) {
			case iter.Seq[T]:
				for innerV := range v {
					if !yield(innerV) {
						return
					}
				}
			case *iter.Seq[T]:
				if v != nil {
					for innerV := range *v {
						if !yield(innerV) {
							return
						}
					}
				}
			case []T:
				for _, innerV := range v {
					if !yield(innerV) {
						return
					}
				}
			case *[]T:
				if v != nil {
					for _, innerV := range *v {
						if !yield(innerV) {
							return
						}
					}
				}
			case T:
				if !yield(v) {
					return
				}
			}
		}
	}
}

func FilterMap[T, R any](s iter.Seq[T], transform func(T) (R, error)) iter.Seq[R] {
	return func(yield func(R) bool) {
		for v := range s {
			if result, err := transform(v); err == nil {
				if !yield(result) {
					return
				}
			}
		}
	}
}

func FilterMap2[K1, V1, K2, V2 any](s iter.Seq2[K1, V1], transform func(K1, V1) (K2, V2, error)) iter.Seq2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k1, v1 := range s {
			if k2, v2, err := transform(k1, v1); err == nil {
				if !yield(k2, v2) {
					return
				}
			}
		}
	}
}
