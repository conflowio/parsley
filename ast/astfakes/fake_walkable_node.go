// Code generated by counterfeiter. DO NOT EDIT.
package astfakes

import (
	"sync"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
)

type FakeWalkableNode struct {
	TokenStub        func() string
	tokenMutex       sync.RWMutex
	tokenArgsForCall []struct{}
	tokenReturns     struct {
		result1 string
	}
	tokenReturnsOnCall map[int]struct {
		result1 string
	}
	TypeStub        func() string
	typeMutex       sync.RWMutex
	typeArgsForCall []struct{}
	typeReturns     struct {
		result1 string
	}
	typeReturnsOnCall map[int]struct {
		result1 string
	}
	ValueStub        func(ctx interface{}) (interface{}, parsley.Error)
	valueMutex       sync.RWMutex
	valueArgsForCall []struct {
		ctx interface{}
	}
	valueReturns struct {
		result1 interface{}
		result2 parsley.Error
	}
	valueReturnsOnCall map[int]struct {
		result1 interface{}
		result2 parsley.Error
	}
	PosStub        func() parsley.Pos
	posMutex       sync.RWMutex
	posArgsForCall []struct{}
	posReturns     struct {
		result1 parsley.Pos
	}
	posReturnsOnCall map[int]struct {
		result1 parsley.Pos
	}
	ReaderPosStub        func() parsley.Pos
	readerPosMutex       sync.RWMutex
	readerPosArgsForCall []struct{}
	readerPosReturns     struct {
		result1 parsley.Pos
	}
	readerPosReturnsOnCall map[int]struct {
		result1 parsley.Pos
	}
	WalkStub        func(f func(n parsley.Node) bool) bool
	walkMutex       sync.RWMutex
	walkArgsForCall []struct {
		f func(n parsley.Node) bool
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

func (fake *FakeWalkableNode) Token() string {
	fake.tokenMutex.Lock()
	ret, specificReturn := fake.tokenReturnsOnCall[len(fake.tokenArgsForCall)]
	fake.tokenArgsForCall = append(fake.tokenArgsForCall, struct{}{})
	fake.recordInvocation("Token", []interface{}{})
	fake.tokenMutex.Unlock()
	if fake.TokenStub != nil {
		return fake.TokenStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.tokenReturns.result1
}

func (fake *FakeWalkableNode) TokenCallCount() int {
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	return len(fake.tokenArgsForCall)
}

func (fake *FakeWalkableNode) TokenReturns(result1 string) {
	fake.TokenStub = nil
	fake.tokenReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeWalkableNode) TokenReturnsOnCall(i int, result1 string) {
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

func (fake *FakeWalkableNode) Type() string {
	fake.typeMutex.Lock()
	ret, specificReturn := fake.typeReturnsOnCall[len(fake.typeArgsForCall)]
	fake.typeArgsForCall = append(fake.typeArgsForCall, struct{}{})
	fake.recordInvocation("Type", []interface{}{})
	fake.typeMutex.Unlock()
	if fake.TypeStub != nil {
		return fake.TypeStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.typeReturns.result1
}

func (fake *FakeWalkableNode) TypeCallCount() int {
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	return len(fake.typeArgsForCall)
}

func (fake *FakeWalkableNode) TypeReturns(result1 string) {
	fake.TypeStub = nil
	fake.typeReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeWalkableNode) TypeReturnsOnCall(i int, result1 string) {
	fake.TypeStub = nil
	if fake.typeReturnsOnCall == nil {
		fake.typeReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.typeReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeWalkableNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	fake.valueMutex.Lock()
	ret, specificReturn := fake.valueReturnsOnCall[len(fake.valueArgsForCall)]
	fake.valueArgsForCall = append(fake.valueArgsForCall, struct {
		ctx interface{}
	}{ctx})
	fake.recordInvocation("Value", []interface{}{ctx})
	fake.valueMutex.Unlock()
	if fake.ValueStub != nil {
		return fake.ValueStub(ctx)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.valueReturns.result1, fake.valueReturns.result2
}

func (fake *FakeWalkableNode) ValueCallCount() int {
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	return len(fake.valueArgsForCall)
}

func (fake *FakeWalkableNode) ValueArgsForCall(i int) interface{} {
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	return fake.valueArgsForCall[i].ctx
}

func (fake *FakeWalkableNode) ValueReturns(result1 interface{}, result2 parsley.Error) {
	fake.ValueStub = nil
	fake.valueReturns = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeWalkableNode) ValueReturnsOnCall(i int, result1 interface{}, result2 parsley.Error) {
	fake.ValueStub = nil
	if fake.valueReturnsOnCall == nil {
		fake.valueReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 parsley.Error
		})
	}
	fake.valueReturnsOnCall[i] = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeWalkableNode) Pos() parsley.Pos {
	fake.posMutex.Lock()
	ret, specificReturn := fake.posReturnsOnCall[len(fake.posArgsForCall)]
	fake.posArgsForCall = append(fake.posArgsForCall, struct{}{})
	fake.recordInvocation("Pos", []interface{}{})
	fake.posMutex.Unlock()
	if fake.PosStub != nil {
		return fake.PosStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.posReturns.result1
}

func (fake *FakeWalkableNode) PosCallCount() int {
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	return len(fake.posArgsForCall)
}

func (fake *FakeWalkableNode) PosReturns(result1 parsley.Pos) {
	fake.PosStub = nil
	fake.posReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeWalkableNode) PosReturnsOnCall(i int, result1 parsley.Pos) {
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
	fake.readerPosArgsForCall = append(fake.readerPosArgsForCall, struct{}{})
	fake.recordInvocation("ReaderPos", []interface{}{})
	fake.readerPosMutex.Unlock()
	if fake.ReaderPosStub != nil {
		return fake.ReaderPosStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.readerPosReturns.result1
}

func (fake *FakeWalkableNode) ReaderPosCallCount() int {
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
	return len(fake.readerPosArgsForCall)
}

func (fake *FakeWalkableNode) ReaderPosReturns(result1 parsley.Pos) {
	fake.ReaderPosStub = nil
	fake.readerPosReturns = struct {
		result1 parsley.Pos
	}{result1}
}

func (fake *FakeWalkableNode) ReaderPosReturnsOnCall(i int, result1 parsley.Pos) {
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

func (fake *FakeWalkableNode) Walk(f func(n parsley.Node) bool) bool {
	fake.walkMutex.Lock()
	ret, specificReturn := fake.walkReturnsOnCall[len(fake.walkArgsForCall)]
	fake.walkArgsForCall = append(fake.walkArgsForCall, struct {
		f func(n parsley.Node) bool
	}{f})
	fake.recordInvocation("Walk", []interface{}{f})
	fake.walkMutex.Unlock()
	if fake.WalkStub != nil {
		return fake.WalkStub(f)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.walkReturns.result1
}

func (fake *FakeWalkableNode) WalkCallCount() int {
	fake.walkMutex.RLock()
	defer fake.walkMutex.RUnlock()
	return len(fake.walkArgsForCall)
}

func (fake *FakeWalkableNode) WalkArgsForCall(i int) func(n parsley.Node) bool {
	fake.walkMutex.RLock()
	defer fake.walkMutex.RUnlock()
	return fake.walkArgsForCall[i].f
}

func (fake *FakeWalkableNode) WalkReturns(result1 bool) {
	fake.WalkStub = nil
	fake.walkReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeWalkableNode) WalkReturnsOnCall(i int, result1 bool) {
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
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	fake.valueMutex.RLock()
	defer fake.valueMutex.RUnlock()
	fake.posMutex.RLock()
	defer fake.posMutex.RUnlock()
	fake.readerPosMutex.RLock()
	defer fake.readerPosMutex.RUnlock()
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

var _ ast.WalkableNode = new(FakeWalkableNode)
