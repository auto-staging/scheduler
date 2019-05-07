// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ASGModelAPI is an autogenerated mock type for the ASGModelAPI type
type ASGModelAPI struct {
	mock.Mock
}

// DescribeInstancesForTagsAndAction provides a mock function with given fields: repository, branch, action
func (_m *ASGModelAPI) DescribeInstancesForTagsAndAction(repository string, branch string, action string) ([]*string, error) {
	ret := _m.Called(repository, branch, action)

	var r0 []*string
	if rf, ok := ret.Get(0).(func(string, string, string) []*string); ok {
		r0 = rf(repository, branch, action)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(repository, branch, action)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StartASGInstances provides a mock function with given fields: instanceIDs
func (_m *ASGModelAPI) StartASGInstances(instanceIDs []*string) error {
	ret := _m.Called(instanceIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*string) error); ok {
		r0 = rf(instanceIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StopASGInstances provides a mock function with given fields: instanceIDs
func (_m *ASGModelAPI) StopASGInstances(instanceIDs []*string) error {
	ret := _m.Called(instanceIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*string) error); ok {
		r0 = rf(instanceIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}