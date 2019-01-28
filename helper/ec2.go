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
