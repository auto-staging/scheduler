package model

import (
	"log"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// RDSModelAPI is an interface including all RDS model functions
type RDSModelAPI interface {
	GetRDSClusterForTags(repository, branch string) (*string, *string, error)
	StopRDSCluster(clusterARN, clusterStatus *string) (bool, error)
	StartRDSCluster(clusterARN, clusterStatus *string) (bool, error)
}

// RDSModel is a struct including the AWS SDK RDS interface, all RDS model functions are called on this struct and the included AWS SDK RDS service
type RDSModel struct {
	rdsiface.RDSAPI
}

// NewRDSModel takes the AWS SDK RDS Interface as parameter and returns the pointer to an RDSModel struct, on which all RDS model functions can be called
func NewRDSModel(svc rdsiface.RDSAPI) *RDSModel {
	return &RDSModel{
		RDSAPI: svc,
	}
}

// GetRDSClusterForTags returns the ARN and the status of the Cluster found for the given repository and branch tag values.
// If an error occurs, the error gets logged and then returned.
func (rdsmodel *RDSModel) GetRDSClusterForTags(repository, branch string) (*string, *string, error) {
	result, err := rdsmodel.RDSAPI.DescribeDBClusters(nil)
	if err != nil {
		return nil, nil, err
	}

	// Check tags for each Cluster
	for i := range result.DBClusters {
		clusterARN := result.DBClusters[i].DBClusterArn
		clusterStatus := result.DBClusters[i].Status

		// Get tags for Cluster
		result, err := rdsmodel.RDSAPI.ListTagsForResource(&rds.ListTagsForResourceInput{
			ResourceName: clusterARN,
		})
		if err != nil {
			return nil, nil, err
		}
		tagMap := map[string]string{}
		for a := range result.TagList {
			tagMap[*result.TagList[a].Key] = *result.TagList[a].Value
		}

		if tagMap["repository"] == repository && tagMap["branch_raw"] == branch {
			log.Printf("Found cluster %s matching the tags with status %s \n", *clusterARN, *clusterStatus)
			return clusterARN, clusterStatus, nil
		}
	}
	log.Println("Found no matching RDS Cluster")
	return nil, nil, nil
}

// StopRDSCluster stops the RDS Cluster for the given Cluster ARN and status. It returns true, if the state of the Cluster was changed and false if not.
// If an error occurs, the error gets logged and then returned.
func (rdsmodel *RDSModel) StopRDSCluster(clusterARN, clusterStatus *string) (bool, error) {
	if *clusterStatus == "available" {
		log.Println("Stopping RDS CLUSTER")
		_, err := rdsmodel.RDSAPI.StopDBCluster(&rds.StopDBClusterInput{
			DBClusterIdentifier: clusterARN,
		})
		if err != nil {
			log.Println(err)
			return false, err
		}
		return true, nil
	}
	log.Println("RDS - No action required")
	return false, nil
}

// StartRDSCluster starts the RDS Cluster for the given Cluster ARN and status. It returns true, if the state of the Cluster was changed and false if not.
// If an error occurs, the error gets logged and then returned.
func (rdsmodel *RDSModel) StartRDSCluster(clusterARN, clusterStatus *string) (bool, error) {
	if *clusterStatus == "stopped" {
		log.Println("Starting RDS CLUSTER")
		_, err := rdsmodel.RDSAPI.StartDBCluster(&rds.StartDBClusterInput{
			DBClusterIdentifier: clusterARN,
		})
		if err != nil {
			log.Println(err)
			return false, err
		}
		return true, nil
	}
	log.Println("RDS - No action required")
	return false, nil
}
