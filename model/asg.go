package model

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

// ASGModelAPI is an interface including all ASG model functions
type ASGModelAPI interface {
	DescribeInstancesForTagsAndAction(repository, branch, action string) ([]*string, error)
	StartASGInstances(instanceIDs []*string) error
	StopASGInstances(instanceIDs []*string) error
}

// ASGModel is a struct including the AWS SDK ASG interface, all ASG model functions are called on this struct and the included AWS SDK ASG service
type ASGModel struct {
	autoscalingiface.AutoScalingAPI
}

// NewASGModel takes the AWS SDK ASG Interface as parameter and returns the pointer to an ASGModel struct, on which all ASG model functions can be called
func NewASGModel(svc autoscalingiface.AutoScalingAPI) *ASGModel {
	return &ASGModel{
		AutoScalingAPI: svc,
	}
}

func (asgModel *ASGModel) DescribeAutoScalingGroupsForTagsAndAction(repository, branch, action string) (*string, error) {
	asgs, _ := asgModel.AutoScalingAPI.DescribeAutoScalingGroups(nil)

	for _, asg := range asgs.AutoScalingGroups {
		foundBranch := false
		foundRepository := false
		for _, tag := range asg.Tags {
			switch *tag.Key {
			case "branch_raw":
				if *tag.Value == branch {
					foundBranch = true
				}
			case "repository":
				if *tag.Value == repository {
					foundRepository = true
				}
			}
		}
		if foundBranch && foundRepository && *asg.MinSize != 0 && action == "stop" {
			return asg.AutoScalingGroupName, nil
		}
		if foundBranch && foundRepository && *asg.MinSize == 0 && action == "start" {
			return asg.AutoScalingGroupName, nil
		}
	}

	return nil, nil
}

func (asgModel *ASGModel) SetASGMinToPreviousValue(asgName *string) error {
	log.Println("Starting ASG")
	_, err := asgModel.AutoScalingAPI.UpdateAutoScalingGroup(&autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: asgName,
		MinSize:              aws.Int64(1),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (asgModel *ASGModel) SetASGMinToZero(asgName *string) error {
	log.Println("Stopping ASG")
	_, err := asgModel.AutoScalingAPI.UpdateAutoScalingGroup(&autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: asgName,
		MinSize:              aws.Int64(0),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
