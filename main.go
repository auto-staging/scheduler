package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/auto-staging/scheduler/helper"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/auto-staging/scheduler/types"
)

type services struct {
	rdsiface.RDSAPI
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
		RDSAPI:          svcRDS,
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
	resultRDS, err := base.RDSAPI.DescribeDBClusters(nil)
	if err != nil {
		return err
	}

	// Check tags for each Cluster
	for i := range resultRDS.DBClusters {
		clusterARN := resultRDS.DBClusters[i].DBClusterArn
		clusterStatus := resultRDS.DBClusters[i].Status

		fmt.Println("Current cluster status = " + *clusterStatus)

		// Get tags for resource
		resultRDS, err := base.RDSAPI.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: clusterARN,
		})
		if err != nil {
			return err
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
					_, err := base.RDSAPI.StopDBCluster(&rds.StopDBClusterInput{
						DBClusterIdentifier: clusterARN,
					})
					if err != nil {
						return err
					}
					base.StatusHelperAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "stopped")
				} else {
					log.Println("RDS - No action required")
				}

			case "start":
				if *clusterStatus == "stopped" {
					log.Println("Starting RDS CLUSTER")
					_, err := base.RDSAPI.StartDBCluster(&rds.StartDBClusterInput{
						DBClusterIdentifier: clusterARN,
					})
					if err != nil {
						return err
					}
					base.StatusHelperAPI.SetStatusForEnvironment(cwEvent.Repository, cwEvent.Branch, "running")
				} else {
					log.Println("RDS - No action required")
				}

			}
		}
	}

	return nil
}
