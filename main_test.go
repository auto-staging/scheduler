package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/auto-staging/scheduler/mocks"
	"github.com/auto-staging/scheduler/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangeEC2State(t *testing.T) {
	svcEC2HelperAPI := new(mocks.EC2HelperAPI)
	svcEC2HelperAPI.On("DescribeInstancesForTagsAndAction", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]*string{
		aws.String("i-1234567890abcdef0"),
		aws.String("i-1234567890abcdef1"),
	}, nil)

	svcStatusHelperAPI := new(mocks.StatusHelperAPI)
	svcStatusHelperAPI.On("SetStatusForEnvironment", mock.AnythingOfType("string")).Return(nil)

	base := services{
		EC2HelperAPI:    svcEC2HelperAPI,
		StatusHelperAPI: svcStatusHelperAPI,
	}

	err := base.changeEC2State(types.Event{})

	assert.Nil(t, err, "Expected no error")
}
