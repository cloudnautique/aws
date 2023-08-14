package cloudformation

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/sirupsen/logrus"
)

func DeployStack(c *Client, stackName, template string) error {
	if template == "" {
		return fmt.Errorf("template is empty")
	}

	logrus.Infof("Deploying: %s", stackName)

	updateStackWaiter := cloudformation.NewStackUpdateCompleteWaiter(c.client)

	stack, err := GetStack(c, stackName)
	if err != nil && stack.Exists {
		return err
	}

	changeSetOutput, err := createAndWaitForChangeset(c, stack, template)
	if err != nil || changeSetOutput == nil {
		return err
	}

	// Logging
	logrus.Info("ChangeSet:")
	describeChangeSetOutput, err := c.client.DescribeChangeSet(c.ctx, &cloudformation.DescribeChangeSetInput{
		ChangeSetName: changeSetOutput.Id,
		StackName:     aws.String(stackName),
	})
	if err != nil {
		return err
	}

	for _, change := range describeChangeSetOutput.Changes {
		logrus.Infof("  %s: %s", change.ResourceChange.Action, *change.ResourceChange.LogicalResourceId)
	}

	// this is the apply
	if _, err := c.client.ExecuteChangeSet(c.ctx, &cloudformation.ExecuteChangeSetInput{
		ChangeSetName: changeSetOutput.Id,
		StackName:     aws.String(stackName),
	}); err != nil {
		return err
	}

	if err := updateStackWaiter.Wait(c.ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}, time.Minute*60); err != nil {
		return err
	}

	logrus.Info("Stack Created/Updated")
	return nil
}

func createAndWaitForChangeset(c *Client, stack *CfnStack, template string) (*cloudformation.CreateChangeSetOutput, error) {
	createChangeSetWaiter := cloudformation.NewChangeSetCreateCompleteWaiter(c.client)

	changeSetType := types.ChangeSetTypeCreate
	if stack.Exists {
		changeSetType = types.ChangeSetTypeUpdate
	}

	changeSetOutput, err := c.client.CreateChangeSet(c.ctx, &cloudformation.CreateChangeSetInput{
		ChangeSetName: aws.String(fmt.Sprintf("%s-%d", stack.StackName, time.Now().Unix())),
		StackName:     aws.String(stack.StackName),
		TemplateBody:  aws.String(template),
		Capabilities: []types.Capability{
			types.CapabilityCapabilityIam,
			types.CapabilityCapabilityNamedIam,
		},
		ChangeSetType: changeSetType,
	})
	if err != nil {
		return nil, err
	}

	if err := createChangeSetWaiter.Wait(c.ctx, &cloudformation.DescribeChangeSetInput{
		ChangeSetName: changeSetOutput.Id,
		StackName:     aws.String(stack.StackName),
	}, time.Second*30); err != nil {
		output, err := c.client.DescribeChangeSet(c.ctx, &cloudformation.DescribeChangeSetInput{
			ChangeSetName: changeSetOutput.Id,
			StackName:     aws.String(stack.StackName),
		})
		if output.Status == "FAILED" && (strings.Contains(aws.ToString(output.StatusReason), "didn't contain changes") || strings.Contains(aws.ToString(output.StatusReason), "No updates are to be performed")) {
			logrus.Info("No changes to be made")
			return nil, nil
		}
		return nil, err
	}

	return changeSetOutput, nil
}
