package main

import (
	"github.com/acorn-io/services/aws/libs/common"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/sirupsen/logrus"
)

type DynamodbTableStackProps struct {
	StackProps awscdk.StackProps
	UserTags   map[string]string `json:"tags"`
	TableName  string            `json:"tableName"`
}

func NewDynamoStack(scope constructs.Construct, id string, props *DynamodbTableStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	db := awsdynamodb.NewTable(stack, jsii.String("Table"), &awsdynamodb.TableProps{
		TableName: jsii.String(props.TableName),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("TableARN"), &awscdk.CfnOutputProps{
		Value: db.TableArn(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := common.NewAcornTaggedApp(nil)
	stackProps := &DynamodbTableStackProps{
		StackProps: *common.NewAWSCDKStackProps(),
	}

	if err := common.NewConfig(&stackProps); err != nil {
		logrus.Fatal(err)
	}

	common.AppendScopedTags(app, stackProps.UserTags)

	NewDynamoStack(app, "DynamoStack", stackProps)

	app.Synth(nil)
}
