// Code generated by counterfeiter. DO NOT EDIT.
package parsleyfakes

import (
	"sync"

	"github.com/conflowio/parsley/parsley"
)

type FakeError struct {
	CauseStub        func() error
	causeMutex       sync.RWMutex
	causeArgsForCall []struct {
	}
	causeReturns struct {
		result1 error
	}
	causeReturnsOnCall map[int]struct {
		result1 error
	}
	ErrorStub        func() string
	errorMutex       sync.RWMutex
	errorArgsForCall []struct {
	}
	errorReturns struct {
		result1 string
	}
	errorReturnsOnCall map[int]struct {
		result1 string
	}
	PosStub        func() parsley.Pos
	posMutex       sync.RWMutex
	posArgsForCall []struct {
	}
	posReturns struct {
		result1 parsley.Pos
	}
	posReturnsOnCall map[int]struct {
		result1 parsley.Pos
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeError) Cause() error {
	fake.causeMutex.Lock()
	ret, specificReturn := fake.causeReturnsOnCall[len(fake.causeArgsForCall)]
	fake.causeArgsForCall = append(fake.causeArgsForCall, struct {
	}{})
	fake.recordInvocation("Cause", []interface{}{})
	fake.causeMutex.Unlock()
	if fake.CauseStub != nil {
		return fake.CauseStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.causeReturns
	return fakeReturns.result1
}

func (fake *FakeError) CauseCallCount() int {
	fake.causeMutex.RLock()
	defer fake.causeMutex.RUnlock()
	return len(fake.causeArgsForCall)
}

func (fake *FakeError) CauseCalls(stub func() error) {
	fake.causeMutex.Lock()
	defer fake.causeMutex.Unlock()
	fake.CauseStub = stub
}

func (fake *FakeError) CauseReturns(result1 error) {
	fake.causeMutex.Lock()
	defer fake.causeMutex.Unlock()
	fake.CauseStub = nil
	fake.causeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeError) CauseReturnsOnCall(i int, result1 error) {
	fake.causeMutex.Lock()
	defer fake.causeMutex.Unlock()
	fake.CauseStub = nil
	if fake.causeReturnsOnCall == nil {
		fake.causeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.causeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeError) Error() string {
	fake.errorMutex.Lock()
	ret, specificReturn := fake.errorReturnsOnCall[len(fake.errorArgsForCall)]
	fake.errorArgsForCall = append(fake.errorArgsForCall, struct {
	}{})
	fake.recordInvocation("Error", []interface{}{})
	fake.errorMutex.Unlock()
	if fake.ErrorStub != nil {
		return fake.ErrorStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.errorReturns
	return fakeReturns.result1
}

func (fake *FakeError) ErrorCallCount() int {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return len(fake.errorArgsForCall)
}

func (fake *FakeError) ErrorCalls(stub func() string) {
	fake.errorMutex.Lock()
	defer fake.errorMutex.Unlock()
	fake.ErrorStub = stub
}

func (fake *FakeError) ErrorReturns(result1 string) {
	fake.errorMutex.Lock()
	defer fake.errorMutex.Unlock()
	fake.ErrorStub = nil
	fake.errorReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeError) ErrorReturnsOnCall(i int, result1 string) {
	fake.errorMutex.Lock()
	defer fake.errorMutex.Unlock()
	fake.ErrorStub = nil
	if fake.errorReturnsOnCall == nil {
		fake.errorReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.errorReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeError) Pos() parsley.Pos {
	fake.posMutex.Lock()
	ret, specificReturn := fake.posReturnsOnCall[len(fake.posArgsForCall)]
	fake.posArgsForCall = append(fake.posArgsForCall, struct {
	}{})
	fake.recordInvocation("Pos", []interface{}{})
	fake.posMutex.Unlock()
	if fake.PosStub != nil {
		return fake.PosStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.posReturns
	return fakeReturns.result1
}

func (fake *FakeError) PosCallCount() int {
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	return len(fake.posArgsForCall)
}

func (fake *FakeError) PosCalls(stub func() parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = stub
}

func (fake *FakeError) PosReturns(result1 parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = nil
	fake.posReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeError) PosReturnsOnCall(i int, result1 parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = nil
	if fake.posReturnsOnCall == nil {
		fake.posReturnsOnCall = make(map[int]struct {
			result1 parsley.Pos
		})
	}
	fake.posReturnsOnCall[i] = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeError) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.causeMutex.RLock()
	defer fake.causeMutex.RUnlock()
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeError) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ parsley.Error = new(FakeError)
