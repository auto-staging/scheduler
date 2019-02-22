package model

import (
	"log"

	"github.com/auto-staging/scheduler/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// StatusModelAPI is an interface including all Status model functions
type StatusModelAPI interface {
	SetStatusForEnvironment(repository, branch, status string) error
}

// StatusModel is a struct including the AWS SDK DynamoDB interface, all status change functions are called on this struct and the included AWS SDK DynamoDB service
type StatusModel struct {
	dynamodbiface.DynamoDBAPI
}

// NewStatusModel takes the AWS SDK DynamoDB Interface as parameter and returns the pointer to an StatusModel struct, on which status change model functions can be called
func NewStatusModel(svc dynamodbiface.DynamoDBAPI) *StatusModel {
	return &StatusModel{
		DynamoDBAPI: svc,
	}
}

// SetStatusForEnvironment updates the status for the Environment given in the parameters to the status given in the parameters.
// If an error occurs the error gets logged and the returned.
func (statusModel *StatusModel) SetStatusForEnvironment(repository, branch, status string) error {
	updateStruct := types.StatusUpdate{
		Status: status,
	}
	update, err := dynamodbattribute.MarshalMap(updateStruct)
	if err != nil {
		log.Println(err)
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("auto-staging-environments"),
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"), // Workaround reserved keywoard issue
		},
		Key: map[string]*dynamodb.AttributeValue{
			"repository": {
				S: aws.String(repository),
			},
			"branch": {
				S: aws.String(branch),
			},
		},
		UpdateExpression:          aws.String("SET #status = :status"),
		ExpressionAttributeValues: update,
		ConditionExpression:       aws.String("attribute_exists(repository) AND attribute_exists(branch)"),
	}
	_, err = statusModel.DynamoDBAPI.UpdateItem(input)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
