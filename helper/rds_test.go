package helper

import (
	"errors"
	"testing"

	"github.com/auto-staging/scheduler/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRDSClusterForTags(t *testing.T) {
	clusterArn := aws.String("arn:aws:rds:eu-west-1:123456789012:db:mysql-db")
	clusterStatus := aws.String("available")

	svc := new(mocks.RDSAPI)
	svc.On("DescribeDBClusters", mock.Anything).Return(&rds.DescribeDBClustersOutput{
		DBClusters: []*rds.DBCluster{
			&rds.DBCluster{
				DBClusterArn: clusterArn,
				Status:       clusterStatus,
			},
		},
	}, nil)

	svc.On("ListTagsForResource", mock.AnythingOfType("*rds.ListTagsForResourceInput")).Return(&rds.ListTagsForResourceOutput{
		TagList: []*rds.Tag{
			&rds.Tag{
				Key:   aws.String("repository"),
				Value: aws.String("repo"),
			},
			&rds.Tag{
				Key:   aws.String("branch_raw"),
				Value: aws.String("branch"),
			},
		},
	}, nil)

	rdsHelper := RDSHelper{
		RDSAPI: svc,
	}

	resultArn, resultStatus, err := rdsHelper.GetRDSClusterForTags("repo", "branch")
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, resultArn, clusterArn, "Expected defined clusterARN")
	assert.Equal(t, resultStatus, clusterStatus, "Expected defined clusterStatus")
}

func TestGetRDSClusterForTagsNoCluster(t *testing.T) {
	clusterArn := aws.String("arn:aws:rds:eu-west-1:123456789012:db:mysql-db")
	clusterStatus := aws.String("available")

	svc := new(mocks.RDSAPI)
	svc.On("DescribeDBClusters", mock.Anything).Return(&rds.DescribeDBClustersOutput{
		DBClusters: []*rds.DBCluster{
			&rds.DBCluster{
				DBClusterArn: clusterArn,
				Status:       clusterStatus,
			},
		},
	}, nil)

	svc.On("ListTagsForResource", mock.AnythingOfType("*rds.ListTagsForResourceInput")).Return(&rds.ListTagsForResourceOutput{
		TagList: []*rds.Tag{
			&rds.Tag{
				Key:   aws.String("repository"),
				Value: aws.String("repo"),
			},
			&rds.Tag{
				Key:   aws.String("branch_raw"),
				Value: aws.String("branch"),
			},
		},
	}, nil)

	rdsHelper := RDSHelper{
		RDSAPI: svc,
	}

	resultArn, resultStatus, err := rdsHelper.GetRDSClusterForTags("repo", "no_branch")
	assert.Nil(t, err, "Expected no error")
	assert.Nil(t, resultArn, "Expected resultArn to be empty")
	assert.Nil(t, resultStatus, "Expected resultStatus to be empty")
}

func TestGetRDSClusterForTagsDescribeError(t *testing.T) {
	errorMsg := errors.New("Test error")

	svc := new(mocks.RDSAPI)
	svc.On("DescribeDBClusters", mock.Anything).Return(&rds.DescribeDBClustersOutput{}, errorMsg)

	rdsHelper := RDSHelper{
		RDSAPI: svc,
	}

	resultArn, resultStatus, err := rdsHelper.GetRDSClusterForTags("repo", "no_branch")
	assert.Error(t, err, "Expected error")
	assert.Equal(t, errorMsg, err, "Error message didn't match the given one")
	assert.Nil(t, resultArn, "Expected resultArn to be empty")
	assert.Nil(t, resultStatus, "Expected resultStatus to be empty")
}

func TestGetRDSClusterForTagsDescribeTagsError(t *testing.T) {
	clusterArn := aws.String("arn:aws:rds:eu-west-1:123456789012:db:mysql-db")
	clusterStatus := aws.String("available")
	errorMsg := errors.New("Test error")

	svc := new(mocks.RDSAPI)
	svc.On("DescribeDBClusters", mock.Anything).Return(&rds.DescribeDBClustersOutput{
		DBClusters: []*rds.DBCluster{
			&rds.DBCluster{
				DBClusterArn: clusterArn,
				Status:       clusterStatus,
			},
		},
	}, nil)

	svc.On("ListTagsForResource", mock.AnythingOfType("*rds.ListTagsForResourceInput")).Return(&rds.ListTagsForResourceOutput{}, errorMsg)

	rdsHelper := RDSHelper{
		RDSAPI: svc,
	}

	resultArn, resultStatus, err := rdsHelper.GetRDSClusterForTags("repo", "no_branch")
	assert.Error(t, err, "Expected error")
	assert.Equal(t, errorMsg, err, "Error message didn't match the given one")
	assert.Nil(t, resultArn, "Expected resultArn to be empty")
	assert.Nil(t, resultStatus, "Expected resultStatus to be empty")
}
