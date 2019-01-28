package helper

import (
	"errors"
	"testing"

	"github.com/auto-staging/scheduler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetStatusForEnvironment(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)
	svc.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, nil)

	err := SetStatusForEnvironment(svc, "", "", "running")
	assert.Nil(t, err, "Expected no error")
}

func TestSetStatusForEnvironmentError(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)
	svc.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, errors.New("Test error"))

	err := SetStatusForEnvironment(svc, "", "", "running")
	assert.Error(t, err, "Expected error")
}
