package main

import (
	"fmt"
	"reflect"
)

// import "fmt"

// Set is a generic implementation of a set using a map.
type Set[T comparable] struct {
	emap map[T]struct{}
}

func (s *Set[T]) String() string {
	str := ""
	for fn := range s.emap {
		Ty := reflect.TypeOf(fn)
		if Ty.Kind() == reflect.Func {
			// TODO: FIXME
			// p := reflect.ValueOf(fn) //.UnsafePointer()
			str += fmt.Sprintf("  %v\n", fn)
		} else {
			str += fmt.Sprintf("  %v\n", fn)
		}
	}
	return str
}

// NewSet creates a new empty set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{emap: make(map[T]struct{})}
}

// Add inserts an element into the set.
func (s *Set[T]) Add(value T) {
	s.emap[value] = struct{}{}
}

// Remove deletes an element from the set.
func (s *Set[T]) Remove(value T) {
	delete(s.emap, value)
}

// Contains checks if an element is in the set.
func (s *Set[T]) Contains(value T) bool {
	_, exists := s.emap[value]
	return exists
}

// Size returns the number of elements in the set.
func (s *Set[T]) Size() int {
	return len(s.emap)
}

// Union returns a new set that is the union of the current set and another set.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for key := range s.emap {
		result.Add(key)
	}
	for key := range other.emap {
		result.Add(key)
	}
	return result
}

// Intersection returns a new set that is the intersection of the current set and another set.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for key := range s.emap {
		if other.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

// Difference returns a new set that is the difference of the current set and another set.
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for key := range s.emap {
		if !other.Contains(key) {
			result.Add(key)
		}
	}
	return result
}

// Elements returns a slice of all elements in the set.
func (s *Set[T]) Elements() []T {
	result := make([]T, 0, len(s.emap))
	for key := range s.emap {
		result = append(result, key)
	}
	return result
}

func (s *Set[T]) Clear() {
	for key := range s.emap {
		s.Remove(key)
	}
}

func NewSetFromSlice[T comparable](slc []T) *Set[T] {
	s := NewSet[T]()
	for i := range slc {
		s.Add(slc[i])
	}
	return s
}

func NewSetFromMapKeys[T comparable, A any](m map[T]A) *Set[T] {
	s := NewSet[T]()
	for k := range m {
		s.Add(k)
	}
	return s
}
