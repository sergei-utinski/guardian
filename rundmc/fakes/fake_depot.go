// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/guardian/rundmc"
)

type FakeDepot struct {
	CreateStub        func(handle string) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		handle string
	}
	createReturns struct {
		result1 error
	}
}

func (fake *FakeDepot) Create(handle string) error {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		handle string
	}{handle})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(handle)
	} else {
		return fake.createReturns.result1
	}
}

func (fake *FakeDepot) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeDepot) CreateArgsForCall(i int) string {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].handle
}

func (fake *FakeDepot) CreateReturns(result1 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

var _ rundmc.Depot = new(FakeDepot)
