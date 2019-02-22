package model

import (
	"errors"
	"testing"

	"github.com/auto-staging/scheduler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewStatusModel(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)

	model := NewStatusModel(svc)

	assert.NotEmpty(t, model, "Expected not empty")
	assert.Equal(t, svc, model.DynamoDBAPI, "DynamoDB service from model is not matching the one used as parameter")
}

func TestSetStatusForEnvironment(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)
	svc.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, nil)

	statusHelper := StatusModel{
		DynamoDBAPI: svc,
	}

	err := statusHelper.SetStatusForEnvironment("", "", "running")
	assert.Nil(t, err, "Expected no error")
}

func TestSetStatusForEnvironmentError(t *testing.T) {
	svc := new(mocks.DynamoDBAPI)
	svc.On("UpdateItem", mock.AnythingOfType("*dynamodb.UpdateItemInput")).Return(nil, errors.New("Test error"))

	statusHelper := StatusModel{
		DynamoDBAPI: svc,
	}

	err := statusHelper.SetStatusForEnvironment("", "", "running")
	assert.Error(t, err, "Expected error")
}
