package helper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// EC2HelperAPI is an interface including all EC2 helper functions
type EC2HelperAPI interface {
	DescribeInstancesForTagsAndAction(repository, branch, action string) ([]*string, error)
	StartEC2Instances(instanceIDs []*string) error
	StopEC2Instances(instanceIDs []*string) error
}

// EC2Helper is a struct including the AWS SDK EC2 interface, all EC2 Helper functions are called on this struct and the included AWS SDK EC2 service
type EC2Helper struct {
	ec2iface.EC2API
}

// NewEC2Helper takes the AWS SDK EC2 Interface as parameter and returns the pointer to an EC2Helper struct, on which all EC2 Helper functions can be called
func NewEC2Helper(svc ec2iface.EC2API) *EC2Helper {
	return &EC2Helper{
		EC2API: svc,
	}
}

// DescribeInstancesForTagsAndAction takes a repository name, a branch name and an action (which can be "start" or "stop"). The function filters all EC2 Instances by
// repository and branch_raw tag and then writes all instanceIDs of instances to the *string array, which must get adapted based on the given action.
// If an error occurs, it gets logged and then returned
func (ec2Helper *EC2Helper) DescribeInstancesForTagsAndAction(repository, branch, action string) ([]*string, error) {
	result, err := ec2Helper.EC2API.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:repository"),
				Values: []*string{aws.String(repository)},
			},
			{
				Name:   aws.String("tag:branch_raw"),
				Values: []*string{aws.String(branch)},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return []*string{}, err
	}

	instanceIDs := []*string{}
	for i := range result.Reservations {
		fmt.Printf("Found instance with id = %s and state = %s \n", *result.Reservations[i].Instances[0].InstanceId, *result.Reservations[i].Instances[0].State.Name)
		if *result.Reservations[i].Instances[0].State.Name == "running" && action == "stop" {
			instanceIDs = append(instanceIDs, result.Reservations[i].Instances[0].InstanceId)
		}
		if *result.Reservations[i].Instances[0].State.Name == "stopped" && action == "start" {
			instanceIDs = append(instanceIDs, result.Reservations[i].Instances[0].InstanceId)
		}
	}

	return instanceIDs, nil
}

// StartEC2Instances starts all EC2 instances given in the instanceIDs array by using the AWS SDK.
// If an error occurs, it gets logged and then returned
func (ec2Helper *EC2Helper) StartEC2Instances(instanceIDs []*string) error {
	log.Println("Starting EC2")
	startResult, err := ec2Helper.EC2API.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Changed state from %s to %s \n", *startResult.StartingInstances[0].PreviousState.Name, *startResult.StartingInstances[0].CurrentState.Name)
	return nil
}

// StopEC2Instances stops all EC2 instances given in the instanceIDs array by using the AWS SDK.
// If an error occurs, it gets logged and then returned
func (ec2Helper *EC2Helper) StopEC2Instances(instanceIDs []*string) error {
	log.Println("Stopping EC2")
	stopResult, err := ec2Helper.EC2API.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Changed state from %s to %s \n", *stopResult.StoppingInstances[0].PreviousState.Name, *stopResult.StoppingInstances[0].CurrentState.Name)
	return nil
}
