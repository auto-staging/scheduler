package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/auto-staging/scheduler/model"
	"github.com/auto-staging/scheduler/types"
)

var version string
var commitHash string
var branch string
var buildTime string

type services struct {
	model.RDSModelAPI
	model.StatusModelAPI
	model.EC2ModelAPI
	model.ASGModelAPI
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

	if cwEvent.Operation == "VERSION" {
		return returnVersionInformation()
	}

	svcEC2 := ec2.New(sess)
	svcRDS := rds.New(sess)
	svcASG := autoscaling.New(sess)
	svcDynamoDB := dynamodb.New(sess)

	svcBase := services{
		RDSModelAPI:    model.NewRDSModel(svcRDS),
		EC2ModelAPI:    model.NewEC2Model(svcEC2),
		StatusModelAPI: model.NewStatusModel(svcDynamoDB),
		ASGModelAPI:    model.NewASGModel(svcASG),
	}

	err = svcBase.changeASGState(cwEvent)
	if err != nil {
		return "", err
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
	log.Printf("version - %s | branch - %s | commit hash - %s | build time - %s \n", version, branch, commitHash, buildTime)

	lambda.Start(Handler)
}

func returnVersionInformation() (string, error) {
	componentVersion := types.SingleComponentVersion{
		Name:       "scheduler",
		Version:    version,
		CommitHash: commitHash,
		Branch:     branch,
		BuildTime:  buildTime,
	}

	body, err := json.Marshal(componentVersion)
	if err != nil {
		log.Println("Error marshaling version information, " + err.Error())
		return fmt.Sprint("{\"message\" : \"Internal server error\"}"), err
	}

	return string(body), nil
}

func (base *services) changeASGState(cwEvent types.Event) error {
	autoscalingGroup, err := base.ASGModelAPI.DescribeAutoScalingGroupForTagsAndAction(cwEvent.Repository, cwEvent.Branch, cwEvent.Action)
	if err != nil {
		return err
	}

	if autoscalingGroup != nil {
		switch cwEvent.Action {
		case "stop":
			err = base.ASGModelAPI.SetASGMinToZero(autoscalingGroup)
			if err != nil {
				return err
			}
			err = base.StatusModelAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "stopped")
			if err != nil {
				return err
			}

		case "start":
			err = base.ASGModelAPI.SetASGMinToPreviousValue(autoscalingGroup)
			if err != nil {
				return err
			}
			err = base.StatusModelAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "running")
			if err != nil {
				return err
			}

		}
	} else {
		log.Println("ASG - No action required")
	}

	return nil
}

func (base *services) changeEC2State(cwEvent types.Event) error {
	instanceIDs, err := base.EC2ModelAPI.DescribeInstancesForTagsAndAction(cwEvent.Repository, cwEvent.Branch, cwEvent.Action)
	if err != nil {
		return err
	}

	if len(instanceIDs) > 0 {
		switch cwEvent.Action {
		case "stop":
			err = base.EC2ModelAPI.StopEC2Instances(instanceIDs)
			if err != nil {
				return err
			}
			err = base.StatusModelAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "stopped")
			if err != nil {
				return err
			}

		case "start":
			err = base.EC2ModelAPI.StartEC2Instances(instanceIDs)
			if err != nil {
				return err
			}
			err = base.StatusModelAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "running")
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
	clusterARN, clusterStatus, err := base.RDSModelAPI.GetRDSClusterForTags(cwEvent.Repository, cwEvent.Branch)
	if err != nil {
		return err
	}
	if clusterARN == nil {
		// No matching cluster found, nothing to do
		return nil
	}

	switch cwEvent.Action {
	case "stop":
		changed, err := base.RDSModelAPI.StopRDSCluster(clusterARN, clusterStatus)
		if err != nil {
			return err
		}
		if changed {
			err := base.StatusModelAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "stopped")
			if err != nil {
				return err
			}
		}

	case "start":
		changed, err := base.RDSModelAPI.StartRDSCluster(clusterARN, clusterStatus)
		if err != nil {
			return err
		}
		if changed {
			err := base.StatusModelAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "running")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
