package helper

import (
	"log"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

type RDSHelperAPI interface {
	GetRDSClusterForTags(repository, branch string) (*string, *string, error)
	StopRDSCluster(clusterARN, clusterStatus *string) (bool, error)
	StartRDSCluster(clusterARN, clusterStatus *string) (bool, error)
}

type RDSHelper struct {
	rdsiface.RDSAPI
}

func NewRDSHelper(svc rdsiface.RDSAPI) *RDSHelper {
	return &RDSHelper{
		RDSAPI: svc,
	}
}

func (rdshelper *RDSHelper) GetRDSClusterForTags(repository, branch string) (*string, *string, error) {
	result, err := rdshelper.RDSAPI.DescribeDBClusters(nil)
	if err != nil {
		return nil, nil, err
	}

	// Check tags for each Cluster
	for i := range result.DBClusters {
		clusterARN := result.DBClusters[i].DBClusterArn
		clusterStatus := result.DBClusters[i].Status

		// Get tags for Cluster
		result, err := rdshelper.RDSAPI.ListTagsForResource(&rds.ListTagsForResourceInput{
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

func (rdshelper *RDSHelper) StopRDSCluster(clusterARN, clusterStatus *string) (bool, error) {
	if *clusterStatus == "available" {
		log.Println("Stopping RDS CLUSTER")
		_, err := rdshelper.RDSAPI.StopDBCluster(&rds.StopDBClusterInput{
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

func (rdshelper *RDSHelper) StartRDSCluster(clusterARN, clusterStatus *string) (bool, error) {
	if *clusterStatus == "stopped" {
		log.Println("Starting RDS CLUSTER")
		_, err := rdshelper.RDSAPI.StartDBCluster(&rds.StartDBClusterInput{
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
