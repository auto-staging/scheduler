package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/auto-staging/scheduler/helper"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/auto-staging/scheduler/types"
)

type services struct {
	helper.RDSHelperAPI
	helper.StatusHelperAPI
	helper.EC2HelperAPI
}

// Handler is the main function called by lambda.Start, it starts / stops EC2 Instances and RDS Clusters based on the information in the eventJSON.
// Since the Lambda function is invoked by CloudWatchEvents rules it uses json.RawMessage as parameter.
func Handler(eventJSON json.RawMessage) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cwEvent := types.Event{}
	err := json.Unmarshal(eventJSON, &cwEvent)
	if err != nil {
		return "", err
	}

	fmt.Println(cwEvent)

	svcEC2 := ec2.New(sess)
	svcRDS := rds.New(sess)
	svcDynamoDB := dynamodb.New(sess)

	svcBase := services{
		RDSHelperAPI:    helper.NewRDSHelper(svcRDS),
		EC2HelperAPI:    helper.NewEC2Helper(svcEC2),
		StatusHelperAPI: helper.NewStatusHelper(svcDynamoDB),
	}

	err = svcBase.changeEC2State(cwEvent)
	if err != nil {
		return "", err
	}

	err = svcBase.changeRDSState(cwEvent)
	if err != nil {
		return "", err
	}

	return "{ \"message\": \"success\" }", nil
}

func main() {
	lambda.Start(Handler)
}

func (base *services) changeEC2State(cwEvent types.Event) error {
	instanceIDs, err := base.EC2HelperAPI.DescribeInstancesForTagsAndAction(cwEvent.Repository, cwEvent.Branch, cwEvent.Action)
	if err != nil {
		return err
	}

	if len(instanceIDs) > 0 {
		switch cwEvent.Action {
		case "stop":
			err = base.EC2HelperAPI.StopEC2Instances(instanceIDs)
			if err != nil {
				return err
			}
			err = base.StatusHelperAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "stopped")
			if err != nil {
				return err
			}

		case "start":
			err = base.EC2HelperAPI.StartEC2Instances(instanceIDs)
			if err != nil {
				return err
			}
			err = base.StatusHelperAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "running")
			if err != nil {
				return err
			}

		}
	} else {
		log.Println("EC2 - No action required")
	}

	return nil
}

func (base *services) changeRDSState(cwEvent types.Event) error {
	clusterARN, clusterStatus, err := base.RDSHelperAPI.GetRDSClusterForTags(cwEvent.Repository, cwEvent.Branch)
	if err != nil {
		return err
	}
	if *clusterARN == "" {
		// No matching cluster found, nothing to do
		return nil
	}

	switch cwEvent.Action {
	case "stop":
		changed, err := base.RDSHelperAPI.StopRDSCluster(clusterARN, clusterStatus)
		if err != nil {
			return err
		}
		if changed {
			base.StatusHelperAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "stopped")
		}

	case "start":
		changed, err := base.RDSHelperAPI.StartRDSCluster(clusterARN, clusterStatus)
		if err != nil {
			return err
		}
		if changed {
			base.StatusHelperAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "running")
		}
	}

	return nil
}
