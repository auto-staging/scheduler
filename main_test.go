package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/auto-staging/scheduler/mocks"
	"github.com/auto-staging/scheduler/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//
// EC2 Tests
//

func TestChangeEC2StateStart(t *testing.T) {
	instanceIDs := []*string{
		aws.String("i-1234567890abcdef0"),
		aws.String("i-1234567890abcdef1"),
	}

	cwEvent := types.Event{
		Action:     "start",
		Branch:     "branch",
		Repository: "repo",
	}

	svcEC2HelperAPI := new(mocks.EC2HelperAPI)
	svcEC2HelperAPI.On("DescribeInstancesForTagsAndAction", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(instanceIDs, nil)

	svcEC2HelperAPI.On("StartEC2Instances", mock.AnythingOfType("[]*string")).Return(nil)

	svcStatusHelperAPI := new(mocks.StatusHelperAPI)
	svcStatusHelperAPI.On("SetStatusForEnvironment", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	base := services{
		EC2HelperAPI:    svcEC2HelperAPI,
		StatusHelperAPI: svcStatusHelperAPI,
	}

	err := base.changeEC2State(cwEvent)

	assert.Nil(t, err, "Expected no error")
	svcEC2HelperAPI.AssertCalled(t, "StartEC2Instances", instanceIDs)
	svcStatusHelperAPI.AssertCalled(t, "SetStatusForEnvironment", cwEvent.Repository, cwEvent.Branch, "running")
}

func TestChangeEC2StateStop(t *testing.T) {
	instanceIDs := []*string{
		aws.String("i-1234567890abcdef0"),
		aws.String("i-1234567890abcdef1"),
	}

	cwEvent := types.Event{
		Action:     "stop",
		Branch:     "branch",
		Repository: "repo",
	}

	svcEC2HelperAPI := new(mocks.EC2HelperAPI)
	svcEC2HelperAPI.On("DescribeInstancesForTagsAndAction", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(instanceIDs, nil)

	svcEC2HelperAPI.On("StopEC2Instances", mock.AnythingOfType("[]*string")).Return(nil)

	svcStatusHelperAPI := new(mocks.StatusHelperAPI)
	svcStatusHelperAPI.On("SetStatusForEnvironment", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	base := services{
		EC2HelperAPI:    svcEC2HelperAPI,
		StatusHelperAPI: svcStatusHelperAPI,
	}

	err := base.changeEC2State(cwEvent)

	assert.Nil(t, err, "Expected no error")
	svcEC2HelperAPI.AssertCalled(t, "StopEC2Instances", instanceIDs)
	svcStatusHelperAPI.AssertCalled(t, "SetStatusForEnvironment", cwEvent.Repository, cwEvent.Branch, "stopped")
}

func TestChangeEC2StateNoInstances(t *testing.T) {
	cwEvent := types.Event{
		Action:     "stop",
		Branch:     "branch",
		Repository: "repo",
	}

	svcEC2HelperAPI := new(mocks.EC2HelperAPI)
	svcEC2HelperAPI.On("DescribeInstancesForTagsAndAction", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]*string{}, nil)

	base := services{
		EC2HelperAPI: svcEC2HelperAPI,
	}

	err := base.changeEC2State(cwEvent)

	assert.Nil(t, err, "Expected no error")
	svcEC2HelperAPI.AssertCalled(t, "DescribeInstancesForTagsAndAction", cwEvent.Repository, cwEvent.Branch, cwEvent.Action)
}

func TestChangeEC2StateDescribeError(t *testing.T) {
	errorMsg := errors.New("Test error")
	svcEC2HelperAPI := new(mocks.EC2HelperAPI)
	svcEC2HelperAPI.On("DescribeInstancesForTagsAndAction", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return([]*string{}, errorMsg)

	base := services{
		EC2HelperAPI: svcEC2HelperAPI,
	}

	err := base.changeEC2State(types.Event{})

	assert.Error(t, err, "Expected error")
	assert.Equal(t, errorMsg, err, "Error didn't match given error")
}

func TestChangeEC2StateStopError(t *testing.T) {
	errorMsg := errors.New("Test error")

	instanceIDs := []*string{
		aws.String("i-1234567890abcdef0"),
		aws.String("i-1234567890abcdef1"),
	}

	cwEvent := types.Event{
		Action:     "stop",
		Branch:     "branch",
		Repository: "repo",
	}

	svcEC2HelperAPI := new(mocks.EC2HelperAPI)
	svcEC2HelperAPI.On("DescribeInstancesForTagsAndAction", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(instanceIDs, nil)

	svcEC2HelperAPI.On("StopEC2Instances", mock.AnythingOfType("[]*string")).Return(errorMsg)

	base := services{
		EC2HelperAPI: svcEC2HelperAPI,
	}

	err := base.changeEC2State(cwEvent)

	assert.Error(t, err, "Expected error")
	assert.Equal(t, errorMsg, err, "Error didn't match given error")
}

func TestChangeEC2StateStopStatusError(t *testing.T) {
	errorMsg := errors.New("Test error")

	instanceIDs := []*string{
		aws.String("i-1234567890abcdef0"),
		aws.String("i-1234567890abcdef1"),
	}

	cwEvent := types.Event{
		Action:     "stop",
		Branch:     "branch",
		Repository: "repo",
	}

	svcEC2HelperAPI := new(mocks.EC2HelperAPI)
	svcEC2HelperAPI.On("DescribeInstancesForTagsAndAction", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(instanceIDs, nil)

	svcEC2HelperAPI.On("StopEC2Instances", mock.AnythingOfType("[]*string")).Return(nil)

	svcStatusHelperAPI := new(mocks.StatusHelperAPI)
	svcStatusHelperAPI.On("SetStatusForEnvironment", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errorMsg)

	base := services{
		EC2HelperAPI:    svcEC2HelperAPI,
		StatusHelperAPI: svcStatusHelperAPI,
	}

	err := base.changeEC2State(cwEvent)

	assert.Error(t, err, "Expected error")
	assert.Equal(t, errorMsg, err, "Error didn't match given error")
}

func TestChangeEC2StateStartError(t *testing.T) {
	errorMsg := errors.New("Test error")

	instanceIDs := []*string{
		aws.String("i-1234567890abcdef0"),
		aws.String("i-1234567890abcdef1"),
	}

	cwEvent := types.Event{
		Action:     "start",
		Branch:     "branch",
		Repository: "repo",
	}

	svcEC2HelperAPI := new(mocks.EC2HelperAPI)
	svcEC2HelperAPI.On("DescribeInstancesForTagsAndAction", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(instanceIDs, nil)

	svcEC2HelperAPI.On("StartEC2Instances", mock.AnythingOfType("[]*string")).Return(errorMsg)

	base := services{
		EC2HelperAPI: svcEC2HelperAPI,
	}

	err := base.changeEC2State(cwEvent)

	assert.Error(t, err, "Expected error")
	assert.Equal(t, errorMsg, err, "Error didn't match given error")
}

func TestChangeEC2StateStartStatusError(t *testing.T) {
	errorMsg := errors.New("Test error")

	instanceIDs := []*string{
		aws.String("i-1234567890abcdef0"),
		aws.String("i-1234567890abcdef1"),
	}

	cwEvent := types.Event{
		Action:     "start",
		Branch:     "branch",
		Repository: "repo",
	}

	svcEC2HelperAPI := new(mocks.EC2HelperAPI)
	svcEC2HelperAPI.On("DescribeInstancesForTagsAndAction", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(instanceIDs, nil)

	svcEC2HelperAPI.On("StartEC2Instances", mock.AnythingOfType("[]*string")).Return(nil)

	svcStatusHelperAPI := new(mocks.StatusHelperAPI)
	svcStatusHelperAPI.On("SetStatusForEnvironment", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errorMsg)

	base := services{
		EC2HelperAPI:    svcEC2HelperAPI,
		StatusHelperAPI: svcStatusHelperAPI,
	}

	err := base.changeEC2State(cwEvent)

	assert.Error(t, err, "Expected error")
	assert.Equal(t, errorMsg, err, "Error didn't match given error")
}
