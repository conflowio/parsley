// Code generated by counterfeiter. DO NOT EDIT.
package parsleyfakes

import (
	"sync"

	"github.com/conflowio/parsley/parsley"
)

type FakeWalkableNode struct {
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
	ReaderPosStub        func() parsley.Pos
	readerPosMutex       sync.RWMutex
	readerPosArgsForCall []struct {
	}
	readerPosReturns struct {
		result1 parsley.Pos
	}
	readerPosReturnsOnCall map[int]struct {
		result1 parsley.Pos
	}
	SchemaStub        func() interface{}
	schemaMutex       sync.RWMutex
	schemaArgsForCall []struct {
	}
	schemaReturns struct {
		result1 interface{}
	}
	schemaReturnsOnCall map[int]struct {
		result1 interface{}
	}
	TokenStub        func() string
	tokenMutex       sync.RWMutex
	tokenArgsForCall []struct {
	}
	tokenReturns struct {
		result1 string
	}
	tokenReturnsOnCall map[int]struct {
		result1 string
	}
	WalkStub        func(func(n parsley.Node) bool) bool
	walkMutex       sync.RWMutex
	walkArgsForCall []struct {
		arg1 func(n parsley.Node) bool
	}
	walkReturns struct {
		result1 bool
	}
	walkReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeWalkableNode) Pos() parsley.Pos {
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

func (fake *FakeWalkableNode) PosCallCount() int {
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	return len(fake.posArgsForCall)
}

func (fake *FakeWalkableNode) PosCalls(stub func() parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = stub
}

func (fake *FakeWalkableNode) PosReturns(result1 parsley.Pos) {
	fake.posMutex.Lock()
	defer fake.posMutex.Unlock()
	fake.PosStub = nil
	fake.posReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeWalkableNode) PosReturnsOnCall(i int, result1 parsley.Pos) {
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

func (fake *FakeWalkableNode) ReaderPos() parsley.Pos {
	fake.readerPosMutex.Lock()
	ret, specificReturn := fake.readerPosReturnsOnCall[len(fake.readerPosArgsForCall)]
	fake.readerPosArgsForCall = append(fake.readerPosArgsForCall, struct {
	}{})
	fake.recordInvocation("ReaderPos", []interface{}{})
	fake.readerPosMutex.Unlock()
	if fake.ReaderPosStub != nil {
		return fake.ReaderPosStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.readerPosReturns
	return fakeReturns.result1
}

func (fake *FakeWalkableNode) ReaderPosCallCount() int {
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	return len(fake.readerPosArgsForCall)
}

func (fake *FakeWalkableNode) ReaderPosCalls(stub func() parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = stub
}

func (fake *FakeWalkableNode) ReaderPosReturns(result1 parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = nil
	fake.readerPosReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeWalkableNode) ReaderPosReturnsOnCall(i int, result1 parsley.Pos) {
	fake.readerPosMutex.Lock()
	defer fake.readerPosMutex.Unlock()
	fake.ReaderPosStub = nil
	if fake.readerPosReturnsOnCall == nil {
		fake.readerPosReturnsOnCall = make(map[int]struct {
			result1 parsley.Pos
		})
	}
	fake.readerPosReturnsOnCall[i] = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeWalkableNode) Schema() interface{} {
	fake.schemaMutex.Lock()
	ret, specificReturn := fake.schemaReturnsOnCall[len(fake.schemaArgsForCall)]
	fake.schemaArgsForCall = append(fake.schemaArgsForCall, struct {
	}{})
	fake.recordInvocation("Schema", []interface{}{})
	fake.schemaMutex.Unlock()
	if fake.SchemaStub != nil {
		return fake.SchemaStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.schemaReturns
	return fakeReturns.result1
}

func (fake *FakeWalkableNode) SchemaCallCount() int {
	fake.schemaMutex.RLock()
	defer fake.schemaMutex.RUnlock()
	return len(fake.schemaArgsForCall)
}

func (fake *FakeWalkableNode) SchemaCalls(stub func() interface{}) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = stub
}

func (fake *FakeWalkableNode) SchemaReturns(result1 interface{}) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = nil
	fake.schemaReturns = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeWalkableNode) SchemaReturnsOnCall(i int, result1 interface{}) {
	fake.schemaMutex.Lock()
	defer fake.schemaMutex.Unlock()
	fake.SchemaStub = nil
	if fake.schemaReturnsOnCall == nil {
		fake.schemaReturnsOnCall = make(map[int]struct {
			result1 interface{}
		})
	}
	fake.schemaReturnsOnCall[i] = struct {
		result1 interface{}
	}{result1}
}

func (fake *FakeWalkableNode) Token() string {
	fake.tokenMutex.Lock()
	ret, specificReturn := fake.tokenReturnsOnCall[len(fake.tokenArgsForCall)]
	fake.tokenArgsForCall = append(fake.tokenArgsForCall, struct {
	}{})
	fake.recordInvocation("Token", []interface{}{})
	fake.tokenMutex.Unlock()
	if fake.TokenStub != nil {
		return fake.TokenStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.tokenReturns
	return fakeReturns.result1
}

func (fake *FakeWalkableNode) TokenCallCount() int {
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	return len(fake.tokenArgsForCall)
}

func (fake *FakeWalkableNode) TokenCalls(stub func() string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = stub
}

func (fake *FakeWalkableNode) TokenReturns(result1 string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = nil
	fake.tokenReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeWalkableNode) TokenReturnsOnCall(i int, result1 string) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = nil
	if fake.tokenReturnsOnCall == nil {
		fake.tokenReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.tokenReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeWalkableNode) Walk(arg1 func(n parsley.Node) bool) bool {
	fake.walkMutex.Lock()
	ret, specificReturn := fake.walkReturnsOnCall[len(fake.walkArgsForCall)]
	fake.walkArgsForCall = append(fake.walkArgsForCall, struct {
		arg1 func(n parsley.Node) bool
	}{arg1})
	fake.recordInvocation("Walk", []interface{}{arg1})
	fake.walkMutex.Unlock()
	if fake.WalkStub != nil {
		return fake.WalkStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.walkReturns
	return fakeReturns.result1
}

func (fake *FakeWalkableNode) WalkCallCount() int {
	fake.walkMutex.RLock()
	defer fake.walkMutex.RUnlock()
	return len(fake.walkArgsForCall)
}

func (fake *FakeWalkableNode) WalkCalls(stub func(func(n parsley.Node) bool) bool) {
	fake.walkMutex.Lock()
	defer fake.walkMutex.Unlock()
	fake.WalkStub = stub
}

func (fake *FakeWalkableNode) WalkArgsForCall(i int) func(n parsley.Node) bool {
	fake.walkMutex.RLock()
	defer fake.walkMutex.RUnlock()
	argsForCall := fake.walkArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeWalkableNode) WalkReturns(result1 bool) {
	fake.walkMutex.Lock()
	defer fake.walkMutex.Unlock()
	fake.WalkStub = nil
	fake.walkReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeWalkableNode) WalkReturnsOnCall(i int, result1 bool) {
	fake.walkMutex.Lock()
	defer fake.walkMutex.Unlock()
	fake.WalkStub = nil
	if fake.walkReturnsOnCall == nil {
		fake.walkReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.walkReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeWalkableNode) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	fake.schemaMutex.RLock()
	defer fake.schemaMutex.RUnlock()
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	fake.walkMutex.RLock()
	defer fake.walkMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeWalkableNode) recordInvocation(key string, args []interface{}) {
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

var _ parsley.WalkableNode = new(FakeWalkableNode)
