// Code generated by counterfeiter. DO NOT EDIT.
package clientfakes

import (
	"sync"

	"github.com/cloudfoundry/honeycomb-ginkgo-reporter/honeycomb/client"
)

type FakeClient struct {
	SendEventStub        func(data interface{}, globalTags interface{}, customTags interface{}) error
	sendEventMutex       sync.RWMutex
	sendEventArgsForCall []struct {
		data       interface{}
		globalTags interface{}
		customTags interface{}
	}
	sendEventReturns struct {
		result1 error
	}
	sendEventReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) SendEvent(data interface{}, globalTags interface{}, customTags interface{}) error {
	fake.sendEventMutex.Lock()
	ret, specificReturn := fake.sendEventReturnsOnCall[len(fake.sendEventArgsForCall)]
	fake.sendEventArgsForCall = append(fake.sendEventArgsForCall, struct {
		data       interface{}
		globalTags interface{}
		customTags interface{}
	}{data, globalTags, customTags})
	fake.recordInvocation("SendEvent", []interface{}{data, globalTags, customTags})
	fake.sendEventMutex.Unlock()
	if fake.SendEventStub != nil {
		return fake.SendEventStub(data, globalTags, customTags)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.sendEventReturns.result1
}

func (fake *FakeClient) SendEventCallCount() int {
	fake.sendEventMutex.RLock()
	defer fake.sendEventMutex.RUnlock()
	return len(fake.sendEventArgsForCall)
}

func (fake *FakeClient) SendEventArgsForCall(i int) (interface{}, interface{}, interface{}) {
	fake.sendEventMutex.RLock()
	defer fake.sendEventMutex.RUnlock()
	return fake.sendEventArgsForCall[i].data, fake.sendEventArgsForCall[i].globalTags, fake.sendEventArgsForCall[i].customTags
}

func (fake *FakeClient) SendEventReturns(result1 error) {
	fake.SendEventStub = nil
	fake.sendEventReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) SendEventReturnsOnCall(i int, result1 error) {
	fake.SendEventStub = nil
	if fake.sendEventReturnsOnCall == nil {
		fake.sendEventReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendEventReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sendEventMutex.RLock()
	defer fake.sendEventMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
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

var _ client.Client = new(FakeClient)