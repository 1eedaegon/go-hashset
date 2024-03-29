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
 1. 함수타입은 별도의 uuid와 포인터로 보관되어야한다.
 2. 구조체타입도 uuid로 할까 마샬링으로 할까 했는데 uuid가 좋아보이넹

6. 같은 set에 대해 동시작업이 원자성을 보장받아야한다.
*/

// func testGongurrency(numOfGoroutine, numOfEashWork int, work func(prefix string, num int)) {
// 	var wg sync.WaitGroup
// 	for nthGoroutine := 0; nthGoroutine < numOfGoroutine; nthGoroutine++ {
// 		wg.Add(1)
// 		go func(n int) {
// 			defer wg.Done()
// 			for nthWork := 0; nthWork < numOfEashWork; nthWork++ {
// 				work("testing-", nthWork)
// 			}

//			}(nthGoroutine)
//		}
//		wg.Wait()
//	}
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

// func TestConcurrentRemoveElement(t *testing.T) {}
// func TestMembershipCheck(t *testing.T)         {}
// func TestDuplicate(t *testing.T)               {}
// func TestConvertToSet(t *testing.T)            {}
// func TestConvertToSlice(t *testing.T)          {}
// func TestUnion(t *testing.T)                   {}
// func TestIntersection(t *testing.T)            {}
// func TestDifference(t *testing.T)              {}
// func TestFunctionElement(t *testing.T)         {}
// func TestStructElement(t *testing.T)           {}
