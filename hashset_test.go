package hashset

import (
	"encoding/json"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

/*
2. slice가오면 set으로 변경되어야한다.
3. set을 slice로 변경할 수 있어야한다.
4. set(집합) 연산이 가능해야한다.
 1. intersection
 2. union
 3. difference
5. 거의 모든타입 지원이 되어야한다.
6. 같은 set에 대해 동시작업이 원자성을 보장받아야한다.
*/

// Basic operations 0. Initialize
func TestInitializeFromArguments(t *testing.T) {
	s := New(1, 2, "3")
	require.NotEmpty(t, s)
	require.Equal(t, 3, s.Len())

	s.Add("3") // It lengths is 3, because Set already has "3"
	require.NotEqual(t, 4, s.Len())

	s1 := New([]int{1, 2, 3})
	require.Equal(t, 3, s1.Len())
	require.True(t, s1.Contains(2))
}

// Basic operations: 1. Add element
func TestAddElement(t *testing.T) {
	s := New()
	require.NotEmpty(t, s)
	require.Equal(t, 0, s.Len())
	s.Add("Add")
	s.Add(" multiple")
	s.Add("values of")
	s.Add(3)
	s.Add("like this:")
	s.Add([]int{1, 2, 3})
	require.Equal(t, 6, s.Len())
}

// Basic operations: 2. Remove element
func TestRemoveElement(t *testing.T) {
	s := New("1", "2", 3)
	require.NotEmpty(t, s)
	require.Equal(t, 3, s.Len())
	s.Remove("1")
	s.Remove("multiple")
	s.Remove("3")
	require.Equal(t, 2, s.Len())
}

// Basic operations: 3. Check membership in set
func TestMembershipCheck(t *testing.T) {
	caseFunction := func() bool { return true }
	s := New("1", "2", "3", 1, caseFunction)
	require.NotEmpty(t, s)
	require.Equal(t, 5, s.Len())

	require.True(t, s.Contains("1"))
	require.True(t, s.Contains("2"))
	require.True(t, s.Contains("3"))
	require.False(t, s.Contains(2))
	require.True(t, s.Contains(caseFunction))
}

func TestDuplicate(t *testing.T) {
	// Same address testing
	caseFunction := func() bool { return true }
	caseMap := map[string]int{}
	caseSlice := []int{1, 2, 3}
	caseSliceTwo := []string{"1", "a", "b"}

	s := New(caseFunction, caseMap, caseSlice, caseSliceTwo)
	require.Equal(t, 8, s.Len())
	s.Add(caseFunction)
	s.Add(caseMap)
	s.Add(caseSlice)
	s.Add(caseSliceTwo)
	require.Equal(t, 10, s.Len())

	// However, if the same signature arguments address are different, it is not a duplicate
	caseFunction2 := func() bool { return true }
	caseMap2 := map[string]int{}
	caseSlice2 := []int{1, 2, 3}
	caseSliceTwo2 := []string{"1", "a", "b"}
	s.Add(caseFunction2)
	s.Add(caseMap2)
	s.Add(caseSlice2)
	s.Add(caseSliceTwo2)
	// require.NotEqual(t, 4, s.Len())
	require.Equal(t, 14, s.Len())

	// Comparable types are duplicated checked.
	s.Add(1)
	s.Add(2)
	s.Add("caseSlice2")
	s.Add("caseSliceTwo2")
	require.Equal(t, 16, s.Len())
	s.Add(1)
	s.Add(2)
	s.Add("caseSlice2")
	s.Add("caseSliceTwo2")
	require.Equal(t, 16, s.Len())
}

func TestConvertToSet(t *testing.T) {
	caseSlice := []int{1, 2, 3}
	caseSliceTwo := []string{"1", "a", "b"}
	s := New(caseSlice, caseSliceTwo)
	require.Equal(t, 6, s.Len())
	require.True(t, s.Contains(1))
	require.True(t, s.Contains("1"))
	require.True(t, s.Contains("b"))
}

func TestConvertToSlice(t *testing.T) {
	caseSlice := []int{1, 2, 3}
	caseSliceTwo := []string{"1", "a", "b"}
	// Splitting slices
	s := New(caseSlice, caseSliceTwo)
	require.Equal(t, 6, s.Len())
	// Converting set to slice
	arr := s.ToSlice()
	require.Equal(t, 6, len(arr))
	require.True(t, reflect.ValueOf(arr).Kind() == reflect.Slice)
	require.Contains(t, arr, 1)
	require.Contains(t, arr, 3)
	require.Contains(t, arr, "1")
	require.Contains(t, arr, "b")
}

func TestUnion(t *testing.T) {
	caseSlice := []int{1, 2, 3}
	caseSliceTwo := []int{3, 4, 5}
	s1 := New(caseSlice)
	require.Equal(t, 3, s1.Len())
	s2 := New(caseSliceTwo)
	require.Equal(t, 3, s2.Len())

	// Union numbers
	s3 := s1.Union(s2)
	require.Equal(t, 5, s3.Len())
	require.True(t, s3.Contains(1))
	require.True(t, s3.Contains(5))
	require.True(t, s3.Contains(3))

	// Union other types
	caseSliceThree := []interface{}{1, 4, 5, "1", "5"}
	s4 := New(caseSliceThree)
	s5 := s3.Union(s4)
	require.Equal(t, 7, s5.Len())
	require.True(t, s5.Contains(1))
	require.True(t, s5.Contains(5))
	require.True(t, s5.Contains("5"))
}

func TestIntersection(t *testing.T) {

}
func TestDifference(t *testing.T)      {}
func TestFunctionElement(t *testing.T) {}
func TestStructElement(t *testing.T)   {}

func TestMarshalJSON(t *testing.T) {

	caseSlice := []int{1, 2, 3}
	caseSliceTwo := []string{"1", "a", "b"}
	// Splitting slices
	s := New(caseSlice, caseSliceTwo)
	ms, err := json.Marshal(s)
	require.NoError(t, err)
	s2 := New()
	err = json.Unmarshal(ms, s2)
	require.NoError(t, err)
}
func TestConcurrentAddElement10Goroutine100000Loop(t *testing.T) {
	var wg sync.WaitGroup
	s := New()
	numOfGoroutine := 10
	numOfLoop := 100000
	totalExpectElement := numOfGoroutine * numOfLoop

	for nthGoroutine := 0; nthGoroutine < numOfGoroutine; nthGoroutine++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for nthWork := 0; nthWork < numOfLoop; nthWork++ {
				s.Add(strconv.Itoa(n) + ".testing-" + strconv.Itoa(nthWork))
			}

		}(nthGoroutine)
	}
	wg.Wait()
	require.Equal(t, totalExpectElement, s.Len())
}
func TestConcurrentAddElement100000Goroutine10Loop(t *testing.T) {
	var wg sync.WaitGroup
	s := New()
	numOfGoroutine := 100000
	numOfLoop := 10
	totalExpectElement := numOfGoroutine * numOfLoop

	for nthGoroutine := 0; nthGoroutine < numOfGoroutine; nthGoroutine++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for nthWork := 0; nthWork < numOfLoop; nthWork++ {
				s.Add(strconv.Itoa(n) + ".testing-" + strconv.Itoa(nthWork))
			}

		}(nthGoroutine)
	}
	wg.Wait()
	require.Equal(t, totalExpectElement, s.Len())
}

