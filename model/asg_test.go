package model

import (
	"errors"
	"strconv"
	"testing"

	"github.com/auto-staging/scheduler/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewASGModel(t *testing.T) {
	svc := new(mocks.AutoScalingAPI)
	model := NewASGModel(svc)

	assert.NotEmpty(t, model, "Expected not empty")
	assert.Equal(t, svc, model.AutoScalingAPI, "ASG service from model is not matching the one used as parameter")
}

func TestGetPreviousMinValueOfASG(t *testing.T) {
	svc := new(mocks.AutoScalingAPI)
	model := NewASGModel(svc)

	expectedMinSize := 2

	svc.On("DescribeAutoScalingGroups", mock.AnythingOfType("*autoscaling.DescribeAutoScalingGroupsInput")).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []*autoscaling.Group{
			&autoscaling.Group{
				AutoScalingGroupName: aws.String("testASG"),
				Tags: []*autoscaling.TagDescription{
					&autoscaling.TagDescription{
						Key:   aws.String("minSize"),
						Value: aws.String(strconv.Itoa(expectedMinSize)),
					},
				},
			},
		},
	}, nil)

	min, err := model.GetPreviousMinValueOfASG(aws.String("testASG"))
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, expectedMinSize, min)
}

func TestGetPreviousMinValueOfASGNoASGFound(t *testing.T) {
	svc := new(mocks.AutoScalingAPI)
	model := NewASGModel(svc)

	svc.On("DescribeAutoScalingGroups", mock.AnythingOfType("*autoscaling.DescribeAutoScalingGroupsInput")).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []*autoscaling.Group{},
	}, nil)

	_, err := model.GetPreviousMinValueOfASG(aws.String("testASG"))
	assert.Error(t, err)
	assert.Equal(t, errors.New("found no autoscaling group for testASG"), err)
}

func TestGetPreviousMinValueOfASGNoInteger(t *testing.T) {
	svc := new(mocks.AutoScalingAPI)
	model := NewASGModel(svc)

	svc.On("DescribeAutoScalingGroups", mock.AnythingOfType("*autoscaling.DescribeAutoScalingGroupsInput")).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []*autoscaling.Group{
			&autoscaling.Group{
				AutoScalingGroupName: aws.String("testASG"),
				Tags: []*autoscaling.TagDescription{
					&autoscaling.TagDescription{
						Key:   aws.String("minSize"),
						Value: aws.String("hello"),
					},
				},
			},
		},
	}, nil)

	_, err := model.GetPreviousMinValueOfASG(aws.String("testASG"))
	assert.Error(t, err)
}

func TestGetPreviousMinValueOfASGAwsError(t *testing.T) {
	svc := new(mocks.AutoScalingAPI)
	model := NewASGModel(svc)

	svc.On("DescribeAutoScalingGroups", mock.AnythingOfType("*autoscaling.DescribeAutoScalingGroupsInput")).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []*autoscaling.Group{},
	}, errors.New("aws-error"))

	_, err := model.GetPreviousMinValueOfASG(aws.String("testASG"))
	assert.Error(t, err)
	assert.Equal(t, errors.New("aws-error"), err)
}

// SetASGMinToZero

func TestSetASGMinToZero(t *testing.T) {
	asgName := "testASG"
	svc := new(mocks.AutoScalingAPI)
	checkInput := func(input *autoscaling.UpdateAutoScalingGroupInput) error {
		if *input.AutoScalingGroupName != asgName {
			t.Error("Exptected asg name to be " + asgName + ", was " + *input.AutoScalingGroupName)
			t.FailNow()
			return errors.New("")
		}
		if *input.MinSize != 0 {
			t.Error("Exptected minSize to be 0, was " + strconv.Itoa(int(*input.MinSize)))
			t.FailNow()
			return errors.New("")
		}
		return nil
	}
	svc.On("UpdateAutoScalingGroup", mock.AnythingOfType("*autoscaling.UpdateAutoScalingGroupInput")).Return(nil, checkInput)

	model := NewASGModel(svc)
	err := model.SetASGMinToZero(aws.String("testASG"))

	assert.Nil(t, err)
}

func TestSetASGMinToZeroAwsError(t *testing.T) {
	svc := new(mocks.AutoScalingAPI)
	svc.On("UpdateAutoScalingGroup", mock.AnythingOfType("*autoscaling.UpdateAutoScalingGroupInput")).Return(nil, errors.New("aws-error"))

	model := NewASGModel(svc)
	err := model.SetASGMinToZero(aws.String("testASG"))

	assert.Error(t, err)
	assert.Equal(t, errors.New("aws-error"), err)
}

// DescribeAutoScalingGroupForTagsAndAction

func TestDescribeAutoScalingGroupForTagsAndActionAwsError(t *testing.T) {
	svc := new(mocks.AutoScalingAPI)
	svc.On("DescribeAutoScalingGroups", mock.AnythingOfType("*autoscaling.DescribeAutoScalingGroupsInput")).Return(nil, errors.New("aws-error"))

	model := NewASGModel(svc)
	asgName, err := model.DescribeAutoScalingGroupForTagsAndAction("repo", "branch", "action")

	assert.Error(t, err)
	assert.Nil(t, asgName)
	assert.Equal(t, errors.New("aws-error"), err)
}
