// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/guardian/rundmc/runrunc"
)

type FakeIdGetter struct {
	GetIDsStub        func(containerID string, user string) (uint32, uint32, error)
	getIDsMutex       sync.RWMutex
	getIDsArgsForCall []struct {
		containerID string
		user        string
	}
	getIDsReturns struct {
		result1 uint32
		result2 uint32
		result3 error
	}
}

func (fake *FakeIdGetter) GetIDs(containerID string, user string) (uint32, uint32, error) {
	fake.getIDsMutex.Lock()
	fake.getIDsArgsForCall = append(fake.getIDsArgsForCall, struct {
		containerID string
		user        string
	}{containerID, user})
	fake.getIDsMutex.Unlock()
	if fake.GetIDsStub != nil {
		return fake.GetIDsStub(containerID, user)
	} else {
		return fake.getIDsReturns.result1, fake.getIDsReturns.result2, fake.getIDsReturns.result3
	}
}

func (fake *FakeIdGetter) GetIDsCallCount() int {
	fake.getIDsMutex.RLock()
	defer fake.getIDsMutex.RUnlock()
	return len(fake.getIDsArgsForCall)
}

func (fake *FakeIdGetter) GetIDsArgsForCall(i int) (string, string) {
	fake.getIDsMutex.RLock()
	defer fake.getIDsMutex.RUnlock()
	return fake.getIDsArgsForCall[i].containerID, fake.getIDsArgsForCall[i].user
}

func (fake *FakeIdGetter) GetIDsReturns(result1 uint32, result2 uint32, result3 error) {
	fake.GetIDsStub = nil
	fake.getIDsReturns = struct {
		result1 uint32
		result2 uint32
		result3 error
	}{result1, result2, result3}
}

var _ runrunc.IdGetter = new(FakeIdGetter)
