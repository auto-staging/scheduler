package helper

import (
	"errors"
	"testing"

	"github.com/auto-staging/scheduler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStatusHelper(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)

	helper := NewStatusHelper(svc)

	assert.NotEmpty(t, helper, "Expected not empty")
	assert.Equal(t, svc, helper.DynamoDBAPI, "DynamoDB service from helper is not matching the one used as parameter")
}

func TestSetStatusForEnvironment(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)
	svc.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, nil)

	statusHelper := StatusHelper{
		DynamoDBAPI: svc,
	}

	err := statusHelper.SetStatusForEnvironment("", "", "running")
	assert.Nil(t, err, "Expected no error")
}

func TestSetStatusForEnvironmentError(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)
	svc.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, errors.New("Test error"))

	statusHelper := StatusHelper{
		DynamoDBAPI: svc,
	}

	err := statusHelper.SetStatusForEnvironment("", "", "running")
	assert.Error(t, err, "Expected error")
}
