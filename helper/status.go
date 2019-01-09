package helper

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/auto-staging/scheduler/types"
)

// SetStatusForEnvironment updates the status for the Environment given in the parameters to the status given in the parameters.
// If an error occurs the error gets logged and the returned.
func SetStatusForEnvironment(repository, branch, status string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		log.Println(err)
		return err
	}

	svc := dynamodb.New(sess)
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

	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
