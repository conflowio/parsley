// Code generated by counterfeiter. DO NOT EDIT.
package parsleyfakes

import (
	"sync"

	"github.com/conflowio/parsley/parsley"
)

type FakeInterpreter struct {
	EvalStub        func(interface{}, parsley.NonTerminalNode) (interface{}, parsley.Error)
	evalMutex       sync.RWMutex
	evalArgsForCall []struct {
		arg1 interface{}
		arg2 parsley.NonTerminalNode
	}
	evalReturns struct {
		result1 interface{}
		result2 parsley.Error
	}
	evalReturnsOnCall map[int]struct {
		result1 interface{}
		result2 parsley.Error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeInterpreter) Eval(arg1 interface{}, arg2 parsley.NonTerminalNode) (interface{}, parsley.Error) {
	fake.evalMutex.Lock()
	ret, specificReturn := fake.evalReturnsOnCall[len(fake.evalArgsForCall)]
	fake.evalArgsForCall = append(fake.evalArgsForCall, struct {
		arg1 interface{}
		arg2 parsley.NonTerminalNode
	}{arg1, arg2})
	stub := fake.EvalStub
	fakeReturns := fake.evalReturns
	fake.recordInvocation("Eval", []interface{}{arg1, arg2})
	fake.evalMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeInterpreter) EvalCallCount() int {
	fake.evalMutex.RLock()
	defer fake.evalMutex.RUnlock()
	return len(fake.evalArgsForCall)
}

func (fake *FakeInterpreter) EvalCalls(stub func(interface{}, parsley.NonTerminalNode) (interface{}, parsley.Error)) {
	fake.evalMutex.Lock()
	defer fake.evalMutex.Unlock()
	fake.EvalStub = stub
}

func (fake *FakeInterpreter) EvalArgsForCall(i int) (interface{}, parsley.NonTerminalNode) {
	fake.evalMutex.RLock()
	defer fake.evalMutex.RUnlock()
	argsForCall := fake.evalArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeInterpreter) EvalReturns(result1 interface{}, result2 parsley.Error) {
	fake.evalMutex.Lock()
	defer fake.evalMutex.Unlock()
	fake.EvalStub = nil
	fake.evalReturns = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeInterpreter) EvalReturnsOnCall(i int, result1 interface{}, result2 parsley.Error) {
	fake.evalMutex.Lock()
	defer fake.evalMutex.Unlock()
	fake.EvalStub = nil
	if fake.evalReturnsOnCall == nil {
		fake.evalReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 parsley.Error
		})
	}
	fake.evalReturnsOnCall[i] = struct {
		result1 interface{}
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeInterpreter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.evalMutex.RLock()
	defer fake.evalMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeInterpreter) recordInvocation(key string, args []interface{}) {
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

var _ parsley.Interpreter = new(FakeInterpreter)
