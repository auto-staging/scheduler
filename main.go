package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"

	"gitlab.com/auto-staging/scheduler/types"

	"github.com/aws/aws-sdk-go/aws/session"
)

func Handler(eventJson json.RawMessage) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)

	cwEvent := types.Event{
		Action:     "start",
		Branch:     "feat/test",
		Repository: "auto-staging-demo-app",
	}

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
		log.Fatal(err)
	}

	instanceIDs := []*string{}
	for i := range result.Reservations {
		fmt.Printf("Found instance with id = %s and state = %s \n", *result.Reservations[i].Instances[0].InstanceId, *result.Reservations[i].Instances[0].State.Name)
		instanceIDs = append(instanceIDs, result.Reservations[i].Instances[0].InstanceId)
	}

	switch cwEvent.Action {
	case "stop":
		log.Println("Stopping EC2")
		stopResult, err := svc.StopInstances(&ec2.StopInstancesInput{
			InstanceIds: instanceIDs,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(stopResult)
	case "start":
		log.Println("Starting EC2")
		startResult, err := svc.StartInstances(&ec2.StartInstancesInput{
			InstanceIds: instanceIDs,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(startResult)
	}

	svcRDS := rds.New(sess)

	resultRDS, err := svcRDS.DescribeDBClusters(nil)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
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
				if *clusterStatus != "available" {
					log.Println("Cluster must be in available state to execute stop")
					return nil
				}
				stopResult, err := svcRDS.StopDBCluster(&rds.StopDBClusterInput{
					DBClusterIdentifier: clusterARN,
				})
				if err != nil {
					log.Fatal(err)
				}
				log.Println(stopResult)

			case "start":
				if *clusterStatus != "stopped" {
					log.Println("Cluster must be in stopped state to execute start")
					return nil
				}
				log.Println("Starting RDS CLUSTER")
				stopResult, err := svcRDS.StartDBCluster(&rds.StartDBClusterInput{
					DBClusterIdentifier: clusterARN,
				})
				if err != nil {
					log.Fatal(err)
				}
				log.Println(stopResult)
			}
		}
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
