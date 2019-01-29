package helper

import (
	"log"

	"github.com/auto-staging/scheduler/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type StatusHelperAPI interface {
	SetStatusForEnvironment(repository, branch, status string) error
}

type StatusHelper struct {
	dynamodbiface.DynamoDBAPI
}

func NewStatusHelper(svc dynamodbiface.DynamoDBAPI) *StatusHelper {
	return &StatusHelper{
		DynamoDBAPI: svc,
	}
}

// SetStatusForEnvironment updates the status for the Environment given in the parameters to the status given in the parameters.
// If an error occurs the error gets logged and the returned.
func (statusHelper *StatusHelper) SetStatusForEnvironment(repository, branch, status string) error {
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
	_, err = statusHelper.DynamoDBAPI.UpdateItem(input)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
