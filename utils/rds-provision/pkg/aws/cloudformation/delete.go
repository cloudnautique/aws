package cloudformation

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/sirupsen/logrus"
)

func Delete(c *Client, stackName string) error {
	deleteStackWaiter := cloudformation.NewStackDeleteCompleteWaiter(c.client)

	logrus.Infof("Deleting stack %s", stackName)

	c.client.DeleteStack(c.ctx, &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	})

	return deleteStackWaiter.Wait(c.ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}, time.Minute*60)
}
