// Code generated by mockery v2.0.4. DO NOT EDIT.

package mocks

import (
	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"

	mock "github.com/stretchr/testify/mock"
)

// Unmarshaler is an autogenerated mock type for the Unmarshaler type
type Unmarshaler struct {
	mock.Mock
}

// UnmarshalDynamoDBAttributeValue provides a mock function with given fields: _a0
func (_m *Unmarshaler) UnmarshalDynamoDBAttributeValue(_a0 *dynamodb.AttributeValue) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*dynamodb.AttributeValue) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}