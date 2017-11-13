// Code generated by counterfeiter. DO NOT EDIT.
package gardenerfakes

import (
	"sync"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/gardener"
	"code.cloudfoundry.org/lager"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type FakeVolumizer struct {
	CreateStub        func(log lager.Logger, spec garden.ContainerSpec) (specs.Spec, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		log  lager.Logger
		spec garden.ContainerSpec
	}
	createReturns struct {
		result1 specs.Spec
		result2 error
	}
	createReturnsOnCall map[int]struct {
		result1 specs.Spec
		result2 error
	}
	DestroyStub        func(log lager.Logger, handle string) error
	destroyMutex       sync.RWMutex
	destroyArgsForCall []struct {
		log    lager.Logger
		handle string
	}
	destroyReturns struct {
		result1 error
	}
	destroyReturnsOnCall map[int]struct {
		result1 error
	}
	MetricsStub        func(log lager.Logger, handle string, privileged bool) (garden.ContainerDiskStat, error)
	metricsMutex       sync.RWMutex
	metricsArgsForCall []struct {
		log        lager.Logger
		handle     string
		privileged bool
	}
	metricsReturns struct {
		result1 garden.ContainerDiskStat
		result2 error
	}
	metricsReturnsOnCall map[int]struct {
		result1 garden.ContainerDiskStat
		result2 error
	}
	GCStub        func(log lager.Logger) error
	gCMutex       sync.RWMutex
	gCArgsForCall []struct {
		log lager.Logger
	}
	gCReturns struct {
		result1 error
	}
	gCReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVolumizer) Create(log lager.Logger, spec garden.ContainerSpec) (specs.Spec, error) {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		log  lager.Logger
		spec garden.ContainerSpec
	}{log, spec})
	fake.recordInvocation("Create", []interface{}{log, spec})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(log, spec)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createReturns.result1, fake.createReturns.result2
}

func (fake *FakeVolumizer) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeVolumizer) CreateArgsForCall(i int) (lager.Logger, garden.ContainerSpec) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].log, fake.createArgsForCall[i].spec
}

func (fake *FakeVolumizer) CreateReturns(result1 specs.Spec, result2 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 specs.Spec
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumizer) CreateReturnsOnCall(i int, result1 specs.Spec, result2 error) {
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 specs.Spec
			result2 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 specs.Spec
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumizer) Destroy(log lager.Logger, handle string) error {
	fake.destroyMutex.Lock()
	ret, specificReturn := fake.destroyReturnsOnCall[len(fake.destroyArgsForCall)]
	fake.destroyArgsForCall = append(fake.destroyArgsForCall, struct {
		log    lager.Logger
		handle string
	}{log, handle})
	fake.recordInvocation("Destroy", []interface{}{log, handle})
	fake.destroyMutex.Unlock()
	if fake.DestroyStub != nil {
		return fake.DestroyStub(log, handle)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.destroyReturns.result1
}

func (fake *FakeVolumizer) DestroyCallCount() int {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return len(fake.destroyArgsForCall)
}

func (fake *FakeVolumizer) DestroyArgsForCall(i int) (lager.Logger, string) {
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	return fake.destroyArgsForCall[i].log, fake.destroyArgsForCall[i].handle
}

func (fake *FakeVolumizer) DestroyReturns(result1 error) {
	fake.DestroyStub = nil
	fake.destroyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeVolumizer) DestroyReturnsOnCall(i int, result1 error) {
	fake.DestroyStub = nil
	if fake.destroyReturnsOnCall == nil {
		fake.destroyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.destroyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeVolumizer) Metrics(log lager.Logger, handle string, privileged bool) (garden.ContainerDiskStat, error) {
	fake.metricsMutex.Lock()
	ret, specificReturn := fake.metricsReturnsOnCall[len(fake.metricsArgsForCall)]
	fake.metricsArgsForCall = append(fake.metricsArgsForCall, struct {
		log        lager.Logger
		handle     string
		privileged bool
	}{log, handle, privileged})
	fake.recordInvocation("Metrics", []interface{}{log, handle, privileged})
	fake.metricsMutex.Unlock()
	if fake.MetricsStub != nil {
		return fake.MetricsStub(log, handle, privileged)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.metricsReturns.result1, fake.metricsReturns.result2
}

func (fake *FakeVolumizer) MetricsCallCount() int {
	fake.metricsMutex.RLock()
	defer fake.metricsMutex.RUnlock()
	return len(fake.metricsArgsForCall)
}

func (fake *FakeVolumizer) MetricsArgsForCall(i int) (lager.Logger, string, bool) {
	fake.metricsMutex.RLock()
	defer fake.metricsMutex.RUnlock()
	return fake.metricsArgsForCall[i].log, fake.metricsArgsForCall[i].handle, fake.metricsArgsForCall[i].privileged
}

func (fake *FakeVolumizer) MetricsReturns(result1 garden.ContainerDiskStat, result2 error) {
	fake.MetricsStub = nil
	fake.metricsReturns = struct {
		result1 garden.ContainerDiskStat
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumizer) MetricsReturnsOnCall(i int, result1 garden.ContainerDiskStat, result2 error) {
	fake.MetricsStub = nil
	if fake.metricsReturnsOnCall == nil {
		fake.metricsReturnsOnCall = make(map[int]struct {
			result1 garden.ContainerDiskStat
			result2 error
		})
	}
	fake.metricsReturnsOnCall[i] = struct {
		result1 garden.ContainerDiskStat
		result2 error
	}{result1, result2}
}

func (fake *FakeVolumizer) GC(log lager.Logger) error {
	fake.gCMutex.Lock()
	ret, specificReturn := fake.gCReturnsOnCall[len(fake.gCArgsForCall)]
	fake.gCArgsForCall = append(fake.gCArgsForCall, struct {
		log lager.Logger
	}{log})
	fake.recordInvocation("GC", []interface{}{log})
	fake.gCMutex.Unlock()
	if fake.GCStub != nil {
		return fake.GCStub(log)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.gCReturns.result1
}

func (fake *FakeVolumizer) GCCallCount() int {
	fake.gCMutex.RLock()
	defer fake.gCMutex.RUnlock()
	return len(fake.gCArgsForCall)
}

func (fake *FakeVolumizer) GCArgsForCall(i int) lager.Logger {
	fake.gCMutex.RLock()
	defer fake.gCMutex.RUnlock()
	return fake.gCArgsForCall[i].log
}

func (fake *FakeVolumizer) GCReturns(result1 error) {
	fake.GCStub = nil
	fake.gCReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeVolumizer) GCReturnsOnCall(i int, result1 error) {
	fake.GCStub = nil
	if fake.gCReturnsOnCall == nil {
		fake.gCReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.gCReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeVolumizer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.destroyMutex.RLock()
	defer fake.destroyMutex.RUnlock()
	fake.metricsMutex.RLock()
	defer fake.metricsMutex.RUnlock()
	fake.gCMutex.RLock()
	defer fake.gCMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeVolumizer) recordInvocation(key string, args []interface{}) {
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

var _ gardener.Volumizer = new(FakeVolumizer)