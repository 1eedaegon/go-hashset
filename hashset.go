package hashset

import (
	"reflect"
	"sync"
)

// Set represents a thread-safe collection of unique elements.

type Set struct {
	mu   sync.RWMutex         // Guards access to the internal hash map.
	hash map[interface{}]bool // Stores elements as keys with a boolean value.
}

// New initializes and returns a new Set with optional initial elements.
func New(initialValue ...interface{}) *Set {
	s := &Set{
		hash: make(map[interface{}]bool),
	}
	for _, iv := range initialValue {
		v := reflect.ValueOf(iv)
		if v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				// TODO: Optimize destructuring slice
				s.Add(v.Index(i).Interface())
			}
		} else {
			s.Add(iv)
		}
	}
	return s
}

// Add inserts an element into the Set. If the element is already present, it does nothing.
func (s *Set) Add(element interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !IsComparable(element) {
		element = MakeComparable(element)
	}
	s.hash[element] = true
}

// Remove deletes an element from the Set. If the element is not present, it does nothing.
func (s *Set) Remove(element interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !IsComparable(element) {
		element = MakeComparable(element)
	}
	delete(s.hash, element)
}

// Contains checks if an element is present in the Set.
func (s *Set) Contains(element interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !IsComparable(element) {
		element = MakeComparable(element)
	}
	_, exists := s.hash[element]
	return exists
}

// Difference returns a new Set containing elements present in the original Set but not in the given Set.
func (s *Set) Difference(set *Set) *Set {
	diff := make(map[interface{}]bool)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, _ := range s.hash {
		if _, exists := set.hash[k]; !exists {
			diff[k] = true
		}
	}
	return &Set{hash: diff}
}

// Do applies a given function to each element of the Set. This can be used for iterating over the Set.
func (s *Set) Do(f func(interface{})) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, _ := range s.hash {
		f(k)
	}
}

// Intersection returns a new Set containing elements that are present in both Sets.
func (s *Set) Intersection(set *Set) *Set {
	intersect := make(map[interface{}]bool)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, _ := range s.hash {
		if _, exists := set.hash[k]; exists {
			intersect[k] = true
		}
	}
	return &Set{hash: intersect}
}

// Len returns the number of elements in the Set.
func (s *Set) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.hash)
}

// SubsetOf checks if the Set is a subset of the given Set.
func (s *Set) SubsetOf(set *Set) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.Len() > set.Len() {
		return false
	}
	for k, _ := range s.hash {
		if _, exists := set.hash[k]; !exists {
			return false
		}
	}
	return true
}

// Union returns a new Set containing all elements that are present in either Set.
func (s *Set) Union(set *Set) *Set {
	union := make(map[interface{}]bool)
	s.mu.Lock()
	for k, _ := range s.hash {
		union[k] = true
	}
	s.mu.Unlock()

	s.mu.Lock()
	for k, _ := range set.hash {
		union[k] = true
	}
	s.mu.Unlock()
	return &Set{hash: union}
}

// ToSlice function returns converted slice from this set
func (s *Set) ToSlice() []interface{} {
	uniTypeSlice := make([]interface{}, 0)
	s.mu.RLock()
	for key := range s.hash {
		uniTypeSlice = append(uniTypeSlice, key)
	}
	s.mu.RUnlock()
	return uniTypeSlice
}

// MakeComparable returns pointer(address) not comparable types: slice, map, function
func MakeComparable(element interface{}) interface{} {
	/*
		Not comparable types: slice, map, function
	*/
	elementType := reflect.TypeOf(element)
	switch elementType.Kind() {
	case reflect.Slice, reflect.Map, reflect.Func:
		return reflect.ValueOf(element).Pointer()
	default:
		return element
	}
}

// Powerful assertion comparable type by generic on compile time
func IsComparable[T comparable](element T) bool {
	defer func() bool {
		if r := recover(); r != nil {
			return false
		}
		return true
	}()
	return element == element
}
