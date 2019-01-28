package helper

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

func DescribeInstancesForTagsAndAction(svc ec2iface.EC2API, repository, branch, action string) ([]*string, error) {
	result, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
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

func StartEC2Instances(svc ec2iface.EC2API, instanceIDs []*string) error {
	log.Println("Starting EC2")
	startResult, err := svc.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Changed state from %s to %s \n", *startResult.StartingInstances[0].PreviousState.Name, *startResult.StartingInstances[0].CurrentState.Name)
	return nil
}

func StopEC2Instances(svc ec2iface.EC2API, instanceIDs []*string) error {
	log.Println("Stopping EC2")
	stopResult, err := svc.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Changed state from %s to %s \n", *stopResult.StoppingInstances[0].PreviousState.Name, *stopResult.StoppingInstances[0].CurrentState.Name)
	return nil
}
