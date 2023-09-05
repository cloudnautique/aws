package main

import (
	"github.com/acorn-io/services/aws/libs/common"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/sirupsen/logrus"
)

type MyStackProps struct {
	awscdk.StackProps
	TopicName string            `json:"topicName,omitempty"`
	UserTags  map[string]string `json:"tags,omitempty"`
}

// Need this info to get to awsiam.IPrincipal
type principalFromAcornfileJson struct {
	PrincipalType string `json:"principalType,omitempty"`
	Identity      string `json:"identity,omitempty"`
}

func NewSNSStack(scope constructs.Construct, id string, props *MyStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	topic := awssns.NewTopic(stack, jsii.String("Topic"), &awssns.TopicProps{})

	awscdk.NewCfnOutput(stack, jsii.String("TopicARN"), &awscdk.CfnOutputProps{
		Value: topic.TopicArn(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("TopicName"), &awscdk.CfnOutputProps{
		Value: topic.TopicName(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := common.NewAcornTaggedApp(nil)
	stackProps := &MyStackProps{
		StackProps: *common.NewAWSCDKStackProps(),
	}

	if err := common.NewConfig(stackProps); err != nil {
		logrus.Fatal(err)
	}

	common.AppendScopedTags(app, stackProps.UserTags)
	NewSNSStack(app, "snsStack", stackProps)

	app.Synth(nil)
}
