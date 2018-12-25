// Code generated by counterfeiter. DO NOT EDIT.
package parsleyfakes

import (
	"sync"

	"github.com/opsidian/parsley/parsley"
)

type FakeNodeTransformer struct {
	TransformNodeStub        func(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error)
	transformNodeMutex       sync.RWMutex
	transformNodeArgsForCall []struct {
		userCtx interface{}
		node    parsley.Node
	}
	transformNodeReturns struct {
		result1 parsley.Node
		result2 parsley.Error
	}
	transformNodeReturnsOnCall map[int]struct {
		result1 parsley.Node
		result2 parsley.Error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeNodeTransformer) TransformNode(userCtx interface{}, node parsley.Node) (parsley.Node, parsley.Error) {
	fake.transformNodeMutex.Lock()
	ret, specificReturn := fake.transformNodeReturnsOnCall[len(fake.transformNodeArgsForCall)]
	fake.transformNodeArgsForCall = append(fake.transformNodeArgsForCall, struct {
		userCtx interface{}
		node    parsley.Node
	}{userCtx, node})
	fake.recordInvocation("TransformNode", []interface{}{userCtx, node})
	fake.transformNodeMutex.Unlock()
	if fake.TransformNodeStub != nil {
		return fake.TransformNodeStub(userCtx, node)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.transformNodeReturns.result1, fake.transformNodeReturns.result2
}

func (fake *FakeNodeTransformer) TransformNodeCallCount() int {
	fake.transformNodeMutex.RLock()
	defer fake.transformNodeMutex.RUnlock()
	return len(fake.transformNodeArgsForCall)
}

func (fake *FakeNodeTransformer) TransformNodeArgsForCall(i int) (interface{}, parsley.Node) {
	fake.transformNodeMutex.RLock()
	defer fake.transformNodeMutex.RUnlock()
	return fake.transformNodeArgsForCall[i].userCtx, fake.transformNodeArgsForCall[i].node
}

func (fake *FakeNodeTransformer) TransformNodeReturns(result1 parsley.Node, result2 parsley.Error) {
	fake.TransformNodeStub = nil
	fake.transformNodeReturns = struct {
		result1 parsley.Node
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeNodeTransformer) TransformNodeReturnsOnCall(i int, result1 parsley.Node, result2 parsley.Error) {
	fake.TransformNodeStub = nil
	if fake.transformNodeReturnsOnCall == nil {
		fake.transformNodeReturnsOnCall = make(map[int]struct {
			result1 parsley.Node
			result2 parsley.Error
		})
	}
	fake.transformNodeReturnsOnCall[i] = struct {
		result1 parsley.Node
		result2 parsley.Error
	}{result1, result2}
}

func (fake *FakeNodeTransformer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.transformNodeMutex.RLock()
	defer fake.transformNodeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeNodeTransformer) recordInvocation(key string, args []interface{}) {
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

var _ parsley.NodeTransformer = new(FakeNodeTransformer)