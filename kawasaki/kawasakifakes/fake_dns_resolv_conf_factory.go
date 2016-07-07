// This file was generated by counterfeiter
package kawasakifakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/guardian/kawasaki"
)

type FakeDnsResolvConfFactory struct {
	CreateDNSResolvConfigurerStub        func(bundlePath string, cfg kawasaki.NetworkConfig) kawasaki.DnsResolvConfigurer
	createDNSResolvConfigurerMutex       sync.RWMutex
	createDNSResolvConfigurerArgsForCall []struct {
		bundlePath string
		cfg        kawasaki.NetworkConfig
	}
	createDNSResolvConfigurerReturns struct {
		result1 kawasaki.DnsResolvConfigurer
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDnsResolvConfFactory) CreateDNSResolvConfigurer(bundlePath string, cfg kawasaki.NetworkConfig) kawasaki.DnsResolvConfigurer {
	fake.createDNSResolvConfigurerMutex.Lock()
	fake.createDNSResolvConfigurerArgsForCall = append(fake.createDNSResolvConfigurerArgsForCall, struct {
		bundlePath string
		cfg        kawasaki.NetworkConfig
	}{bundlePath, cfg})
	fake.recordInvocation("CreateDNSResolvConfigurer", []interface{}{bundlePath, cfg})
	fake.createDNSResolvConfigurerMutex.Unlock()
	if fake.CreateDNSResolvConfigurerStub != nil {
		return fake.CreateDNSResolvConfigurerStub(bundlePath, cfg)
	} else {
		return fake.createDNSResolvConfigurerReturns.result1
	}
}

func (fake *FakeDnsResolvConfFactory) CreateDNSResolvConfigurerCallCount() int {
	fake.createDNSResolvConfigurerMutex.RLock()
	defer fake.createDNSResolvConfigurerMutex.RUnlock()
	return len(fake.createDNSResolvConfigurerArgsForCall)
}

func (fake *FakeDnsResolvConfFactory) CreateDNSResolvConfigurerArgsForCall(i int) (string, kawasaki.NetworkConfig) {
	fake.createDNSResolvConfigurerMutex.RLock()
	defer fake.createDNSResolvConfigurerMutex.RUnlock()
	return fake.createDNSResolvConfigurerArgsForCall[i].bundlePath, fake.createDNSResolvConfigurerArgsForCall[i].cfg
}

func (fake *FakeDnsResolvConfFactory) CreateDNSResolvConfigurerReturns(result1 kawasaki.DnsResolvConfigurer) {
	fake.CreateDNSResolvConfigurerStub = nil
	fake.createDNSResolvConfigurerReturns = struct {
		result1 kawasaki.DnsResolvConfigurer
	}{result1}
}

func (fake *FakeDnsResolvConfFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createDNSResolvConfigurerMutex.RLock()
	defer fake.createDNSResolvConfigurerMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeDnsResolvConfFactory) recordInvocation(key string, args []interface{}) {
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

var _ kawasaki.DnsResolvConfFactory = new(FakeDnsResolvConfFactory)
