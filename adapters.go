package adapters

import (
	"iter"
)

// Filter returns a new sequence containing only the elements that satisfy the predicate.
//
// Example:
//
//	numbers := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//	evenNumbers := adapters.Filter(numbers, func(n int) bool { return n%2 == 0 })
//	// evenNumbers will yield: 2, 4, 6, 8, 10
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

// Map applies a transformation function to each element in the sequence.
//
// Example:
//
//	numbers := slices.Values([]int{1, 2, 3, 4, 5})
//	squared := adapters.Map(numbers, func(n int) int { return n * n })
//	// squared will yield: 1, 4, 9, 16, 25
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

// Reduce applies a reducer function to all elements of the sequence,
// accumulating the result into a single value.
//
// Example:
//
//	numbers := slices.Values([]int{1, 2, 3, 4, 5})
//	sum := adapters.Reduce(numbers, 0, func(acc, n int) int { return acc + n })
//	// sum will be 15
func Reduce[T, R any](s iter.Seq[T], initial R, reducer func(R, T) R) R {
	result := initial
	for v := range s {
		result = reducer(result, v)
	}
	return result
}

// Take returns a new sequence with at most n elements from the original sequence.
//
// Example:
//
//	numbers := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//	firstFive := adapters.Take(numbers, 5)
//	// firstFive will yield: 1, 2, 3, 4, 5
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

// Skip returns a new sequence that skips the first n elements of the original sequence.
//
// Example:
//
//	numbers := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//	afterFive := adapters.Skip(numbers, 5)
//	// afterFive will yield: 6, 7, 8, 9, 10
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

// Zip combines two sequences into a single sequence of pairs.
//
// Example:
//
//	numbers := slices.Values([]int{1, 2, 3, 4, 5})
//	letters := slices.Values([]string{"a", "b", "c", "d", "e"})
//	zipped := adapters.Zip(numbers, letters)
//	// zipped will yield: (1, "a"), (2, "b"), (3, "c"), (4, "d"), (5, "e")
func Zip[T, U any](s1 iter.Seq[T], s2 iter.Seq[U]) iter.Seq[struct {
	First  T
	Second U
}] {
	return func(yield func(struct {
		First  T
		Second U
	}) bool) {
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
			if !yield(struct {
				First  T
				Second U
			}{v1, v2}) {
				return
			}
		}
	}
}

// FlatMap applies a transformation to each element and flattens the result.
//
// Example:
//
//	nestedNumbers := slices.Values([][]int{{1, 2}, {3, 4}, {5, 6}})
//	flattened := adapters.FlatMap(nestedNumbers, func(slice []int) iter.Seq[int] {
//		return slices.Values(slice)
//	})
//	// flattened will yield: 1, 2, 3, 4, 5, 6
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

// Flatten flattens a sequence of various types into a single sequence.
// It can handle regular sequences, nullable sequences, and single values.
//
// Example:
//
//	mixedSeq := slices.Values([]any{
//		[]int{1, 2, 3},
//		iter.Seq[int](func(yield func(int) bool) {
//			yield(4)
//			yield(5)
//		}),
//		6,
//		&[]int{7, 8},
//		[]int{9, 10},
//	})
//	flattenedMixed := adapters.Flatten[int](mixedSeq)
//	// flattenedMixed will yield: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
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

// FilterMap applies a transformation function to each element in the sequence,
// keeping only the successful results and discarding errors.
//
// Example:
//
//	numbers := slices.Values([]int{1, 2, 3, 4, 5})
//	squared := adapters.FilterMap(numbers, func(n int) (int, error) {
//		if n == 3 {
//			return 0, fmt.Errorf("skipping 3")
//		}
//		return n * n, nil
//	})
//	// squared will yield: 1, 4, 16, 25
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
