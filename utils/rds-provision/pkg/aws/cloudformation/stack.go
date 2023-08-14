package cloudformation

import (
	"strings"

	awsCfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/goformation"
	"github.com/awslabs/goformation/cloudformation"
	"github.com/sirupsen/logrus"
)

type CfnStack struct {
	StackName string
	Exists    bool
	Current   types.Stack
}

func GetStack(c *Client, stackName string) (*CfnStack, error) {
	logrus.Infof("Getting stack: %s", stackName)
	cStack := &CfnStack{
		StackName: stackName,
		Exists:    false,
	}

	s, err := c.client.DescribeStacks(c.ctx, &awsCfn.DescribeStacksInput{
		StackName: &stackName,
	})
	if err != nil && strings.Contains(err.Error(), "does not exist") {
		logrus.Infof("Stack %s does not exist", stackName)
		return cStack, err
	}

	if len(s.Stacks) > 0 {
		cStack.Exists = true
		cStack.Current = s.Stacks[0]
	}

	return cStack, err
}

func (s *CfnStack) GetCurrentTemplate(c *Client) (*cloudformation.Template, error) {
	current, err := c.client.GetTemplate(c.ctx, &awsCfn.GetTemplateInput{
		StackName: &s.StackName,
	})
	if err != nil {
		return nil, err
	}
	return ParseYAMLCFN([]byte(*current.TemplateBody))
}

func ParseYAMLCFN(template []byte) (*cloudformation.Template, error) {
	return goformation.ParseYAML(template)
}