func TestConcurrentExclusiveLock4Goroutine100000Loop(t *testing.T) {

	var wg sync.WaitGroup

	s := New()
	numOfGoroutine := 4
	numOfLoop := 100000
	for idx := 0; idx < numOfLoop; idx++ {
		for nthGoroutine := 0; nthGoroutine < numOfGoroutine; nthGoroutine++ {
			key := strconv.Itoa(idx) + ".testing-" + strconv.Itoa(nthGoroutine/2)
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				if n%2 == 0 {
					for s.Contains(key) {
						time.Sleep(time.Nanosecond)
					}
					s.Add(key)
				} else {
					for !s.Contains(key) {
						time.Sleep(time.Nanosecond)
					}
					s.Remove(key)
				}
			}(nthGoroutine)
		}
		wg.Wait()
	}
	wg.Wait()
	require.Equal(t, 0, s.Len())
}

func TestConcurrentRemoveElement(t *testing.T) {
	var wg sync.WaitGroup

	s := New()
	numOfGoroutine := 1000
	numOfLoop := 10
	totalExpectElement := numOfGoroutine * numOfLoop

	for nthGoroutine := 0; nthGoroutine < numOfGoroutine; nthGoroutine++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for nthWork := 0; nthWork < numOfLoop; nthWork++ {
				s.Add(strconv.Itoa(n) + ".testing-" + strconv.Itoa(nthWork))
			}

		}(nthGoroutine)
	}
	wg.Wait()
	require.Equal(t, totalExpectElement, s.Len())

	for nthGoroutine := 0; nthGoroutine < numOfGoroutine; nthGoroutine++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for nthWork := 0; nthWork < numOfLoop; nthWork++ {
				s.Remove(strconv.Itoa(n) + ".testing-" + strconv.Itoa(nthWork))
			}

		}(nthGoroutine)
	}
	wg.Wait()
	require.Equal(t, 0, s.Len())
}

func TestConcurrentMembershipCheck(t *testing.T) {
	var wg sync.WaitGroup

	s := New()
	numOfGoroutine := 10000
	numOfLoop := 10
	totalExpectElement := numOfGoroutine * numOfLoop

	for nthGoroutine := 0; nthGoroutine < numOfGoroutine; nthGoroutine++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for nthWork := 0; nthWork < numOfLoop; nthWork++ {
				s.Add(strconv.Itoa(n) + ".testing-" + strconv.Itoa(nthWork))
			}

		}(nthGoroutine)
	}
	wg.Wait()

	require.Equal(t, totalExpectElement, s.Len())
	for nthGoroutine := 0; nthGoroutine < numOfGoroutine; nthGoroutine++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for nthWork := 0; nthWork < numOfLoop; nthWork++ {
				s.Add(strconv.Itoa(n) + ".testing-phase-two-" + strconv.Itoa(nthWork))
				s.Remove(strconv.Itoa(n) + ".testing-" + strconv.Itoa(nthWork))
			}

		}(nthGoroutine)
	}
	wg.Wait()
}
