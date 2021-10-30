// Code generated by counterfeiter. DO NOT EDIT.
package parsleyfakes

import (
	"sync"

	"github.com/conflowio/parsley/parsley"
)

type FakeFile struct {
	LenStub        func() int
	lenMutex       sync.RWMutex
	lenArgsForCall []struct {
	}
	lenReturns struct {
		result1 int
	}
	lenReturnsOnCall map[int]struct {
		result1 int
	}
	PosStub        func(int) parsley.Pos
	posMutex       sync.RWMutex
	posArgsForCall []struct {
		arg1 int
	}
	posReturns struct {
		result1 parsley.Pos
	}
	posReturnsOnCall map[int]struct {
		result1 parsley.Pos
	}
	PositionStub        func(int) parsley.Position
	positionMutex       sync.RWMutex
	positionArgsForCall []struct {
		arg1 int
	}
	positionReturns struct {
		result1 parsley.Position
	}
	positionReturnsOnCall map[int]struct {
		result1 parsley.Position
	}
	SetOffsetStub        func(int)
	setOffsetMutex       sync.RWMutex
	setOffsetArgsForCall []struct {
		arg1 int
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeFile) Len() int {
	fake.lenMutex.Lock()
	ret, specificReturn := fake.lenReturnsOnCall[len(fake.lenArgsForCall)]
	fake.lenArgsForCall = append(fake.lenArgsForCall, struct {
	}{})
	stub := fake.LenStub
	fakeReturns := fake.lenReturns
	fake.recordInvocation("Len", []interface{}{})
	fake.lenMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeFile) LenCallCount() int {
	fake.lenMutex.RLock()
	defer fake.lenMutex.RUnlock()
	return len(fake.lenArgsForCall)
}

func (fake *FakeFile) LenCalls(stub func() int) {
	fake.lenMutex.Lock()
	defer fake.lenMutex.Unlock()
	fake.LenStub = stub
}

func (fake *FakeFile) LenReturns(result1 int) {
	fake.lenMutex.Lock()
	defer fake.lenMutex.Unlock()
	fake.LenStub = nil
	fake.lenReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeFile) LenReturnsOnCall(i int, result1 int) {
	fake.lenMutex.Lock()
	defer fake.lenMutex.Unlock()
	fake.LenStub = nil
	if fake.lenReturnsOnCall == nil {
		fake.lenReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.lenReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeFile) Pos(arg1 int) parsley.Pos {
	fake.posMutex.Lock()
	ret, specificReturn := fake.posReturnsOnCall[len(fake.posArgsForCall)]
	fake.posArgsForCall = append(fake.posArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.PosStub
	fakeReturns := fake.posReturns
	fake.recordInvocation("Pos", []interface{}{arg1})
	fake.posMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeFile) PosCallCount() int {
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	return len(fake.posArgsForCall)
}

func (fake *FakeFile) PosCalls(stub func(int) parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = stub
}

func (fake *FakeFile) PosArgsForCall(i int) int {
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	argsForCall := fake.posArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeFile) PosReturns(result1 parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = nil
	fake.posReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeFile) PosReturnsOnCall(i int, result1 parsley.Pos) {
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

func (fake *FakeFile) Position(arg1 int) parsley.Position {
	fake.positionMutex.Lock()
	ret, specificReturn := fake.positionReturnsOnCall[len(fake.positionArgsForCall)]
	fake.positionArgsForCall = append(fake.positionArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.PositionStub
	fakeReturns := fake.positionReturns
	fake.recordInvocation("Position", []interface{}{arg1})
	fake.positionMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeFile) PositionCallCount() int {
	fake.positionMutex.RLock()
	defer fake.positionMutex.RUnlock()
	return len(fake.positionArgsForCall)
}

func (fake *FakeFile) PositionCalls(stub func(int) parsley.Position) {
	fake.positionMutex.Lock()
	defer fake.positionMutex.Unlock()
	fake.PositionStub = stub
}

func (fake *FakeFile) PositionArgsForCall(i int) int {
	fake.positionMutex.RLock()
	defer fake.positionMutex.RUnlock()
	argsForCall := fake.positionArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeFile) PositionReturns(result1 parsley.Position) {
	fake.positionMutex.Lock()
	defer fake.positionMutex.Unlock()
	fake.PositionStub = nil
	fake.positionReturns = struct {
		result1 parsley.Position
	}{result1}
}

func (fake *FakeFile) PositionReturnsOnCall(i int, result1 parsley.Position) {
	fake.positionMutex.Lock()
	defer fake.positionMutex.Unlock()
	fake.PositionStub = nil
	if fake.positionReturnsOnCall == nil {
		fake.positionReturnsOnCall = make(map[int]struct {
			result1 parsley.Position
		})
	}
	fake.positionReturnsOnCall[i] = struct {
		result1 parsley.Position
	}{result1}
}

func (fake *FakeFile) SetOffset(arg1 int) {
	fake.setOffsetMutex.Lock()
	fake.setOffsetArgsForCall = append(fake.setOffsetArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.SetOffsetStub
	fake.recordInvocation("SetOffset", []interface{}{arg1})
	fake.setOffsetMutex.Unlock()
	if stub != nil {
		fake.SetOffsetStub(arg1)
	}
}

func (fake *FakeFile) SetOffsetCallCount() int {
	fake.setOffsetMutex.RLock()
	defer fake.setOffsetMutex.RUnlock()
	return len(fake.setOffsetArgsForCall)
}

func (fake *FakeFile) SetOffsetCalls(stub func(int)) {
	fake.setOffsetMutex.Lock()
	defer fake.setOffsetMutex.Unlock()
	fake.SetOffsetStub = stub
}

func (fake *FakeFile) SetOffsetArgsForCall(i int) int {
	fake.setOffsetMutex.RLock()
	defer fake.setOffsetMutex.RUnlock()
	argsForCall := fake.setOffsetArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeFile) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.lenMutex.RLock()
	defer fake.lenMutex.RUnlock()
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	fake.positionMutex.RLock()
	defer fake.positionMutex.RUnlock()
	fake.setOffsetMutex.RLock()
	defer fake.setOffsetMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeFile) recordInvocation(key string, args []interface{}) {
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

var _ parsley.File = new(FakeFile)
