// This file was generated by counterfeiter
package gcpfakes

import (
	"sync"

	"github.com/pivotal-cf/cliaas/iaas/gcp"
	compute "google.golang.org/api/compute/v1"
)

type FakeClientAPI struct {
	CreateVMStub        func(instance compute.Instance) error
	createVMMutex       sync.RWMutex
	createVMArgsForCall []struct {
		instance compute.Instance
	}
	createVMReturns struct {
		result1 error
	}
	createVMReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteVMStub        func(instanceName string) error
	deleteVMMutex       sync.RWMutex
	deleteVMArgsForCall []struct {
		instanceName string
	}
	deleteVMReturns struct {
		result1 error
	}
	deleteVMReturnsOnCall map[int]struct {
		result1 error
	}
	GetVMInfoStub        func(filter gcp.Filter) (*compute.Instance, error)
	getVMInfoMutex       sync.RWMutex
	getVMInfoArgsForCall []struct {
		filter gcp.Filter
	}
	getVMInfoReturns struct {
		result1 *compute.Instance
		result2 error
	}
	getVMInfoReturnsOnCall map[int]struct {
		result1 *compute.Instance
		result2 error
	}
	StopVMStub        func(instanceName string) error
	stopVMMutex       sync.RWMutex
	stopVMArgsForCall []struct {
		instanceName string
	}
	stopVMReturns struct {
		result1 error
	}
	stopVMReturnsOnCall map[int]struct {
		result1 error
	}
	CreateImageStub        func(tarball string) (string, error)
	createImageMutex       sync.RWMutex
	createImageArgsForCall []struct {
		tarball string
	}
	createImageReturns struct {
		result1 string
		result2 error
	}
	createImageReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClientAPI) CreateVM(instance compute.Instance) error {
	fake.createVMMutex.Lock()
	ret, specificReturn := fake.createVMReturnsOnCall[len(fake.createVMArgsForCall)]
	fake.createVMArgsForCall = append(fake.createVMArgsForCall, struct {
		instance compute.Instance
	}{instance})
	fake.recordInvocation("CreateVM", []interface{}{instance})
	fake.createVMMutex.Unlock()
	if fake.CreateVMStub != nil {
		return fake.CreateVMStub(instance)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.createVMReturns.result1
}

func (fake *FakeClientAPI) CreateVMCallCount() int {
	fake.createVMMutex.RLock()
	defer fake.createVMMutex.RUnlock()
	return len(fake.createVMArgsForCall)
}

func (fake *FakeClientAPI) CreateVMArgsForCall(i int) compute.Instance {
	fake.createVMMutex.RLock()
	defer fake.createVMMutex.RUnlock()
	return fake.createVMArgsForCall[i].instance
}

func (fake *FakeClientAPI) CreateVMReturns(result1 error) {
	fake.CreateVMStub = nil
	fake.createVMReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClientAPI) CreateVMReturnsOnCall(i int, result1 error) {
	fake.CreateVMStub = nil
	if fake.createVMReturnsOnCall == nil {
		fake.createVMReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createVMReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClientAPI) DeleteVM(instanceName string) error {
	fake.deleteVMMutex.Lock()
	ret, specificReturn := fake.deleteVMReturnsOnCall[len(fake.deleteVMArgsForCall)]
	fake.deleteVMArgsForCall = append(fake.deleteVMArgsForCall, struct {
		instanceName string
	}{instanceName})
	fake.recordInvocation("DeleteVM", []interface{}{instanceName})
	fake.deleteVMMutex.Unlock()
	if fake.DeleteVMStub != nil {
		return fake.DeleteVMStub(instanceName)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteVMReturns.result1
}

func (fake *FakeClientAPI) DeleteVMCallCount() int {
	fake.deleteVMMutex.RLock()
	defer fake.deleteVMMutex.RUnlock()
	return len(fake.deleteVMArgsForCall)
}

func (fake *FakeClientAPI) DeleteVMArgsForCall(i int) string {
	fake.deleteVMMutex.RLock()
	defer fake.deleteVMMutex.RUnlock()
	return fake.deleteVMArgsForCall[i].instanceName
}

func (fake *FakeClientAPI) DeleteVMReturns(result1 error) {
	fake.DeleteVMStub = nil
	fake.deleteVMReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClientAPI) DeleteVMReturnsOnCall(i int, result1 error) {
	fake.DeleteVMStub = nil
	if fake.deleteVMReturnsOnCall == nil {
		fake.deleteVMReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteVMReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClientAPI) GetVMInfo(filter gcp.Filter) (*compute.Instance, error) {
	fake.getVMInfoMutex.Lock()
	ret, specificReturn := fake.getVMInfoReturnsOnCall[len(fake.getVMInfoArgsForCall)]
	fake.getVMInfoArgsForCall = append(fake.getVMInfoArgsForCall, struct {
		filter gcp.Filter
	}{filter})
	fake.recordInvocation("GetVMInfo", []interface{}{filter})
	fake.getVMInfoMutex.Unlock()
	if fake.GetVMInfoStub != nil {
		return fake.GetVMInfoStub(filter)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getVMInfoReturns.result1, fake.getVMInfoReturns.result2
}

func (fake *FakeClientAPI) GetVMInfoCallCount() int {
	fake.getVMInfoMutex.RLock()
	defer fake.getVMInfoMutex.RUnlock()
	return len(fake.getVMInfoArgsForCall)
}

func (fake *FakeClientAPI) GetVMInfoArgsForCall(i int) gcp.Filter {
	fake.getVMInfoMutex.RLock()
	defer fake.getVMInfoMutex.RUnlock()
	return fake.getVMInfoArgsForCall[i].filter
}

func (fake *FakeClientAPI) GetVMInfoReturns(result1 *compute.Instance, result2 error) {
	fake.GetVMInfoStub = nil
	fake.getVMInfoReturns = struct {
		result1 *compute.Instance
		result2 error
	}{result1, result2}
}

func (fake *FakeClientAPI) GetVMInfoReturnsOnCall(i int, result1 *compute.Instance, result2 error) {
	fake.GetVMInfoStub = nil
	if fake.getVMInfoReturnsOnCall == nil {
		fake.getVMInfoReturnsOnCall = make(map[int]struct {
			result1 *compute.Instance
			result2 error
		})
	}
	fake.getVMInfoReturnsOnCall[i] = struct {
		result1 *compute.Instance
		result2 error
	}{result1, result2}
}

func (fake *FakeClientAPI) StopVM(instanceName string) error {
	fake.stopVMMutex.Lock()
	ret, specificReturn := fake.stopVMReturnsOnCall[len(fake.stopVMArgsForCall)]
	fake.stopVMArgsForCall = append(fake.stopVMArgsForCall, struct {
		instanceName string
	}{instanceName})
	fake.recordInvocation("StopVM", []interface{}{instanceName})
	fake.stopVMMutex.Unlock()
	if fake.StopVMStub != nil {
		return fake.StopVMStub(instanceName)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.stopVMReturns.result1
}

func (fake *FakeClientAPI) StopVMCallCount() int {
	fake.stopVMMutex.RLock()
	defer fake.stopVMMutex.RUnlock()
	return len(fake.stopVMArgsForCall)
}

func (fake *FakeClientAPI) StopVMArgsForCall(i int) string {
	fake.stopVMMutex.RLock()
	defer fake.stopVMMutex.RUnlock()
	return fake.stopVMArgsForCall[i].instanceName
}

func (fake *FakeClientAPI) StopVMReturns(result1 error) {
	fake.StopVMStub = nil
	fake.stopVMReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClientAPI) StopVMReturnsOnCall(i int, result1 error) {
	fake.StopVMStub = nil
	if fake.stopVMReturnsOnCall == nil {
		fake.stopVMReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.stopVMReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClientAPI) CreateImage(tarball string) (string, error) {
	fake.createImageMutex.Lock()
	ret, specificReturn := fake.createImageReturnsOnCall[len(fake.createImageArgsForCall)]
	fake.createImageArgsForCall = append(fake.createImageArgsForCall, struct {
		tarball string
	}{tarball})
	fake.recordInvocation("CreateImage", []interface{}{tarball})
	fake.createImageMutex.Unlock()
	if fake.CreateImageStub != nil {
		return fake.CreateImageStub(tarball)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createImageReturns.result1, fake.createImageReturns.result2
}

func (fake *FakeClientAPI) CreateImageCallCount() int {
	fake.createImageMutex.RLock()
	defer fake.createImageMutex.RUnlock()
	return len(fake.createImageArgsForCall)
}

func (fake *FakeClientAPI) CreateImageArgsForCall(i int) string {
	fake.createImageMutex.RLock()
	defer fake.createImageMutex.RUnlock()
	return fake.createImageArgsForCall[i].tarball
}

func (fake *FakeClientAPI) CreateImageReturns(result1 string, result2 error) {
	fake.CreateImageStub = nil
	fake.createImageReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeClientAPI) CreateImageReturnsOnCall(i int, result1 string, result2 error) {
	fake.CreateImageStub = nil
	if fake.createImageReturnsOnCall == nil {
		fake.createImageReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.createImageReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeClientAPI) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createVMMutex.RLock()
	defer fake.createVMMutex.RUnlock()
	fake.deleteVMMutex.RLock()
	defer fake.deleteVMMutex.RUnlock()
	fake.getVMInfoMutex.RLock()
	defer fake.getVMInfoMutex.RUnlock()
	fake.stopVMMutex.RLock()
	defer fake.stopVMMutex.RUnlock()
	fake.createImageMutex.RLock()
	defer fake.createImageMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeClientAPI) recordInvocation(key string, args []interface{}) {
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

var _ gcp.ClientAPI = new(FakeClientAPI)