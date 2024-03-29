package hashset

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

/*
0. 기본 연산
 1. 추가
 2. 삭제
 3. 요소 확인

1. 값의 중복이 없어야한다.
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

// Basic operations: 3.
// func TestMembershipCheck(t *testing.T) {}

// func TestDuplicate(t *testing.T)               {}
// func TestConvertToSet(t *testing.T)            {}
// func TestConvertToSlice(t *testing.T)          {}
// func TestUnion(t *testing.T)                   {}
// func TestIntersection(t *testing.T)            {}
// func TestDifference(t *testing.T)              {}
// func TestFunctionElement(t *testing.T)         {}
// func TestStructElement(t *testing.T)           {}
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

func TestConcurrentExclusiveLock100000Loop(t *testing.T) {
	var wg sync.WaitGroup

	s := New()
	numOfGoroutine := 2
	numOfLoop := 100000
	for idx := 0; idx < numOfLoop; idx++ {
		key := strconv.Itoa(idx) + ".testing-" + strconv.Itoa(idx)
		for nthGoroutine := 0; nthGoroutine < numOfGoroutine; nthGoroutine++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				if n == 0 {
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

func TestMembershipCheck(t *testing.T) {
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
