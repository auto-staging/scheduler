package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gitlab.com/auto-staging/scheduler/helper"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"

	"gitlab.com/auto-staging/scheduler/types"

	"github.com/aws/aws-sdk-go/aws/session"
)

// Handler is the main function called by lambda.Start, it starts / stops EC2 Instances and RDS Clusters based on the information in the eventJSON.
// Since the Lambda function is invoked by CloudWatchEvents rules it uses json.RawMessage as parameter.
func Handler(eventJSON json.RawMessage) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)

	cwEvent := types.Event{}
	err := json.Unmarshal(eventJSON, &cwEvent)
	if err != nil {
		return "", err
	}

	fmt.Println(cwEvent)

	result, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:repository"),
				Values: []*string{aws.String(cwEvent.Repository)},
			},
			{
				Name:   aws.String("tag:branch_raw"),
				Values: []*string{aws.String(cwEvent.Branch)},
			},
		},
	})
	if err != nil {
		return "", err
	}

	instanceIDs := []*string{}
	for i := range result.Reservations {
		fmt.Printf("Found instance with id = %s and state = %s \n", *result.Reservations[i].Instances[0].InstanceId, *result.Reservations[i].Instances[0].State.Name)
		if *result.Reservations[i].Instances[0].State.Name == "running" && cwEvent.Action == "stop" {
			instanceIDs = append(instanceIDs, result.Reservations[i].Instances[0].InstanceId)
		}
		if *result.Reservations[i].Instances[0].State.Name == "stopped" && cwEvent.Action == "start" {
			instanceIDs = append(instanceIDs, result.Reservations[i].Instances[0].InstanceId)
		}
	}

	if len(instanceIDs) > 0 {
		switch cwEvent.Action {
		case "stop":
			log.Println("Stopping EC2")
			stopResult, err := svc.StopInstances(&ec2.StopInstancesInput{
				InstanceIds: instanceIDs,
			})
			if err != nil {
				return "", err
			}
			log.Printf("Changed state from %s to %s \n", *stopResult.StoppingInstances[0].PreviousState.Name, *stopResult.StoppingInstances[0].CurrentState.Name)
			helper.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "stopped")

		case "start":
			log.Println("Starting EC2")
			startResult, err := svc.StartInstances(&ec2.StartInstancesInput{
				InstanceIds: instanceIDs,
			})
			if err != nil {
				return "", err
			}
			log.Printf("Changed state from %s to %s \n", *startResult.StartingInstances[0].PreviousState.Name, *startResult.StartingInstances[0].CurrentState.Name)
			helper.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "running")

		}
	} else {
		log.Println("EC2 - No action required")
	}

	svcRDS := rds.New(sess)

	resultRDS, err := svcRDS.DescribeDBClusters(nil)
	if err != nil {
		return "", err
	}

	// Check tags for each Cluster
	for i := range resultRDS.DBClusters {
		clusterARN := resultRDS.DBClusters[i].DBClusterArn
		clusterStatus := resultRDS.DBClusters[i].Status

		fmt.Println("Current cluster status = " + *clusterStatus)

		// Get tags for resource
		resultRDS, err := svcRDS.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: clusterARN,
		})
		if err != nil {
			return "", err
		}
		tagMap := map[string]string{}
		for a := range resultRDS.TagList {
			tagMap[*resultRDS.TagList[a].Key] = *resultRDS.TagList[a].Value
		}

		if tagMap["repository"] == cwEvent.Repository && tagMap["branch_raw"] == cwEvent.Branch {
			// Found matching Custer
			fmt.Printf("Found cluster %s matching the tags \n", *clusterARN)
			switch cwEvent.Action {
			case "stop":
				log.Println("Stopping RDS CLUSTER")
				if *clusterStatus == "available" {
					_, err := svcRDS.StopDBCluster(&rds.StopDBClusterInput{
						DBClusterIdentifier: clusterARN,
					})
					if err != nil {
						return "", err
					}
					helper.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "stopped")
				} else {
					log.Println("RDS - No action required")
				}

			case "start":
				if *clusterStatus == "stopped" {
					log.Println("Starting RDS CLUSTER")
					_, err := svcRDS.StartDBCluster(&rds.StartDBClusterInput{
						DBClusterIdentifier: clusterARN,
					})
					if err != nil {
						return "", err
					}
					helper.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "running")
				} else {
					log.Println("RDS - No action required")
				}

			}
		}
	}

	return "{ \"message\": \"success\" }", nil
}

func main() {
	lambda.Start(Handler)
}
