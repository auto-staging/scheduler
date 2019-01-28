package helper

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/auto-staging/scheduler/mocks"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDescribeInstancesStopAction(t *testing.T) {
	svc := new(mocks.EC2API)
	svc.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			&ec2.Reservation{
				Instances: []*ec2.Instance{
					&ec2.Instance{
						InstanceId: aws.String("i-1234567890abcdef0"),
						State: &ec2.InstanceState{
							Code: aws.Int64(16),
							Name: aws.String("running"),
						},
					},
				},
			},
			&ec2.Reservation{
				Instances: []*ec2.Instance{
					&ec2.Instance{
						InstanceId: aws.String("i-1234567890abcdef1"),
						State: &ec2.InstanceState{
							Code: aws.Int64(80),
							Name: aws.String("stopped"),
						},
					},
				},
			},
		},
	}, nil)

	result, err := DescribeInstancesForTagsAndAction(svc, "", "", "stop")
	assert.Nil(t, err, "Expected no error")
	assert.Len(t, result, 1, "Expect one instance")
	assert.Equal(t, *result[0], "i-1234567890abcdef0", "Expected i-1234567890abcdef0")
}

func TestDescribeInstancesStartAction(t *testing.T) {
	svc := new(mocks.EC2API)
	svc.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			&ec2.Reservation{
				Instances: []*ec2.Instance{
					&ec2.Instance{
						InstanceId: aws.String("i-1234567890abcdef0"),
						State: &ec2.InstanceState{
							Code: aws.Int64(16),
							Name: aws.String("running"),
						},
					},
				},
			},
			&ec2.Reservation{
				Instances: []*ec2.Instance{
					&ec2.Instance{
						InstanceId: aws.String("i-1234567890abcdef1"),
						State: &ec2.InstanceState{
							Code: aws.Int64(80),
							Name: aws.String("stopped"),
						},
					},
				},
			},
		},
	}, nil)

	result, err := DescribeInstancesForTagsAndAction(svc, "", "", "start")
	assert.Nil(t, err, "Expected no error")
	assert.Len(t, result, 1, "Expect one instance")
	assert.Equal(t, *result[0], "i-1234567890abcdef1", "Expected i-1234567890abcdef1")
}

func TestDescribeInstancesError(t *testing.T) {
	svc := new(mocks.EC2API)
	svc.On("DescribeInstances", mock.AnythingOfType("*ec2.DescribeInstancesInput")).Return(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{},
	}, errors.New("Test error"))

	result, err := DescribeInstancesForTagsAndAction(svc, "", "", "start")
	assert.Error(t, err, "Expected error")
	assert.Len(t, result, 0, "Expected no instance")
}

func TestStartEC2Instances(t *testing.T) {
	svc := new(mocks.EC2API)
	svc.On("StartInstances", mock.AnythingOfType("*ec2.StartInstancesInput")).Return(&ec2.StartInstancesOutput{
		StartingInstances: []*ec2.InstanceStateChange{
			&ec2.InstanceStateChange{
				CurrentState: &ec2.InstanceState{
					Code: aws.Int64(16),
					Name: aws.String("running"),
				},
				PreviousState: &ec2.InstanceState{
					Code: aws.Int64(80),
					Name: aws.String("stopped"),
				},
				InstanceId: aws.String("i-1234567890abcdef0"),
			},
		},
	}, nil)

	err := StartEC2Instances(svc, []*string{})
	assert.Nil(t, err, "Expected no error")
}

func TestStartEC2InstancesError(t *testing.T) {
	svc := new(mocks.EC2API)
	svc.On("StartInstances", mock.AnythingOfType("*ec2.StartInstancesInput")).Return(&ec2.StartInstancesOutput{}, errors.New("Test error"))

	err := StartEC2Instances(svc, []*string{})
	assert.Error(t, err, "Expected error")
}

func TestStopEC2Instances(t *testing.T) {
	svc := new(mocks.EC2API)
	svc.On("StopInstances", mock.AnythingOfType("*ec2.StopInstancesInput")).Return(&ec2.StopInstancesOutput{
		StoppingInstances: []*ec2.InstanceStateChange{
			&ec2.InstanceStateChange{
				CurrentState: &ec2.InstanceState{
					Code: aws.Int64(80),
					Name: aws.String("stopped"),
				},
				PreviousState: &ec2.InstanceState{
					Code: aws.Int64(16),
					Name: aws.String("running"),
				},
				InstanceId: aws.String("i-1234567890abcdef0"),
			},
		},
	}, nil)

	err := StopEC2Instances(svc, []*string{})
	assert.Nil(t, err, "Expected no error")
}

func TestStopEC2InstancesError(t *testing.T) {
	svc := new(mocks.EC2API)
	svc.On("StopInstances", mock.AnythingOfType("*ec2.StopInstancesInput")).Return(&ec2.StopInstancesOutput{}, errors.New("Test error"))

	err := StopEC2Instances(svc, []*string{})
	assert.Error(t, err, "Expected error")
}
